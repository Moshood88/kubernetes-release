/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spdx

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/release-utils/util"
)

func TestBuildIDString(t *testing.T) {
	cases := []struct {
		seeds    []string
		expected string
	}{
		{[]string{"1234"}, "1234"},
		{[]string{"abc"}, "abc"},
		{[]string{"ABC"}, "ABC"},
		{[]string{"ABC", "123"}, "ABC-123"},
		{[]string{"Hello:bye", "123"}, "Hello-bye-123"},
		{[]string{"Hello^bye", "123"}, "Hellobye-123"},
		{[]string{"Hello:bye", "123", "&^%&$"}, "Hello-bye-123"},
	}
	for _, tc := range cases {
		require.Equal(t, tc.expected, buildIDString(tc.seeds...))
	}

	// If we do not pass any seeds, func should return an UUID
	// which is 36 chars long
	require.Len(t, buildIDString(), 36)

	// Same thing for only invalid chars
	require.Len(t, buildIDString("&^$&^%"), 36)
}

func TestUnitExtractTarballTmp(t *testing.T) {
	tar := writeTestTarball(t)
	require.NotNil(t, tar)
	defer os.Remove(tar.Name())

	sut := NewSPDX()
	_, err := sut.ExtractTarballTmp("lsdjkflskdjfl")
	require.NotNil(t, err)
	dir, err := sut.ExtractTarballTmp(tar.Name())
	require.Nil(t, err, "extracting file")
	defer os.RemoveAll(dir)

	require.True(t, util.Exists(filepath.Join(dir, "/text.txt")), "checking directory")
	require.True(t, util.Exists(filepath.Join(dir, "/subdir/text.txt")), "checking subdirectory")
	require.True(t, util.Exists(dir), "checking directory")

	// Check files
}

func TestReadArchiveManifest(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "sample-manifest-*.json")
	require.Nil(t, err)
	defer os.Remove(f.Name())
	require.Nil(t, os.WriteFile(
		f.Name(), []byte(sampleManifest), os.FileMode(0o644),
	), "writing test manifest file")

	sut := spdxDefaultImplementation{}
	_, err = sut.ReadArchiveManifest("laksdjlakjsdlkjsd")
	require.NotNil(t, err)
	manifest, err := sut.ReadArchiveManifest(f.Name())
	require.Nil(t, err)
	require.Equal(
		t, "386bcf5c63de46c7066c42d4ae1c38af0689836e88fed37d1dca2d484b343cf5.json",
		manifest.ConfigFilename,
	)
	require.Equal(t, 1, len(manifest.RepoTags))
	require.Equal(t, "k8s.gcr.io/kube-apiserver-amd64:v1.22.0-alpha.1", manifest.RepoTags[0])
	require.Equal(t, 3, len(manifest.LayerFiles))
	for i, fname := range []string{
		"23e140cb8e03a12cba4ac571d9a7143cf5e2e9b72de3b33ce3243b4f7ad6a188/layer.tar",
		"48dd73ececdf0f52a174ad33a469145824713bd2b73c6257ce1ba8502003ad4e/layer.tar",
		"d397673d78556210baa112013c960cb95a3fd452e5c4a2ead2b26e5a458cd87f/layer.tar",
	} {
		require.Equal(t, fname, manifest.LayerFiles[i])
	}
}

func TestPackageFromTarball(t *testing.T) {
	tar := writeTestTarball(t)
	require.NotNil(t, tar)
	defer os.Remove(tar.Name())

	sut := spdxDefaultImplementation{}
	_, err := sut.PackageFromTarball("lsdkjflksdjflk", &TarballOptions{})
	require.NotNil(t, err)
	pkg, err := sut.PackageFromTarball(tar.Name(), &TarballOptions{})
	require.Nil(t, err)
	require.NotNil(t, pkg)

	require.NotNil(t, pkg.Checksum)
	_, ok := pkg.Checksum["SHA256"]
	require.True(t, ok, "checking if sha256 checksum is set")
	_, ok = pkg.Checksum["SHA512"]
	require.True(t, ok, "checking if sha512 checksum is set")
	require.Equal(t, "5e75826e1baf84d5c5b26cc8fc3744f560ef0288c767f1cbc160124733fdc50e", pkg.Checksum["SHA256"])
	require.Equal(t, "f3b48a64a3d9db36fff10a9752dea6271725ddf125baf7026cdf09a2c352d9ff4effadb75da31e4310bc1b2513be441c86488b69d689353128f703563846c97e", pkg.Checksum["SHA512"])
}

func TestExternalDocRef(t *testing.T) {
	cases := []struct {
		DocRef    ExternalDocumentRef
		StringVal string
	}{
		{ExternalDocumentRef{ID: "", URI: "", Checksums: map[string]string{}}, ""},
		{ExternalDocumentRef{ID: "", URI: "http://example.com/", Checksums: map[string]string{"SHA256": "d3b53860aa08e5c7ea868629800eaf78856f6ef3bcd4a2f8c5c865b75f6837c8"}}, ""},
		{ExternalDocumentRef{ID: "test-id", URI: "", Checksums: map[string]string{"SHA256": "d3b53860aa08e5c7ea868629800eaf78856f6ef3bcd4a2f8c5c865b75f6837c8"}}, ""},
		{ExternalDocumentRef{ID: "test-id", URI: "http://example.com/", Checksums: map[string]string{}}, ""},
		{
			ExternalDocumentRef{
				ID: "test-id", URI: "http://example.com/", Checksums: map[string]string{"SHA256": "d3b53860aa08e5c7ea868629800eaf78856f6ef3bcd4a2f8c5c865b75f6837c8"},
			},
			"DocumentRef-test-id http://example.com/ SHA256: d3b53860aa08e5c7ea868629800eaf78856f6ef3bcd4a2f8c5c865b75f6837c8",
		},
	}
	for _, tc := range cases {
		require.Equal(t, tc.StringVal, tc.DocRef.String())
	}
}

func TestExtDocReadSourceFile(t *testing.T) {
	// Create a known testfile
	f, err := os.CreateTemp("", "")
	require.Nil(t, err)
	require.Nil(t, os.WriteFile(f.Name(), []byte("Hellow World"), os.FileMode(0o644)))
	defer os.Remove(f.Name())

	ed := ExternalDocumentRef{}
	require.NotNil(t, ed.ReadSourceFile("/kjfhg/skjdfkjh"))
	require.Nil(t, ed.ReadSourceFile(f.Name()))
	require.NotNil(t, ed.Checksums)
	require.Equal(t, len(ed.Checksums), 1)
	require.Equal(t, "5f341d31f6b6a8b15bc4e6704830bf37f99511d1", ed.Checksums["SHA1"])
}

func writeTestTarball(t *testing.T) *os.File {
	// Create a testdir
	tar, err := os.CreateTemp(os.TempDir(), "test-tar-*.tar.gz")
	require.Nil(t, err)

	tardata, err := base64.StdEncoding.DecodeString(testTar)
	require.Nil(t, err)

	reader := bytes.NewReader(tardata)
	zipreader, err := gzip.NewReader(reader)
	require.Nil(t, err)

	bindata, err := ioutil.ReadAll(zipreader)
	require.Nil(t, err)

	require.Nil(t, os.WriteFile(
		tar.Name(), bindata, os.FileMode(0o644)), "writing test tar file",
	)
	return tar
}

func TestRelationshipRender(t *testing.T) {
	host := NewPackage()
	host.BuildID("TestHost")
	peer := NewFile()
	peer.BuildID("TestPeer")
	dummyref := "SPDXRef-File-6c0c16be41af1064ee8fd2328b17a0a778dd5e52"

	cases := []struct {
		Rel      Relationship
		MustErr  bool
		Rendered string
	}{
		{
			// Relationships with a full peer object have to render
			Relationship{FullRender: false, Type: DEPENDS_ON, Peer: peer},
			false, fmt.Sprintf("Relationship: %s DEPENDS_ON %s\n", host.SPDXID(), peer.SPDXID()),
		},
		{
			// Relationships with a remote reference
			Relationship{FullRender: false, Type: DEPENDS_ON, Peer: peer, PeerExtReference: "Remote"},
			false, fmt.Sprintf("Relationship: %s DEPENDS_ON DocumentRef-Remote:%s\n", host.SPDXID(), peer.SPDXID()),
		},
		{
			// Relationships without a full object, but
			// with a set reference must render
			Relationship{FullRender: false, PeerReference: dummyref, Type: DEPENDS_ON},
			false, fmt.Sprintf("Relationship: %s DEPENDS_ON %s\n", host.SPDXID(), dummyref),
		},
		{
			// Relationships without a object and without a set reference
			// must return an error
			Relationship{FullRender: false, Type: DEPENDS_ON}, true, "",
		},
		{
			// Relationships with a peer object withouth id should err
			Relationship{FullRender: false, Peer: &File{}, Type: DEPENDS_ON}, true, "",
		},
		{
			// Relationships with only a a peer reference that should render
			// in full should err
			Relationship{FullRender: true, PeerReference: dummyref, Type: DEPENDS_ON}, true, "",
		},
		{
			// Relationships without a type should err
			Relationship{FullRender: false, PeerReference: dummyref}, true, "",
		},
	}

	for _, tc := range cases {
		res, err := tc.Rel.Render(host)
		if tc.MustErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
			require.Equal(t, tc.Rendered, res)
		}
	}

	// Full rednering should not be the same as non full render
	nonFullRender, err := cases[0].Rel.Render(host)
	require.Nil(t, err)
	cases[0].Rel.FullRender = true
	fullRender, err := cases[0].Rel.Render(host)
	require.Nil(t, err)
	require.NotEqual(t, nonFullRender, fullRender)

	// Finally, rendering with a host objectwithout an ID should err
	_, err = cases[0].Rel.Render(&File{})
	require.NotNil(t, err)
}

var testTar = `H4sICPIFo2AAA2hlbGxvLnRhcgDt1EsKwjAUBdCMXUXcQPuS5rMFwaEraDGgUFpIE3D5puAPRYuD
VNR7Jm+QQh7c3hQly44Sa/U4hdV0O8+YUKTJkLRCMhKk0zHX+VdjLA6h9pyz6Ju66198N3H+pYpy
iM1273P+Bm/lX4mUvyQlkP8cLvkHdwhFOIQMd4wBG6Oe5y/1Xf6VNhXjlGGXB3+e/yY2O9e2PV/H
xvnOBTcsF59eCmZT5Cz+yXT/5bX/pMb3P030fw4rlB8AAAAAAAAAAAAAAAAA4CccAXRRwL4AKAAA
`

var sampleManifest = `[{"Config":"386bcf5c63de46c7066c42d4ae1c38af0689836e88fed37d1dca2d484b343cf5.json","RepoTags":["k8s.gcr.io/kube-apiserver-amd64:v1.22.0-alpha.1"],"Layers":["23e140cb8e03a12cba4ac571d9a7143cf5e2e9b72de3b33ce3243b4f7ad6a188/layer.tar","48dd73ececdf0f52a174ad33a469145824713bd2b73c6257ce1ba8502003ad4e/layer.tar","d397673d78556210baa112013c960cb95a3fd452e5c4a2ead2b26e5a458cd87f/layer.tar"]}]
`

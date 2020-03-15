/*
Copyright 2019 The Kubernetes Authors.

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

package document

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"

	"k8s.io/release/pkg/notes/options"
	"k8s.io/release/pkg/release"
)

func TestFileMetadata(t *testing.T) {
	// Given
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	for _, file := range []string{
		"kubernetes-client-darwin-386.tar.gz",
		"kubernetes-client-darwin-amd64.tar.gz",
		"kubernetes-client-linux-386.tar.gz",
		"kubernetes-client-linux-amd64.tar.gz",
		"kubernetes-client-linux-arm.tar.gz",
		"kubernetes-client-linux-arm64.tar.gz",
		"kubernetes-client-linux-ppc64le.tar.gz",
		"kubernetes-client-linux-s390x.tar.gz",
		"kubernetes-client-windows-386.tar.gz",
		"kubernetes-client-windows-amd64.tar.gz",
		"kubernetes-node-linux-amd64.tar.gz",
		"kubernetes-node-linux-arm.tar.gz",
		"kubernetes-node-linux-arm64.tar.gz",
		"kubernetes-node-linux-ppc64le.tar.gz",
		"kubernetes-node-linux-s390x.tar.gz",
		"kubernetes-node-windows-amd64.tar.gz",
		"kubernetes-server-linux-amd64.tar.gz",
		"kubernetes-server-linux-arm.tar.gz",
		"kubernetes-server-linux-arm64.tar.gz",
		"kubernetes-server-linux-ppc64le.tar.gz",
		"kubernetes-server-linux-s390x.tar.gz",
		"kubernetes-src.tar.gz",
		"kubernetes.tar.gz",
	} {
		require.Nil(t, ioutil.WriteFile(
			filepath.Join(dir, file), []byte{1, 2, 3}, os.FileMode(0644),
		))
	}

	metadata, err := fetchMetadata(dir, "http://test.com", "test-release")
	require.Nil(t, err)

	expected := &FileMetadata{
		Source: []File{
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes.tar.gz", URL: "http://test.com/test-release/kubernetes.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-src.tar.gz", URL: "http://test.com/test-release/kubernetes-src.tar.gz"},
		},
		Client: []File{
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-darwin-386.tar.gz", URL: "http://test.com/test-release/kubernetes-client-darwin-386.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-darwin-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-client-darwin-amd64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-386.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-386.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-amd64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-arm.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-arm.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-arm64.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-arm64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-ppc64le.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-ppc64le.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-linux-s390x.tar.gz", URL: "http://test.com/test-release/kubernetes-client-linux-s390x.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-windows-386.tar.gz", URL: "http://test.com/test-release/kubernetes-client-windows-386.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-client-windows-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-client-windows-amd64.tar.gz"},
		},
		Server: []File{
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-server-linux-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-server-linux-amd64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-server-linux-arm.tar.gz", URL: "http://test.com/test-release/kubernetes-server-linux-arm.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-server-linux-arm64.tar.gz", URL: "http://test.com/test-release/kubernetes-server-linux-arm64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-server-linux-ppc64le.tar.gz", URL: "http://test.com/test-release/kubernetes-server-linux-ppc64le.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-server-linux-s390x.tar.gz", URL: "http://test.com/test-release/kubernetes-server-linux-s390x.tar.gz"},
		},
		Node: []File{
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-linux-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-node-linux-amd64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-linux-arm.tar.gz", URL: "http://test.com/test-release/kubernetes-node-linux-arm.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-linux-arm64.tar.gz", URL: "http://test.com/test-release/kubernetes-node-linux-arm64.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-linux-ppc64le.tar.gz", URL: "http://test.com/test-release/kubernetes-node-linux-ppc64le.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-linux-s390x.tar.gz", URL: "http://test.com/test-release/kubernetes-node-linux-s390x.tar.gz"},
			{Checksum: "27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29", Name: "kubernetes-node-windows-amd64.tar.gz", URL: "http://test.com/test-release/kubernetes-node-windows-amd64.tar.gz"},
		},
	}
	require.Equal(t, metadata, expected)
}

func TestRenderMarkdownTemplateGoldenFile(t *testing.T) {
	tests := []struct {
		name           string
		templateSpec   string
		userTemplate   bool
		hasDownloads   bool
		wantGoldenFile string
	}{
		{
			"render default template and downloads",
			options.FormatSpecDefaultGoTemplate,
			false,
			true,
			"document.md.golden",
		},
		{
			"render default template and no downloads",
			options.FormatSpecDefaultGoTemplate,
			true,
			false,
			"document_without_downloads.md.golden",
		},
		{
			"render user-specified template and downloads",
			"go-template:user-template.tmpl",
			true,
			true,
			"document.md.golden",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			dir, err := ioutil.TempDir("", "")
			require.NoError(t, err, "Creating tmpDir")
			defer os.RemoveAll(dir)

			if tc.hasDownloads {
				setupTestDir(t, dir)
			}

			doc := Document{
				NotesWithActionRequired: []string{"If an API changes and no one documented it, did it really happen?"},
				Notes: NoteCollection{
					NoteCategory{Kind: KindAPIChange, NoteEntries: &Notes{"This might make people sad...or happy."}},
					NoteCategory{Kind: KindBug, NoteEntries: &Notes{"This will likely get you promoted."}},
					NoteCategory{Kind: KindCleanup, NoteEntries: &Notes{"This usually does not get you promoted but it should."}},
					NoteCategory{Kind: KindDeprecation, NoteEntries: &Notes{"Delorted."}},
					NoteCategory{Kind: KindDesign, NoteEntries: &Notes{"Change the world."}},
					NoteCategory{Kind: KindDocumentation, NoteEntries: &Notes{"There was  a library in Alexandria, Egypt once."}},
					NoteCategory{Kind: KindFailingTest, NoteEntries: &Notes{"Please run presubmit test!"}},
					NoteCategory{Kind: KindFeature, NoteEntries: &Notes{"This will get you promoted."}},
					NoteCategory{Kind: KindFlake, NoteEntries: &Notes{"This *should* get you promoted."}},
					NoteCategory{Kind: KindBugCleanupFlake, NoteEntries: &Notes{"This should definitely get you promoted."}},
					NoteCategory{Kind: KindUncategorized, NoteEntries: &Notes{"Someone somewhere did the world a great justice."}},
				},
				PreviousRevision: "v1.16.0",
				CurrentRevision:  "v1.16.1",
			}

			templateSpec := tc.templateSpec
			if tc.userTemplate {
				// Write out the default template to simulate reading an actual template.
				p := filepath.Join(dir, strings.Split(tc.templateSpec, ":")[1])
				templateSpec = fmt.Sprintf("go-template:%s", p)

				require.NoError(t, ioutil.WriteFile(p, []byte(defaultReleaseNotesTemplate), 0664))
			}

			goldenFile, err := bazel.Runfile(filepath.Join("testdata", tc.wantGoldenFile))
			require.NoError(t, err, "Locating runfiles are you using bazel test?")

			b, err := ioutil.ReadFile(goldenFile)
			require.NoError(t, err, "Reading golden file %q", goldenFile)
			expected := string(b)

			// When
			got, err := doc.RenderMarkdownTemplate(release.ProductionBucket, dir, templateSpec)
			require.NoError(t, err, "Unexpected error rendering document")

			// Then
			require.Empty(t, diff.Diff(expected, got), "Unexpected diff")
		})
	}
}

func TestRenderMarkdownTemplateFailure(t *testing.T) {
	tests := []struct {
		name             string
		templateSpec     string
		templateContents string
		templateExist    bool
	}{
		{
			"given template exist but is empty",
			"go-template:empty.tmpl",
			"",
			true,
		},
		{
			"given bad template spec",
			"wrong-prefix:template.tmpl",
			"",
			true,
		},
		{
			"given bad template contents",
			"go-template:bad.tmpl",
			"# This template will not parse: {{}",
			true,
		},
		{
			"given non-existent template",
			"go-template:non-exist.tmpl",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "")
			require.Nil(t, err)
			defer os.RemoveAll(dir)

			if tt.templateExist {
				fileName := strings.Split(tt.templateSpec, ":")[1]
				p := filepath.Join(dir, fileName)
				require.Nil(t, ioutil.WriteFile(p, []byte(tt.templateContents), 0664))
			}

			doc := Document{}
			_, err = doc.RenderMarkdownTemplate("", "", tt.templateSpec)
			require.Error(t, err, "Unexpected success")
		})
	}
}

func TestCreateDownloadsTable(t *testing.T) {
	// Given
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)
	setupTestDir(t, dir)

	// When
	output := &strings.Builder{}
	require.Nil(t, CreateDownloadsTable(
		output, release.ProductionBucket, dir, "v1.16.0", "v1.16.1",
	))

	// Then
	require.Equal(t, `# v1.16.1

[Documentation](https://docs.k8s.io)

## Downloads for v1.16.1

filename | sha512 hash
-------- | -----------
[kubernetes.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-src.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-src.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`

### Client Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-client-darwin-386.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-darwin-386.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-darwin-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-darwin-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-386.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-386.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-arm.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-arm.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-arm64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-arm64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-ppc64le.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-linux-s390x.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-linux-s390x.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-windows-386.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-windows-386.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-client-windows-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-client-windows-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`

### Server Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-server-linux-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-server-linux-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-server-linux-arm.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-server-linux-arm.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-server-linux-arm64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-server-linux-arm64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-server-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-server-linux-ppc64le.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-server-linux-s390x.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-server-linux-s390x.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`

### Node Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-node-linux-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-linux-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-node-linux-arm.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-linux-arm.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-node-linux-arm64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-linux-arm64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-node-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-linux-ppc64le.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-node-linux-s390x.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-linux-s390x.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`
[kubernetes-node-windows-amd64.tar.gz](https://dl.k8s.io/v1.16.1/kubernetes-node-windows-amd64.tar.gz) | `+"`"+`27864cc5219a951a7a6e52b8c8dddf6981d098da1658d96258c870b2c88dfbcb51841aea172a28bafa6a79731165584677066045c959ed0f9929688d04defc29`+"`"+`

## Changelog since v1.16.0

`, output.String())
}

func TestSortKinds(t *testing.T) {
	input := NotesByKind{
		"cleanup":                       nil,
		"api-change":                    nil,
		"deprecation":                   nil,
		"documentation":                 nil,
		"Other (Bug, Cleanup or Flake)": nil,
		"failing-test":                  nil,
		"design":                        nil,
		"flake":                         nil,
		"bug":                           nil,
		"feature":                       nil,
		"Uncategorized":                 nil,
	}
	res := sortKinds(input)
	require.Equal(t, res, kindPriority)
}

// setupTestDir adds basic test files to a given directory.
func setupTestDir(t *testing.T, dir string) {
	for _, file := range []string{
		"kubernetes-client-darwin-386.tar.gz",
		"kubernetes-client-darwin-amd64.tar.gz",
		"kubernetes-client-linux-386.tar.gz",
		"kubernetes-client-linux-amd64.tar.gz",
		"kubernetes-client-linux-arm.tar.gz",
		"kubernetes-client-linux-arm64.tar.gz",
		"kubernetes-client-linux-ppc64le.tar.gz",
		"kubernetes-client-linux-s390x.tar.gz",
		"kubernetes-client-windows-386.tar.gz",
		"kubernetes-client-windows-amd64.tar.gz",
		"kubernetes-node-linux-amd64.tar.gz",
		"kubernetes-node-linux-arm.tar.gz",
		"kubernetes-node-linux-arm64.tar.gz",
		"kubernetes-node-linux-ppc64le.tar.gz",
		"kubernetes-node-linux-s390x.tar.gz",
		"kubernetes-node-windows-amd64.tar.gz",
		"kubernetes-server-linux-amd64.tar.gz",
		"kubernetes-server-linux-arm.tar.gz",
		"kubernetes-server-linux-arm64.tar.gz",
		"kubernetes-server-linux-ppc64le.tar.gz",
		"kubernetes-server-linux-s390x.tar.gz",
		"kubernetes-src.tar.gz",
		"kubernetes.tar.gz",
	} {
		require.Nil(t, ioutil.WriteFile(
			filepath.Join(dir, file), []byte{1, 2, 3}, os.FileMode(0644),
		))
	}
}

func TestNoteCollection(t *testing.T) {
	tests := []struct {
		name       string
		collection NoteCollection
		wantOutput string
	}{
		{
			"render collection with higher priority kind appearing first",
			NoteCollection{
				NoteCategory{
					Kind:        KindDeprecation,
					NoteEntries: &Notes{"change 1"},
				},
				NoteCategory{
					Kind:        KindAPIChange,
					NoteEntries: &Notes{"change 2"},
				},
			},
			"\n# Deprecation\n  * change 1\n# API Change\n  * change 2",
		},
		{
			"render collection with lower priority kind appearing first",
			NoteCollection{
				NoteCategory{
					Kind:        KindAPIChange,
					NoteEntries: &Notes{"change 2"},
				},
				NoteCategory{
					Kind:        KindDeprecation,
					NoteEntries: &Notes{"change 1"},
				},
			},
			"\n# Deprecation\n  * change 1\n# API Change\n  * change 2",
		},
		{
			"render with equal priority kinds",
			NoteCollection{
				NoteCategory{
					Kind:        KindDeprecation,
					NoteEntries: &Notes{"change 1"},
				},
				NoteCategory{
					Kind:        KindDeprecation,
					NoteEntries: &Notes{"change 2"},
				},
			},
			"\n# Deprecation\n  * change 1\n# Deprecation\n  * change 2",
		},
	}

	const goTemplate = `
{{- range . }}
# {{ .Kind | prettyKind}}
{{range $note := .NoteEntries}}  * {{$note}}{{end -}}
{{end -}}
`

	tmpl, err := template.New("markdown").
		Funcs(template.FuncMap{"prettyKind": prettyKind}).
		Parse(goTemplate)
	require.NoError(t, err, "parsing template")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got strings.Builder
			tt.collection.Sort(kindPriority)
			require.NoError(t, tmpl.Execute(&got, tt.collection), "rendering template")

			require.Equal(t, tt.wantOutput, got.String())
		})
	}
}

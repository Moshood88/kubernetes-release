/*
Copyright 2020 The Kubernetes Authors.

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

package object

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/release/pkg/gcp"
	"k8s.io/utils/pointer"
)

var (
	// GcsPrefix url prefix for google cloud storage buckets
	GcsPrefix      = "gs://"
	concurrentFlag = "-m"
	recursiveFlag  = "-r"
	noClobberFlag  = "-n"
)

type GCS struct {
	// TODO: Implement store
	opts *GCSOptions
}

func NewGCS(opts *GCSOptions) *GCS {
	return &GCS{opts}
}

// GCSOptions are the main options to pass to `GCS`.
type GCSOptions struct {
	// TODO: Populate fields
	// gsutil options
	Concurrent *bool
	Recursive  *bool
	NoClobber  *bool

	// local options
	// AllowMissing allows a copy operation to be skipped if the source or
	// destination does not exist. This is useful for scenarios where copy
	// operations happen in a loop/channel, so a single "failure" does not block
	// the entire operation.
	AllowMissing *bool
}

// TODO: Consider a method to set options

func NewDefaultGCS() *GCS {
	return &GCS{DefaultGCSCopyOptions}
}

// DefaultGCSCopyOptions have the default options for the GCS copy action
var DefaultGCSCopyOptions = &GCSOptions{
	Concurrent:   pointer.BoolPtr(true),
	Recursive:    pointer.BoolPtr(true),
	NoClobber:    pointer.BoolPtr(true),
	AllowMissing: pointer.BoolPtr(true),
}

// CopyToGCS copies a local directory to the specified GCS path
func (g *GCS) CopyToGCS(src, gcsPath string) error {
	logrus.Infof("Copying %s to GCS (%s)", src, gcsPath)
	gcsPath, gcsPathErr := g.NormalizeGCSPath(gcsPath)
	if gcsPathErr != nil {
		return errors.Wrap(gcsPathErr, "normalize GCS path")
	}

	_, err := os.Stat(src)
	if err != nil {
		logrus.Info("Unable to get local source directory info")

		if *g.opts.AllowMissing {
			logrus.Infof("Source directory (%s) does not exist. Skipping GCS upload.", src)
			return nil
		}

		return errors.New("source directory does not exist")
	}

	return g.bucketCopy(src, gcsPath)
}

// CopyToLocal copies a GCS path to the specified local directory
func (g *GCS) CopyToLocal(gcsPath, dst string) error {
	logrus.Infof("Copying GCS (%s) to %s", gcsPath, dst)
	gcsPath, gcsPathErr := g.NormalizeGCSPath(gcsPath)
	if gcsPathErr != nil {
		return errors.Wrap(gcsPathErr, "normalize GCS path")
	}

	return g.bucketCopy(gcsPath, dst)
}

// CopyBucketToBucket copies between two GCS paths.
func (g *GCS) CopyBucketToBucket(src, dst string) error {
	logrus.Infof("Copying %s to %s", src, dst)

	src, srcErr := g.NormalizeGCSPath(src)
	if srcErr != nil {
		return errors.Wrap(srcErr, "normalize GCS path")
	}

	dst, dstErr := g.NormalizeGCSPath(dst)
	if dstErr != nil {
		return errors.Wrap(dstErr, "normalize GCS path")
	}

	return g.bucketCopy(src, dst)
}

func (g *GCS) bucketCopy(src, dst string) error {
	args := []string{}

	if *g.opts.Concurrent {
		logrus.Debug("Setting GCS copy to run concurrently")
		args = append(args, concurrentFlag)
	}

	args = append(args, "cp")
	if *g.opts.Recursive {
		logrus.Debug("Setting GCS copy to run recursively")
		args = append(args, recursiveFlag)
	}
	if *g.opts.NoClobber {
		logrus.Debug("Setting GCS copy to not clobber existing files")
		args = append(args, noClobberFlag)
	}

	args = append(args, src, dst)

	if err := gcp.GSUtil(args...); err != nil {
		return errors.Wrap(err, "gcs copy")
	}

	return nil
}

// GetReleasePath returns a GCS path to retrieve builds from or push builds to
//
// Expected destination format:
//   gs://<bucket>/<gcsRoot>[/fast][/<version>]
func (g *GCS) GetReleasePath(
	bucket, gcsRoot, version string,
	fast bool) (string, error) {
	gcsPath, err := g.getPath(
		bucket,
		gcsRoot,
		version,
		"release",
		fast,
	)
	if err != nil {
		return "", errors.Wrap(err, "normalize GCS path")
	}

	logrus.Infof("Release path is %s", gcsPath)
	return gcsPath, nil
}

// GetMarkerPath returns a GCS path where version markers should be stored
//
// Expected destination format:
//   gs://<bucket>/<gcsRoot>
func (g *GCS) GetMarkerPath(
	bucket, gcsRoot string) (string, error) {
	gcsPath, err := g.getPath(
		bucket,
		gcsRoot,
		"",
		"marker",
		false,
	)
	if err != nil {
		return "", errors.Wrap(err, "normalize GCS path")
	}

	logrus.Infof("Version marker path is %s", gcsPath)
	return gcsPath, nil
}

// GetReleasePath returns a GCS path to retrieve builds from or push builds to
//
// Expected destination format:
//   gs://<bucket>/<gcsRoot>[/fast][/<version>]
// TODO: Support "release" buildType
func (g *GCS) getPath(
	bucket, gcsRoot, version, pathType string,
	fast bool) (string, error) {
	if gcsRoot == "" {
		return "", errors.New("GCS root must be specified")
	}

	gcsPathParts := []string{}

	gcsPathParts = append(gcsPathParts, bucket, gcsRoot)

	if pathType == "release" {
		if fast {
			gcsPathParts = append(gcsPathParts, "fast")
		}

		if version != "" {
			gcsPathParts = append(gcsPathParts, version)
		}
	} else if pathType == "marker" {
	} else {
		return "", errors.New("a GCS path type must be specified")
	}

	// Ensure any constructed GCS path is prefixed with `gs://`
	return g.NormalizeGCSPath(gcsPathParts...)
}

// NormalizeGCSPath takes a GCS path and ensures that the `GcsPrefix` is
// prepended to it.
// TODO: Should there be an append function for paths to prevent multiple calls
//       like in build.checkBuildExists()?
func (g *GCS) NormalizeGCSPath(gcsPathParts ...string) (string, error) {
	gcsPath := ""

	// Ensure there is at least one element in the gcsPathParts slice before
	// trying to construct a path
	if len(gcsPathParts) == 0 {
		return "", errors.New("must contain at least one path part")
	} else if len(gcsPathParts) == 1 {
		if gcsPathParts[0] == "" {
			return "", errors.New("path should not be an empty string")
		}

		gcsPath = gcsPathParts[0]
	} else {
		var emptyParts int

		for i, part := range gcsPathParts {
			if part == "" {
				emptyParts++
			}

			if i == 0 {
				continue
			}

			if strings.Contains(part, "gs:/") {
				return "", errors.New("one of the GCS path parts contained a `gs:/`, which may suggest a filepath.Join() error in the caller")
			}

			if i == len(gcsPathParts)-1 && emptyParts == len(gcsPathParts) {
				return "", errors.New("all paths provided were empty")
			}
		}

		gcsPath = filepath.Join(gcsPathParts...)
	}

	// Strip `gs://` if it was included in gcsPathParts
	gcsPath = strings.TrimPrefix(gcsPath, GcsPrefix)

	// Strip `gs:/` if:
	// - `gs://` was included in gcsPathParts
	// - gcsPathParts had more than element
	// - filepath.Join() was called somewhere in a caller's logic
	gcsPath = strings.TrimPrefix(gcsPath, "gs:/")

	// Strip `/`
	// This scenario may never happen, but let's catch it, just in case
	gcsPath = strings.TrimPrefix(gcsPath, "/")

	gcsPath = GcsPrefix + gcsPath

	isNormalized := g.IsPathNormalized(gcsPath)
	if !isNormalized {
		return gcsPath, errors.New("unknown error while trying to normalize GCS path")
	}

	return gcsPath, nil
}

// IsPathNormalized determines if a GCS path is prefixed with `gs://`.
// Use this function as pre-check for any gsutil/GCS functions that manipulate
// GCS bucket contents.
func (g *GCS) IsPathNormalized(gcsPath string) bool {
	var errCount int

	if !strings.HasPrefix(gcsPath, GcsPrefix) {
		logrus.Errorf("GCS path (%s) should be prefixed with `gs://`", gcsPath)
		errCount++
	}

	strippedPath := strings.TrimPrefix(gcsPath, GcsPrefix)
	if strings.Contains(strippedPath, "gs:/") {
		logrus.Errorf("GCS path (%s) should be prefixed with `gs:/`", gcsPath)
		errCount++
	}

	// TODO: Add logic to handle invalid path characters

	if errCount > 0 {
		return false
	}

	return true
}

// RsyncRecursive runs `gsutil rsync` in recursive mode. The caller of this
// function has to ensure that the provided paths are prefixed with gs:// if
// necessary (see `NormalizeGCSPath()`).
func (g *GCS) RsyncRecursive(src, dst string) error {
	return errors.Wrap(
		gcp.GSUtil(concurrentFlag, "rsync", recursiveFlag, src, dst),
		"running gsutil rsync",
	)
}

// PathExists returns true if the specified GCS path exists.
func (g *GCS) PathExists(gcsPath string) (bool, error) {
	if !g.IsPathNormalized(gcsPath) {
		return false, errors.New("cannot run `gsutil ls` GCS path does not begin with `gs://`")
	}

	err := gcp.GSUtil(
		"ls",
		gcsPath,
	)
	if err != nil {
		return false, err
	}

	logrus.Infof("Found %s", gcsPath)
	return true, nil
}

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

package anago

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/release/pkg/build"
	"k8s.io/release/pkg/build/make"
	"k8s.io/release/pkg/changelog"
	"k8s.io/release/pkg/gcp/gcb"
	"k8s.io/release/pkg/git"
	"k8s.io/release/pkg/release"
)

// stageClient is a client for staging releases.
//counterfeiter:generate . stageClient
type stageClient interface {
	// Submit can be used to submit a Google Cloud Build (GCB) job.
	Submit() error

	// Validate if the provided `ReleaseOptions` are correctly set.
	ValidateOptions() error

	// CheckPrerequisites verifies that a valid GITHUB_TOKEN environment
	// variable is set. It also checks for the existence and version of
	// required packages and if the correct Google Cloud project is set. A
	// basic hardware check will ensure that enough disk space is available,
	// too.
	CheckPrerequisites() error

	// SetBuildCandidate discovers the release branch, parent branch (if
	// available) and build version for this release.
	SetBuildCandidate() error

	// GenerateReleaseVersion discovers the next versions to be released.
	GenerateReleaseVersion() error

	// PrepareWorkspace verifies that the working directory is in the desired
	// state. This means that the build directory is cleaned up and the checked
	// out repository is in a clean state.
	PrepareWorkspace() error

	// TagRepository creates all necessary git objects by tagging the
	// repository for the provided `versions` the main version `versionPrime`
	// and the `parentBranch`.
	TagRepository() error

	// Build runs 'make cross-in-a-container' by using the latest kubecross
	// container image. This step also build all necessary release tarballs.
	Build() error

	// GenerateChangelog builds the CHANGELOG-x.y.md file and commits it
	// into the local repository.
	GenerateChangelog() error

	// StageArtifacts copies the build artifacts to a Google Cloud Bucket.
	StageArtifacts() error
}

// DefaultStage is the default staging implementation used in production.
type DefaultStage struct {
	impl    stageImpl
	options *StageOptions
	state   *StageState
}

// NewDefaultStage creates a new defaultStage instance.
func NewDefaultStage(options *StageOptions) *DefaultStage {
	return &DefaultStage{&defaultStageImpl{}, options, nil}
}

// SetImpl can be used to set the internal stage implementation.
func (d *DefaultStage) SetImpl(impl stageImpl) {
	d.impl = impl
}

// SetState fixes the current state. Mainly used for passing
// arbitrary values during testing
func (d *DefaultStage) SetState(state *StageState) {
	d.state = state
}

// defaultStageImpl is the default internal stage client implementation.
type defaultStageImpl struct{}

// stageImpl is the implementation of the stage client.
//counterfeiter:generate . stageImpl
type stageImpl interface {
	Submit(options *gcb.Options) error
	PrepareWorkspaceStage() error
	GenerateReleaseVersion(
		releaseType, version, branch string, branchFromMaster bool,
	) (*release.Versions, error)
	ConfigureGlobalDefaultUserAndEmail() error
	OpenRepo(repoPath string) (*git.Repo, error)
	RevParse(repo *git.Repo, rev string) (string, error)
	HasBranch(repo *git.Repo, branch string) (bool, error)
	Checkout(repo *git.Repo, rev string, args ...string) error
	CurrentBranch(repo *git.Repo) (string, error)
	CommitEmpty(repo *git.Repo, msg string) error
	Tag(repo *git.Repo, name, message string) error
	CheckReleaseBucket(options *build.Options) error
	MakeCross(version string) error
	GenerateChangelog(options *changelog.Options) error
	StageLocalSourceTree(
		options *build.Options, workDir, buildVersion string,
	) error
	StageLocalArtifacts(options *build.Options) error
	PushReleaseArtifacts(
		options *build.Options, srcPath, gcsPath string,
	) error
	PushContainerImages(options *build.Options) error
}

func (d *defaultStageImpl) Submit(options *gcb.Options) error {
	return gcb.New(options).Submit()
}

func (d *defaultStageImpl) PrepareWorkspaceStage() error {
	if err := release.PrepareWorkspaceStage(gitRoot); err != nil {
		return err
	}
	return os.Chdir(gitRoot)
}

func (d *defaultStageImpl) GenerateReleaseVersion(
	releaseType, version, branch string, branchFromMaster bool,
) (*release.Versions, error) {
	return release.GenerateReleaseVersion(
		releaseType, version, branch, branchFromMaster,
	)
}

func (d *defaultStageImpl) ConfigureGlobalDefaultUserAndEmail() error {
	return git.ConfigureGlobalDefaultUserAndEmail()
}

func (d *defaultStageImpl) OpenRepo(repoPath string) (*git.Repo, error) {
	return git.OpenRepo(repoPath)
}

func (d *defaultStageImpl) RevParse(repo *git.Repo, rev string) (string, error) {
	return repo.RevParse(rev)
}

func (d *defaultStageImpl) HasBranch(repo *git.Repo, branch string) (bool, error) {
	return repo.HasBranch(branch)
}

func (d *defaultStageImpl) Checkout(repo *git.Repo, rev string, args ...string) error {
	return repo.Checkout(rev, args...)
}

func (d *defaultStageImpl) CurrentBranch(repo *git.Repo) (string, error) {
	return repo.CurrentBranch()
}

func (d *defaultStageImpl) CommitEmpty(repo *git.Repo, msg string) error {
	return repo.CommitEmpty(msg)
}

func (d *defaultStageImpl) Tag(repo *git.Repo, name, message string) error {
	return repo.Tag(name, message)
}

func (d *defaultStageImpl) MakeCross(version string) error {
	return make.New().MakeCross(version)
}

func (d *defaultStageImpl) GenerateChangelog(options *changelog.Options) error {
	return changelog.New(options).Run()
}

func (d *defaultStageImpl) CheckReleaseBucket(
	options *build.Options,
) error {
	return build.New(options).CheckReleaseBucket()
}

func (d *defaultStageImpl) StageLocalSourceTree(
	options *build.Options, workDir, buildVersion string,
) error {
	return build.New(options).StageLocalSourceTree(workDir, buildVersion)
}

func (d *defaultStageImpl) StageLocalArtifacts(
	options *build.Options,
) error {
	return build.New(options).StageLocalArtifacts()
}

func (d *defaultStageImpl) PushReleaseArtifacts(
	options *build.Options, srcPath, gcsPath string,
) error {
	return build.New(options).PushReleaseArtifacts(srcPath, gcsPath)
}

func (d *defaultStageImpl) PushContainerImages(
	options *build.Options,
) error {
	return build.New(options).PushContainerImages()
}

func (d *DefaultStage) Submit() error {
	options := gcb.NewDefaultOptions()
	options.Stage = true
	options.NoMock = d.options.NoMock
	options.Branch = d.options.ReleaseBranch
	options.ReleaseType = d.options.ReleaseType
	options.NoAnago = true
	return d.impl.Submit(options)
}

func (d *DefaultStage) ValidateOptions() error {
	// Call options, validate. The validation returns the initial
	// state of the stage process
	state, err := d.options.Validate()
	if err != nil {
		return errors.Wrap(err, "validating options")
	}
	d.state = state
	return nil
}

func (d *DefaultStage) CheckPrerequisites() error { return nil }

func (d *DefaultStage) SetBuildCandidate() error {
	// TODO: the parent branch has to be returned by the SetBuildCandidate
	// method. It should be empty (releases cut from master) or
	// git.DefaultBranch / "master" (releases cut from release branches).
	//
	// d.state.parentBranch = XXXXX
	d.state.parentBranch = ""
	return nil
}

func (d *DefaultStage) GenerateReleaseVersion() error {
	versions, err := d.impl.GenerateReleaseVersion(
		d.options.ReleaseType,
		d.options.BuildVersion,
		d.options.ReleaseBranch,
		d.state.parentBranch == git.DefaultBranch,
	)
	if err != nil {
		return errors.Wrap(err, "generating release versions for stage")
	}
	// Set the versions on the state
	d.state.versions = versions
	return nil
}

func (d *DefaultStage) PrepareWorkspace() error {
	if err := d.impl.PrepareWorkspaceStage(); err != nil {
		return errors.Wrap(err, "prepare workspace")
	}
	return nil
}

func (d *DefaultStage) TagRepository() error {
	logrus.Info("Configuring git user and email")
	if err := d.impl.ConfigureGlobalDefaultUserAndEmail(); err != nil {
		return errors.Wrap(err, "configure git user and email")
	}

	repo, err := d.impl.OpenRepo(gitRoot)
	if err != nil {
		return errors.Wrap(err, "open Kubernetes repository")
	}

	for _, version := range d.state.versions.Ordered() {
		logrus.Infof("Preparing version %s", version)

		// Ensure that the tag not already exists
		if _, err := d.impl.RevParse(repo, version); err == nil {
			return errors.Errorf("tag %s already exists", version)
		}

		commit := d.state.semverBuildVersion.Build[0]
		if d.state.parentBranch != "" {
			logrus.Infof("Parent branch provided: %s", d.state.parentBranch)

			if version == d.state.versions.Prime() {
				logrus.Infof("Version %s is the prime version", version)
				logrus.Infof(
					"Creating or checking out release branch %s",
					d.options.ReleaseBranch,
				)

				hasBranch, err := d.impl.HasBranch(
					repo, d.options.ReleaseBranch,
				)
				if err != nil {
					return errors.Wrap(err, "check if repository has branch")
				}
				logrus.Infof("Branch already exist: %v", hasBranch)

				if !hasBranch {
					logrus.Infof(
						"Creating release branch %s from commit %s",
						d.options.ReleaseBranch, commit,
					)
					if err := d.impl.Checkout(
						repo, "-b", d.options.ReleaseBranch, commit,
					); err != nil {
						return errors.Wrap(err, "create new release branch")
					}
				} else {
					logrus.Infof(
						"Checking out release branch %s since it already exist",
						d.options.ReleaseBranch,
					)
					if err := d.impl.Checkout(
						repo, d.options.ReleaseBranch,
					); err != nil {
						return errors.Wrap(err, "checkout release branch")
					}
				}
			} else {
				logrus.Infof(
					"Version %s it not the prime, checking out parent branch",
					version,
				)
				if err := d.impl.Checkout(repo, d.state.parentBranch); err != nil {
					return errors.Wrap(err, "checkout parent branch")
				}
			}
		} else {
			logrus.Infof("Checking out commit %s", commit)
			if err := d.impl.Checkout(repo, commit); err != nil {
				return errors.Wrap(err, "checkout release commit")
			}
		}

		// `branch == ""` in case we checked out a commit directly, which is
		// then in detached head state.
		branch, err := d.impl.CurrentBranch(repo)
		if err != nil {
			return errors.Wrap(err, "get current branch")
		}
		logrus.Infof("Current branch is %q", branch)

		// For release branches, we create an empty release commit to avoid
		// potential ambiguous 'git describe' logic between the official
		// release, 'x.y.z' and the next beta of that release branch,
		// 'x.y.(z+1)-beta.0'.
		//
		// We avoid doing this empty release commit on 'master', as:
		//   - there is a potential for branch conflicts as upstream/master
		//     moves ahead
		//   - we're checking out a git ref, as opposed to a branch, which
		//     means the tag will detached from 'upstream/master'
		//
		// A side-effect of the tag being detached from 'master' is the primary
		// build job (ci-kubernetes-build) will build as the previous alpha,
		// instead of the assumed tag. This causes the next anago run against
		// 'master' to fail due to an old build version.
		//
		// Example: 'v1.18.0-alpha.2.663+df908c3aad70be'
		//          (should instead be:
		//			 'v1.18.0-alpha.3.<commits-since-tag>+<commit-ish>')
		//
		// ref:
		//   - https://github.com/kubernetes/release/issues/1020
		//   - https://github.com/kubernetes/release/pull/1030
		//   - https://github.com/kubernetes/release/issues/1080
		//   - https://github.com/kubernetes/kubernetes/pull/88074
		if strings.HasPrefix(branch, "release-") {
			logrus.Infof("Creating empty release commit for tag %s", version)
			if err := d.impl.CommitEmpty(
				repo,
				fmt.Sprintf("Release commit for Kubernetes %s", version),
			); err != nil {
				return errors.Wrap(err, "create empty release commit")
			}
		}

		// Do the actual tag
		logrus.Infof("Tagging version %s", version)
		if err := d.impl.Tag(
			repo,
			version,
			fmt.Sprintf(
				"Kubernetes %s release %s", d.options.ReleaseType, version,
			),
		); err != nil {
			return errors.Wrap(err, "tag version")
		}
	}
	return nil
}

func (d *DefaultStage) Build() error {
	for _, version := range d.state.versions.Ordered() {
		if err := d.impl.MakeCross(version); err != nil {
			return errors.Wrap(err, "build artifacts")
		}
	}
	return nil
}

func (d *DefaultStage) GenerateChangelog() error {
	branch := d.options.ReleaseBranch
	if d.state.parentBranch != "" {
		branch = d.state.parentBranch
	}
	return d.impl.GenerateChangelog(&changelog.Options{
		RepoPath:     gitRoot,
		Tag:          d.state.versions.Prime(),
		Branch:       branch,
		Bucket:       d.options.Bucket(),
		HTMLFile:     filepath.Join(workspaceDir, "src/release-notes.html"),
		Dependencies: true,
		Tars: filepath.Join(
			gitRoot,
			fmt.Sprintf("%s-%s", release.BuildDir, d.state.versions.Prime()),
			release.ReleaseTarsPath,
		),
	})
}

func (d *DefaultStage) StageArtifacts() error {
	for _, version := range d.state.versions.Ordered() {
		logrus.Infof("Staging artifacts for version %s", version)
		buildDir := filepath.Join(
			gitRoot, fmt.Sprintf("%s-%s", release.BuildDir, version),
		)
		bucket := d.options.Bucket()
		containerRegistry := d.options.ContainerRegistry()
		pushBuildOptions := &build.Options{
			Bucket:                     bucket,
			BuildDir:                   buildDir,
			Registry:                   containerRegistry,
			Version:                    version,
			AllowDup:                   true,
			ValidateRemoteImageDigests: true,
		}
		if err := d.impl.CheckReleaseBucket(pushBuildOptions); err != nil {
			return errors.Wrap(err, "check release bucket access")
		}

		// Stage the local source tree
		if err := d.impl.StageLocalSourceTree(
			pushBuildOptions,
			workspaceDir,
			d.options.BuildVersion,
		); err != nil {
			return errors.Wrap(err, "staging local source tree")
		}

		// Stage local artifacts and write checksums
		if err := d.impl.StageLocalArtifacts(pushBuildOptions); err != nil {
			return errors.Wrap(err, "staging local artifacts")
		}
		gcsPath := filepath.Join("stage", d.options.BuildVersion, version)

		// Push gcs-stage to GCS
		if err := d.impl.PushReleaseArtifacts(
			pushBuildOptions,
			filepath.Join(buildDir, release.GCSStagePath, version),
			filepath.Join(gcsPath, release.GCSStagePath, version),
		); err != nil {
			return errors.Wrap(err, "pushing release artifacts")
		}

		// Push container release-images to GCS
		if err := d.impl.PushReleaseArtifacts(
			pushBuildOptions,
			filepath.Join(buildDir, release.ImagesPath),
			filepath.Join(gcsPath, release.ImagesPath),
		); err != nil {
			return errors.Wrap(err, "pushing release artifacts")
		}

		// Push container images into registry
		if err := d.impl.PushContainerImages(pushBuildOptions); err != nil {
			return errors.Wrap(err, "pushing container images")
		}
	}

	noMockFlag := ""
	if d.options.NoMock {
		noMockFlag = "--nomock"
	}

	logrus.Infof(
		"To release this staged build, run:\n\n"+
			"$ krel gcbmgr --no-anago --release "+
			"--type %s "+
			"--branch %s "+
			"--build-version=%s %s",
		d.options.ReleaseType,
		d.options.ReleaseBranch,
		d.options.BuildVersion,
		noMockFlag,
	)
	return nil
}

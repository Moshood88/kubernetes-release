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

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"k8s.io/release/cmd/schedule-builder/model"
	"sigs.k8s.io/release-utils/log"
	"sigs.k8s.io/yaml"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "schedule-builder --config-path path/to/schedule.yaml --type <release>/or/<patch>[--output-file <filename.md>]",
	Short:             "schedule-builder generate a human readable format of the Kubernetes release schedule",
	Example:           "schedule-builder --config-path /home/user/kubernetes/sig-release/releases/schedule.yaml --type patch",
	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRunE: initLogging,
	RunE: func(*cobra.Command, []string) error {
		return run(opts)
	},
}

type options struct {
	configPath string
	outputFile string
	logLevel   string
	typeFile   string
}

var opts = &options{}

const (
	configPathFlag = "config-path"
	outputFileFlag = "output-file"
	typeFlag       = "type"
)

var requiredFlags = []string{
	configPathFlag,
	typeFlag,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&opts.configPath,
		configPathFlag,
		"",
		"path where can find the schedule.yaml file",
	)

	rootCmd.PersistentFlags().StringVar(
		&opts.outputFile,
		outputFileFlag,
		"",
		"name of the file that save the schedule to. If not set it will just output to the stdout.",
	)

	rootCmd.PersistentFlags().StringVar(
		&opts.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)

	rootCmd.PersistentFlags().StringVar(
		&opts.typeFile,
		typeFlag,
		"patch",
		"type of file to be produced - release cycle schedule or patch schedule. To be set to `release` or `patch` and respective yaml needs to be supplied with --config-path",
	)

	for _, flag := range requiredFlags {
		if err := rootCmd.MarkPersistentFlagRequired(flag); err != nil {
			logrus.Fatal(err)
		}
	}
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(opts.logLevel)
}

func run(opts *options) error {
	if err := opts.SetAndValidate(); err != nil {
		return fmt.Errorf("validating schedule-path options: %w", err)
	}

	logrus.Infof("Reading the schedule file %s...", opts.configPath)
	data, err := os.ReadFile(opts.configPath)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}

	var (
		patchSchedule   model.PatchSchedule
		releaseSchedule model.ReleaseSchedule
		scheduleOut     string
	)

	logrus.Info("Parsing the schedule...")

	switch opts.typeFile {
	case "patch":
		if err := yaml.UnmarshalStrict(data, &patchSchedule); err != nil {
			return fmt.Errorf("failed to decode the file: %w", err)
		}

		logrus.Info("Generating the markdown output...")
		scheduleOut = parseSchedule(patchSchedule)

	case "release":
		if err := yaml.UnmarshalStrict(data, &releaseSchedule); err != nil {
			return fmt.Errorf("failed to decode the file: %w", err)
		}

		logrus.Info("Generating the markdown output...")
		scheduleOut = parseReleaseSchedule(releaseSchedule)

	default:
		return errors.New("type must be either `release` or `patch`")
	}

	if opts.outputFile != "" {
		logrus.Infof("Saving schedule to a file %s.", opts.outputFile)
		//nolint:gosec // TODO(gosec): G306: Expect WriteFile permissions to be
		// 0600 or less
		err := os.WriteFile(opts.outputFile, []byte(scheduleOut), 0o644)
		if err != nil {
			return fmt.Errorf("failed to save schedule to the file: %w", err)
		}
		logrus.Info("File saved")
	}

	return nil
}

// SetAndValidate sets some default options and verifies if options are valid
func (o *options) SetAndValidate() error {
	logrus.Info("Validating schedule-path options...")

	if o.configPath == "" {
		return fmt.Errorf("need to set the config-path")
	}

	return nil
}

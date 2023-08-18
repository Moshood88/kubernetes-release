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

// nolint
// This file is intended for legacy support and should not be linted
package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetKubeadmConfig(t *testing.T) {
	testcases := []struct {
		version      string
		expectConfig string
		expectErr    bool
	}{
		{
			"not-a-real-version",
			"",
			true,
		},
		{
			"1.22.0",
			"post-1.10/10-kubeadm.conf",
			false,
		},
	}

	for _, tc := range testcases {
		v := version{
			Version: tc.version,
		}
		kubeadmConfig, err := getKubeadmKubeletConfigFile(v)

		if err != nil {
			if !tc.expectErr {
				t.Errorf("getKubeadmConfig(%s) returned unwanted error: %v", tc.version, err)
			}
		} else {
			if kubeadmConfig != tc.expectConfig {
				t.Errorf("getKubeadmConfig(%s) got %q, wanted %q", tc.version, kubeadmConfig, tc.expectConfig)
			}
		}
	}
}

func TestGetKubeadmDependencies(t *testing.T) {
	testcases := []struct {
		name    string
		version string
		deps    []string
	}{
		{
			name:    "minimum supported kubernetes",
			version: "1.19.0",
			deps: []string{
				"kubelet (>= 1.19.0)",
				"kubectl (>= 1.19.0)",
				"kubernetes-cni (>= 0.8.7)",
				"${misc:Depends}",
				"cri-tools (>= 1.25.0)",
			},
		},
		{
			name:    "latest stable minor kubernetes",
			version: "1.22.0",
			deps: []string{
				"kubelet (>= 1.19.0)",
				"kubectl (>= 1.19.0)",
				"kubernetes-cni (>= 0.8.7)",
				"${misc:Depends}",
				"cri-tools (>= 1.25.0)",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			v := version{Version: tc.version}
			deps, err := getKubeadmDependencies(v)
			if err != nil {
				t.Fatalf("did not expect an error: %v", err)
			}
			actual := strings.Split(deps, ", ")
			if len(actual) != len(tc.deps) {
				t.Fatalf("Expected %d deps but found %d", len(tc.deps), len(actual))
			}
			if !reflect.DeepEqual(actual, tc.deps) {
				t.Fatalf("expected %q but got %q", tc.deps, actual)
			}
		})
	}
}

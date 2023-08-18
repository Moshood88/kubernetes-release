#!/usr/bin/env bash
# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

docker build -t kubelet-rpm-builder .
echo "Cleaning output directory..."
rm -rf output/*
mkdir -p output

architectures=${1:-""}

docker run -i --rm -v $PWD/output/:/root/rpmbuild/RPMS/ kubelet-rpm-builder "$architectures"

USER=${USER:-$(id -u)}
if [[ $USER != 0 ]]; then
  sudo chown -R "$USER" "$PWD/output"
fi

echo
echo "----------------------------------------"
echo
echo "RPMs written to: "
ls $PWD/output/*/
echo
echo "Yum repodata written to: "
ls $PWD/output/*/repodata/

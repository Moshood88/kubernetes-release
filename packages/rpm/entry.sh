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

# Entrypoint for the build container to create the rpms and yum repodata:
# Usage: ./entry.sh GOARCH/RPMARCH,GOARCH/RPMARCH,....

set -e

declare -a ARCHS

if [ $# -gt 0 ]; then
  IFS=','; ARCHS=($1); unset IFS;
else
  #GOARCH/RPMARCH
  ARCHS=(
    amd64
    arm
    arm64
    ppc64le
    s390x
  )
fi

declare -A GOTORPMARCH=(
    [amd64]=x86_64
    [arm]=armhfp
    [arm64]=aarch64
    [ppc64le]=ppc64le
    [s390x]=s390x
)

for GOARCH in "${ARCHS[@]}"; do
  RPMARCH=${GOTORPMARCH[$GOARCH]}
  SRC_PATH="/root/rpmbuild/SOURCES/${RPMARCH}"
  mkdir -p ${SRC_PATH}
  cp -r /root/rpmbuild/SPECS/* ${SRC_PATH}
  echo "Building RPM's for ${GOARCH}....."
  sed -i "s/\%global ARCH.*/\%global ARCH ${GOARCH}/" ${SRC_PATH}/kubelet.spec
  # Download sources if not already available
  cd ${SRC_PATH} && spectool -gf kubelet.spec
  /usr/bin/rpmbuild --target ${RPMARCH} --define "_sourcedir ${SRC_PATH}" -bb ${SRC_PATH}/kubelet.spec
  mkdir -p /root/rpmbuild/RPMS/${RPMARCH}
  createrepo -o /root/rpmbuild/RPMS/${RPMARCH}/ /root/rpmbuild/RPMS/${RPMARCH}
done

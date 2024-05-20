#!/usr/bin/env bash

# Copyright 2024 The Kubesysadm Authors.
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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
DIFFROOT="${SCRIPT_ROOT}/api"
DIFFMANIFESTSROOT="${SCRIPT_ROOT}/config"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp"
TMP_DIFFMANIFESTSROOT="${TMP_DIFFROOT}/config"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp"
CONTROLLER_TOOLS_VERSION="v0.14.0"
CONTROLLER_GEN=${SCRIPT_ROOT}/bin/controller-gen-${CONTROLLER_TOOLS_VERSION}
DEFINED_PACKAGE_LIST="v1beta1"


cleanup() {
  rm -rf "${TMP_DIFFROOT}"
}

trap "cleanup" EXIT SIGINT

cleanup

mkdir -p "${TMP_DIFFROOT}/api"
mkdir -p "${TMP_DIFFROOT}/config"
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}/api"
mkdir -p ${TMP_DIFFROOT}
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}"
cd ${SCRIPT_ROOT}

. ${SCRIPT_ROOT}/hack/update-gencode.sh
echo "diffing ${DIFFROOT} against freshly generated codegen"
ret=0
diff -Naupr "${DIFFROOT}" "${TMP_DIFFROOT}/api" || ret=$?
diff -Naupr "${DIFFROOT}" "${TMP_DIFFROOT}" || ret=$?
if [[ $ret -eq 0 ]]
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please run hack/update-gencode.sh"
  exit 1
fi

echo "diffing manifests  against freshly generated codegen"
ret=0
diff -Naupr "${DIFFMANIFESTSROOT}/crd/bases" "${TMP_DIFFMANIFESTSROOT}/crd/bases" || ret=$?
if [[ $ret -eq 0 ]]
then
  echo "manifests up to date."
else
  echo "manifests is out of date. Please run hack/update-gencode.sh"
  exit 1
fi


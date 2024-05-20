#!/bin/bash

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
# set -o nounset
set -o pipefail


SCRIPT_ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
CONTROLLER_TOOLS_VERSION="v0.14.0"
CONTROLLER_GEN=${SCRIPT_ROOT}/bin/controller-gen-${CONTROLLER_TOOLS_VERSION}
PACKAGE_LIST="v1beta1"

if [ "X${DEFINED_PACKAGE_LIST}" ==  "X" ]; then
  PACKAGE_LIST="v1beta1"
else
  PACKAGE_LIST=${DEFINED_PACKAGE_LIST}
fi

if [ "X${TMP_DIFFROOT}" == "X" ]; then
   PKGROOT="${SCRIPT_ROOT}/api"
   MANIFESTSROOT="${SCRIPT_ROOT}/config"
else
   PKGROOT="${TMP_DIFFROOT}/api"
   MANIFESTSROOT="${TMP_DIFFROOT}/config"
fi

chmod +x ${CONTROLLER_GEN}
for p in ${PACKAGE_LIST}
do
  ${CONTROLLER_GEN} object:headerFile="hack/boilerplate.go.txt" paths="${PKGROOT}/${p}"
  if [ $? -ne 0 ]; then
     echo "generating code for ${p} error"
     exit 1
  fi

  ${CONTROLLER_GEN} rbac:roleName=manager-role crd webhook paths="${PKGROOT}/${p}" output:crd:artifacts:config="${MANIFESTSROOT}/crd/bases"
  if [ $? -ne 0 ]; then
     echo "generating manifests for ${p} error"
     exit 1
  fi
done

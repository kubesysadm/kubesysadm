#!/bin/bash

# Copyright 2024 The kubeSysadm Authors.
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


# check if kubectl installed
function check-prerequisites {
  echo "Checking prerequisites"
  which kubectl >/dev/null 2>&1
  if [[ $? -ne 0 ]]; then
    echo -e "\033[31mERROR\033[0m: kubectl not installed"
    exit 1
  else
    echo -n "Found kubectl, version: " && kubectl version --short --client
  fi
}

function kind-up-cluster {
  check-kind

  echo "Running kind: [kind create cluster ${CLUSTER_CONTEXT} ${KIND_OPT}]"
  kind create cluster ${CLUSTER_CONTEXT} ${KIND_OPT}

  echo
  check-images

  echo
  echo "Loading docker images into kind cluster"
  kind load docker-image ${IMAGE_PREFIX}/kubesysadm-controller-manager:${TAG} ${CLUSTER_CONTEXT}
}

# check if kind installed
function check-kind {
  echo "Checking kind"
  which kind >/dev/null 2>&1
  if [[ $? -ne 0 ]]; then
    echo "Installing kind ..."
    go install sigs.k8s.io/kind@v0.15.0
  else
    echo -n "Found kind, version: " && kind version
  fi
}

# check if the required images exist
function check-images {
  echo "Checking whether the required images exist"
  docker image inspect "${IMAGE_PREFIX}/kubesysadm-controller-manager:${TAG}" > /dev/null
  if [[ $? -ne 0 ]]; then
    echo -e "\033[31mERROR\033[0m: ${IMAGE_PREFIX}/kubesysadm-controller-manager:${TAG} does not exist"
    exit 1
  fi

}

# install helm if not installed
function install-helm {
  echo "Checking helm"
  which helm >/dev/null 2>&1
  if [[ $? -ne 0 ]]; then
    echo "Installing helm via script"
    HELM_TEMP_DIR=$(mktemp -d)
    curl -fsSL -o ${HELM_TEMP_DIR}/get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 ${HELM_TEMP_DIR}/get_helm.sh && ${HELM_TEMP_DIR}/get_helm.sh
  else
    echo -n "Found helm, version: " && helm version
  fi
}
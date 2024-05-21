#!/bin/bash

#
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
#

export KS_ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/..
export LOG_LEVEL=3
export CLEANUP_CLUSTER=${CLEANUP_CLUSTER:-1}

NAMESPACE=${NAMESPACE:-kube-sysadm}
CLUSTER_NAME=${CLUSTER_NAME:-integration}

export CLUSTER_CONTEXT="--name ${CLUSTER_NAME}"
export SHOW_KUBESYSADM_LOGS=${SHOW_KUBESYSADM_LOGS:-1}
export CONTROLLER_IMAGE_NAME=${CONTROLLER_IMAGE:-kubesysadm-controller}
export CONTROLLER_IMAGE_TAG=${IMAGE_TAG:-latest}
export CONTROLLER_IMAGE_PREFIX=${IMAGE_PREFIX:-kubesysadm}

export KIND_OPT=${KIND_OPT:="--config ${KS_ROOT}/hack/e2e-kind-config.yaml"}

echo $* | grep -E -q "\-\-help|\-h"
if [[ $? -eq 0 ]]; then
  echo "Customize the kind-cluster name:

    export CLUSTER_NAME=<custom cluster name>  # default: integration

Customize kind options other than --name:

    export KIND_OPT=<kind options>

Disable displaying kubesysadm component logs:
    export SHOW_KUBESYSADM_LOGS=0
"
  exit 0
fi

if [[ $CLEANUP_CLUSTER -eq 1 ]]; then
    trap cleanup EXIT
fi

source "${KS_ROOT}/hack/lib/install.sh"

function get-k8s-server-version {
    echo $(kubectl version --short=true | grep Server | sed "s/.*: v//" | tr "." " ")
}

function install-kubesysadm {
  install-helm

  # judge crd version
  serverVersion=($(get-k8s-server-version))
  major=${serverVersion[0]}
  minor=${serverVersion[1]}
  crd_version="v1beta1"
  # if k8s version less than v1.18, crd version use v1beta
  # if [ "$major" -le "1" ]; then
  #  if [ "$minor" -lt "18" ]; then
  #    crd_version="v1beta1"
  #  fi
  # fi

  echo "Ensure create namespace"
  kubectl apply -f installer/namespace.yaml

  echo "Install kubesysadm chart"
  helm install ${CLUSTER_NAME} installer/helm/kubesysadm --namespace ${NAMESPACE} --kubeconfig ${KUBECONFIG} \
    --set image.imagePrefix=${CONTROLLER_IMAGE_PREFIX} \
    --set image.controllerManager=${CONTROLLER_IMAGE_NAME} \
    --set image.controllerManagerTag=${CONTROLLER_IMAGE_TAG} \
    --wait
}

function uninstall-kubesysadm {
  helm uninstall ${CLUSTER_NAME} -n ${NAMESPACE}
}

function generate-log {
    echo "Generating kubesysadm log files"
    kubectl logs deployment/${CLUSTER_NAME}-controller-manager -n ${NAMESPACE} > kubesysadm-controller-manager.log
}

function show-log() {
  log_files=("kubesysadm-controller-manager.log" )
  for log_file in "${log_files[@]}"; do
    if [ -f "$log_file" ]; then
      echo "Showing ${log_file}..."
      cat "$log_file"
    else
      echo "${log_file} not found"
    fi
  done
}


# clean up
function cleanup {
  uninstall-kubesysadm

  echo "Running kind: [kind delete cluster ${CLUSTER_CONTEXT}]"
  kind delete cluster ${CLUSTER_CONTEXT}

  if [[ ${SHOW_KUBESYSADM_LOGS} -eq 1 ]]; then
    show-log
  fi
}

check-prerequisites
kind-up-cluster

if [[ -z ${KUBECONFIG+x} ]]; then
    export KUBECONFIG="${HOME}/.kube/config"
fi

install-kubesysadm

# Run e2e test
cd ${KS_ROOT}

install-ginkgo-if-not-exist

case ${E2E_TYPE} in
"ALL")
    echo "Running e2e..."

    ;;
"KSCTL")
    echo "Running ksctl e2e suite..."
    KUBECONFIG=${KUBECONFIG} KIND_CLUSTER=${CLUSTER_NAME} ginkgo -r --slow-spec-threshold='30s' --progress ./test/e2e/
    ;;
esac

if [[ $? -ne 0 ]]; then
  generate-log
  exit 1
fi
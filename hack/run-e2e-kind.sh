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
CLUSTER_NAME=${CLUSTER_NAME:-integration}\

export CLUSTER_CONTEXT="--name ${CLUSTER_NAME}"

export KIND_OPT=${KIND_OPT:="--config ${KS_ROOT}/hack/e2e-kind-config.yaml"}

echo $* | grep -E -q "\-\-help|\-h"
if [[ $? -eq 0 ]]; then
  echo "Customize the kind-cluster name:

    export CLUSTER_NAME=<custom cluster name>  # default: integration

Customize kind options other than --name:

    export KIND_OPT=<kind options>

Disable displaying volcano component logs:

    export SHOW_VOLCANO_LOGS=0
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

  echo "Install volcano chart with crd version $crd_version"
  helm install ${CLUSTER_NAME} installer/helm/kubesysadm --namespace ${NAMESPACE} --kubeconfig ${KUBECONFIG} \
    --set basic.image_pull_policy=IfNotPresent \
    --set basic.image_tag_version=${TAG} \
    --set basic.scheduler_config_file=config/volcano-scheduler-ci.conf \
    --set basic.crd_version=${crd_version} \
    --wait
}


check-prerequisites
kind-up-cluster

if [[ -z ${KUBECONFIG+x} ]]; then
    export KUBECONFIG="${HOME}/.kube/config"
fi

install-kubesysadm

# Run e2e test
cd ${KS_ROOT}
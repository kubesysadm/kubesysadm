<a href="https://www.sysadm.cn">
    <img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/kubesysadm.png"/>
</a>


---
[![Go Report Card](https://goreportcard.com/badge/github.com/kubesysadm/kubesysadm)](https://goreportcard.com/report/github.com/kubesysadm/kubesysadm)
[![RepoSize](https://img.shields.io/github/repo-size/kubesysadm/kubesysadm.svg)](https://github.com/kubesysadm/kubesysadm)
[![Release](https://img.shields.io/github/release/kubesysadm/kubesysadm.svg)](https://github.com/kubesysadm/kubesysadm/releases)
[![LICENSE](https://img.shields.io/github/license/kubesysadm/kubesysadm.svg)](https://github.com/kubesysadm/kubesysadm/blob/main/LICENSE)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9148/badge)](https://www.bestpractices.dev/projects/9148)

# kubesysadm(Kubernetes System Administration)
[kubesysadm](https://kubesysadm.sysadm.cn/) is a management tool for managing kubernete's resource. It is short for 
KUBErnete SYStem ADMInistration. And it provides a suite mechanisms and methods to manage resource of kubernetes.

Kubesysadm is based on kubernetes Operator technoloygy. And it focuses on automated operations and maintenance of kubernetes clusters. 


## Features
Now the features of kubesysadm described as the following:
- Automatically restart the workload gracefully
  
  we know that the workload(Deployment, statefulSet, DaemaonSet) does not automatically restart gracefully when the configMap/secret
  referenced by the workload changes. This results in updates to configMap not being applied to the workload in a timely manner.
  
  kubesysadm can monitor the configMaps/secrets referenced by workloads based on user-configured rules.  Kubesysadm will restart 
  the workload gracefully when it finds the configMap/secret referenced by the workload has changed.

- Automatically cleaning no-running Pods
  We know that no-running(Such as Completed, Error) Pods in K8S cluster are not be deleted automatically. Kubesysadm can 
  delete the no-running pods according to the rules which be configured by user.

- Pod cleaning rules configuration
  User can configure the rules for deleting no-running

- Monitoring rule configuration
  users can configure the monitoring rules.

  
## Quick Start Guide

### Prerequisites
- kubectl version v1.12+ with CRD support.
- Access to a Kubernetes v1.12+ cluster.

### Install CRD and instance of Kubesysadm into  cluster
Install Kubesysadm on an existing Kubernetes cluster. This way is both available for x86_64 and arm64 architecture.
```
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/install.yaml
```

Enjoy! Kubesysadm will create the following resources in the cluster
```
namespace/kubesysadm-system created
customresourcedefinition.apiextensions.k8s.io/cmmonitors.monitoring.sysadm.cn created
customresourcedefinition.apiextensions.k8s.io/podcleanrules.monitoring.sysadm.cn created
serviceaccount/kubesysadm-controller-manager created
role.rbac.authorization.k8s.io/kubesysadm-leader-election-role created
clusterrole.rbac.authorization.k8s.io/kubesysadm-cmmonitor-editor-role created
clusterrole.rbac.authorization.k8s.io/kubesysadm-cmmonitor-viewer-role created
clusterrole.rbac.authorization.k8s.io/kubesysadm-manager-role created
clusterrole.rbac.authorization.k8s.io/kubesysadm-metrics-reader created
clusterrole.rbac.authorization.k8s.io/kubesysadm-proxy-role created
rolebinding.rbac.authorization.k8s.io/kubesysadm-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/kubesysadm-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/kubesysadm-proxy-rolebinding created
service/kubesysadm-controller-manager-metrics-service created
deployment.apps/kubesysadm-controller-manager created
```
Check whether kubesysadm controller is running by the following command
```
kubectl get po -n kubesysadm-system
```

The output of the above command like the following showing
```
NAME                                             READY   STATUS    RESTARTS   AGE
kubesysadm-controller-manager-5f78865594-b6gsc   2/2     Running   0          2m30s
```

### Configure configMap monitoring rules
In the following example we will to do the following things:
- Create a namespace named test-kubesysadm
- Create a configMap named cm-env in test-kubesysadm namespace
- Create a configMap named cm-mount in test-kubesysadm namespace
- Create configMap monitoring rules named cm-env and cm-mount
- Create a deployment named test-kubesysadm in test-kubesysadm namespace which using cm-env and cm-mount configMap

We create the above resource by the following commands:
```azure
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/create_ns.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/cm-env.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/cm-mount.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/cmmonitor_cm-env.yaml    
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/cmmonitor_cm-mount.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/configMap/deploy.yaml
```

Now we try to check the results
We get the pods' status by the following command
```azure
  kubectl get po -n test-kubesysadm
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/mornitoring_cm_started.png"/>

Then we try to change cm-env/cm-mount configMap by the following command. After that, we will find that the pods of
deployment test-kubesysadm will be restared like the following image shown.
```azure
  kubectl edit cm -n test-kubesysadm cm-env/cm-mount 
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/mornitoring_cm_restart.png"/>

### Configure secret monitoring rules
Like configMap, in the following example we will to do the following things:
- Create a namespace named test-kubesysadm
- Create a secret named secret-env in test-kubesysadm namespace
- Create a secret named secret-mount in test-kubesysadm namespace
- Create secret monitoring rules named secret-env and secret-mount
- Create a deployment named test-kubesysadm in test-kubesysadm namespace which using secret-env and secret-mount secret

We create the above resource by the following commands:
```azure
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/create_ns.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/secret_env.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/secret_mount.yaml 
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/cmmonitor_secret-env.yaml    
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/cmmonitor_secret-mount.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/secret/deploy.yaml
```

Now we try to check the results
We get the pods' status by the following command
```azure
  kubectl get po -n test-kubesysadm
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/mornitoring_cm_started.png"/>

Then we try to change secret-env/secret-mount secret by the following command. After that, we will find that the pods of
deployment test-kubesysadm will be restared like the following image shown.
```azure
  kubectl edit secret -n test-kubesysadm secret-env/secret-mount
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/mornitoring_cm_restart.png"/>


### Configure Pod Cleaning rules
In the following example we will to do the following things:
- Create a namespace named test-kubesysadm
- Create a Job named job1
- Create a PodCleanRule named cleanpods. The age value of rule is 300 and namespace is test-kubesysadm. That meaning is 
  PodCleanManager will clean the pods which created before 5 minutes and in no-running/no-pending status.

We create the above resource by the following commands:
```azure
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/podclean/create_ns.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/podclean/podcleanrule.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/installer/podclean/job.yaml
```

Now we try to check the results
We get the pods' status by the following command
```azure
  kubectl get po -n test-kubesysadm
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/job_pods.png"/>

We found that all pods in "Completed" status in test-kubesysadm namespace have be deleted after 5 minutes when we re-run 
the above command. And we get the log message like the following image shown when we run the following command.
```azure
   kubectl logs -n kubesysadm-system kubesysadm-controller-manager-5f78865594-b6gsc
```

<img src="https://raw.githubusercontent.com/kubesysadm/kubesysadm/main/docs/images/podcleanmessage.png"/>

## Contributing/贡献
### English
All those who are interested in the kubesysadm project are welcome to contribute to kubesysadm.
We encourage you to communicate in English, but do not exclude Chinese.
The [Contributor Guide](https://raw.githubusercontent.com/kubesysadm/community/contribute.md) provides detailed instruction on how to get your ideas and bug fixes seen and accepted, including:

1. How to [find something to work on](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#find-something-to-work-on)
1. How to [create a pull request](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#creating-pull-requests)
1. How to [code review](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#code-review)

If you're interested in being a contributor and want to get involved in
developing the Kubesysadm code, please see [contribute](https://raw.githubusercontent.com/kubesysadm/community/contribute.md) for
details on submitting patches and the contribution workflow.

More information can be found via the [community](https://github.com/kubesysadm/community)

### 中文
我们非常欢迎所有对kubesysadm项目感兴趣的人为kubesysadm做出贡献。
我们鼓励您用英文进行沟通，但是不排斥中文。
[贡献者向导](https://raw.githubusercontent.com/kubesysadm/community/contribute.md) 为您提供了一个详细的说明，以便我们更容易的接受您的想法或者您为修正Bug而做的贡献，它包括：

1. [如何找到适合自已做的事情](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#find-something-to-work-on)
2. [如何创建一个PR](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#creating-pull-requests)
3. [如何Review代码](https://raw.githubusercontent.com/kubesysadm/community/contribute.md#code-review)

如果你有兴趣成为一个贡献者，并希望参与Kubesysadm代码的开发，请参阅提交补丁和贡献工作流程细节
[contribute](https://raw.githubusercontent.com/kubesysadm/community/contribute.md) 

更多信息情参阅我们的社区 [community](https://github.com/kubesysadm/community)

## License

Copyright 2024 Wayne Wang<net_use@bzhy.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


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


- Monitoring rule configuration
  users can configure the monitoring rules.

  
## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kubesysadm:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kubesysadm:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/kubesysadm:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/kubesysadm/<tag or branch>/dist/install.yaml
```

## Contributing

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [community](https://github.com/kubesysadm/community)

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


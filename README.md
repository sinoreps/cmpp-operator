[![Build Status][1]][2]

## Overview

CMPP-Operator is the tool to simplify the deployment and management of http-cmpp-proxy clients on your Kubernetes cluster.

It is designed on top of the kubernetes CRD capabilities and runs as a series of deployment resources.. Although currently CMPPProxy is the only supported model, CRDs for more CMPP functionalities can be extended soon. With CMPPProxy as a CRD, configurations for one CMPP account is kept within one k8s resource object, acting as the single source of truth. To handle multiple CMPP account configurations, you can create multiple CMPPProxies, without the need to touch deployments/pods/services under the hood, thus keeping your hands clean.

## Get started

### Prerequisites
* Access to a Kubernetes cluster, version 1.7 or later. 
* kubectl installed

### Installation
#### Clone the repo
```bash
$ git clone https://github.com/sinoreps/cmpp-operator.git && cd cmpp-operator
```
#### Setup Service Account
```bash
$ kubectl create -f deploy/service_account.yaml
```
#### Setup RBAC
```bash
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
```

#### Setup the CRD
```bash
$ kubectl create -f deploy/crds/cmpp_v1alpha1_cmppproxy_crd.yaml
```
#### Deploy the cmpp-operator
```bash
$ kubectl create -f deploy/operator.yaml
```

### Add a CMPPProxy
```bash
$ kubectl create -f -<<EOF
apiVersion: cmpp.io/v1alpha1
kind: CMPPProxy
metadata:
  name: example-cmppproxy
spec:
  image: sinoreps/cmppproxy:latest
  account: "333"
  password: "0555"
  numConnections: 1
  serverAddr: "127.0.0.1"
  serviceCode: "9999"
  enterpriseCode: "044022"
EOF
```

#### check the newly created pods
```bash
$ kubectl get pods -w
```
The deployment and its pods will be created in your currently active namespace.

## Contributing

Feel free to [file an issue][4] if you encounter issues, or create [pull requests][6].


[1]: https://travis-ci.org/sinoreps/cmpp-operator.svg?branch=master
[2]: https://travis-ci.org/sinoreps/cmpp-operator
[4]: https://github.com/sinoreps/cmpp-operator/issues
[6]: https://github.com/sinoreps/cmpp-operator/pulls

# Alexa Service

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-contrib/alexa-service)
[![Build Status](https://travis-ci.org/keptn-contrib/alexa-service.svg?branch=master)](https://travis-ci.org/keptn-contrib/alexa-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-contrib/alexa-service)](https://goreportcard.com/report/github.com/keptn-contrib/alexa-service)

The *alexa-service* is a [Keptn](https://keptn.sh) service for sending Keptn events to Amazon Alexa.

# Setup

## Deploy in your Kubernetes cluster

First create a Kubernetes secret containing your Alexa webhook url and token.
You can therefore use the file `deploy/secret.yaml` as template, which afterwards has to be applied:
```console
kubectl apply -f deploy/secret.yaml
``` 

To deploy the current version of the *alexa-service* in your Keptn Kubernetes cluster,
use the files `deploy/service.yaml` and `deploy/distributor.yaml`
from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
kubectl apply -f deploy/distributor.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *prometheus-service*, use the file `deploy/service.yaml` and `deploy/distributor.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
kubectl delete -f deploy/distributor.yaml
```


# Contributions

You are welcome to contribute using Pull Requests against the **master** branch. Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md).

# Travis-CI setup

Travis is configured with CI to automatically build docker images for pull requests and commits. The  pipeline can be viewed at https://travis-ci.org/keptn-contrib/alexa-service.
The travis pipeline needs to be configured with the `REGISTRY_USER` and `REGISTRY_PASSWORD` variables. 

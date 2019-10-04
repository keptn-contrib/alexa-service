# Alexa Service

The *alexa-service* is a Keptn service for sending Keptn events to Amazon Alexa.

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
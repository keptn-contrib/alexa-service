# Alexa Notification Service

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn-sandbox/alexa-notification-service)
[![Build Status](https://travis-ci.org/keptn-sandbox/alexa-notification-service.svg?branch=master)](https://travis-ci.org/keptn-sandbox/alexa-notification-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn-sandbox/alexa-notification-service)](https://goreportcard.com/report/github.com/keptn-sandbox/alexa-notification-service)

| Authors | Keptn Version | Comment |
| ------ | --------------| -------- |
| [@alipatton10](https://github.com/alipatton10) | 0.7.0 | Initial Release |

The *alexa-notification-service* is a [Keptn](https://keptn.sh) service for sending Keptn events to an Amazon Alexa. This skill is primarilly designed for demoing the capabilities of keptn. When certain keptn events are triggered a push notification will be sent to your Amazon Alexa account. You can then ask the Alexa to read your notification describing the current state. The table below describes the current supported events and the messages you should expect to receive

| Keptn Event | Default State | Sample Message | Comment |
| ------ | --------------| -------- | -------- |
| configuration-change | Disabled | New Keptn event detected. Configuration change, has been reported for carts in staging. | This message is static with the exception of the service and stage names which are dynamically populated. |
| deployment-finished | Enabled | New Keptn event detected. Deployment finished, has been reported for carts in staging. | This message is static with the exception of the service and stage names which are dynamically populated. |
| tests-finished | Disabled | New Keptn event detected. Tests finished, has been reported for carts in staging | This message is static with the exception of the service and stage names which are dynamically populated. |
| evaluation-done | Enabled | New Keptn event detected. Evaluation done, has been reported for carts in staging. The result of the evaluation was pass. Promoting artifact to next stage. | The message will call out the result of the evalualtion and the service and stage names. If it has passed and is not in production then the message will also include Promoting artifact to next stage, if it failes the message will include \"The artifact will not be promoted from staging to next stage.\" |
| problem-distributor | Enabled | New problem reported by Dynatrace. P.I.D. 803 . Failure rate increase . The impact is failure rate increase on items controller staging. | Supports Open, Closed and State change problems from dynatrace. The message will differ dependon on this.

# Setup

## Preperation

This services utilises a third party Alexa skill [Notify Me](https://www.thomptronics.com/about/notify-me). In order to use it you must first add the skill to your alexa account. To get started, just enable the skill from the Alexa app or from the Amazon website, giving it permission to send notifications to your Alexa and linking it to your Amazon account so we can send a unique access code to the email address associated with that account.  Next, just say, "Alexa, open Notify Me" and the skill will introduce itself and send your access code via email. Please note you will not recieve your access code until you launch the skill for the first time.

If it seems to be taking too long to receive the access code, please check your SPAM filter.  The email will come from from notifyme@notifymyecho.com via amazonses.com. Remember: launching the skill \(\"Alexa, open Notify Me\" or \"Alexa, launch the Notify Me skill\"\) is what triggers the email with the access code.

Once you have your access code modify the `deploy/secret.yaml` file and your access code to line 8 after "token: ". Ensure to leave a space after :

**Optional** As detailed above 2 notification types are disabled by default. If you wish to change this then please uncomment them in the file `deploy/distributor.yaml`

## Deploy in your Kubernetes cluster

First create a Kubernetes secret containing your Notify Me access token.
Ensure you have added your access token to the file `deploy/secret.yaml`
```console
kubectl apply -f deploy/secret.yaml -n keptn
``` 

To deploy the current version of the *alexa-notification-service* in your Keptn Kubernetes cluster,
use the files `deploy/service.yaml` and `deploy/distributor.yaml`
from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
kubectl apply -f deploy/distributor.yaml
```

To test our your skill trigger a new a

## Delete in your Kubernetes cluster

To delete a deployed *prometheus-service*, use the file `deploy/service.yaml` and `deploy/distributor.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
kubectl delete -f deploy/distributor.yaml
```

# Limitations

1. Notify Me is available in all areas that allow English-speaking skills. As of this writing, these locations include the United States, Canada, United Kingdom, Australia, and India. Your Amazon Alexa account must be registered in one of these locations.
2. Amazon enforces a per-user limit of 5 notifications in a 5 minute period. This is why by default we disable some of the notifications. If you exceed this limit the notificatiosn are surpressed and will not be send again once the timer resets.

# Contributions

You are welcome to contribute using Pull Requests against the **master** branch. Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md).

# Travis-CI setup

Travis is configured with CI to automatically build docker images for pull requests and commits. The  pipeline can be viewed at https://travis-ci.org/github/keptn-sandbox/alexa-notification-service.
The travis pipeline needs to be configured with the `REGISTRY_USER` and `REGISTRY_PASSWORD` variables. 

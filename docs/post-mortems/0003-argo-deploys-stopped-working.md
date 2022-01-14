## Summary

2022-01-11 ArgoCD stopped doing deployments. Applications would be in an endless "Syncronizing state", without being able to
deploy new changes pushed to Git. This lasted for 24 hours, then it started working again.

## Timeline

2022-01-11
* 00:00 ArgoCD stops doing automatic deployments.
* 09:20 A user notifies us on Slack about the issue. [Thread link](https://oslokommune.slack.com/archives/CV9EGL9UG/p1641889201058000).

2022-01-12
* 00:00 ArgoCD deployments starts working again by itself. 

## Impact

* ArgoCD deployments stopped working for 24 hours. Only one team notified us about this problem, but in theory all our teams
should have had this problem.

## Cause(s)

The error message from ArgoCD (the argocd-repo-server pod), was:

time="2022-01-11T00:03:33Z" level=error msg="finished unary call with code Unknown" error="unknown error: ERROR: You're using an
RSA key with SHA-1, which is no longer allowed. Please use a newer client or a different key type." grpc.code=Unknown
grpc.method=GenerateManifest grpc.request.deadline="2022-01-11T00:04:31Z" grpc.service=repository.RepoServerService
grpc.start_time="2022-01-11T00:03:31Z" grpc.time_ms=1359.753 span.kind=server system=grpc

The issue is we're using the SHA-1 hashing algorithm for our GitHub deploy keys. GitHub deactivated support for SHA-1 keys for 24
hours the 11th january 2022, in a so-called "brownout", to alert clients still using these keys to replace them. See: [1] and [2]
and

The deploy key is created by Okctl (see
[keypair.go](https://github.com/oslokommune/okctl/blob/2173be74b104abd36b9b9b17b12b3f3d9f41fed2/pkg/keypair/keypair.go)), and is a
SSH-key using SHA-1. The public key resides in every IAC-repository in GitHub, while the private key is injected as a secret to
ArgoCD. ArgoCD uses this to pull changes from GitHub. So when ArgoCD used the key to pull the IAC repository, GitHub responded
with the above error message.

## Solution

We need to update Okctl to use a new key format for the deploy key. [1] suggests ed25519, referring
to https://security.stackexchange.com/a/144044/29078, which appears to be a solid choice. Then we need to roll out this change,
preferably through an upgrade.

## What went well


## What went wrong

* We were not aware in any way of this breaking change from GitHub. They linked blog post above is the only message I could find
about it. We have now subscribed to this blog, but the blog is filled with non-relevant noise, and we're already following
other blogs from Kubernetes and AWS, which adds to the noise.

## Where were we lucky

* GitHub used a "brownout" for this, meaning they pushed the change only for 1 day, for then to restore support for the SHA-1 key
  format. The intention is to force clients to notice the issue. They will implement the actual, permanent change 2022-03-15, so
  we need to come up with a fix in good time before this, so our teams can roll out a fix before this happens. 

## Mitigations

<!--
Describe potential steps to prevent this event from happening again, be it technical, processes, etc.

Example:
- Request that proactive messages will be sendt in smaller batches, ideally during regular work hours, for load to be
    better distributed.
- Create tests that cover creation and authentication of users
-->

* We have subscribed to the GitHub blog in a Slack bot that forwards posts to #kjøremiljø-feeds. I advise all members of our team
to also have a personal system of following RSS news from GitHub and other important sources like AWS and Kubernetes.
https://newsblur.com is one option.
  * This may not be enough though, so suggestions are welcome.

## Action items

[KM498](https://trello.com/c/9GxBElfD/498-argocd-stops-working-2022-03-15-due-to-old-format-of-deploy-key) - ArgoCD stops working
2022-03-15 due to old format of deploy key

## Supporting information and documentation

* [1] https://github.blog/2021-09-01-improving-git-protocol-security-github/#when-are-these-changes-effective
* [2] https://github.com/argoproj/argo-cd/issues/7600
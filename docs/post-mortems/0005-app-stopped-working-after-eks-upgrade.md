## Summary
<!--
Describe the reason and the consequences of the event as short and concise as possible.

Example:
The backend of Oslo Nøkkelen were unavailable for 31 minutes due to an increase in traffic that happened after a
proaktiv melding.
-->

Team Booking and Kjøremiljø together upgraded Booking's development environment from EKS 1.19 to 1.20, using the guide https://github.com/oslokommune/okctl-upgrade/blob/64920127a5cd68a9af9f2e0b2396d00fb8ef03f6/gists/bump-eks-to-1-20/README.md.

After running the final step, the backend application become unresponsive. The issue was outdated application configuration, specifically using these values in the application's configuration:
* `dnsPolicy: Default`
* missing an EKS control plane security group as part of the application's `SecurityGroupPolicy`.

We don't know why, but using these settings in EKS 1.19 worked fine, but not in 1.20. 

## Timeline

### 2022-06-13

Timestamps are approximate.

* 15:15: We completed the final step in the upgrade guide. Some time after, we discovered that the frontend application didn't work.

* 15:30: We saw that `kubectl describe pod booking-server-...` contained warnings looking like this:

```
Events:
Type     Reason                  Age   From                     Message
  ----     ------                  ----  ----                     -------
Normal   Scheduled               29m   default-scheduler        Successfully assigned booking/booking-server-56974858b9-27zcw to ip-xxxxxxxxx.xxxxxxxx.compute.internal
Normal   SecurityGroupRequested  29m   vpc-resource-controller  Pod will get the following Security Groups [sg-xxxxxxxxxxxx]
Warning  FailedCreatePodSandBox  29m   kubelet                  Failed to create pod sandbox: rpc error: code = Unknown desc = failed to set up sandbox container "5249d574a203073df639fd024ca0480740d39cc532f2ca7d88ca828ec865656b" network for pod "booking-server-56974858b9-27zcw": networkPlugin cni failed to set up pod "booking-server-56974858b9-27zcw_booking" network: add cmd: failed to assign an IP address to container
Normal   ResourceAllocated       29m   vpc-resource-controller  Allocated [{"eniId":"eni-xxxxxxxxxxxxxx","ifAddress":"0a:xxxxxxxxxxxxx","privateIp":"192.168.xxxxxxxxxxx","vlanId":1,"subnetCidr":"192.168.xxxxxxxxxxxx/yyy"}] to the pod
Normal   SandboxChanged          29m   kubelet                  Pod sandbox changed, it will be killed and re-created.
Normal   Pulled                  29m   kubelet                  Container image "xxxxxxxxxxxxxx.dkr.ecr.xxxxxxxxx.amazonaws.com/booking-server:xxxxxxxxxxxxxxxx" already present on machine
Normal   Created                 29m   kubelet                  Created container booking-server
Normal   Started                 29m   kubelet                  Started container booking-server
```

* 16:00 AWS' support contact offered to help debugging.

* 16:19 Kjøremiljø met with AWS' contact. During this session, we realized that the SecurityGroupPolicy for the failing application was missing a SecurityGroup. Kjøremiljø, having fixed issues with this [before](https://trello.com/c/y8VyNwRy/427-dns-doesnt-work-for-pods-in-securitygroups), remembered that applications using databases need two security groups, but the failing application only had one.

* 16:35 We tested setting the correct security group, and verified that the application worked.

## Impact
<!--
Describe the consequences this had for the organization.

Example:
- Approximately 30.000 users downloaded and experienced an error using the app.
- Users already logged in experienced some minor delay in the app
-->

Booking's backend in development stopped working between 15:15 and 16:35. 

Booking's frontend returned a web page, but it wasn't working properly.

## Cause(s)
<!--
Describe relevant factors that played a part in causing the event.

Example:
The error happened due to a combination of:
- A spike of new users
- A bug in the backend code for user creation
-->

The error happened due to a combination of using EKS 1.20 together with application configuration (Kubernetes manifests) that were outdated.

## Solution
<!--
Describe how the problem was fixed.

Example:
- Manually scaled up number of backend instances letting regular traffic stay unaffected
- Created a fix for the database schema
-->

We updated to the correct application configuration:
* removed `dnsPolicy: true` from all three application's Deployment.
* Added `terminationGracePeriodSeconds: 30` to the backend Deployment.
* Added the correct cluster security group to the backend's SecurityGroupPolicy.

## What went well
<!--
Describe what went well trying to handle the event.

Example:
- Alerting mechanisms worked brilliantly when errors started comming in
- Deploying the database update was fast
-->

* We knew by memory what a correct application configuration was. This would also be possible to aquire by running `okctl apply application` with an existing Application manifest, or generating a new on with `okctl scaffold application`. 
* We ran the EKS upgrade in the development environment before the production environment
* In the AWS account where we did the upgrade, we only has "Basic" support level, which probably would have been to slow to use. Luckily, someone on Slack was available to tell how to upgrade it - though we ended up fixing it before this became necessary.  

## What went wrong
<!--
Describe what went wrong trying to handle the event.

Example:
- Had to manually roll back the database while testing fixes
- We missed the notification regarding Keycloak upgrade
-->

* We didn't check that the application configuration was up-to-date before running through the EKS upgrade guide.
* Not specifically wrong, but trying to update the Kubernetes application manifests in one environment (dev) while not affecting production is a pain with Kustomize - when the YAML you really want to change is in folder that is common for both environments (`base`). You have to use overlays, whose syntax is totally different from the base files.

## Where were we lucky
<!--
Describe briefly the situations related to this event where we got lucky.

Example:
- Updating the database schema fixed the problem
- Manual horizontal scale up ensured a small amount of users was affected
-->

* The person in Kjøremiljø who ran the EKS upgrade, was the same one as implemented the [security group issue](https://trello.com/c/y8VyNwRy/427-dns-doesnt-work-for-pods-in-securitygroups), and therefore knew in detail what was wrong. 
* The AWS contact for Oslo kommune was available to help.

## Mitigations
<!--
Describe potential steps to prevent this event from happening again, be it technical, processes, etc.

Example:
- Request that proactive messages will be sendt in smaller batches, ideally during regular work hours, for load to be
    better distributed.
- Create tests that cover creation and authentication of users
-->

We added more pre-checks to the EKS upgrade guide: https://github.com/oslokommune/okctl-upgrade/commit/bf64dc14d5f995bc69071a1e26e1070436eab005

## Action items
<!--
Describe and reference specific action items

Example:
[KM123](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add backoff for login in the iOS app
[KM221](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add tests in the backend for creating a user and authenticating
-->

The mitigation above should be enough.

## Supporting information and documentation

<!--
Add relevant additional information or documentation.

Example: images of metrics that show when something went wrong, relevant pieces of logs, etc
-->

Upgrade guide: https://github.com/oslokommune/okctl-upgrade/blob/64920127a5cd68a9af9f2e0b2396d00fb8ef03f6/gists/bump-eks-to-1-20/README.md

Slack thread: https://oslokommune.slack.com/archives/C01JPNWHF7H/p1655128125086249?thread_ts=1653295843.929539&cid=C01JPNWHF7H

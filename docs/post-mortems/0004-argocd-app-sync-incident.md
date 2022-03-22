## Summary
<!--
Describe the reason and the consequences of the event as short and concise as possible.

Example:
The backend of Oslo Nøkkelen were unavailable for 31 minutes due to an increase in traffic that happened after a
proaktiv melding.
-->

Two things happened:

1. When running the ArgoCD app sync upgrade (okctl-upgrade_0.0.88.activate-argo-app-sync-v2) the upgrade ended up
    migrating files concerning other environments not specified by `-c` as well
2. When running `delete application`, the command ended up deleting the namespace where other applications were running
    causing the other applications to go down.

## Timeline
<!--
Describe the relevant activites in a timeline format. Remember activities
leading up to the event being triggered. Make sure to link to graphs, logs and
other relevant information sources.

2019-12-30
    23:30 A proactive message got sendt out to all the citizens of Oslo
2019-12-31
    08:23 Oslonøkkelen had been downloaded 30.000 times, distributed evenly across iOS and Android, during the past 3 hours
    09:09 DOWNTIME START - backend fails due to the increased load
    09:14 Received Slack notifications due to an increase in 500 status codes in the backend
    09:16 Initiated an investigation into the issue
    09:18 EVENT START - Team creates the Slack channel #citykey-incident-backend to better cooperate on the issue
    09:20 Manually scaling up the number of backend instances
    09:24 Found the error. The backend fails when it tries to read an item in the database that doesn't exist. Only happens upon a users first login
    09:28 Pull request with a fix created: http://github.com/oslokommune/something/something/pr/298
    09:34 Pull request 298 merged and deployed to production
    09:39 PROBLEM FIXED - Amount of 500 status code requests are decreasing
    09:40 DOWNTIME END - Amount of 500 status codes are down to zero
    10:10 EVENT END - Reaches criteria of 30 minutes with normal activity
-->
### 2022-03-07

#### Incident 1

1. Team booking ran `okctl upgrade` from version 0.0.82, 0.0.92 being the latest.
2. After a successful upgrade, we noticed synchronization errors in ArgoCD.
3. Researching these errors, we noticed that the app sync upgrade had moved all `argocd-application.yaml` files in
    `infrastructure/applications` to `infrastructure/<relevant cluster/argocd/applications/<application name>` which
    included `argocd-application.yaml` files from _all_ environments.
4. After manually reverting changes and migrating them, ArgoCD once again was OK.

#### Incident 2

1. Team booking ran `okctl delete application -f <app manifest>` to delete an old application which hasn't been in use
    for a while.
2. Delete application naively purged all resources related to the old application, where one of them was the
    namespace manifest. The old application shared namespace with new and in use application causing these applications
    to go down.

## Impact
<!--
Describe the consequences this had for the organization.

Example:
- Approximately 30.000 users downloaded and experienced an error using the app.
- Users already logged in experienced some minor delay in the app
-->

Team Booking's applications in the development environment went down for approximately five minutes while ArgoCD was
realizing that there were other applications that had specified the same namespace as well, and therefore brought it
back up.

## Cause(s)
<!--
Describe relevant factors that played a part in causing the event.

Example:
The error happened due to a combination of:
- A spike of new users
- A bug in the backend code for user creation
-->

### Incident 1

A bug in the ArgoCD app sync upgrade causing the upgrade to ignore what environment the files were related to, migrating
them all. One can speculate in why we didn't pick this up before. Maybe due to Team Kjøremiljø developers generally only
having one dev environment available?

### Incident 2

Fault in the feature design.

## Solution
<!--
Describe how the problem was fixed.

Example:
- Manually scaled up number of backend instances letting regular traffic stay unaffected
- Created a fix for the database schema
-->

### Incident 1

Manually reverted and migrated relevant ArgoCD application manifests.

### Incident 2

Nothing. ArgoCD realised other applications were dependent on the namespace and that it was missing, so ArgoCD
automatically brought it back up.

## What went well
<!--
Describe what went well trying to handle the event.

Example:
- Alerting mechanisms worked brilliantly when errors started comming in
- Deploying the database update was fast
-->

* Started upgrading the development environment first.
* ArgoCD ensured the namespace got back up after we deleted it.
* ArgoCD application manifests did not have the "Purge" flag set, meaning even if the ArgoCD application manifests
    disappeared from the git repository, the production apps were not affected.

## What went wrong
<!--
Describe what went wrong trying to handle the event.

Example:
- Had to manually roll back the database while testing fixes
- We missed the notification regarding Keycloak upgrade
-->


## Where were we lucky
<!--
Describe briefly the situations related to this event where we got lucky.

Example:
- Updating the database schema fixed the problem
- Manual horizontal scale up ensured a small amount of users was affected
-->

* Lucky Booking requested support before trying the upgrade.
* Lucky that the person who created the upgrade and the new feature were present.

## Mitigations
<!--
Describe potential steps to prevent this event from happening again, be it technical, processes, etc.

Example:
- Request that proactive messages will be sendt in smaller batches, ideally during regular work hours, for load to be
    better distributed.
- Create tests that cover creation and authentication of users
-->

- Potentially more in depth design QA before implementation
- Spend more time envisioning worst case, and create a test preventing it

## Action items
<!--
Describe and reference specific action items

Example:
[KM123](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add backoff for login in the iOS app
[KM221](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add tests in the backend for creating a user and authenticating
-->

[KM620](https://trello.com/c/AyCdNq4y) - 0.0.88 upgrade fails: dev and prod files get mixed

[KM621](https://trello.com/c/wjF1XyfK) - `okctl delete` application removes the whole namespace

## Supporting information and documentation
<!--
Add relevant additional information or documentation.

Example: images of metrics that show when something went wrong, relevant pieces of logs, etc
-->

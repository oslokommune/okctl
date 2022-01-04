## Summary
<!--
Describe the reason and the consequences of the event as short and concise as possible.

Example:
The backend of Oslo Nøkkelen were unavailable for 31 minutes due to an increase in traffic that happened after a
proaktiv melding.
-->
Users were unable to deploy updates to their application through ArgoCD.

## Timeline
<!-- 
Describe the relevant activites in a timeline format. Remember activities leading up to the event being triggered.

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

### 2021-09-1

| time  | description                                                    |
|-------|----------------------------------------------------------------|
| --:-- | Github releases a blog post describing the changes and the fix |

### 2021-11-17

| time  | description                                                                                            |
|-------|--------------------------------------------------------------------------------------------------------|
| --:-- | ArgoCD releases a version with updated known_hosts                                                     |
| 03:35 | Nightly fails due to `cloning repository: ssh: handshake failed: knownhosts: key mismatch`<sub>2</sub> |
| 16:26 | User reports error in #kjøremiljø-support<sub>1</sub>                                                  |
| 16:29 | User agrees that we can wait until the day after to look into it                                       |

### 2021-11-18

| time  | description                                                                                  |
|-------|----------------------------------------------------------------------------------------------|
| 08:33 | The team identifies the problem - Github rotated the certificate for their servers           |
| 08:53 | The team identifies a solution, and announces it and an impending upgrade - upgrading ArgoCD |
| 12:27 | Another user posts a temporary fix - using `ssh-keyscan github.com`                          |

### 2021-12-06

| time  | description                                                                                                                                   |
|-------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| 15:53 | The team releases an update for the okctl binary ([0.0.75](https://github.com/oslokommune/okctl/releases/tag/v0.0.75)) with the newest ArgoCD |

## Impact

<!--
Describe the consequences this had for the organization.

Example:
- Approximately 30.000 users downloaded and experienced an error using the app.
- Users already logged in experienced some minor delay in the app
-->
Most of our users were not able to use ArgoCD to deploy updated versions of their apps.

## Cause(s)
<!--
Describe relevant factors that played a part in causing the event.

Example:
The error happened due to a combination of:
- A spike of new users
- A bug in the backend code for user creation
-->
- Github rotating their keys
- The team not picking up the change (it was announced)
- The dependency ArgoCD had fallen too far behind version-wise

## Solution
<!--
Describe how the problem was fixed.

Example:
- Manually scaled up number of backend instances letting regular traffic stay unaffected
- Created a fix for the database schema
-->
- Referred users to the user's fix <sub>4</sub>
- Updated ArgoCD for new clusters and with it a suboptimal fix
- Created upgrade for ArgoCD

## What went well
<!--
Describe what went well trying to handle the event.

Example:
- Alerting mechanisms worked brilliantly when errors started comming in
- Deploying the database update was fast
-->
- Quickly figured out the issue
- Quickly identified a fix
- Quickly pushed a fix thanks to one of the users<sub>4</sub>, in other words; our efforts in concentrating knowledge
    and code are working.

## What could we have done better
<!--
Describe what went wrong trying to handle the event.

Example:
- Had to manually roll back the database while testing fixes
-->
- We could have picked up the breaking change when it was announced, and prevented any problems
- The ArgoCD values file had to be converted to YAML from struct spaghetti before working on the fix
- Due to the ArgoCD version lagging behind, we had to migrate and test more for the upgrade

## Where were we lucky
<!--
Describe briefly the situations related to this event where we got lucky.

Example:
- Updating the database schema fixed the problem
- Manual horizontal scale up ensured a small amount of users was affected
-->
- Lucky that one of the other users found and shared a temporary fix

## Action items
<!--
Describe and reference specific action items

Example:
[KM123](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add backoff for login in the iOS app
[KM221](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add tests in the backend for creating a user and authenticating
-->
- [KM484](https://trello.com/c/KFY52VE3) - Add Github blog to #kjøremiljø-feeds
- [KM485](https://trello.com/c/uZawm9Sz) - Discuss strategy for keeping components up-to-date
 
## Supporting information and documentation
<!--
Add relevant additional information or documentation.

Example: images of metrics that show when something went wrong, relevant pieces of logs, etc
-->

1) [#kjøremiljø-support request](https://oslokommune.slack.com/archives/CV9EGL9UG/p1637162801019900)
2) [Failed nightly](https://oslokommune.slack.com/archives/C018EFCSC4W/p1637116522239900)
3) [Github blog article describing the situation](https://github.blog/2021-09-01-improving-git-protocol-security-github/)
4) [external user's fix](https://oslokommune.slack.com/archives/CV9EGL9UG/p1637234863021300?thread_ts=1637162801.019900&cid=CV9EGL9UG)
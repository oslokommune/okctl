## Summary
<!--
Describe the reason and the consequences of the event as short and concise as possible.

Example:
The backend of Oslo Nøkkelen were unavailable for 31 minutes due to an increase in traffic that happened after a
proaktiv melding.
-->
Team Barnehagepris had to postpone batch processing operations for seven hours due to suspended logging.

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

| Timeline 20.04.22 | Event                                                                                                                                     |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| 11:28             | Team Barnehagepris was made aware of the problem                                                                                          |
| 11:38             | Team Barnehagepris [reached out](https://oslokommune.slack.com/archives/CV9EGL9UG/p1650447510612329) to Team Kjøremiljø about the problem |
| 11:56             | Team Kjøremiljø initiated support by suggesting different possible solutions                                                              |
| 19:33             | PROBLEM FIXED - processing resumes.                                                                                                       |

## Impact
<!--
Describe the consequences this had for the organization.

Example:
- Approximately 30.000 users downloaded and experienced an error using the app.
- Users already logged in experienced some minor delay in the app
-->
- During the annual creation of re-applications for citizens with reduced payment the processing had to be halte for
  hours.
- Approximately two hours of log data got lost

## Cause(s)
<!--
Describe relevant factors that played a part in causing the event.

Example:
The error happened due to a combination of:
- A spike of new users
- A bug in the backend code for user creation
-->
- Log retention policies for Loki had not been configured

## Solution
<!--
Describe how the problem was fixed.

Example:
- Manually scaled up number of backend instances letting regular traffic stay unaffected
- Created a fix for the database schema
-->
The following steps were done as a temporary solution

- Created a bigger volume for logs (from 10GB to 50GB) and migrated old logs to the new volume
- Configured log retention to suit needs

As a more permanent fix, we'll use S3 and DynamoDB as a persistence layer for Loki as described in [KM676](https://trello.com/c/QIgNLSh1)

## What went well
<!--
Describe what went well trying to handle the event.

Example:
- Alerting mechanisms worked brilliantly when errors started comming in
- Deploying the database update was fast
-->

- Great cooperation between Team Barnehagepris and Team Kjøremiljø
- Impressive will to learn, execute and a commendable attitude in Team Barnehagepris

## What went wrong
<!--
Describe what went wrong trying to handle the event.

Example:
- Had to manually roll back the database while testing fixes
- We missed the notification regarding Keycloak upgrade
-->
- Development environment did not match production environment causing delays during the research and test phase. The
    log persistence patch had only been applied in production.
 
## Mitigations
<!--
Describe potential steps to prevent this event from happening again, be it technical, processes, etc.

Example:
- Request that proactive messages will be sendt in smaller batches, ideally during regular work hours, for load to be
    better distributed.
- Create tests that cover creation and authentication of users
-->

One or more of the following mitigations will work:

- Configure log retention for Loki
- Use S3 and DynamoDB for logs and indexes
- Use managed Loki

## Action items
<!--
Describe and reference specific action items

Example:
[KM123](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add backoff for login in the iOS app
[KM221](https://trello.com/c/nJpSSGCs/455-post-mortem-the-github-cert-issue) - Add tests in the backend for creating a user and authenticating
-->
[KM676](https://trello.com/c/QIgNLSh1)

## Supporting information and documentation
<!--
Add relevant additional information or documentation.

Example: images of metrics that show when something went wrong, relevant pieces of logs, etc
-->
[Thread in #kjøremiljø-support](https://oslokommune.slack.com/archives/CV9EGL9UG/p1650447510612329)
[Start of event in #team-barnehagepris](https://oslokommune.slack.com/archives/C02CLDZDD0A/p1650446906306059)

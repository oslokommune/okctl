## Summary

A team's logging was discovered not work 2022-09-21, the duration is unknown. This happened because Loki,
which is responsible for forwarding logs to AWS S3, did not have correct IAM permissions.

## Timeline

2022-09-21

* 14:08 The affected team notifies Kjøremiljø that they cannot find any logs in Grafana.
* 14:14 Kjøremiljø responds, and starts digging into the issue.
* 15:45 Kjøremiljø finds the culprit: Loki's pod is missing a DeleteTable permission in IAM. This is
the same as happened to another team, see
  https://oslokommune.slack.com/archives/CV9EGL9UG/p1662034566815389?thread_ts=1662033302.833479&cid=CV9EGL9UG
* 16:06 Kjøremiljø runs upgrade [0.0.104-fix-loki-delete-table-issue](https://github.com/oslokommune/okctl-upgrade/releases/tag/0.0.104%2Bfix-loki-delete-table-issue) which fixes the issue in the production cluster.
  
  Side note: For some reason, the pod `loki-0` needs to be deleted with kubect's `--force` flag, no idea
  why. 
* 16:30 When attempting to find the same error in their own reference clusters, Kjøremiljø notices some
  weird behavior. There is another permission error there for CreateTable: `User: arn:aws:sts::xxxxxxxxxxxxx:assumed-role/okctl-okctl-reference-loki/1663833091140436960 is not authorized to perform: dynamodb:CreateTable on resource: arn:aws:dynamodb:eu-west-1:xxxxxxxxxxxxx:table/index_19139 because no identity-based policy allows the dynamodb:CreateTable action
  status code: 400`
* 16:35 End of work day, work continues next day.

2022-09-22
* 09:30 Kjøremiljø runs the 0.0.104 upgrade in the team's development cluster.
* 10:04 Kjøremiljø verifies in Grafana that logging in the team's development and production clusters now
  fixed.
* 14:01 Kjøremiljø concludes that the CreateTable permission didn't happen in the team's cluster, only our
  own clusters.

## Impact

* Logging stopped working for the team
  * in the production Kubernetes cluster between unknown time and 2022-09-21 16:06
  * in the development Kubernetes cluster between unknown time and 2022-09-22 10:04

## Cause(s)

Loki's pod is missing an IAM permission ("DeleteTable") to be able to delete DynamoDB tables, which is used
by Loki for storing log indexes.

## Solution

Okctl upgrade
[0.0.104-fix-loki-delete-table-issue](https://github.com/oslokommune/okctl-upgrade/releases/tag/0.0.104%2Bfix-loki-delete-table-issue)
solved the issue.

Side note: It was not possible to run `okctl upgrade` to for applying this upgrade, because upgrade
[0.0.102-eks-1-21](https://github.com/oslokommune/okctl-upgrade/releases/tag/0.0.102%2Beks-1-21) fails if
EKS 1.21 is not installed. At the time, the team was using EKS 1.20. Kjøremiljø and the team had attempted
to upgrade to EKS 1.21 the week before, but it failed, and Kjøremiljø haven't had time to fix that
[issue](https://github.com/oslokommune/okctl/issues/1038) yet.  The solution was to download and run
the 0.0.104 upgrade directly.

## What went well

* The team forwards all "ERROR" log lines to Slack, so even if Loki failed, they still got the most
  important logs.

## What went wrong

* The 0.0.104 patch had been ready for 15 days (since 2022-09-06), so the problem could have been avoided
if this upgrade had been run. 
  
  It was unfortunately not possible for the team to apply this upgrade using `okctl upgrade`, as
described above. Kjøremiljø would have had to been involved to be able to apply this upgrade, as teams
don't know how to apply upgrades directly.
  
* Kjøremiljø uses the Slack channel #okctl-viktig to notify teams about important stuff that teams should
  know about. In this case, Kjøremiljø had not notified teams to run this critical upgrade. Although teams
  are expected to stay updated on Okctl, the Okctl change log only show changes in Okctl, and not the
  upgrades themselves.
  
## Other information

The team didn't notice anything was wrong for many days because they don't usually read logs directly. They
instead rely on alerts notifying them on Slack if something goes wrong.
Detalis: https://oslokommune.slack.com/archives/CV9EGL9UG/p1656683913778469?thread_ts=1656678663.318469&cid=CV9EGL9UG

## Where were we lucky

The team forwards important logs to Slack.

## Mitigations

* Kjøremiljø have previously decided to put important messages about okctl and upgrades on #okctl-viktig,
  so we need to ensure that this is not forgotten.
* On a higher level, focus on getting the new golden path ready. Then this issue will become irrelevant.

## Action items

* https://github.com/oslokommune/okctl/issues/1038 - Research and fix the issue
* https://github.com/oslokommune/kjoremiljo-admin-tasks/issues/56 - Ensure teams are notified about important upgrades

## Supporting information and documentation

* Communication with the team: https://oslokommune.slack.com/archives/CV9EGL9UG/p1663762107259069
* Thread with technical details: https://oslokommune.slack.com/archives/C018EFCSC4W/p1663768362162369  
* Communication with other teams about the same issue:
  * https://oslokommune.slack.com/archives/CV9EGL9UG/p1661499598067769
  * https://oslokommune.slack.com/archives/CV9EGL9UG/p1662033302833479

## Summary

Most of the pods running in the development environment of team Booking was down for around 1h 30min. 

The cause is the same as an earlier
[incident](https://docs.google.com/document/d/1_NRW32xvu7YgEakl1JnLf2nDRbOPe8exrzgnwPAAj30/edit#heading=h.mwz2ai1qrpl4).

## Timeline

* 2022-01-07 12:39 - 15:30
  * 12:39 A user from team Booking notifies us on Slack that their development environment is down. We immediately start
investigating.

  * We identify the cause to be missing resource quotas. This has been fixed in with a
    [resource quota patch](https://github.com/oslokommune/okctl-upgrade/blob/main/gists/resource-quota-patch/resource-quota-patch.sh),
    but it has not been run in this cluster, so we run it.

  * We identify that the patch for
    [persistent logs for Loki](https://github.com/oslokommune/okctl-upgrade/tree/main/gists/loki-persistent-storage) has not
    been run. This complicates putting resource quotas on the Loki pods, because logs will disappear as soon as the quotas
    are applied, because the pod will restart.

  * In AWS console, we set Auto Scaling Groups min and max capacity from 1 to 0, and then back again. At the time, we had an
  idea it would fix issues with the node. In hindsight, this was probably not necessary.

  * We run the resource quota patch to set resource quotas on most pods except Loki.

  * We create a backup of Loki logs in development, and test that it's possible to restore. We then apply the persistent
    storage patch for Loki, and restore the logs. We then run the resource quota patch again, this time enabling it for Loki also.

  * We then run the resource quota patch, without Loki, in production. As it's getting late, we delay the rest of the fix to
    after the weekend.

* 2022-01-10
  * We repeat the process we did for the development environment, i.e. taking backup of logs, running log persistence patch,
    restoring backups, and running resource quota patch for Loki.

## Impact

* Development environment went down for approximately 1 hour and 30 minutes.

## Cause(s)

* When the load on a node is too high, Kubernetes is supposed to spin up a new node to handle the increased load. This does
  however not work when pods don't have resource quotas set. We don't know exactly why the existing node stops working. However,
  we know that our resource quota patch fixes the issue, so the cause is not running this patch.

## Solution

* We ran the following patches both in development and production environments
  * https://github.com/oslokommune/okctl-upgrade/blob/main/gists/loki-persistent-storage/loki-values.yaml
  * https://github.com/oslokommune/okctl-upgrade/blob/main/gists/resource-quota-patch/resource-quota-patch.sh

## What went well

* We responded fast.
* We were able to quickly identify the issue.
* We had patches to fix the problem readily available.
* Communication between us and the Booking team went well (using Whereby).

## What went wrong

* None of our critical patches had been applied to the any of the environments.
  * Our notifications about these patches had got lost between all other messages on our Slack support channel.  
  * For this team, the person who had most know-how about the Okctl infrastructure has recently changed team. Also see comment in
  [Supporting information and documentation](#supporting-information-and-documentation).

## Where were we lucky

* The affected environment was development and not production. 
* The environment went down during work hours.

## Mitigations
<!--
Describe potential steps to prevent this event from happening again, be it technical, processes, etc.

Example:
- Request that proactive messages will be sendt in smaller batches, ideally during regular work hours, for load to be
    better distributed.
- Create tests that cover creation and authentication of users
-->

* We have decided to put all messages about upgrades, patches and other important messages on okctl.io. We'll create an RSS feed
for posts, and forward those to Slack on a separate channel (TODO update with decision). With these steps, people
have multiple options for getting information about their infrastructure.

## Supporting information and documentation

* For new applications made with `okctl apply applications`, resource quotas are set by default. However, in theory, if somebody
spins up lots of pods wihout resource quotas, this could happen again.
* A final thought: This is what can happen when we're establishing a culture of "you own it, you run it", where people run their
own infrastructure. Also, most of our teams have more experience with developement than operations. We're in the start of this
culture change, and it's not unlikely to have a few bumps in the road. As people will get used to this culture, we expect this to
not be a significant problem.


## Summary
A team's logging stopped working between 2022-06-22 and 2022-07-01. This happened because Loki, which is responsible for forwarding logs to AWS S3, was running a too new and unsupported version.

## Timeline
2022-07-01

* 14:31 The affected team notifies Kjøremiljø that they cannot find any logs in Grafana.
* 14:37 Kjøremiljø responds, and starts digging into the issue.
* 14:41 Kjøremiljø confirms that Grafana cannot find any logs in the team's development cluster.
* 15:45 Kjøremiljø finds the culprit: Loki is running version 2.4.1, which is a too new version (it should be 2.1.0).
* 16:06 Kjøremiljø suggest to reinstall Loki, which will delete logs preceding 2022-06-22, which is approved by the team.
* 16:15 Kjøremiljø reinstalls Loki in the development cluster (which sets the Loki version to 2.1.0).
* 16:23 Kjøremiljø confirms that logs are working in the development cluster.
* 16:36 Kjøremiljø starts reinstalling Loki in the production cluster as well.
* 16:55 Kjøremiljø confirms logs are now working in production.

## Impact
* Logging stopped working for the team
  * in the production Kubernetes cluster between 2022-06-22 00:00 and 2022-07-01 16:55
  * in the development Kubernetes cluster between  2022-06-13 00:00 and 2022-07-01 16:36

## Cause(s)
We don't know the exact cause yet. This post mortem will be updated if we find out, see [action items](#Action-items).

However, we do know that logging stopped working the day after we did the following maintenance on the cluster

* Running `okctl upgrade`, which for both the development and production cluster included the upgrades `0.0.95.persist-loki` and `0.0.96.remote-state-versioning`. 

* Upgraded EKS control plane and nodes from 1.19 to 1.20, causing all pods in the `monitoring` namespace (including Loki) to be stopped on old nodes and re-run at new nodes.

None of these should affect the version of Loki. Loki's version is specified in the Statefulset `loki` in the `monitoring` namespace, and we don't have code changing this.

`0.0.95.persist-loki` is touching the Loki configuration, but it should not affect the version of the docker image being used.

Lastly, both `okctl upgrade` and EKS upgrade have been run for other clusters without any issues. We do however see that two teams are running unexpected versions of the Loki chart. See supporting links at the bottom of this document.

## Solution
The solution was to change Loki version to the same version as Okctl `0.0.98` (the Okctl cluster version) installs, which is `2.1.0` .

We did this by uninstalling Loki and Promtail, and reinstalling again, using okctl:

* We edited the cluster manifest to use `loki: false` and `promtail: false`. Promtail was probably not needed, but we did just because it is related, and doesn't hurt to reinstall as well.

```bash
okctl apply cluster -f cluster-manifest.yaml
```

This uninstalled Loki and Promtail. Then we set  `loki: true` and `promtail: true` and ran

```bash
okctl apply cluster -f cluster-manifest.yaml
```

This installed Loki with the correct version, `2.1.0`.

## What went well
Nothing in particular.

## What went wrong
Nothing in particular.

## Other information

The team didn't notice anything was wrong for many days because they don't usually read logs directly. They instead rely on alerts notifying them on Slack if something goes wrong. Detalis: https://oslokommune.slack.com/archives/CV9EGL9UG/p1656683913778469?thread_ts=1656678663.318469&cid=CV9EGL9UG

## Where were we lucky
The team didn't need to preserve old logs contained in Loki's persistent volume in Kubernetes. If this had been the case, we would have had to backup the volume and restore it after reinstalling Loki, which would have added to the recovery time and complexity of recovering.

## Mitigations
* There is currently no system (monitoring and alerting) in place for noticing if logging in an Okctl EKS cluster is working. If Okctl was not going to be discontinued, Kjøremiljø should implement some alerting system to ensure logging is always up.
* We have confirmed that no other Okctl cluster is running Loki 2.4.x or newer. Some clusters are running 2.3.0, but that should not be a problem.

## Action items
* [22Q2-31](https://trello.com/c/K4QykQh0/31-okctio-beskriv-i-common-issues-hvordan-fikse-logging-som-slutter-%C3%A5-funke-pga-loki-24x) - Describe in common issues in our documentation how to fix logs if they stop working due to this issue.
* [22Q2-32](https://trello.com/c/plYjrcfb/32-diskuter-hvordan-loki-i-et-cluster-kunne-f%C3%A5-versjon-241) - Discuss what causes Loki to use version 2.4.1, and why some teams are using unexpected versions of the Loki helm chart.

## Supporting information and documentation
* Communcation with team - https://oslokommune.slack.com/archives/CV9EGL9UG/p1656678663318469
* Loki version in all clusters - https://oslokommune.slack.com/archives/CV9EGL9UG/p1656687494513269
* Issue making Loki 2.4.x not work - https://github.com/grafana/loki/issues/5021




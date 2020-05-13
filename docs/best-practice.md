# Observability

We use [grafana](https://grafana.com/) as an entrypoint to observability of the cluster and applications within it.

## Logs

We use [loki](https://grafana.com/oss/loki/) and [promtail](https://github.com/grafana/loki/tree/master/docs/clients/promtail) for gathering cluster and application logs.

- Logs should not contain any secrets or user identifiable data

## Metrics

We use [prometheus](https://prometheus.io/) for gathering cluster and application metrics. This is achieved by deploying a [prometheus operator](https://github.com/coreos/prometheus-operator) to the cluster, which makes it easier to configure service monitoring.

Prometheus is pull-based, this means it needs to be configured to know where and how often it should fetch metrics from a given application. Our current recommendation, for achieving good security and separation is to follow these guidelines:

- The application metrics will be exported on a port separate, e.g, `localhost:9020/metrics`

## Tracing

TBD

## Ready / Live

These checks should run on the same port as the application
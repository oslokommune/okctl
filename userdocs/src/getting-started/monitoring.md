With `okctl` we create an observability stack in the cluster that provides metrics, traces and logs from the Kubernetes cluster, relevant AWS resources and the applications running in the cluster.

## Observability stack

The observability stack relies on [Grafana](https://grafana.com/oss/grafana/) at its core. Using the [data sources](https://grafana.com/docs/grafana/latest/datasources/) capability of Grafana we can easily integrate with a variety of backends. The ones we currently support are:

- [AWS CloudWatch](https://grafana.com/docs/grafana/latest/datasources/cloudwatch/) for AWS resources, including EKS control plane logs
- [Loki](https://grafana.com/docs/grafana/latest/datasources/loki/) for logs
- [Prometheus](https://grafana.com/docs/grafana/latest/datasources/prometheus/) for metrics
- [Tempo](https://grafana.com/docs/grafana/latest/datasources/tempo/) for traces

These backends provide us with the basic building blocks we required to build a fully functional observability stack.

### Declarative configuration

We love declarative configuration, being able to check _everything_ into git is the best thing to happen since sliced bread. We use declarative configuration to add dashboards, alerts, and scrapers to Grafana and Prometheus. By using `ConfigMap`s with annotations to add dashboards to Grafana we can easily track these resources in our version control system also, similarly for `ServiceMonitor` type for Prometheus.

### Alerting

We will eventually use the [AlertManager](https://prometheus.io/docs/alerting/latest/alertmanager/) for setting up alerts, feel free to do so now, but we haven't started looking at this in-depth yet.

### Prometheus

**NB:** We only support prometheus for metrics currently, as such, you need to ensure that your application has a metrics endpoint that can be scraped by the `ServiceMonitor` you setup.

We integrate [Prometheus](https://prometheus.io) into the cluster by using [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus), where [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) is used to automatically start scraping an application for metrics.

The full list of available [CustomResourceDefinitions](https://github.com/prometheus-operator/prometheus-operator#customresourcedefinitions) provides a good overview of the capabilities provided by this operator. 

For setting up monitoring of your application, we recommend reading [this guide](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md). The most relevant part is the setup of the `ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-app
  labels:
    team: frontend
spec:
  selector:
    matchLabels:
      app: example-app
  endpoints:
  - port: web
```

Once you have setup a `ServiceMonitor` for your application you can login to the `Grafana` dashboard and use the query explorer using the `Prometheus` datasource and start searching for the metrics you have defined.

### Loki

We scrape the logs of all pods in the Kubernetes cluster using [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) and send these on to Loki. For details on how best to setup your logging for use with Loki, we recommend reading the [documentation](https://grafana.com/docs/loki/latest/). Essentially, these logs will be available, and can be queried from the query explorer in Grafana when the Loki datasource is selected.

### Tempo

[Tempo](https://grafana.com/oss/tempo/) receives the traces from your application and makes them available for querying. For more details on how to use Tempo, we recommend reading the [documentation](https://grafana.com/docs/grafana/latest/datasources/loki/).

### Grafana

We have set up [Grafana](https://grafana.com/oss/grafana/) for you, and multiple datasources. Once you have found the metrics, traces or logs you are interested in following you can setup a dashboard for easy viewing. The easiest way of creating a dashboard is through the [UI](https://grafana.com/docs/grafana/latest/getting-started/getting-started/#step-3-create-a-dashboard). Once you are satisfied with the result, you can export it as json and add it to grafana via a declarative config.

We achieve this, because we have enabled a [sidecar for dashboards](https://github.com/grafana/helm-charts/tree/main/charts/grafana#sidecar-for-dashboards). In essence, you define a `ConfigMap` like so:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sample-grafana-dashboard
  labels:
     grafana_dashboard: "1"
data:
  k8s-dashboard.json: |-
  [...]
```

The important part is the `grafana_dashboard` label, also, please remember that the name of the dashboard, e.g.: `k8s-dashboard.json` needs to be unique, if you use the same name everywhere they will overwrite each other. 
 
### Roadmap

- [x] Collecting logs, metrics and traces
- [ ] Setting up alarms and alerts
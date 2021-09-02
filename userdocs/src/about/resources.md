<style>
.components-map-container {
    display: flex;
    justify-content: space-between;

    font-size: 12pt;
}

.components-map-column {
    flex-basis: 32%;

    font-size: .7rem;
}

.centered-title {
    text-align: center;
}

.component-title {
    margin-top: 0.5em !important;
    margin-bottom: 0 !important;
}

.component-icon {
    margin: 0 auto;
    display: block;

    height: 29px !important;
}

ul {
    margin-top: 0 !important;
    margin-left: 0 !important;
}

li {
    margin-left: 15px !important;
}
</style>


# Overview

## AWS resources used

The following is a map of all the resources okctl will set up

<div style="display: flex;">
    <img alt="cloud resources" src="/img/Resources.png">
</div>

## AWS resource dependencies

The following shows how relations between the different resources

<div style="display: flex;">
    <img alt="resource dependencies" src="/img/Dependencies.png">
</div>

## Kubernetes installed compontents

The following shows how relations between the different resources

<div class="components-map-container">
    <div class="components-map-column">
        <h3 class="centered-title">Monitoring</h3>
        <div>
            <object class="component-icon" data="/img/Icon-Prometheus.svg" type="image/svg+xml"></object>
            <h4 class="component-title">Prometheus</h4>
            <span>Time series database</span>
            <ul>
                <li>Enables scaping of metrics from pods</li>
                <li>Retrieves log data from Loki</li>
                <li>Retrieves trace data from Tempo</li>
                <li>Provides data for Grafana</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-grafana.svg" type="image/svg+xml"></object>
            <h4 class="component-title">Grafana</h4>
            <span>Data visualizer</span>
            <ul>
                <li>Reads log, metric and trace data from Prometheus</li>
                <li>Provides powerful tools for data visualization</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-loki.svg" type="image/svg+xml"></object>
            <h4 class="component-title">Loki</h4>
            <span>Log aggregator</span>
            <ul>
                <li>Acts as a data source for Prometheus</li>
                <li>Accepts log data from Promtail</li>
            </ul>
        </div>

        <div>
            <h4 class="component-title">Promtail</h4>
            <span>Log scraper</span>
            <ul>
                <li>Scrapes and labels log data from pods</li>
                <li>Pushes log data to Loki</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-Tempo.svg" type="image/svg+xml"></object>
            <h4 class="component-title">Tempo</h4>
            <span>Tracing backend</span>
            <ul>
                <li>Acts as a datasource for Prometheus</li>
                <li>Supports OpenTelemetry, Jaeger, Zipkin
            </ul>
        </div>
    </div>

    <div class="components-map-column">
        <h3 class="centered-title">Kube System</h3>

        <div>
            <object class="component-icon" data="/img/Icon-AWS-LoadBalancer.svg" type="image/svg+xml"></object>
            <h4 class="component-title">AWS Load Balancer</h4>
            <span>Ingress controller</span>
            <ul>
                <li>Provides load balancers based on Kubernetes Ingress'</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-aws_autoscaler.svg" type="image/svg+xml"></object>
            <h4 class="component-title">AutoScaler</h4>
            <span>Horizontal cluster scaler</span>
            <ul>
                <li>Adds and removes cluster nodes depending on load</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-CSI.svg" type="image/svg+xml"></object>
            <h4 class="component-title">EBS CSI Controller</h4>
            <span>Persistent volume provider</span>
            <ul>
                <li>Creates AWS Elastic Block Store based on Persistent Volume Claims and connects them to pods</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-aws_external_dns.svg" type="image/svg+xml"></object>
            <h4 class="component-title">External DNS</h4>
            <span>Domain administration controller</span>
            <ul>
                <li>Configures Route53 entries based on Kubernetes Ingress'</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-AWS-ExternalSecrets.svg" type="image/svg+xml"></object>
            <h4 class="component-title">External Secrets</h4>
            <span>Secrets controller</span>
            <ul>
                <li>Fetches secrets from AWS Parameter Store and Secrets Manager, then injects them into the cluster</li>
            </ul>
        </div>
    </div>

    <div class="components-map-column">
        <h3 class="centered-title">ArgoCD</h3>

        <div>
            <object class="component-icon" data="/img/Icon-ArgoCD.svg" type="image/svg+xml"></object>
            <h4 class="component-title">ArgoCD</h4>
            <span>Continuous Deployment Provider</span>
            <ul>
                <li>Synchronizes Kubernetes state with a Git repository</li>
                <li>Enables rollback of state based on commit history</li>
                <li>Provides superficial administration of cluster applications</li>
            </ul>
        </div>

        <div>
            <object class="component-icon" data="/img/Icon-Dex.svg" type="image/svg+xml"></object>
            <h4 class="component-title">Dex</h4>
            <span>A pluggable OAuth2 handler</span>
            <ul>
                <li>Acts as a mediator for Cognito</li>
                <li>Handles authentication for ArgoCD</li>
                <li>Handles authentication for Grafana</li>
            </ul>
        </div>
    </div>
</div>

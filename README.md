# typesense-prometheus-exporter

`typesense-prometheus-exporter` is a lightweight Prometheus exporter designed to expose metrics from a Typesense cluster 
for monitoring and alerting purposes. The exporter collects metrics from the Typesense `/metrics.json` and 
`/stats.json` endpoints and presents them in a Prometheus-compatible format.

### **Usage**

#### **Running Locally**

1. Clone the repository:

   ```bash
   git clone https://github.com/your-fork/typesense-prometheus-exporter.git
   cd typesense-prometheus-exporter
   ```

2. Build the exporter:
   ```bash
   make build
   ```

3. Run the binary with the required environment variables:

   ```bash
   LOG_LEVEL=0 TYPESENSE_API_KEY=your-api-key \
   TYPESENSE_HOST=your-host TYPESENSE_PORT=8108 \
   METRICS_PORT=8908 TYPESENSE_PROTOCOL=http \
   TYPESENSE_CLUSTER=your-cluster-name \
   ./cmd/typesense-prometheus-exporter
   ```

#### **Running in Docker**

1. Deploy typesense-prometheus-exporter as a stand-alone stack with docker-compose:

```bash
version: '3.8'

services:
  typesense-prometheus-exporter:
    image: akyriako78/typesense-prometheus-exporter:0.1.7
    container_name: typesense-prometheus-exporter
    environment:
      LOG_LEVEL: "0"
      TYPESENSE_API_KEY: "${TYPESENSE_API_KEY}" # Use an .env file or environment variable for secrets
      TYPESENSE_HOST: "ts.example.com"
      TYPESENSE_PORT: "8108"
      TYPESENSE_PROTOCOL: "http"
      TYPESENSE_CLUSTER: "ts"
    ports:
      - "8908:8908"
```

2. Open http://localhost:8908 in your browser:

![image](https://github.com/user-attachments/assets/c2ccdfe3-1c37-49f0-acda-6b44950c2096)


#### **Running in Kubernetes**

0. If you are not having a Prometheus instance, deploy one with kube-prometheus-stack

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install {RELEASE_NAME} prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

> [!IMPORTANT]
> Note down the `RELEASE_NAME` you've used, we are going to need it later.

1. Deploy the exporter as a `Deployment` in your Kubernetes cluster:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
   name: ts-prometheus-exporter
   namespace: default
spec:
   replicas: 1
   selector:
      matchLabels:
         app: ts-prometheus-exporter
   template:
      metadata:
         labels:
            app: ts-prometheus-exporter
      spec:
         containers:
            - name: typesense-prometheus-exporter
              image: akyriako78/typesense-prometheus-exporter:0.1.7
              env:
                 - name: LOG_LEVEL
                   value: "0"
                 - name: TYPESENSE_API_KEY
                   valueFrom:
                      secretKeyRef:
                         name: ts-api-key
                         key: typesense-api-key
                 - name: TYPESENSE_HOST
                   value: ts-svc.default.svc.cluster.local
                 - name: TYPESENSE_PORT
                   value: "8108"
                 - name: TYPESENSE_PROTOCOL
                   value: "http"
                 - name: TYPESENSE_CLUSTER
                   value: ts
              ports:
                 - containerPort: 8908
```

2. Expose a `Service` for the exporter (Prometheus will scrape this endpoint)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: ts-prometheus-exporter-svc
  namespace: default
  labels:
    app: ts-prometheus-exporter
spec:
  selector:
    app: ts-prometheus-exporter
  ports:
    - name: metrics
      port: 8908
      targetPort: 8908
```

> [!IMPORTANT]
> If you set a custom `METRICS_PORT` remember to adjust `containerPort` in **Deployment** and `port` and `targetPort` in **Service**
> accordingly. Check [Configuration](#configuration) for more info.

3. Provision a `ServiceMonitor`, so this endpoint is picked up automatically from Prometheus as a target:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: ts-servicemonitor
  namespace: default
  labels:
    release: {RELEASE_NAME}
spec:
  selector:
    matchLabels:
      app: ts-prometheus-exporter
  namespaceSelector:
    matchNames:
      - default
  endpoints:
    - port: metrics
      path: /metrics
      interval: 15s
      scheme: http
  jobLabel: app
```

> [!CAUTION]
> 1. Put the `ServiceMonitor` in any namespace you want (either the one hosting your workload or the one hosting Prometheus)
> but make sure `spec.namespaceSelector.matchNames` is pointing to the namespace your `Service` is deployed at.
> 2. Add an extra label in `metadata.Labels`, with key `release` and value the the **release name** you provided when you deployed
> your **kube-prometheus-stack** Helm Chart. Without this label, Prometheus will **not** automatically pick up this endpoint as 
> one of its target. It is not well documented in Prometheus documentation, but it's the _secret sauce_.

![image](https://github.com/user-attachments/assets/084f6fe9-2342-427f-9686-d63461507afa)

![image](https://github.com/user-attachments/assets/04a03c85-5b86-4f37-ada6-9f300a0a811d)

#### Import Grafana Dashboards

Open your Grafana installation and import the dashboards found in **assets/grafana**. There is one for metrics and one for stats.

![image](https://github.com/user-attachments/assets/606e182a-867f-4c62-8668-e6cdc5d2ddb0)

### **Configuration**

The `typesense-prometheus-exporter` is configured via environment variables. Below is a table of the available configuration options:

| **Variable**         | **Type** | **Default** | **Required** | **Description**                                                     |
|----------------------|----------|-------------|--------------|---------------------------------------------------------------------|
| `LOG_LEVEL`          | `int`    | `0`         | No           | (debug) `-4`, (info) `0` , (warn) `4` , (error) `8`                 |
| `TYPESENSE_API_KEY`  | `string` | -           | Yes          | The API key for accessing the Typesense cluster.                    |
| `TYPESENSE_HOST`     | `string` | -           | Yes          | The host address of the Typesense instance.                         |
| `TYPESENSE_PORT`     | `uint`   | `8108`      | No           | The port number of the Typesense API endpoint.                      |
| `METRICS_PORT`       | `uint`   | `8908`      | No           | The port number for serving the Prometheus metrics endpoint.        |
| `TYPESENSE_PROTOCOL` | `string` | `http`      | No           | Protocol used for communication with Typesense (`http` or `https`). |
| `TYPESENSE_CLUSTER`  | `string` | -           | Yes          | The name of the Typesense cluster, used for labeling metrics.       |

### **Metrics**
The exporter gathers various metrics from the Typesense `/metrics.json` endpoint, including:
- **CPU Utilization**: Per-core and overall CPU usage percentages.
- **Memory Usage**: Active, allocated, and retained memory statistics.
- **Disk Usage**: Total and used disk space.
- **Network Activity**: Total bytes sent and received.
- **Typesense-specific Metrics**: Fragmentation ratios, mapped memory, and more.

> [!NOTE]
> - Each **metric** is labeled with `typesense_cluster` as the name of the Typesense cluster you want to fetch metrics from.
> - Each **stat** is labeled with `typesense_cluster` as the name of the Typesense cluster you want to fetch stats from,
> and additionally with `typesense_request` for any metrics reporting back on individual requests.
> - All FQDNs for Prometheus Descriptors collected from **metrics** are prefixed with `typesense_metrics_` 
> - All FQDNs for Prometheus Descriptors collected from **stats** are prefixed with `typesense_stats_`

### **Build and Push Docker Image**

You can build and push the Docker image using the provided `Makefile`.

```bash
# Build the Docker image
make docker-build REGISTRY=your-registry IMAGE_NAME=typesense-prometheus-exporter TAG=latest
```

```bash
# Push the Docker image to the registry
make docker-push REGISTRY=your-registry IMAGE_NAME=typesense-prometheus-exporter TAG=latest
```

Ensure the `REGISTRY`, `IMAGE_NAME`, and `TAG` variables are properly set.

### **License**
This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

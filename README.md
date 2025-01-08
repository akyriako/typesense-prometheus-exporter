# typesense-prometheus-exporter

`typesense-prometheus-exporter` is a lightweight Prometheus exporter designed to expose metrics from a Typesense cluster for monitoring and alerting purposes. The exporter collects metrics from the Typesense `/metrics.json` endpoint and presents them in a Prometheus-compatible format, enriched with Kubernetes-specific labels.

### **Features**
- Fetches and exposes key performance and resource utilization metrics from Typesense clusters.
- Supports Kubernetes environments with labels for `namespace` and `typesense_cluster` for better observability.
- Fully configurable through environment variables.

### **Usage**

#### **Running Locally**

1. Clone the repository:

   ```bash
   git clone https://github.com/akyriako/typesense-prometheus-exporter.git
   cd typesense-prometheus-exporter
   ```

2. Build the exporter:
   ```bash
   make build
   ```

3. Run the binary with the required environment variables:

   ```bash
   LOG_LEVEL=1 TYPESENSE_API_KEY=your-api-key \
   TYPESENSE_HOST=your-host TYPESENSE_PORT=8108 \
   METRICS_PORT=8908 TYPESENSE_PROTOCOL=http \
   POD_NAMESPACE=default TYPESENSE_CLUSTER=your-cluster-name \
   ./cmd/typesense-prometheus-exporter
   ```

#### **Running in Kubernetes**
Deploy the exporter as a pod in your Kubernetes cluster:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: typesense-prometheus-exporter
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: typesense-prometheus-exporter
  template:
    metadata:
      labels:
        app: typesense-prometheus-exporter
    spec:
      containers:
      - name: exporter
        image: your-registry/typesense-prometheus-exporter:latest
        env:
        - name: LOG_LEVEL
          value: "0"
        - name: TYPESENSE_API_KEY
          valueFrom:
            secretKeyRef:
              name: typesense-api-key
              key: api-key
        - name: TYPESENSE_HOST
          value: "your-typesense-host"
        - name: TYPESENSE_PORT
          value: "8108"
        - name: METRICS_PORT
          value: "8908"
        - name: TYPESENSE_PROTOCOL
          value: "http"
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: TYPESENSE_CLUSTER
          value: "your-cluster-name"
        ports:
        - containerPort: 8908
```

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
| `POD_NAMESPACE`      | `string` | `~empty`    | No           | The Kubernetes namespace where the pod is running.                  |
| `TYPESENSE_CLUSTER`  | `string` | -           | Yes          | The name of the Typesense cluster, used for labeling metrics.       |

### **Metrics**
The exporter gathers various metrics from the Typesense `/metrics.json` endpoint, including:
- **CPU Utilization**: Per-core and overall CPU usage percentages.
- **Memory Usage**: Active, allocated, and retained memory statistics.
- **Disk Usage**: Total and used disk space.
- **Network Activity**: Total bytes sent and received.
- **Typesense-specific Metrics**: Fragmentation ratios, mapped memory, and more.

Each metric is labeled with:
- `namespace`: The Kubernetes namespace where the exporter is running.
- `typesense_cluster`: The name of the Typesense cluster.

### **Build and Push Docker Image**

You can build and push the Docker image using the provided `Makefile`.

```bash
# Build the Docker image
make docker-build REGISTRY=myregistry.io IMAGE_NAME=typesense-prometheus-exporter TAG=latest
```

```bash
# Push the Docker image to the registry
make docker-push REGISTRY=myregistry.io IMAGE_NAME=typesense-prometheus-exporter TAG=latest
```

Ensure the `REGISTRY`, `IMAGE_NAME`, and `TAG` variables are properly set.

### **License**
This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
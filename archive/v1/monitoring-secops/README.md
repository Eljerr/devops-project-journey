# 🔐 Security & Monitoring Stack (SecOps)

This directory contains the central observability and alerting stack deployed via **Docker Compose** on a dedicated Proxmox LXC container (`vmid: 203`, provisioned by Terraform in `core-infrastructure/proxmox-k3s/`).

---

## 🏛️ Architecture Overview

```text
[ K3s Worker / Master Nodes ]              [ Dedicated Monitoring LXC ]
   ├── Node Exporter (DaemonSet) ──────┐          ├── Prometheus (Port 9090)
   └── kube-state-metrics ─────────────┼─────────►├── Alertmanager (Port 9093)
                                       │          └── Grafana (Port 3000)
[ Gateway STB ]                        │
   └── Node Exporter / cAdvisor ───────┘
```

> [!NOTE]
> **Difference with `platform-tools/monitoring/`**:
> - `monitoring-secops/` (this folder): The central server stack (Prometheus, Grafana, Alertmanager) that runs in a standalone LXC container.
> - `platform-tools/monitoring/`: Kubernetes-native DaemonSets and kube-state-metrics manifests applied **inside** the K3s cluster so this Prometheus instance can scrape them.

---

## 📦 Services Included

| Service | Port | Description |
|---|---|---|
| **Prometheus** | `9090` | Time-series database scraping metrics from nodes and containers |
| **Grafana** | `3000` | Dashboards and visualization platform |
| **Alertmanager** | `9093` | Alert routing, grouping, and notification dispatch |
| **Node Exporter** | `9100` | Host hardware & OS metrics collector |
| **cAdvisor** | `8080` | Container resource usage and performance metrics |

---

## 📁 Directory Structure

```text
monitoring-secops/
├── 📄 docker-compose.yml         # Main Docker Compose service definitions
├── 📄 alert.rules.yml            # Prometheus alerting rules
├── 📄 alertmanager.yml           # Alertmanager routing and receiver configuration
├── 📁 grafana/
│   ├── 📁 dashboards/            # Pre-provisioned Grafana dashboards
│   └── 📁 provisioning/          # Datasource & dashboard auto-provisioning configs
└── 📁 prometheus/
    └── 📄 prometheus.yml         # Prometheus scrape configs & targets
```

---

## 🚀 Running the Stack

To start the monitoring stack on the host/LXC container:

```bash
# Copy local configuration overrides if necessary:
cp prometheus/prometheus.yml prometheus/prometheus.local.yml
cp alertmanager.yml alertmanager.local.yml

# Start services in detached mode
docker compose up -d

# Verify container statuses
docker compose ps

# View live logs
docker compose logs -f
```

---

## 🛡️ Security Notes

- UFW rules on the LXC should restrict incoming traffic on port 9090 and 9093 to internal network IPs only.
- Grafana default credentials must be updated upon initial startup.

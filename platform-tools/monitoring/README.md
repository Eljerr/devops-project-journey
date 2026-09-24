# 📊 Cluster-Level Monitoring (Kubernetes Manifests)

This directory contains Kubernetes manifests deployed into the K3s cluster to expose internal cluster metrics and host node telemetry to Prometheus.

---

## 🏛️ Context & Difference

- **`platform-tools/monitoring/` (This Directory)**:
  Contains Kubernetes manifests applied inside the cluster (via ArgoCD `argo-monitoring.yml`).
  Runs Node Exporter as a `DaemonSet` on every node and deploys `kube-state-metrics` to listen to the Kubernetes API server.
- **`monitoring-secops/` (Repository Root)**:
  Contains the central monitoring server stack (Prometheus, Grafana, Alertmanager) deployed via Docker Compose on a dedicated Proxmox LXC container.

---

## 📁 Manifests

| File | Type | Purpose |
|---|---|---|
| **`node-exporter-daemonset.yaml`** | `DaemonSet` | Runs an instance of Node Exporter on every K8s node (Master & Workers) using host networking to expose CPU, memory, disk, and network stats on port `9100`. |
| **`kube-state-metrics.yaml`** | `Service` & config | Exposes metrics about the health of Kubernetes objects (deployments, pods, nodes) to Prometheus. |
| **`kustomization.yaml`** | Kustomize file | Bundles the official kube-state-metrics resources with our local daemonset. |

---

## 🚀 Manual Deployment (without ArgoCD)

If you need to deploy or test these manifests manually using `kubectl`:

```bash
kubectl apply -k platform-tools/monitoring/
```

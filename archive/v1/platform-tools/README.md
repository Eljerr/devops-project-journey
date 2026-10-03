# Platform Tools

This directory contains configuration for platform-level tools I set up to manage and deliver applications to the K3s cluster — primarily ArgoCD for GitOps, and Node Exporter for cluster-wide metrics collection.

## ArgoCD

ArgoCD is a declarative, GitOps continuous delivery tool for Kubernetes. It is used to automatically synchronize and deploy our applications directly from this Git repository into our Kubernetes (K3s) cluster.

### Installation

To install ArgoCD into your Kubernetes cluster, you need to create a dedicated namespace and apply the official installation manifest. 

Run the following commands in your terminal:

```bash
# Create a dedicated namespace for argocd
kubectl create namespace argocd

# Install ArgoCD using the official manifest from their stable repository
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

**Command Explanation:**
- `kubectl apply`: Instructs Kubernetes to apply the given configuration.
- `-n argocd`: Specifies that the installation will take place inside the `argocd` namespace.
- `-f https://raw.../install.yaml`: Fetches the installation configuration (manifest) file directly from the official ArgoCD GitHub URL (stable version). This manifest contains all the components (Deployments, Services, Roles, etc.) needed for ArgoCD to run.

### Directory Structure

#### `argocd/`
This folder contains ArgoCD *Application* manifests that dictate how our applications are deployed.

- **`argo-project01.yml`**: This is an `Application` manifest. It tells ArgoCD to monitor the `projects/01-hello-nginx` folder in this GitHub repository (`https://github.com/Eljerr/devops-project-journey.git`). If there are any changes, ArgoCD will automatically deploy (sync) those changes into the Kubernetes cluster (in the `default` namespace). This file also has `automated prune` and `selfHeal` features enabled.

#### App of Apps (Bootstrap)

- **`bootstrap-argocd.yml`**: This is the root "App of Apps" configuration. It continuously monitors the `platform-tools/argocd` folder in this repository. Any new application manifest placed in that folder will be automatically detected and deployed to the Kubernetes cluster by ArgoCD. This enables the installation of new applications to the Kubernetes cluster simply by pushing a YAML file to GitHub—eliminating the need to SSH into the master node or manually run `kubectl apply`.

### GitOps Workflows: Activating & Deactivating Applications

While pure GitOps typically dictates removing files from the repository to uninstall applications, managing application lifecycles can be more flexibly handled directly via the ArgoCD Web UI:

- **Selective Deployment via UI:** Instead of constantly editing or moving files in the repository to turn applications on or off, you can disable the `automated sync` and `prune` features for specific applications.
- **GUI Resource Management:** By relying on the ArgoCD UI, the declarative configurations remain safely stored in Git, but you gain the manual control to `Sync` or `Delete` the application resources on demand. This provides a highly convenient way to pause specific workloads (such as `01-hello-nginx`, `02-nodejs-api`, or `monitoring-secops`) to free up server RAM, and easily reactivate them with a single click whenever necessary, without needing to create new git commits.

---

## Monitoring

### `monitoring/`
Contains a `node-exporter-daemonset.yaml` that deploys Node Exporter as a DaemonSet across **all K8s cluster nodes**. This ensures every node exposes host-level metrics (CPU, memory, disk, network) to be scraped by the Prometheus instance running in the `monitoring-secops` LXC container.

Note: The monitoring stack itself (Prometheus + Grafana) lives in `monitoring-secops/` at the repo root and runs via Docker Compose on a dedicated LXC container — separate from the K8s cluster.

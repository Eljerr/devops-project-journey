# 🐙 ArgoCD Application Manifests

This directory contains declarative ArgoCD `Application` Custom Resource Definitions (CRDs).

These manifests are automatically discovered and reconciled by the root "App of Apps" pattern defined in `platform-tools/bootstrap-argocd.yml`.

---

## 📋 Applications List

| Manifest | Application Name | Target Path in Repo | Target K8s Namespace | Description |
|---|---|---|---|---|
| **`argo-monitoring.yml`** | `monitoring-system` | `platform-tools/monitoring` | `kube-system` | Cluster-level Node Exporter & kube-state-metrics |
| **`argo-project01.yml`** | `hello-nginx` | `projects/01-hello-nginx` | `default` | Basic Nginx static web workload |
| **`argo-project02.yml`** | `project-02-nodejs-api` | `projects/02-nodejs-api` | `devops-prod` | Multi-tier Node.js API + MySQL database workload |

---

## 🔄 App of Apps Synchronization Flow

```text
[ GitHub Repository ]
  └── platform-tools/argocd/*.yml
              │
              ▼ (Scanned continuously)
[ root-app-of-apps ] (Defined in bootstrap-argocd.yml)
      ├──► Creates Application: monitoring-system
      ├──► Creates Application: hello-nginx
      └──► Creates Application: project-02-nodejs-api
```

---

## 🛠️ Adding a New Application

To add a new workload to the K3s cluster:
1. Create your Kubernetes manifests under `projects/your-new-app/`.
2. Add a new `argo-your-new-app.yml` file in this directory following this pattern:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: your-new-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: "https://github.com/Eljerr/devops-project-journey.git"
    targetRevision: HEAD
    path: projects/your-new-app
  destination:
    server: "https://kubernetes.default.svc"
    namespace: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

3. Git commit and push to `main`. ArgoCD will detect the new file and deploy it automatically.

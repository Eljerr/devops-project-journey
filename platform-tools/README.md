# Platform Tools

This directory contains configuration files and manifests for various platform-level tools used to manage, monitor, and deploy our infrastructure and applications.

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

- **`argo-project01.yml`**: This is an `Application` manifest. It tells ArgoCD to monitor the `projects/01-hello-nginx` folder in this GitHub repository (`https://github.com/Eljerr/devops-project-portofolio.git`). If there are any changes, ArgoCD will automatically deploy (sync) those changes into the Kubernetes cluster (in the `default` namespace). This file also has `automated prune` and `selfHeal` features enabled.

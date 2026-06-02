# 🚀 DevOps Project Journey

![CI DevSecOps](https://github.com/Eljerr/devops-project-journey/actions/workflows/ci-devsecops.yml/badge.svg)
![Last Commit](https://img.shields.io/github/last-commit/Eljerr/devops-project-journey)

This is a **personal learning journal** — a record of everything I've explored and built while studying DevOps from scratch. It's not a step-by-step tutorial, but rather a collection of real, working configurations paired with honest notes on how and why things were done a certain way.

If you're on a similar learning path, hopefully the configs, pipeline setups, and the reasoning behind each decision here can serve as a useful reference or starting point for your own experiments.

---

## 🖥️ Infrastructure Overview

This entire project runs on **personal/repurposed hardware** — no cloud provider needed.

| Role | Device | Specs |
|---|---|---|
| **Proxmox Server** | Used Laptop | 2 Cores, 8 GB RAM — hosts K3s cluster inside LXC containers |
| **Gateway** | STB Amlogic S905X4 | 2 GB RAM, Ubuntu OS — reverse proxy & Cloudflare Tunnel entry point |

---

## 🏗️ Repository Structure

```text
devops-project-journey/
├── 📁 core-infrastructure/    # 🏗️ Hardware, Proxmox, & Gateway configurations
│   ├── 📁 gateway-stb/        # Nginx + Cloudflared (Docker Compose)
│   └── 📁 proxmox-k3s/        # Terraform IaC — LXC provisioning for K3s nodes
├── 📁 monitoring-secops/      # 🔐 Security & Monitoring stack (Docker Compose)
├── 📁 platform-tools/         # ⚙️ GitOps, ArgoCD, & Platform-level manifests
│   ├── 📁 argocd/             # ArgoCD Application manifests
│   └── 📁 monitoring/         # Node Exporter DaemonSet for K8s cluster metrics
├── 📁 projects/               # 🌟 Application Workloads & Demonstrations
│   ├── 📁 01-hello-nginx/     # Project 1: K8s Basics & Nginx App
│   └── 📁 02-nodejs-api/      # Project 2: Backend API + MySQL Database
├── 📁 src-app/                # 📦 Node.js API source code & Dockerfile
└── 📄 Jenkinsfile             # Jenkins pipeline template (placeholder)
```

---

## 🛠️ Tech Stack & Tools

These are the primary tools and technologies demonstrated in this repository:

- **Infrastructure & Provisioning:** Proxmox VE, Terraform (telmate/proxmox provider), LXC Containers
- **Containerization:** Docker, Docker Compose, Docker Hub
- **Container Orchestration & GitOps:** Kubernetes (K3s), ArgoCD
- **Applications & Databases:** Node.js (Express), Nginx, MySQL
- **CI/CD & Security:** GitHub Actions, Jenkins (template), CodeQL, Trivy (Aqua Security)
- **Networking:** Cloudflare Tunnels, Kubernetes NodePort Services
- **Monitoring & Observability:** Prometheus, Grafana, Node Exporter, cAdvisor

---

## 🌟 Repository Sections

### 1. [Core Infrastructure](./core-infrastructure/)
Contains the foundational Infrastructure-as-Code (IaC) and configuration.
- **Proxmox K3s:** Terraform scripts to provision K3s nodes (1 Master + 2 Workers) as unprivileged LXC containers on a Proxmox host. The master node is bootstrapped automatically with K3s via SSH provisioner.
- **Gateway STB:** Docker Compose setup running Nginx (reverse proxy), Cloudflared (secure tunnel), Node Exporter, and cAdvisor on an Amlogic S905X4 set-top box. Enables exposing local services to the internet without port forwarding or a static public IP.

### 2. [Platform Tools](./platform-tools/)
Contains configuration for platform-level operational tools.
- **ArgoCD:** Demonstrates GitOps continuous delivery using an "App of Apps" pattern. `bootstrap-argocd.yml` monitors the `platform-tools/argocd/` folder, and any new Application manifest pushed there will be automatically deployed to the K3s cluster.
- **Monitoring (K8s):** A `node-exporter-daemonset.yaml` deploys Node Exporter as a DaemonSet across all K8s cluster nodes, exposing host-level metrics to be scraped by Prometheus.

### 3. [Security & Monitoring Operations](./monitoring-secops/)
Contains the setup for observability and basic security. The monitoring LXC container itself is provisioned via the Terraform script in `core-infrastructure/proxmox-k3s/`. Inside, it runs a full monitoring stack via Docker Compose:
- **Prometheus** — metrics collection & storage
- **Grafana** — visualization dashboards
- **Node Exporter** — host system metrics
- **cAdvisor** — container-level resource metrics

> [!NOTE]
> **Security (UFW):** UFW is recognized as a best practice for hardening security. However, to streamline the initial setup phase and avoid the repetitive process of configuring it individually on every server, strict UFW configuration is deferred for now.

### 4. [Application Projects](./projects/)
Documents the application workloads I deployed to the Kubernetes cluster via ArgoCD.

- **[Hello Nginx (01-hello-nginx)](./projects/01-hello-nginx/):** My first hands-on with Kubernetes basics (Pods, Deployments, Services) — deploying a simple Nginx web server.
- **[Node.js API + Database (02-nodejs-api)](./projects/02-nodejs-api/):** A multi-tier setup connecting a Node.js backend (`mybinichizuru/project-02-api`) to a MySQL database using ConfigMaps and Secrets. The image tag in `backend.yaml` is automatically updated by the CI/CD pipeline.

### 5. [Application Source Code](./src-app/)
Contains the source code for the Node.js API application that is built and deployed through the CI/CD pipeline.

- **`server.js`:** A simple Express.js API (runs on port 80) that exposes a health-check endpoint and reads the `DB_HOST` environment variable injected by Kubernetes.
- **`Dockerfile`:** Containerizes the app using `node:20-alpine` as the base image. Includes `apk update && apk upgrade` to minimize OS-level vulnerabilities before the image is scanned by Trivy.
- **`.trivyignore`:** Lists specific CVEs that are acknowledged but cannot be fixed immediately (e.g., unfixed upstream patches or transitive NPM dependency vulnerabilities), allowing the CI pipeline to proceed.

### 6. CI/CD & DevSecOps (GitHub Actions)
Continuous Integration and DevSecOps practices are integrated using GitHub Actions (`.github/workflows/ci-devsecops.yml`). The pipeline triggers on any push to `main` that modifies the `src-app/` directory.

The full pipeline flow is:

```
Code Push → CodeQL Analysis → Docker Build → Trivy Scan → Push to Docker Hub → GitOps Handoff
```

- **CodeQL Analysis:** Performs static code analysis on the JavaScript source code to detect security vulnerabilities before the image is even built.
- **Docker Build:** Builds the image tagged as `mybinichizuru/project-02-api:v1.0.<run_number>`.
- **Trivy Vulnerability Scan:** Scans the built image for `CRITICAL` and `HIGH` severity vulnerabilities in both OS packages and libraries.
- **Push to Docker Hub:** Pushes the verified image to Docker Hub.
- **GitOps Handoff:** Automatically updates the image tag in `projects/02-nodejs-api/backend.yaml` and commits it back to `main`. This commit is what triggers ArgoCD to detect a drift and automatically sync the new image version to the K3s cluster — completing the full GitOps loop.

> [!NOTE]
> **Trivy Vulnerability Exceptions (`.trivyignore`):** Some vulnerabilities flagged by Trivy could not be resolved immediately due to:
> 1. **Lag in Upstream Patches:** The OS package repository may not yet have released a fix for a newly published CVE even when pulling the latest packages.
> 2. **Transitive NPM Dependencies:** Vulnerabilities in `node_modules` (e.g., `cross-spawn`, `glob`, `tar`) cannot be patched via `apk` — they are locked in `package-lock.json`.
> 3. **Risk Acceptance & Exception Management:** In real-world scenarios, explicitly ignoring unfixable CVEs (`.trivyignore`) is standard practice while continuing to monitor for future patches.

### 7. [Jenkinsfile](./Jenkinsfile)
A Jenkins Declarative Pipeline **template** at the repository root. It defines the standard pipeline stages (`Checkout → Build → Test → Deploy`) aligned with the CI/CD workflow of this project. Currently, the stages contain placeholder `echo` commands and are not yet actively implemented — the primary CI/CD automation is handled by **GitHub Actions**. This file serves as documentation for an alternative Jenkins-based pipeline that can be extended in the future.

---

## 🗺️ Navigating This Repository

1. **Infrastructure:** Review `core-infrastructure/` to see how the base environment is provisioned using Terraform.
2. **Platform:** Check `platform-tools/` to understand how GitOps is configured with ArgoCD (App of Apps pattern).
3. **Monitoring:** See `monitoring-secops/` for the observability stack (LXC provisioned via Terraform, stack deployed via Docker Compose).
4. **Projects:** Explore the `projects/` folder sequentially to see how Kubernetes manifests are structured. With ArgoCD set up, these projects sync automatically. Alternatively, apply them manually:
   ```bash
   kubectl apply -f projects/01-hello-nginx/app-nginx.yaml
   kubectl apply -f projects/02-nodejs-api/
   ```
5. **Application Source:** See `src-app/` for the Node.js API source code, Dockerfile, and `.trivyignore`.
6. **CI/CD Pipelines:** Inspect `.github/workflows/ci-devsecops.yml` to observe how automated Docker builds, CodeQL analysis, Trivy scanning, and GitOps handoff are implemented.

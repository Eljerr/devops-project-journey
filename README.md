# 🚀 DevOps Project Journey

Welcome to the DevOps Project Journey! This repository serves as a comprehensive guide and documentation of practical DevOps implementations, focusing on infrastructure provisioning, GitOps, container orchestration, and application deployment.

This repository is designed to document hands-on configurations and act as a reference for the broader community, as well as a technical knowledge base. 

---

## 🏗️ Repository Structure

Currently, the repository is structured into three main areas:

```text
devops-project-journey/
├── 📁 core-infrastructure/    # 🏗️ Hardware, Proxmox, & Gateway configurations
├── 📁 monitoring-secops/      # 🔐 Security & Monitoring configuration (Prometheus, Grafana, etc.)
├── 📁 platform-tools/         # ⚙️ GitOps, ArgoCD, & Platform-level manifests
└── 📁 projects/               # 🌟 Application Workloads & Demonstrations
    ├── 📁 01-hello-nginx/     # Project 1: K8s Basics & Nginx App
    └── 📁 02-nodejs-api/      # Project 2: Backend API + MySQL Database
```

---

## 🛠️ Tech Stack & Tools

These are the primary tools and technologies demonstrated in this repository:

- **Infrastructure & Provisioning:** Proxmox VE, Terraform, LXC
- **Container Orchestration & GitOps:** Kubernetes (K3s), ArgoCD
- **Applications & Databases:** Nginx, Node.js, MySQL
- **Networking & Security:** Kubernetes Services, Cloudflare Tunnels

---

## 🌟 Repository Sections

### 1. [Core Infrastructure](./core-infrastructure/)
Contains the foundational Infrastructure-as-Code (IaC) and configuration.
- **Proxmox K3s:** Terraform scripts to provision K3s nodes (Master and Workers) within LXC containers on a Proxmox host.
- **Gateway STB:** Docker Compose setup for Nginx and Cloudflared to establish a secure tunnel and act as a reverse proxy into the local network.

### 2. [Platform Tools](./platform-tools/)
Contains configuration for platform-level operational tools.
- **ArgoCD:** Demonstrates GitOps continuous delivery by automatically syncing manifests from this repository into the Kubernetes cluster.

### 3. [Security & Monitoring Operations](./monitoring-secops/)
Contains the setup for observability and basic security.
- **Monitoring LXC:** The infrastructure for monitoring is an LXC container, which is provisioned using the Terraform `main.tf` script located in the `core-infrastructure/proxmox-k3s` folder. Inside, it runs a monitoring stack (Prometheus, Grafana, Node Exporter, cAdvisor) via Docker Compose.
- **Security (UFW):** UFW (Uncomplicated Firewall) is recognized as a best practice for hardening security—such as closing down the server and allowing only specific essential ports. However, to streamline the initial setup phase and avoid the repetitive process of configuring it individually on every server, strict UFW configuration is deferred for now.

### 4. [Application Projects](./projects/)
Showcases individual application workloads deployed to the Kubernetes cluster.

- **[Hello Nginx (01-hello-nginx)](./projects/01-hello-nginx/):** Understand Kubernetes basics (Pods, Deployments, Services) by deploying a simple Nginx web server.
- **[Node.js API + Database (02-nodejs-api)](./projects/02-nodejs-api/):** Deploy a functional multi-tier application connecting a Node.js backend to a MySQL database using ConfigMaps and Secrets.

### 5. CI/CD & DevSecOps (GitHub Actions)
Continuous Integration and DevSecOps practices are integrated into the repository using GitHub Actions (located in `.github/workflows/`).
- **Automated Build & Scan:** The pipeline automatically builds Docker images for the application code (e.g., in `src-app/`) and pushes them to Docker Hub.
- **Vulnerability Scanning (Trivy):** During the build pipeline, images are scanned for critical and high vulnerabilities using Aqua Security's Trivy. 

> [!NOTE]
> **Trivy Vulnerability Exceptions (`.trivyignore`):** During the CI pipeline execution, the Trivy security scan initially flagged several vulnerabilities. Despite mitigating efforts—such as running `apk update && apk upgrade` in the Dockerfile—some vulnerabilities could not be resolved immediately. To unblock the CI pipeline, a `.trivyignore` file is utilized to bypass specific unfixed CVEs. This approach is a standard DevSecOps practice required because:
> 1. **Lag in Upstream Patches:** Even when pulling the latest available packages in the base image, the OS package repository might not yet have released a fix for a newly published CVE.
> 2. **Transitive NPM Dependencies:** Vulnerabilities frequently reside deep within `node_modules` (e.g., packages like `cross-spawn`, `glob`, `tar`). Updating the OS via `apk` does not patch vulnerable NPM packages locked in `package-lock.json`.
> 3. **Risk Acceptance & Exception Management:** In real-world scenarios, if a vulnerability has no available fix from the maintainers, it is often necessary to explicitly ignore it (`.trivyignore`) to allow automated deployments to proceed, while continuing to monitor for future patches.

---

## 🚀 How to Use This Repository

1. **Infrastructure:** Review `core-infrastructure/` to see how the base environment is provisioned using Terraform.
2. **Platform:** Check `platform-tools/` to understand how GitOps is configured with ArgoCD.
3. **Monitoring:** See `monitoring-secops/` for the observability stack (LXC provisioned via Terraform).
4. **Projects:** Explore the `projects/` folder sequentially to see how application manifests are created. If ArgoCD is set up, these projects will be automatically synchronized to the cluster. Alternatively, you can apply them manually (`kubectl apply -f <filename>.yaml`).
5. **CI/CD Pipelines:** Inspect the `.github/workflows/` directory to observe how automated Docker builds and DevSecOps scanning (Trivy) are implemented.

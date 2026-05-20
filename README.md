# 🚀 DevOps Project Portfolio

Welcome to my DevOps Project Portfolio! This repository serves as a showcase of my journey and skills in DevOps, focusing on infrastructure provisioning, GitOps, container orchestration, and application deployment.

This repository is designed to document my practical learning and act as a reference for myself and others. 

---

## 🏗️ Repository Structure

Currently, the repository is structured into three main areas:

```text
devops-project-portofolio/
├── 📁 core-infrastructure/    # 🏗️ Hardware, Proxmox, & Gateway configurations
├── 📁 platform-tools/         # ⚙️ GitOps, ArgoCD, & Platform-level manifests
└── 📁 projects/               # 🌟 Application Portfolio Showcase
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

## 🌟 Portfolio Sections

### 1. [Core Infrastructure](./core-infrastructure/)
Contains the foundational Infrastructure-as-Code (IaC) and configuration.
- **Proxmox K3s:** Terraform scripts to provision K3s nodes (Master and Workers) within LXC containers on a Proxmox host.
- **Gateway STB:** Docker Compose setup for Nginx and Cloudflared to establish a secure tunnel and act as a reverse proxy into the local network.

### 2. [Platform Tools](./platform-tools/)
Contains configuration for platform-level operational tools.
- **ArgoCD:** Demonstrates GitOps continuous delivery by automatically syncing manifests from this repository into the Kubernetes cluster.

### 3. [Featured Projects](./projects/)
Showcases individual application workloads deployed to the Kubernetes cluster.

- **[Hello Nginx (01-hello-nginx)](./projects/01-hello-nginx/):** Understand Kubernetes basics (Pods, Deployments, Services) by deploying a simple Nginx web server.
- **[Node.js API + Database (02-nodejs-api)](./projects/02-nodejs-api/):** Deploy a functional multi-tier application connecting a Node.js backend to a MySQL database using ConfigMaps and Secrets.

---

## 🚀 How to Use This Repository

1. **Infrastructure:** Review `core-infrastructure/` to see how the base environment is provisioned using Terraform.
2. **Platform:** Check `platform-tools/` to understand how GitOps is configured with ArgoCD.
3. **Projects:** Explore the `projects/` folder sequentially to see how application manifests are created. If ArgoCD is set up, these projects will be automatically synchronized to the cluster. Alternatively, you can apply them manually (`kubectl apply -f <filename>.yaml`).

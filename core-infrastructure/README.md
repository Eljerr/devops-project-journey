# Core Infrastructure

This directory contains the configuration and provisioning scripts for the core infrastructure of this project.

## Hardware Specifications

The infrastructure is hosted on the following hardware:

### 1. Proxmox Server (Old Laptop)
- **Hypervisor**: Proxmox VE
- **Hardware Type**: Used Laptop
- **CPU**: 2 Cores
- **RAM**: 8 GB
- **Purpose**: Hosts the K3s Kubernetes cluster (Master and Worker nodes) inside LXC containers.

### 2. Gateway STB (Set-Top Box)
- **Device**: Amlogic S905X4
- **RAM**: 2 GB
- **OS**: Ubuntu
- **Purpose**: Acts as a reverse proxy to securely route traffic from the internet into the local network via a Cloudflare Tunnel.

## Directory Structure

### `gateway-stb/`
Contains the configuration deployed on the STB gateway.
- Uses `docker-compose.yaml` to spin up **Nginx** (for reverse proxying) and **Cloudflared** (for establishing a secure Cloudflare Tunnel).
- Enables exposing local services to the internet without the need for port forwarding or a static public IP.

### `proxmox-k3s/`
Contains Infrastructure-as-Code (IaC) using Terraform to automatically provision the K3s cluster on the Proxmox server.
- `main.tf` and `variables.tf` define the creation of unprivileged LXC containers for the Kubernetes nodes.
- Automatically bootstraps K3s installation via SSH provisioners for the master and worker nodes.

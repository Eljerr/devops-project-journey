# Project 2: Node.js API + MySQL Database

This project demonstrates how to deploy a multi-tier application on Kubernetes, consisting of a Node.js backend API and a MySQL database.

## Objective
To successfully deploy interconnected services, manage configurations, and handle database deployments.

## Concepts Covered
- **Multi-tier Architecture:** Deploying both frontend/API and database layers.
- **ConfigMaps & Secrets:** Managing environment variables and sensitive data (like database passwords).
- **Networking:** Enabling communication between the Node.js backend and the MySQL database within the cluster.

> [!NOTE]
> The database password in this repository is a dummy value for demonstration purposes. In a real-world production environment, I would implement HashiCorp Vault or Kubernetes Sealed Secrets to dynamically inject credentials without committing them to version control.

## Files
- `database.yaml`: Contains the definitions for the MySQL database deployment, service, and necessary secrets/configurations.
- `backend.yaml`: Contains the definitions for the Node.js backend API deployment and service.

"Note: The database password in this repository is a dummy value for demonstration purposes. In a real-world production environment, I would implement HashiCorp Vault or Kubernetes Sealed Secrets to dynamically inject credentials without committing them to version control."

## Usage
To deploy this project to your Kubernetes cluster, run the following commands from this directory:

First, deploy the database:
```bash
kubectl apply -f database.yaml
```

Then, deploy the backend API:
```bash
kubectl apply -f backend.yaml
```

Verify the deployments are running:
```bash
kubectl get pods
kubectl get svc
```

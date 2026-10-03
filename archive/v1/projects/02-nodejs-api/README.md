# Project 2: Node.js API + MySQL Database

A multi-tier setup deploying a Node.js backend API alongside a MySQL database on Kubernetes. This is where things get more interesting — managing inter-service communication, configs, and secrets.

## What I Learned
- **Multi-tier Architecture:** Running both an API layer and a database layer as separate Kubernetes Deployments that communicate via a ClusterIP Service.
- **ConfigMaps & Secrets:** The MySQL password is stored in a `Secret` (`mysql-secret`) and injected into both the database and backend containers via `secretKeyRef` — so no plain-text credentials in the pod spec.
- **Networking:** The backend connects to MySQL using the service DNS name `mysql-service` as `DB_HOST`, which Kubernetes resolves internally within the cluster.

> [!NOTE]
> The database password in this repo (`secret`) is a dummy value for demonstration purposes. In a real production environment, I'd use HashiCorp Vault or Kubernetes Sealed Secrets to inject credentials dynamically without committing them to version control.

## Files
- `database.yaml`: MySQL Secret, Deployment (`mysql:8.0`), and ClusterIP Service on port `3306`.
- `backend.yaml`: Node.js backend Deployment (2 replicas, `mybinichizuru/project-02-api`) and NodePort Service accessible at port `30082`.

## Commands Used
Deploy the database first (backend depends on it):
```bash
kubectl apply -f database.yaml
```

Then deploy the backend:
```bash
kubectl apply -f backend.yaml
```

Verify both are running:
```bash
kubectl get pods
kubectl get svc
```

Access the API at `http://<node-ip>:30082`

# Project 1: Hello Nginx

My first hands-on experiment with Kubernetes — deploying a simple Nginx web server to understand the core building blocks.

## What I Learned
- **Deployment:** How Kubernetes manages stateless workloads. This manifest runs 2 replicas of Nginx (`nginx:1.25`) so the cluster keeps the app running even if one pod dies.
- **Service (NodePort):** How to expose an app outside the cluster. The Nginx service is accessible at port `30085` on any cluster node.

## Files
- `app-nginx.yaml`: Contains both the Deployment (2 replicas) and NodePort Service definitions.

## Commands Used
Deploy to the cluster:
```bash
kubectl apply -f app-nginx.yaml
```

Verify it's running:
```bash
kubectl get pods
kubectl get svc
```

Access the app at `http://<node-ip>:30085`

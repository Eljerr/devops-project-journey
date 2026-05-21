# Project 1: Hello Nginx

This project demonstrates the basics of Kubernetes by deploying a simple Nginx web server.

## Objective
To understand fundamental Kubernetes resources such as Deployments and Services.

## Concepts Covered
- **Deployment:** Managing stateless application workloads (Nginx).
- **Service:** Exposing the Nginx application to be accessible within or outside the cluster.

## Files
- `app-nginx.yaml`: Contains the Kubernetes Deployment and Service definitions for the Nginx application.

## Usage
To deploy this project to your Kubernetes cluster, run the following command from this directory:

```bash
kubectl apply -f app-nginx.yaml
```

Verify the deployment is running:
```bash
kubectl get pods
kubectl get svc
```

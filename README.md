# Canary Deployment

A sample project showcasing various Canary Deployment solutions.

## Medium Articles

- [Canary Deployment in Kubernetes (Part 1) — Simple Canary Deployment using Ingress NGINX](https://medium.com/@jhandguy/canary-deployment-in-kubernetes-part-1-simple-canary-deployment-using-ingress-nginx-f8f5da2b0f38)
- [Canary Deployment in Kubernetes (Part 2) — Automated Canary Deployment using Argo Rollouts](https://medium.com/@jhandguy/canary-deployment-in-kubernetes-part-2-automated-canary-deployment-using-argo-rollouts-8a3550d5a434)
- [Canary Deployment in Kubernetes (Part 3) — Smart Canary Deployment using Argo Rollouts and Prometheus](https://medium.com/@jhandguy/canary-deployment-in-kubernetes-part-3-smart-canary-deployment-using-argo-rollouts-and-47992d72222c)

## Installing

### Using ingress-nginx

```shell
kind create cluster --image kindest/node:v1.23.4 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.0.19 --wait

helm install ingress-nginx --name-template sample-app --create-namespace -n sample-app

helm upgrade sample-app ingress-nginx -n sample-app --reuse-values --set canary.weight=50

kind delete cluster
```

### Using argo-rollouts

```shell
kind create cluster --image kindest/node:v1.23.4 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.0.19 --wait

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --version 2.14.0 --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

kind delete cluster
```

### Using argo-rollouts + prometheus

```shell
kind create cluster --image kindest/node:v1.23.4 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.0.19 --wait

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --version 2.14.0 --wait

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-community/kube-prometheus-stack --name-template prometheus --create-namespace -n prometheus --version 34.8.0 --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app --set prometheus.enabled=true

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

kind delete cluster
```

## Smoke Testing

### Weighted canary

```shell
curl localhost/success -H "Host: sample.app" -v
curl localhost/error -H "Host: sample.app" -v
```

### Always canary

```shell
curl localhost/success -H "Host: sample.app" -H "X-Canary: always" -v
curl localhost/error -H "Host: sample.app" -H "X-Canary: always" -v
```

### Never canary

```shell
curl localhost/success -H "Host: sample.app" -H "X-Canary: never" -v
curl localhost/error -H "Host: sample.app" -H "X-Canary: never" -v
```

## Load Testing

```shell
k6 run k6/script.js
```

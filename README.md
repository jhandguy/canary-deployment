# Canary Deployment

A sample project showcasing various Canary Deployment solutions.

## Blog Posts

- [Canary Deployment in Kubernetes (Part 1) — Simple Canary Deployment using Ingress NGINX](https://jhandguy.github.io/posts/simple-canary-deployment/)
- [Canary Deployment in Kubernetes (Part 2) — Automated Canary Deployment using Argo Rollouts](https://jhandguy.github.io/posts/automated-canary-deployment/)
- [Canary Deployment in Kubernetes (Part 3) — Smart Canary Deployment using Argo Rollouts and Prometheus](https://jhandguy.github.io/posts/smart-canary-deployment/)

## Installing

### Using ingress-nginx

```shell
kind create cluster --image kindest/node:v1.27.3 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.8.3 --wait

helm install sample-app/helm-charts/ingress-nginx --name-template sample-app --create-namespace -n sample-app

helm upgrade sample-app sample-app/helm-charts/ingress-nginx -n sample-app --reuse-values --set canary.weight=50
```

### Using argo-rollouts

```shell
kind create cluster --image kindest/node:v1.27.3 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.8.3 --wait

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --version 2.32.5 --wait

helm install sample-app/helm-charts/argo-rollouts --name-template sample-app --create-namespace -n sample-app

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app
```

### Using argo-rollouts + prometheus

```shell
kind create cluster --image kindest/node:v1.27.3 --config=kind/cluster.yaml

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx/ingress-nginx --name-template ingress-nginx --create-namespace -n ingress-nginx --values kind/ingress-nginx-values.yaml --version 4.8.3 --wait

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --version 2.32.5 --wait

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-community/kube-prometheus-stack --name-template prometheus --create-namespace -n prometheus --version 54.2.2 --wait

helm install sample-app/helm-charts/argo-rollouts --name-template sample-app --create-namespace -n sample-app --set prometheus.enabled=true

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app
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

## Uninstalling

```shell
kind delete cluster
```

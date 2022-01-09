# Canary Deployment

A sample project showcasing various Canary Deployment solutions.

## Using ingress-nginx

```shell
minikube start --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm install ingress-nginx --name-template sample-app --create-namespace -n sample-app

helm upgrade sample-app ingress-nginx -n sample-app --set canary.weight=50

minikube stop && minikube delete
```

## Using argo-rollouts

```shell
minikube start --kubernetes-version=1.21.8 --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

minikube stop && minikube delete
```

## Using argo-rollouts + prometheus

```shell
minikube start --kubernetes-version=1.21.8 --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --wait

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-community/kube-prometheus-stack --name-template prometheus --create-namespace -n prometheus --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app --set prometheus.enabled=true

kubectl argo rollouts dashboard -n argo-rollouts &
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

minikube stop && minikube delete
```

## Smoke Testing

### Weighted canary

```shell
curl $(minikube ip)/success -H "Host: sample.app" -v
curl $(minikube ip)/error -H "Host: sample.app" -v
```

### Always canary

```shell
curl $(minikube ip)/success -H "Host: sample.app" -H "X-Canary: always" -v
curl $(minikube ip)/error -H "Host: sample.app" -H "X-Canary: always" -v
```

### Never canary

```shell
curl $(minikube ip)/success -H "Host: sample.app" -H "X-Canary: never" -v
curl $(minikube ip)/error -H "Host: sample.app" -H "X-Canary: never" -v
```

## Load Testing

```shell
env URL=$(minikube ip) k6 run k6/script.js
```

# canary-deployment

A sample project showcasing different Canary Deployment solutions using ingress-nginx and argo-rollouts.

## Using ingress-nginx

```
minikube start --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm install ingress-nginx --name-template sample-app --create-namespace -n sample-app

helm upgrade sample-app ingress-nginx --create-namespace -n sample-app --set canary.image.tag=latest --set canary.weight=50 --wait

minikube stop && minikube delete
```

## Using argo-rollouts

```
minikube start --kubernetes-version=1.21.8 --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app

kubectl argo rollouts dashboard -n argo-rollouts
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

minikube stop && minikube delete
```

## Using argo-rollouts + prometheus

```
minikube start --kubernetes-version=1.21.8 --addons=ingress $(if [ $(uname) != "Linux" ]; then echo "--vm=true"; fi)

helm repo add argo https://argoproj.github.io/argo-helm
helm install argo/argo-rollouts --name-template argo-rollouts --create-namespace -n argo-rollouts --set dashboard.enabled=true --wait

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-community/kube-prometheus-stack --name-template prometheus --create-namespace -n prometheus --wait

helm install argo-rollouts --name-template sample-app --create-namespace -n sample-app --set prometheus.enabled=true

kubectl argo rollouts dashboard -n argo-rollouts
kubectl argo rollouts set image sample-app sample-app=ghcr.io/jhandguy/canary-deployment/sample-app:latest -n sample-app
kubectl argo rollouts promote sample-app -n sample-app

minikube stop && minikube delete
```

## Testing

### Weighted canary

```
curl $(minikube ip)/success -H "Host: sample.app" -v
curl $(minikube ip)/error -H "Host: sample.app" -v
```

### Always canary

```
curl $(minikube ip)/success -H "Host: sample.app" -H "X-Canary: always" -v
curl $(minikube ip)/error -H "Host: sample.app" -H "X-Canary: always" -v
```

### Never canary

```
curl $(minikube ip)/success -H "Host: sample.app" -H "X-Canary: never" -v
curl $(minikube ip)/error -H "Host: sample.app" -H "X-Canary: never" -v
```

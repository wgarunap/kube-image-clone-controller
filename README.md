# Kubernetes Public Image Cache Controller

## Overview
Kubernetes controller for automatically backup public 
images to user's registry which are used in Deployments and DaemonSets. 
kube-system and controller's namespaces are ignored. 

### Testing On MacOS
##### Install & Start Minikube Cluster 
```shell
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

# Start Cluster
minikube start

# Start Dashboard
minikube dashboard
```

If you're going to use local docker images for testing, use minikube for image building as follows.
1. Set the environment variables with eval $(minikube docker-env)
2. Build the image with the Docker daemon of Minikube (eg docker build -t my-image .)
3. Set the image in the pod spec like the build tag (eg my-image)
4. Set the imagePullPolicy to Never, otherwise Kubernetes will try to download the image.

### Build & Push Your Own Controller
```shell
docker build -t <username>/kube-image-clone-controller:latest .
```

```shell
docker push <username>/kube-image-clone-controller:latest
```

### Pull Public Image from Dockerhub
```shell
docker push testingnew123/kube-image-clone-controller:latest
```

## Quick Start

Create the NameSpace
```shell
kubectl create namespace image-clone-namespace
```
Set the Target Docker Registry Credentials as a Secret
```shell
kubectl create secret --namespace=image-clone-namespace generic docker-registry-credentials \
  --from-literal=docker-server=index.docker.io \
  --from-literal=docker-username=testingnew123 \
  --from-literal=docker-password=xxxx
```

Kubernetes config file is prepared with necessary RBAC config to start.
```shell
kubectl apply -f .kubenates/k8s.yaml
```
### DEMO
[![asciicast](https://asciinema.org/a/qKM2HhpeM1KkHZCIOdgZJ7baH.svg)](https://asciinema.org/a/qKM2HhpeM1KkHZCIOdgZJ7baH)


### Stop Controller
```shell
kubectl delete -f .kubenates/k8s.yaml
```
NOTE: This will remove the namespace and all data associated with image-clone-controller.

### Special Notes and Assumptions
1. Source Repository is properly tagged and previous tags will not be overridden.(once a tag is cloned it will not be cloned again until it remains in the target registry)
2. If `latest` tag is referred in the image it will always clone to the target registry.

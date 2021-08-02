# kube-image-clone-controller

### Overview
Kubernetes controller for automatically backup public images to user's registry which are used in Deployments and DaemonSets. kube-system and controller's namespaces are ignored. 

### Testing 
##### Install & Start Minikube Cluster 
```shell
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

# Start Cluster
minikube start

# Start Dashboard
minikube dashboard
```

### Build Controller
```shell
docker build -t <username>/kube-image-clone-controller:latest .
```

### Push Controller
```shell
docker push <username>/kube-image-clone-controller:latest
```

### Pull Public Kube Image Clone Controller 
```shell
docker push testingnew123/kube-image-clone-controller:latest
```

### Deployment
##### 
```shell
kubectl apply -f .kubenates/k8s.yaml
```

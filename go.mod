module github.com/wgarunap/kube-image-clone-controller

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/docker/cli v20.10.7+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.5.0
	github.com/google/go-containerregistry v0.6.0
	github.com/stretchr/testify v1.7.0
	github.com/wgarunap/goconf v0.5.0
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
)

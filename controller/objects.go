package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object interface {
	// Name is the given name for the object
	Name() string

	//Get return  the client object inside the binded object
	Get() client.Object

	// NewCopy return new deep copy of the object
	NewCopy() Object

	//Containers return list of all available containers inside the object
	Containers() []v1.Container

	// OverrideImage overrides the specific image name which holds the given index
	OverrideImage(containerIndex int, newImage string)
}

type deployment struct {
	Object *appsv1.Deployment
}

func (d *deployment) Name() string {
	return `Deployment`
}

func (d *deployment) Get() client.Object {
	return d.Object
}

func (d *deployment) NewCopy() Object {
	return &deployment{Object: d.Object.DeepCopy()}
}

func (d *deployment) Containers() []v1.Container {
	return d.Object.Spec.Template.Spec.Containers
}

func (d *deployment) OverrideImage(containerIndex int, newImage string) {
	d.Object.Spec.Template.Spec.Containers[containerIndex].Image = newImage
}

func NewDeploymentObject() Object {
	return &deployment{
		Object: &appsv1.Deployment{},
	}
}

type daemonSet struct {
	Object *appsv1.DaemonSet
}

func (d *daemonSet) Name() string {
	return `DaemonSet`
}

func (d *daemonSet) Get() client.Object {
	return d.Object
}

func (d *daemonSet) NewCopy() Object {
	return &daemonSet{Object: d.Object.DeepCopy()}
}

func (d *daemonSet) Containers() []v1.Container {
	return d.Object.Spec.Template.Spec.Containers
}

func (d *daemonSet) OverrideImage(containerIndex int, newImage string) {
	d.Object.Spec.Template.Spec.Containers[containerIndex].Image = newImage
}

func NewDaemonSetObject() Object {
	return &daemonSet{
		Object: &appsv1.DaemonSet{},
	}
}

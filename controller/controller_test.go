package controller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wgarunap/kube-image-clone-controller/config"
	"github.com/wgarunap/kube-image-clone-controller/domain"
	"github.com/wgarunap/kube-image-clone-controller/mocks"
	appsv1 "k8s.io/api/apps/v1"
	cv1 "k8s.io/api/core/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func fakeClient() client.Client {
	return fake.NewClientBuilder().WithObjects(
		&appsv1.Deployment{
			ObjectMeta: mv1.ObjectMeta{Name: "test", Namespace: "something"},
			Spec: appsv1.DeploymentSpec{
				Template: cv1.PodTemplateSpec{
					Spec: cv1.PodSpec{
						Containers: []cv1.Container{{
							Name:  "test_container",
							Image: "test_images/image:v1.0",
						}, {
							Name:  "test_container_222",
							Image: "test_images_222/image2:v1.0",
						}, {
							Name:  "test_container_333",
							Image: "test_container_333/image3:latest",
						}, {
							Name:  "container_clone_error",
							Image: "container_clone_error/image:v1.0",
						}, {
							Name:  "container_no_tag_check",
							Image: "container_no_tag_check/image",
						}, {
							Name:  "dockerhub_container",
							Image: "index.docker.io/dockerhub_container/image",
						}, {
							Name:  "google_container",
							Image: "acr.gcp.io/dockerhub_container/image",
						}},
					},
				},
			},
		},
		&appsv1.DaemonSet{
			ObjectMeta: mv1.ObjectMeta{Name: "test_daemon", Namespace: "something"},
			Spec: appsv1.DaemonSetSpec{
				Template: cv1.PodTemplateSpec{
					Spec: cv1.PodSpec{
						Containers: []cv1.Container{{
							Name:  "test_daemon_container",
							Image: "test_images_daemon/image2:v1.0",
						}},
					},
				},
			},
		},
	).Build()
}

func fakeCloner(t *testing.T) domain.Cloner {
	ctrl := gomock.NewController(t)
	mockcloner := mocks.NewMockCloner(ctrl)
	mockcloner.EXPECT().
		Clone(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	r := &reconcileImage{cfg: config.Config}
	newName, _, _ := r.generateNewImageName("test_images/image:v1.0")
	mockcloner.EXPECT().
		IsExistInClones(context.Background(), gomock.Eq(newName)).
		Return(nil, true).MaxTimes(1)

	mockcloner.EXPECT().
		IsExistInClones(context.Background(), gomock.Any()).
		Return(nil, false).MaxTimes(2)

	return mockcloner
}

func TestReconcileImage_ReconcileDevelopment(t *testing.T) {
	config.Config.TargetRegistry = "localhost:5000"
	config.Config.UserName = "new_user"
	log.SetLogger(zap.New())

	out := []string{
		"localhost:5000/new_user/test_images_image:v1.0",
		"localhost:5000/new_user/test_images_222_image2:v1.0",
		"localhost:5000/new_user/test_container_333_image3:latest",
		"localhost:5000/new_user/container_clone_error_image:v1.0",
		"localhost:5000/new_user/container_no_tag_check_image:latest",
		"localhost:5000/new_user/dockerhub_container_image:latest",
		"localhost:5000/new_user/dockerhub_container_image:latest",
	}

	reconcileOb := &reconcileImage{
		client: fakeClient(),
		object: NewDeploymentObject(),
		cfg:    config.Config,
		cloner: fakeCloner(t),
	}

	namespace := types.NamespacedName{
		Namespace: "something",
		Name:      "test",
	}
	_, err := reconcileOb.Reconcile(context.Background(), reconcile.Request{
		NamespacedName: namespace,
	})
	if err != nil {
		t.Error(err)
	}

	rs := &appsv1.Deployment{}
	err = reconcileOb.client.Get(context.Background(), namespace, rs)
	if err != nil {
		t.Error(err)
	}
	for i, container := range rs.Spec.Template.Spec.Containers {
		require.Equal(t, out[i], container.Image)
	}
}

func TestReconcileImage_ReconcileDaemonSet(t *testing.T) {
	config.Config.TargetRegistry = "localhost:5000"
	config.Config.UserName = "new_user"
	log.SetLogger(zap.New())

	out := []string{
		"localhost:5000/new_user/test_images_daemon_image2:v1.0",
	}

	reconcileOb := &reconcileImage{
		client: fakeClient(),
		object: NewDaemonSetObject(),
		cfg:    config.Config,
		cloner: fakeCloner(t),
	}

	namespace := types.NamespacedName{
		Namespace: "something",
		Name:      "test_daemon",
	}
	_, err := reconcileOb.Reconcile(context.Background(), reconcile.Request{
		NamespacedName: namespace,
	})
	if err != nil {
		t.Error(err)
	}

	rs := &appsv1.DaemonSet{}
	err = reconcileOb.client.Get(context.Background(), namespace, rs)
	if err != nil {
		t.Error(err)
	}
	for i, container := range rs.Spec.Template.Spec.Containers {
		require.Equal(t, out[i], container.Image)
	}
}

func TestGenerateNewImageName(t *testing.T) {
	username := `kgkgkg`
	cloneRegistry := `test_registry:9000`
	config.Config.TargetRegistry = cloneRegistry
	config.Config.UserName = username

	prefix := cloneRegistry + "/" + username
	testData := []struct {
		image        string
		outName      string
		changeStatus bool
	}{{
		image:        "user/image:v1.0.0",
		outName:      prefix + "/user_image:v1.0.0",
		changeStatus: true,
	}, {
		image:        prefix + "/image:v1.0.55",
		outName:      prefix + "/image:v1.0.55",
		changeStatus: false,
	}, {
		image:        "dockerhub/image-asdad:latest",
		outName:      prefix + "/dockerhub_image-asdad:latest",
		changeStatus: true,
	}, {
		image:        "dockerhub/image-asdad",
		outName:      prefix + "/dockerhub_image-asdad:latest",
		changeStatus: true,
	}, {
		image:        "wgarunap/generic-kafka-event-producer",
		outName:      prefix + "/wgarunap_generic-kafka-event-producer:latest",
		changeStatus: true,
	}}

	r := reconcileImage{cfg: config.Config}
	for _, datum := range testData {
		out, changed, err := r.generateNewImageName(datum.image)
		require.Equal(t, nil, err)
		require.Equal(t, datum.outName, out.Name())
		require.Equal(t, datum.changeStatus, changed)
		//println(out.Name())
	}
}

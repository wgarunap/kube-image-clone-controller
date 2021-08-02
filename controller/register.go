package controller

import (
	"github.com/go-logr/logr"
	"github.com/wgarunap/kube-image-clone-controller/config"
	"github.com/wgarunap/kube-image-clone-controller/domain"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Register register the different type of objects with the controller and
// setup the watcher conditions
func Register(ob Object, log logr.Logger, mgr manager.Manager, cloner domain.Cloner) error {
	c, err := controller.New(ob.Name(), mgr, controller.Options{
		Reconciler: &reconcileImage{
			client: mgr.GetClient(),
			object: ob,
			cfg:    config.Config,
			cloner: cloner,
		},
	})
	if err != nil {
		log.Error(err, "unable to set up individual controller")
		return err
	}

	// Watch received Object and based on the given predicate
	err = c.Watch(&source.Kind{Type: ob.Get()}, &handler.EnqueueRequestForObject{},
		predicate.And(
			predicate.NewPredicateFuncs(func(object client.Object) bool {
				switch object.GetNamespace() {
				case "kube-system", "kubernetes-dashboard", "image-clone-namespace":
					return false
				}
				return true
			}),
			predicate.Funcs{
				DeleteFunc: func(event event.DeleteEvent) bool {
					return false
				},
				GenericFunc: func(event event.GenericEvent) bool {
					return false
				},
			},
		),
	)

	if err != nil {
		log.Error(err, "unable to watch "+ob.Get().GetName())
		return err
	}
	return nil
}

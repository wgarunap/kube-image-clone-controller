package controller

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/wgarunap/kube-image-clone-controller/config"
	"github.com/wgarunap/kube-image-clone-controller/domain"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = (*reconcileImage)(nil)

// reconcileImage reconciles ReplicaSets
type reconcileImage struct {
	// Client can be used to retrieve objects from the APIServer.
	client client.Client

	object Object

	cfg config.Conf

	// cloner
	cloner domain.Cloner
}

func (r *reconcileImage) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	logger := log.FromContext(ctx)

	fmt.Println(`Reconcile containers starting`, request.String())

	// Fetch the controller
	rs := r.object.Get()
	err := r.client.Get(ctx, request.NamespacedName, rs)
	if errors.IsNotFound(err) {
		logger.Error(nil, "error finding "+r.object.Name())
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch %s: %+v", r.object.Name(), err)
	}

	newCopy := r.object.NewCopy()
	wg := sync.WaitGroup{}
	errorChan := make(chan error, len(r.object.Containers()))

	for i, container := range r.object.Containers() {
		newName, isChanged, err := r.generateNewImageName(container.Image)
		if err != nil {
			logger.Error(err, "error generating image name")
			return reconcile.Result{}, err
		}

		logger.Info("container process", container.Image, newName.Name())

		if isChanged {
			newCopy.OverrideImage(i, newName.Name())

			wg.Add(1)
			go func(imageName string) {
				defer wg.Done()

				source, _ := name.ParseReference(imageName)
				if source.Identifier() != `latest` {
					_, exist := r.cloner.IsExistInClones(ctx, newName)
					if exist {
						return
					}
				}

				logger.Info(`clonning image`, "newName", newName.Name())

				err := r.cloner.Clone(ctx, source, newName)
				if err != nil {
					errorChan <- err
				}
			}(container.Image)
		}
	}
	wg.Wait()
	close(errorChan)

	var errs []string
	for err := range errorChan {
		errs = append(errs, err.Error())
		logger.Error(err, "error cloning")
	}
	if len(errs) > 0 {
		return reconcile.Result{}, fmt.Errorf("error clonning the docker images: %v", strings.Join(errs, ", "))
	}

	patchObject := client.StrategicMergeFrom(rs)

	// Patch data object
	// NOTE: if used Update instead of patch, it will conflict with the parallel changes and output errors
	err = r.client.Patch(ctx, newCopy.Get(), patchObject) //err = r.client.Update(ctx, rs)
	if err != nil {
		logger.Error(err, "Unable to update the containers")
		return reconcile.Result{}, fmt.Errorf("could not write %s: %+v", r.object.Name(), err)
	}

	logger.Info(`reconcile completed`)
	return reconcile.Result{}, nil
}

func (r *reconcileImage) generateNewImageName(imageName string) (n name.Reference, isChanged bool, err error) {
	source, err := name.ParseReference(imageName)
	if err != nil {
		return nil, false, err
	}
	if strings.Contains(source.String(), r.cfg.TargetRegistry+"/"+r.cfg.UserName) {
		return source, false, nil
	}

	target, err := name.ParseReference(r.cfg.UserName+"/"+strings.ReplaceAll(source.Context().RepositoryStr(), "/", "_"), name.WithDefaultRegistry(r.cfg.TargetRegistry))
	if err != nil {
		return nil, false, err
	}
	return target, true, nil
}

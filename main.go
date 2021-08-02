package main

import (
	appconf "github.com/wgarunap/kube-image-clone-controller/config"
	"github.com/wgarunap/kube-image-clone-controller/imagecloner"
	"os"

	"github.com/wgarunap/goconf"
	imageclonecontroller "github.com/wgarunap/kube-image-clone-controller/controller"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("clone-controller")

	if err := goconf.Load(&appconf.Conf{}); err != nil {
		panic(err)
	}

	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	cloner := imagecloner.NewCloner(entryLog, appconf.Config)

	entryLog.Info("registering deployment controller")
	if err := imageclonecontroller.Register(imageclonecontroller.NewDeploymentObject(), entryLog, mgr, cloner); err != nil {
		entryLog.Error(err, "unable to set up Deployment controller")
		os.Exit(1)
	}

	entryLog.Info("registering daemonSet controller")
	if err := imageclonecontroller.Register(imageclonecontroller.NewDaemonSetObject(), entryLog, mgr, cloner); err != nil {
		entryLog.Error(err, "unable to set up DaemonSet controller")
		os.Exit(1)
	}

	entryLog.Info("starting manager")

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

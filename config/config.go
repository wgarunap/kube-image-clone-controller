package config

import (
	"errors"
	"github.com/caarlos0/env"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Conf is the application config object
type Conf struct {
	TargetRegistry string `env:"TARGET_REGISTRY" envDefault:"index.docker.io"`
	UserName       string `env:"USERNAME"`
	Password       string `env:"PASSWORD"`
}

var Config Conf

//Register will be called for config registration from env variables
func (*Conf) Register() error {
	err := env.Parse(&Config)
	if err != nil {
		log.Log.Error(err, "error loading controller config")
		return err
	}
	return nil
}

// Validate validates the configs
func (*Conf) Validate() (err error) {
	if Config.TargetRegistry == "" {
		log.Log.Error(errors.New("env TARGET_REGISTRY not found"), "unable to find target registry")
	}
	return nil
}

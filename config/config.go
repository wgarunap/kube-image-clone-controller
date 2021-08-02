package config

import (
	"errors"
	"github.com/caarlos0/env"
	"github.com/docker/cli/cli/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Conf struct {
	TargetRegistry  string `env:"TARGET_REGISTRY" envDefault:"index.docker.io"`
	DockerConfigDir string `env:"DOCKER_CONFIG"`
	UserName        string // internal use
}

var Config Conf

func (*Conf) Register() error {
	err := env.Parse(&Config)
	if err != nil {
		log.Log.Error(err, "error loading controller config")
		return err
	}
	return nil
}

func (*Conf) Validate() (err error) {
	if Config.TargetRegistry == "" {
		log.Log.Error(errors.New("env TARGET_REGISTRY not found"), "unable to find target registry")
	}

	Config.UserName, err = username(Config.TargetRegistry, Config.DockerConfigDir)
	if err != nil {
		return err
	}
	return nil
}

func username(target string, dockerConfig string) (string, error) {
	cf, err := config.Load(dockerConfig)
	if err != nil {
		return "", err
	}

	for t, ac := range cf.AuthConfigs {
		if t == target {
			return ac.Username, nil
		}
	}
	return "", errors.New(`username not found for the target`)
}

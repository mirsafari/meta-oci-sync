package config

import (
	"log/slog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	HTTPPort int
	K8S      *kubernetes.Clientset
}

var Clients Config

func (cfg *Config) Load() error {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	slog.Info("In cluster K8S client config created")

	cfg.HTTPPort = 8080
	cfg.K8S = clientset

	return nil
}

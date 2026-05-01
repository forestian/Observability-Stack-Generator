package config

import (
	"fmt"
	"strings"
)

const (
	ProfileDev        = "dev"
	ProfileProduction = "production"

	StorageMinIO = "minio"
	StorageS3    = "s3"
)

type Options struct {
	Name            string
	Namespace       string
	OutputDir       string
	Profile         string
	Storage         string
	StorageExplicit bool
}

type StackConfig struct {
	Name      string
	Namespace string
	OutputDir string
	Profile   string
	Storage   StorageConfig

	Loki  ComponentConfig
	Mimir ComponentConfig
	Tempo ComponentConfig
	Alloy ComponentConfig
}

type StorageConfig struct {
	Type                string
	Endpoint            string
	Region              string
	AccessKeySecretName string
	SecretKeySecretName string
	Insecure            bool
}

type ComponentConfig struct {
	Enabled   bool
	Retention string
	Replicas  int
}

func DefaultOptions() Options {
	return Options{
		Name:      "observability-stack",
		Namespace: "monitoring",
		OutputDir: "./observability-stack",
		Profile:   ProfileDev,
		Storage:   StorageMinIO,
	}
}

func NewStackConfig(opts Options) (StackConfig, error) {
	name := strings.TrimSpace(opts.Name)
	namespace := strings.TrimSpace(opts.Namespace)
	outputDir := strings.TrimSpace(opts.OutputDir)
	profile := strings.TrimSpace(opts.Profile)
	storageType := strings.TrimSpace(opts.Storage)

	if name == "" {
		return StackConfig{}, fmt.Errorf("name must not be empty")
	}
	if namespace == "" {
		return StackConfig{}, fmt.Errorf("namespace must not be empty")
	}
	if outputDir == "" {
		return StackConfig{}, fmt.Errorf("output must not be empty")
	}
	if profile != ProfileDev && profile != ProfileProduction {
		return StackConfig{}, fmt.Errorf("profile must be either %q or %q", ProfileDev, ProfileProduction)
	}

	defaults := defaultsForProfile(profile)
	if storageType == "" || (!opts.StorageExplicit && profile == ProfileProduction) {
		storageType = defaults.storageType
	}
	if storageType != StorageMinIO && storageType != StorageS3 {
		return StackConfig{}, fmt.Errorf("storage must be either %q or %q", StorageMinIO, StorageS3)
	}

	return StackConfig{
		Name:      name,
		Namespace: namespace,
		OutputDir: outputDir,
		Profile:   profile,
		Storage:   newStorageConfig(storageType, namespace),
		Loki: ComponentConfig{
			Enabled:   true,
			Retention: defaults.lokiRetention,
			Replicas:  defaults.replicas,
		},
		Mimir: ComponentConfig{
			Enabled:   true,
			Retention: defaults.mimirRetention,
			Replicas:  defaults.replicas,
		},
		Tempo: ComponentConfig{
			Enabled:   true,
			Retention: defaults.tempoRetention,
			Replicas:  defaults.replicas,
		},
		Alloy: ComponentConfig{
			Enabled:  true,
			Replicas: defaults.replicas,
		},
	}, nil
}

type profileDefaults struct {
	lokiRetention  string
	mimirRetention string
	tempoRetention string
	replicas       int
	storageType    string
}

func defaultsForProfile(profile string) profileDefaults {
	if profile == ProfileProduction {
		return profileDefaults{
			lokiRetention:  "30d",
			mimirRetention: "30d",
			tempoRetention: "72h",
			replicas:       2,
			storageType:    StorageS3,
		}
	}

	return profileDefaults{
		lokiRetention:  "7d",
		mimirRetention: "7d",
		tempoRetention: "24h",
		replicas:       1,
		storageType:    StorageMinIO,
	}
}

func newStorageConfig(storageType, namespace string) StorageConfig {
	storage := StorageConfig{
		Type:                storageType,
		Region:              "us-east-1",
		AccessKeySecretName: "obsgen-object-storage",
		SecretKeySecretName: "obsgen-object-storage",
	}

	if storageType == StorageMinIO {
		storage.Endpoint = fmt.Sprintf("http://minio.%s.svc.cluster.local:9000", namespace)
		storage.Insecure = true
		return storage
	}

	storage.Endpoint = "https://s3.us-east-1.amazonaws.com"
	storage.Insecure = false
	return storage
}

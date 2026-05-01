package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg, err := NewStackConfig(DefaultOptions())
	if err != nil {
		t.Fatalf("NewStackConfig returned error: %v", err)
	}

	if cfg.Name != "observability-stack" {
		t.Fatalf("expected default name, got %q", cfg.Name)
	}
	if cfg.Namespace != "monitoring" {
		t.Fatalf("expected default namespace, got %q", cfg.Namespace)
	}
	if cfg.OutputDir != "./observability-stack" {
		t.Fatalf("expected default output, got %q", cfg.OutputDir)
	}
	if cfg.Profile != ProfileDev {
		t.Fatalf("expected dev profile, got %q", cfg.Profile)
	}
	if cfg.Storage.Type != StorageMinIO {
		t.Fatalf("expected minio storage, got %q", cfg.Storage.Type)
	}
	if cfg.Loki.Retention != "7d" || cfg.Mimir.Retention != "7d" || cfg.Tempo.Retention != "24h" {
		t.Fatalf("unexpected dev retention defaults: loki=%s mimir=%s tempo=%s", cfg.Loki.Retention, cfg.Mimir.Retention, cfg.Tempo.Retention)
	}
	if cfg.Loki.Replicas != 1 || cfg.Mimir.Replicas != 1 || cfg.Tempo.Replicas != 1 || cfg.Alloy.Replicas != 1 {
		t.Fatalf("expected dev replicas to be 1")
	}
}

func TestProductionProfileDefaultsStorageToS3WhenNotExplicit(t *testing.T) {
	opts := DefaultOptions()
	opts.Profile = ProfileProduction

	cfg, err := NewStackConfig(opts)
	if err != nil {
		t.Fatalf("NewStackConfig returned error: %v", err)
	}

	if cfg.Storage.Type != StorageS3 {
		t.Fatalf("expected production default storage s3, got %q", cfg.Storage.Type)
	}
	if cfg.Loki.Retention != "30d" || cfg.Mimir.Retention != "30d" || cfg.Tempo.Retention != "72h" {
		t.Fatalf("unexpected production retention defaults: loki=%s mimir=%s tempo=%s", cfg.Loki.Retention, cfg.Mimir.Retention, cfg.Tempo.Retention)
	}
	if cfg.Loki.Replicas != 2 || cfg.Mimir.Replicas != 2 || cfg.Tempo.Replicas != 2 || cfg.Alloy.Replicas != 2 {
		t.Fatalf("expected production replicas to be 2")
	}
}

func TestProductionProfileAllowsExplicitMinIO(t *testing.T) {
	opts := DefaultOptions()
	opts.Profile = ProfileProduction
	opts.Storage = StorageMinIO
	opts.StorageExplicit = true

	cfg, err := NewStackConfig(opts)
	if err != nil {
		t.Fatalf("NewStackConfig returned error: %v", err)
	}

	if cfg.Storage.Type != StorageMinIO {
		t.Fatalf("expected explicit minio storage, got %q", cfg.Storage.Type)
	}
}

func TestInvalidProfile(t *testing.T) {
	opts := DefaultOptions()
	opts.Profile = "staging"

	if _, err := NewStackConfig(opts); err == nil {
		t.Fatalf("expected invalid profile error")
	}
}

func TestInvalidStorageType(t *testing.T) {
	opts := DefaultOptions()
	opts.Storage = "azure"
	opts.StorageExplicit = true

	if _, err := NewStackConfig(opts); err == nil {
		t.Fatalf("expected invalid storage error")
	}
}

func TestEmptyOutputDirectory(t *testing.T) {
	opts := DefaultOptions()
	opts.OutputDir = " "

	if _, err := NewStackConfig(opts); err == nil {
		t.Fatalf("expected empty output error")
	}
}

func TestEmptyNamespace(t *testing.T) {
	opts := DefaultOptions()
	opts.Namespace = " "

	if _, err := NewStackConfig(opts); err == nil {
		t.Fatalf("expected empty namespace error")
	}
}

package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"observability-stack-generator/internal/config"
)

func TestOutputDirectoryOverwriteBehavior(t *testing.T) {
	outputDir := t.TempDir()
	cfg := testConfig(t, outputDir, config.ProfileDev, config.StorageMinIO, true)

	if _, err := Generate(cfg, Options{}); err == nil {
		t.Fatalf("expected existing output directory to fail without force")
	}

	if _, err := Generate(cfg, Options{Force: true}); err != nil {
		t.Fatalf("expected force generation to succeed: %v", err)
	}
}

func TestTemplateRenderCreatesExpectedFiles(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "demo-stack")
	cfg := testConfig(t, outputDir, config.ProfileDev, config.StorageMinIO, true)

	if _, err := Generate(cfg, Options{}); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	expectedFiles := []string{
		"README.md",
		"values/loki-values.yaml",
		"values/mimir-values.yaml",
		"values/tempo-values.yaml",
		"values/alloy-values.yaml",
		"storage/minio-values.yaml",
		"storage/object-storage-notes.md",
		"examples/install.sh",
		"examples/uninstall.sh",
	}

	for _, expected := range expectedFiles {
		path := filepath.Join(outputDir, filepath.FromSlash(expected))
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected generated file %s: %v", expected, err)
		}
	}

	readme := readFile(t, filepath.Join(outputDir, "README.md"))
	if !strings.Contains(readme, "Grafana Loki") {
		t.Fatalf("expected README to mention Grafana Loki")
	}

	installScript := readFile(t, filepath.Join(outputDir, "examples", "install.sh"))
	if !strings.Contains(installScript, "helm upgrade --install minio minio/minio") {
		t.Fatalf("expected MinIO install command in minio mode")
	}
}

func TestS3GenerationOmitsMinIO(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "demo-s3-stack")
	cfg := testConfig(t, outputDir, config.ProfileProduction, config.StorageS3, true)

	if _, err := Generate(cfg, Options{}); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	minioValuesPath := filepath.Join(outputDir, "storage", "minio-values.yaml")
	if _, err := os.Stat(minioValuesPath); !os.IsNotExist(err) {
		t.Fatalf("expected minio-values.yaml to be omitted for s3 storage")
	}

	installScript := readFile(t, filepath.Join(outputDir, "examples", "install.sh"))
	if strings.Contains(installScript, "minio/minio") {
		t.Fatalf("expected install script to omit MinIO install commands for s3 storage")
	}

	readme := readFile(t, filepath.Join(outputDir, "README.md"))
	if !strings.Contains(readme, "production-oriented starter configuration") {
		t.Fatalf("expected production README phrase")
	}
}

func testConfig(t *testing.T, outputDir, profile, storage string, storageExplicit bool) config.StackConfig {
	t.Helper()

	opts := config.DefaultOptions()
	opts.Name = "demo"
	opts.Namespace = "monitoring"
	opts.OutputDir = outputDir
	opts.Profile = profile
	opts.Storage = storage
	opts.StorageExplicit = storageExplicit

	cfg, err := config.NewStackConfig(opts)
	if err != nil {
		t.Fatalf("NewStackConfig returned error: %v", err)
	}

	return cfg
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

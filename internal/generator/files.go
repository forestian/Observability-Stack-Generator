package generator

import "observability-stack-generator/internal/config"

type fileSpec struct {
	path       string
	template   string
	executable bool
	enabled    func(config.StackConfig) bool
}

func filesForStack() []fileSpec {
	return []fileSpec{
		{path: "README.md", template: "README.md.tmpl"},
		{path: "values/loki-values.yaml", template: "loki-values.yaml.tmpl"},
		{path: "values/mimir-values.yaml", template: "mimir-values.yaml.tmpl"},
		{path: "values/tempo-values.yaml", template: "tempo-values.yaml.tmpl"},
		{path: "values/alloy-values.yaml", template: "alloy-values.yaml.tmpl"},
		{
			path:     "storage/minio-values.yaml",
			template: "minio-values.yaml.tmpl",
			enabled: func(cfg config.StackConfig) bool {
				return cfg.Storage.Type == config.StorageMinIO
			},
		},
		{path: "storage/object-storage-notes.md", template: "object-storage-notes.md.tmpl"},
		{path: "examples/install.sh", template: "install.sh.tmpl", executable: true},
		{path: "examples/uninstall.sh", template: "uninstall.sh.tmpl", executable: true},
	}
}

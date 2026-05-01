package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"observability-stack-generator/internal/config"
	templatefiles "observability-stack-generator/internal/templates"
)

type Options struct {
	Force bool
}

type Result struct {
	OutputDir string
	Files     []string
}

func Generate(cfg config.StackConfig, opts Options) (Result, error) {
	if err := validateOutputDir(cfg.OutputDir, opts.Force); err != nil {
		return Result{}, err
	}

	tmpl, err := template.ParseFS(templatefiles.FS, "*.tmpl")
	if err != nil {
		return Result{}, fmt.Errorf("parse templates: %w", err)
	}

	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return Result{}, fmt.Errorf("create output directory: %w", err)
	}

	result := Result{OutputDir: cfg.OutputDir}
	for _, spec := range filesForStack() {
		if spec.enabled != nil && !spec.enabled(cfg) {
			continue
		}

		content, err := renderTemplate(tmpl, spec.template, cfg)
		if err != nil {
			return Result{}, err
		}

		outputPath := filepath.Join(cfg.OutputDir, filepath.FromSlash(spec.path))
		if err := writeGeneratedFile(outputPath, content, spec.executable, opts.Force); err != nil {
			return Result{}, err
		}

		result.Files = append(result.Files, outputPath)
	}

	return result, nil
}

func validateOutputDir(outputDir string, force bool) error {
	info, err := os.Stat(outputDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("output path %q exists and is not a directory", outputDir)
		}
		if !force {
			return fmt.Errorf("output directory %q already exists; use --force to overwrite generated files", outputDir)
		}
		return nil
	}
	if os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("inspect output directory: %w", err)
}

func renderTemplate(tmpl *template.Template, name string, cfg config.StackConfig) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, cfg); err != nil {
		return nil, fmt.Errorf("render %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

func writeGeneratedFile(path string, content []byte, executable bool, force bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create directory for %s: %w", path, err)
	}

	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("file %q already exists; use --force to overwrite generated files", path)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inspect file %s: %w", path, err)
	}

	mode := os.FileMode(0644)
	if executable {
		mode = 0755
	}

	if err := os.WriteFile(path, content, mode); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

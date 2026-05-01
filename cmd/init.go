package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"observability-stack-generator/internal/config"
	"observability-stack-generator/internal/generator"
)

func newInitCommand() *cobra.Command {
	opts := config.DefaultOptions()
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate starter Helm values and install scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.StorageExplicit = cmd.Flags().Changed("storage")

			stackConfig, err := config.NewStackConfig(opts)
			if err != nil {
				return err
			}

			result, err := generator.Generate(stackConfig, generator.Options{Force: force})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Generated observability stack at %s\n", result.OutputDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", opts.Name, "stack name")
	cmd.Flags().StringVar(&opts.Namespace, "namespace", opts.Namespace, "Kubernetes namespace")
	cmd.Flags().StringVar(&opts.OutputDir, "output", opts.OutputDir, "output directory")
	cmd.Flags().StringVar(&opts.Storage, "storage", opts.Storage, "object storage type: minio or s3")
	cmd.Flags().StringVar(&opts.Profile, "profile", opts.Profile, "profile: dev or production")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite generated files if the output directory exists")

	return cmd
}

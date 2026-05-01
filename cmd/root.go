package cmd

import "github.com/spf13/cobra"

func NewRootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "obsgen",
		Short:         "Generate starter Kubernetes observability stack configuration files",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newInitCommand())
	root.AddCommand(newVersionCommand(version))

	return root
}

func Execute(version string) error {
	return NewRootCommand(version).Execute()
}

// Package cli wires together the Cobra command tree exposed by the goforge binary.
package cli

import "github.com/spf13/cobra"

// version is injected at build time via -ldflags.
var version = "dev"

// Execute runs the root command and returns any execution error.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "goforge",
		Short:         "Scaffold a production-ready Go project",
		Long:          "goforge generates an opinionated Go project layout with the integrations you actually want.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newNewCmd(), newVersionCmd(), newListCmd())
	return root
}

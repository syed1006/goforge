package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new [project]",
		Short: "Create a new Go project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented yet — wired up in a later commit")
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available features (frameworks, databases, …)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented yet — wired up in a later commit")
		},
	}
}

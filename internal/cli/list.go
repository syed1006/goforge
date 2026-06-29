package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/syed1006/goforge/internal/config"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available features",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "frameworks",
			Short: "List supported HTTP frameworks",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return printList(cmd, "frameworks", asStrings(config.Frameworks))
			},
		},
		&cobra.Command{
			Use:   "databases",
			Short: "List supported database drivers",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return printList(cmd, "databases", asStrings(config.Databases))
			},
		},
	)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if err := printList(cmd, "frameworks", asStrings(config.Frameworks)); err != nil {
			return err
		}
		return printList(cmd, "databases", asStrings(config.Databases))
	}
	return cmd
}

func printList(cmd *cobra.Command, header string, items []string) error {
	out := cmd.OutOrStdout()
	if _, err := fmt.Fprintf(out, "%s:\n  - %s\n", header, strings.Join(items, "\n  - ")); err != nil {
		return err
	}
	return nil
}

func asStrings[T ~string](in []T) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return out
}

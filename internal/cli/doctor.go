package cli

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Check describes a single external tool goforge or generated projects can use.
type Check struct {
	Tool     string
	Required bool
	Why      string
	Install  string
	// Probe returns the discovered version line, or "" if not found.
	Probe func() string
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check that the external tools goforge expects are installed",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctor(cmd.OutOrStdout())
		},
	}
}

func runDoctor(out io.Writer) error {
	checks := defaultChecks()
	missingRequired := false
	for _, c := range checks {
		ver := c.Probe()
		switch {
		case ver != "":
			fmt.Fprintf(out, "  ok   %-15s %s\n", c.Tool, oneline(ver))
		case c.Required:
			missingRequired = true
			fmt.Fprintf(out, "  miss %-15s required — %s\n       install: %s\n", c.Tool, c.Why, c.Install)
		default:
			fmt.Fprintf(out, "  warn %-15s optional — %s\n       install: %s\n", c.Tool, c.Why, c.Install)
		}
	}
	if missingRequired {
		return fmt.Errorf("required tools missing")
	}
	return nil
}

func defaultChecks() []Check {
	return []Check{
		{
			Tool: "go", Required: true,
			Why:     "needed to run go get / go mod tidy in generated projects",
			Install: "https://go.dev/dl",
			Probe:   func() string { return firstLine("go", "version") },
		},
		{
			Tool: "gofmt", Required: true,
			Why:     "post-step formatter for generated code",
			Install: "ships with Go",
			Probe:   func() string { return whichVersion("gofmt") },
		},
		{
			Tool: "git", Required: false,
			Why:     "Makefile uses `git describe` for VERSION",
			Install: "https://git-scm.com/downloads",
			Probe:   func() string { return firstLine("git", "--version") },
		},
		{
			Tool: "air", Required: false,
			Why:     "hot-reload dev server (--hot-reload)",
			Install: "go install github.com/air-verse/air@latest",
			Probe:   func() string { return firstLine("air", "-v") },
		},
		{
			Tool: "golangci-lint", Required: false,
			Why:     "lint runner (--lint)",
			Install: "https://golangci-lint.run/welcome/install/",
			Probe:   func() string { return firstLine("golangci-lint", "--version") },
		},
		{
			Tool: "buf", Required: false,
			Why:     "regenerate protobuf bindings (--grpc)",
			Install: "https://buf.build/docs/installation",
			Probe:   func() string { return firstLine("buf", "--version") },
		},
		{
			Tool: "docker", Required: false,
			Why:     "build/run the generated image (--docker)",
			Install: "https://docs.docker.com/get-docker/",
			Probe:   func() string { return firstLine("docker", "--version") },
		},
		{
			Tool: "pre-commit", Required: false,
			Why:     "run hooks from .pre-commit-config.yaml (--lint)",
			Install: "pipx install pre-commit",
			Probe:   func() string { return firstLine("pre-commit", "--version") },
		},
	}
}

func firstLine(name string, args ...string) string {
	if _, err := exec.LookPath(name); err != nil {
		return ""
	}
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return ""
	}
	return oneline(string(out))
}

func whichVersion(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

func oneline(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}
	return s
}

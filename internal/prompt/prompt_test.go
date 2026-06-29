package prompt

import (
	"strings"
	"testing"

	"github.com/syed1006/goforge/internal/config"
)

func TestSummaryIncludesAllFields(t *testing.T) {
	t.Parallel()
	cfg := config.Config{
		ProjectName: "myapi",
		ModulePath:  "github.com/me/myapi",
		GoVersion:   "1.23",
		Framework:   config.FrameworkGin,
		Database:    config.DatabasePostgres,
		GRPC:        true,
		GraphQL:     false,
		HotReload:   true,
		Lint:        true,
		Docker:      true,
		CI:          true,
	}
	s := Summary(cfg)
	for _, want := range []string{"myapi", "github.com/me/myapi", "1.23", "gin", "postgres"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary missing %q:\n%s", want, s)
		}
	}
}

func TestValidators(t *testing.T) {
	t.Parallel()
	if err := validateProjectName("MyAPI"); err == nil {
		t.Error("expected validateProjectName to reject mixed case")
	}
	if err := validateProjectName("myapi"); err != nil {
		t.Errorf("expected validateProjectName to accept myapi: %v", err)
	}
	if err := validateModulePath("notamodule"); err == nil {
		t.Error("expected validateModulePath to reject bare token")
	}
	if err := validateModulePath("github.com/me/myapi"); err != nil {
		t.Errorf("expected validateModulePath to accept github.com path: %v", err)
	}
}

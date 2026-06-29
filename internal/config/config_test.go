package config

import "testing"

func TestValidate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		cfg  Config
		ok   bool
	}{
		{"happy", Config{ProjectName: "myapi", ModulePath: "github.com/me/myapi", GoVersion: "1.23", Framework: FrameworkGin, Database: DatabasePostgres}, true},
		{"bad name", Config{ProjectName: "MyAPI", ModulePath: "github.com/me/myapi", GoVersion: "1.23", Framework: FrameworkGin, Database: DatabaseNone}, false},
		{"bad module", Config{ProjectName: "myapi", ModulePath: "notamodule", GoVersion: "1.23", Framework: FrameworkGin, Database: DatabaseNone}, false},
		{"bad framework", Config{ProjectName: "myapi", ModulePath: "github.com/me/myapi", GoVersion: "1.23", Framework: "spring", Database: DatabaseNone}, false},
		{"bad db", Config{ProjectName: "myapi", ModulePath: "github.com/me/myapi", GoVersion: "1.23", Framework: FrameworkGin, Database: "oracle"}, false},
		{"empty go", Config{ProjectName: "myapi", ModulePath: "github.com/me/myapi", Framework: FrameworkGin, Database: DatabaseNone}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if (err == nil) != tc.ok {
				t.Fatalf("Validate(): got err=%v, want ok=%v", err, tc.ok)
			}
		})
	}
}

func TestSlug(t *testing.T) {
	t.Parallel()
	c := Config{ProjectName: "my_api"}
	if got, want := c.Slug(), "my-api"; got != want {
		t.Fatalf("Slug(): got %q, want %q", got, want)
	}
}

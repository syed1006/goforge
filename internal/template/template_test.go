package template

import (
	"testing"
	"testing/fstest"
)

func TestFuncs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, in, out string
		fn            func(string) string
	}{
		{"pascal kebab", "my-api-service", "MyApiService", pascalCase},
		{"pascal snake", "my_api_service", "MyApiService", pascalCase},
		{"pascal already", "MyApiService", "MyApiService", pascalCase},
		{"camel", "my-api", "myApi", camelCase},
		{"kebab from pascal", "MyApiService", "my-api-service", kebabCase},
		{"snake", "MyApiService", "my_api_service", snakeCase},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fn(tc.in); got != tc.out {
				t.Fatalf("got %q, want %q", got, tc.out)
			}
		})
	}
}

func TestEngineRender(t *testing.T) {
	t.Parallel()
	files := fstest.MapFS{
		"greet.tmpl":          {Data: []byte("hello {{ .Name | pascal }}!")},
		"nested/inner.go.tmpl": {Data: []byte("// pkg {{ kebab .Pkg }}")},
		"notatmpl":            {Data: []byte("ignored")},
	}
	eng, err := New(files)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	got, err := eng.Render("greet", map[string]string{"Name": "my-app"})
	if err != nil {
		t.Fatalf("Render greet: %v", err)
	}
	if string(got) != "hello MyApp!" {
		t.Fatalf("greet: %q", got)
	}

	got, err = eng.Render("nested/inner.go", map[string]string{"Pkg": "FooBar"})
	if err != nil {
		t.Fatalf("Render nested: %v", err)
	}
	if string(got) != "// pkg foo-bar" {
		t.Fatalf("nested: %q", got)
	}

	if _, err := eng.Render("missing", nil); err == nil {
		t.Fatal("expected error for missing template")
	}
}

package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestOneline(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"":                      "",
		"single":                "single",
		"   spaced   ":          "spaced",
		"first\nsecond":         "first",
		"only first line\r\nb":  "only first line",
	}
	for in, want := range cases {
		if got := oneline(in); got != want {
			t.Errorf("oneline(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRunDoctorIncludesEveryCheck(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	// runDoctor returns an error only when a required tool is missing. We don't
	// assert nil here because the test machine might not have `go` on PATH; we
	// just verify the listing covers every check.
	_ = runDoctor(&buf)
	for _, c := range defaultChecks() {
		if !strings.Contains(buf.String(), c.Tool) {
			t.Errorf("doctor output missing %q:\n%s", c.Tool, buf.String())
		}
	}
}

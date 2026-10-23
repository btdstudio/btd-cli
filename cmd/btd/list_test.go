package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

func TestRunListOutput(t *testing.T) {
	old := resolveInstalled
	resolveInstalled = func(name string, tool manifest.Tool) string {
		switch name {
		case "kumi":
			return "v0.1.0"
		case "hachigo":
			return "v0.0.9"
		case "tomo":
			return ""
		}
		return ""
	}
	defer func() { resolveInstalled = old }()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var buf bytes.Buffer
	if err := runList(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	for _, want := range []string{"kumi", "hachigo", "tomo", "ok", "mismatch", "missing", "v0.0.9"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

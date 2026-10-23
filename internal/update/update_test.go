package update

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

type fakeRunner struct {
	calls []string
}

func (f *fakeRunner) Run(dir, name string, args ...string) error {
	f.calls = append(f.calls, name+" "+strings.Join(args, " "))
	return nil
}

func TestPlanFindsChanges(t *testing.T) {
	m := manifest.Default()
	m.Tools["goph"] = manifest.Tool{Kind: manifest.KindGo, Path: "github.com/btdstudio/goph/cmd/goph", Version: "v0.0.1"}

	changes := Plan(m)
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1: %#v", len(changes), changes)
	}
	c := changes[0]
	if c.Tool != "goph" || c.From != "v0.0.1" || c.To != "v0.1.0" {
		t.Fatalf("change = %#v", c)
	}
}

func TestPlanNoChanges(t *testing.T) {
	if changes := Plan(manifest.Default()); len(changes) != 0 {
		t.Fatalf("want no changes, got %#v", changes)
	}
}

func TestRunCheckIsDryRun(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	gpath := filepath.Join(xdg, "btd", "manifest.json")
	if err := manifest.WriteFile(gpath, &manifest.Manifest{
		Tools: map[string]manifest.Tool{
			"goph": {Kind: manifest.KindGo, Path: "github.com/btdstudio/goph/cmd/goph", Version: "v0.0.1"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	m, err := manifest.Load("")
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{}
	changes, err := Run(r, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}
	if len(r.calls) != 0 {
		t.Fatalf("dry run should not install, got %#v", r.calls)
	}
	global, _ := manifest.ReadFile(gpath)
	if global.Tools["goph"].Version != "v0.0.1" {
		t.Fatalf("dry run mutated manifest: %+v", global.Tools["goph"])
	}
}

func TestRunAppliesAndReinstalls(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	gpath := filepath.Join(xdg, "btd", "manifest.json")
	if err := manifest.WriteFile(gpath, &manifest.Manifest{
		Tools: map[string]manifest.Tool{
			"goph": {Kind: manifest.KindGo, Path: "github.com/btdstudio/goph/cmd/goph", Version: "v0.0.1"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	m, err := manifest.Load("")
	if err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{}
	changes, err := Run(r, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}
	if len(r.calls) != 1 || r.calls[0] != "go install github.com/btdstudio/goph/cmd/goph@v0.1.0" {
		t.Fatalf("calls = %#v", r.calls)
	}
	global, _ := manifest.ReadFile(gpath)
	if global.Tools["goph"].Version != "v0.1.0" {
		t.Fatalf("manifest not updated: %+v", global.Tools["goph"])
	}
	if global.UpdatedAt == "" {
		t.Fatal("updatedAt not set")
	}
}

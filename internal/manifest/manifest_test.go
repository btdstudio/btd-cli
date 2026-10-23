package manifest

import (
	"path/filepath"
	"testing"
)

func TestMergePrecedence(t *testing.T) {
	global := &Manifest{
		Tools: map[string]Tool{
			"kumi": {Kind: KindGo, Path: "github.com/btdstudio/kumi/cmd/kumi", Version: "v0.2.0"},
		},
	}
	project := &Manifest{
		Tools: map[string]Tool{
			"kumi": {Version: "v0.3.0"},
		},
	}

	eff := Merge(Default(), global, project)

	if eff.SchemaVersion != 1 {
		t.Fatalf("schemaVersion = %d, want 1", eff.SchemaVersion)
	}
	if eff.Compatible != "btd-0.1.0" {
		t.Fatalf("compatible = %q, want btd-0.1.0", eff.Compatible)
	}

	// Project version wins, but inherits kind/path from earlier layers.
	kumi := eff.Tools["kumi"]
	if kumi.Version != "v0.3.0" {
		t.Fatalf("kumi.version = %q, want v0.3.0", kumi.Version)
	}
	if kumi.Kind != KindGo || kumi.Path != "github.com/btdstudio/kumi/cmd/kumi" {
		t.Fatalf("kumi kind/path not inherited: %+v", kumi)
	}

	// Untouched tools inherit defaults.
	if eff.Tools["tomo"].Ref != "v0.1.0" {
		t.Fatalf("tomo.ref = %q, want v0.1.0", eff.Tools["tomo"].Ref)
	}
	if eff.Tools["hachigo"].Kind != KindGo {
		t.Fatalf("hachigo.kind = %q, want %q", eff.Tools["hachigo"].Kind, KindGo)
	}
}

func TestMergeDoesNotMutateLayers(t *testing.T) {
	def := Default()
	proj := &Manifest{Tools: map[string]Tool{"kumi": {Version: "v9.9.9"}}}
	eff := Merge(def, proj)

	if eff.Tools["kumi"].Version != "v9.9.9" {
		t.Fatalf("kumi.version = %q, want v9.9.9", eff.Tools["kumi"].Version)
	}
	if def.Tools["kumi"].Version != "v0.1.0" {
		t.Fatalf("default mutated: kumi.version = %q", def.Tools["kumi"].Version)
	}
}

func TestReadFileMissingReturnsEmpty(t *testing.T) {
	m, err := ReadFile(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatal(err)
	}
	if m == nil || len(m.Tools) != 0 {
		t.Fatalf("want empty manifest, got %+v", m)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	m := Default()
	m.UpdatedAt = "2026-10-21T00:00:00Z"

	if err := WriteFile(path, m); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tools["hachigo"].Kind != KindGo {
		t.Fatalf("hachigo.kind = %q, want %q", got.Tools["hachigo"].Kind, KindGo)
	}
	if got.UpdatedAt != m.UpdatedAt {
		t.Fatalf("updatedAt = %q, want %q", got.UpdatedAt, m.UpdatedAt)
	}
}

func TestLoadMergesGlobalAndProject(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	if err := WriteFile(filepath.Join(xdg, "btd", "manifest.json"), &Manifest{
		Tools: map[string]Tool{"kumi": {Version: "v0.2.0"}},
	}); err != nil {
		t.Fatal(err)
	}

	proj := t.TempDir()
	if err := WriteFile(ProjectPath(proj), &Manifest{
		Tools: map[string]Tool{"kumi": {Version: "v0.3.0"}},
	}); err != nil {
		t.Fatal(err)
	}

	eff, err := Load(proj)
	if err != nil {
		t.Fatal(err)
	}
	if eff.Tools["kumi"].Version != "v0.3.0" {
		t.Fatalf("kumi.version = %q, want v0.3.0 (project wins)", eff.Tools["kumi"].Version)
	}
	if eff.Tools["tomo"].Ref != "v0.1.0" {
		t.Fatalf("tomo.ref = %q, want v0.1.0 (default inherited)", eff.Tools["tomo"].Ref)
	}
}

func TestGlobalPathHonorsXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/xdg")
	const want = "/custom/xdg/btd/manifest.json"
	if got := GlobalPath(); got != want {
		t.Fatalf("GlobalPath = %q, want %q", got, want)
	}
}

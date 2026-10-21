package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWalksUp(t *testing.T) {
	root := t.TempDir()
	proj := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(filepath.Join(proj, MarkerDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, MarkerDir, MarkerFile), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(proj, "x", "y")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	got, ok := Find(nested)
	if !ok {
		t.Fatal("Find returned ok=false, want true")
	}
	if got != proj {
		t.Fatalf("Find = %q, want %q", got, proj)
	}
}

func TestFindStopsAtNearestMarker(t *testing.T) {
	root := t.TempDir()
	outer := filepath.Join(root, "outer")
	inner := filepath.Join(outer, "inner")
	for _, d := range []string{outer, inner} {
		if err := os.MkdirAll(filepath.Join(d, MarkerDir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, MarkerDir, MarkerFile), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, ok := Find(inner)
	if !ok {
		t.Fatal("Find returned ok=false")
	}
	if got != inner {
		t.Fatalf("Find = %q, want nearest marker %q", got, inner)
	}
}

func TestFindNone(t *testing.T) {
	if _, ok := Find(t.TempDir()); ok {
		t.Fatal("Find returned ok=true with no marker present")
	}
}

func TestFindEmptyUsesCwd(t *testing.T) {
	proj := filepath.Join(t.TempDir(), "proj")
	if err := os.MkdirAll(filepath.Join(proj, MarkerDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(proj, MarkerDir, MarkerFile), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(proj); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(old) }()

	got, ok := Find("")
	want, _ := os.Getwd() // canonical (macOS: /private/var/... vs /var/...)
	if !ok || got != want {
		t.Fatalf("Find(\"\") = %q, %v; want %q, true", got, ok, want)
	}
}

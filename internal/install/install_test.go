package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

type recordedCall struct {
	dir  string
	name string
	args []string
}

type fakeRunner struct {
	calls []recordedCall
}

func (f *fakeRunner) Run(dir, name string, args ...string) error {
	f.calls = append(f.calls, recordedCall{dir: dir, name: name, args: append([]string(nil), args...)})
	return nil
}

func assertCalls(t *testing.T, got []recordedCall, want []recordedCall) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d calls, want %d:\ngot:  %#v\nwant: %#v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i].dir != want[i].dir || got[i].name != want[i].name {
			t.Fatalf("call %d = {dir:%q name:%q}, want {dir:%q name:%q}", i, got[i].dir, got[i].name, want[i].dir, want[i].name)
		}
		if len(got[i].args) != len(want[i].args) {
			t.Fatalf("call %d args = %#v, want %#v", i, got[i].args, want[i].args)
		}
		for j := range want[i].args {
			if got[i].args[j] != want[i].args[j] {
				t.Fatalf("call %d arg %d = %q, want %q", i, j, got[i].args[j], want[i].args[j])
			}
		}
	}
}

func TestInstallGophRunsGoInstall(t *testing.T) {
	r := &fakeRunner{}
	if err := Install(r, manifest.Default(), "goph", Options{}); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "go", args: []string{"install", "github.com/btdstudio/goph/cmd/goph@v0.1.0"}},
	})
}

func TestInstallHachigoGlobalRunsGoInstall(t *testing.T) {
	r := &fakeRunner{}
	if err := Install(r, manifest.Default(), "hachigo", Options{}); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "go", args: []string{"install", "github.com/btdstudio/hachigo/cmd/hachigo@v0.1.0"}},
	})
}

func TestInstallHachigoProjectRunsInit(t *testing.T) {
	old := lookPath
	lookPath = func(string) (string, error) { return "/usr/bin/hachigo", nil }
	defer func() { lookPath = old }()

	r := &fakeRunner{}
	o := Options{ProjectRoot: "/proj"}
	if err := Install(r, manifest.Default(), "hachigo", o); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "hachigo", args: []string{"init", filepath.Join("/proj", "site")}},
	})
}

func TestInstallHachigoProjectEnsuresBinary(t *testing.T) {
	old := lookPath
	lookPath = func(string) (string, error) { return "", os.ErrNotExist }
	defer func() { lookPath = old }()

	r := &fakeRunner{}
	o := Options{ProjectRoot: "/proj", SiteDir: "web"}
	if err := Install(r, manifest.Default(), "hachigo", o); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "go", args: []string{"install", "github.com/btdstudio/hachigo/cmd/hachigo@v0.1.0"}},
		{name: "hachigo", args: []string{"init", filepath.Join("/proj", "web")}},
	})
}

func TestInstallTomoGlobalClonesAndInstalls(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "tomo")
	r := &fakeRunner{}
	o := Options{CacheDir: cache}
	if err := Install(r, manifest.Default(), "tomo", o); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "git", args: []string{"clone", "--depth", "1", "--branch", "v0.1.0", "https://github.com/btdstudio/tomo.git", cache}},
		{name: "node", args: []string{filepath.Join(cache, "scripts", "install.js"), "--global"}},
	})
}

func TestInstallTomoProjectInstallsPerProject(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "tomo")
	r := &fakeRunner{}
	o := Options{CacheDir: cache, ProjectRoot: "/proj", Provider: "opencode"}
	if err := Install(r, manifest.Default(), "tomo", o); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "git", args: []string{"clone", "--depth", "1", "--branch", "v0.1.0", "https://github.com/btdstudio/tomo.git", cache}},
		{name: "node", args: []string{filepath.Join(cache, "scripts", "install.js"), "--project", "/proj", "--provider", "opencode"}},
	})
}

func TestInstallTomoRefreshesExistingClone(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "tomo")
	if err := os.MkdirAll(filepath.Join(cache, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	r := &fakeRunner{}
	o := Options{CacheDir: cache}
	if err := Install(r, manifest.Default(), "tomo", o); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{dir: cache, name: "git", args: []string{"fetch", "--depth", "1", "origin", "v0.1.0"}},
		{dir: cache, name: "git", args: []string{"checkout", "--detach", "v0.1.0"}},
		{name: "node", args: []string{filepath.Join(cache, "scripts", "install.js"), "--global"}},
	})
}

func TestInstallAllRunsEveryTool(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "tomo")
	r := &fakeRunner{}
	if err := Install(r, manifest.Default(), "all", Options{CacheDir: cache}); err != nil {
		t.Fatal(err)
	}
	assertCalls(t, r.calls, []recordedCall{
		{name: "go", args: []string{"install", "github.com/btdstudio/goph/cmd/goph@v0.1.0"}},
		{name: "go", args: []string{"install", "github.com/btdstudio/hachigo/cmd/hachigo@v0.1.0"}},
		{name: "git", args: []string{"clone", "--depth", "1", "--branch", "v0.1.0", "https://github.com/btdstudio/tomo.git", cache}},
		{name: "node", args: []string{filepath.Join(cache, "scripts", "install.js"), "--global"}},
	})
}

func TestInstallUnknownToolErrors(t *testing.T) {
	r := &fakeRunner{}
	if err := Install(r, manifest.Default(), "nope", Options{}); err == nil {
		t.Fatal("want error for unknown tool, got nil")
	}
}

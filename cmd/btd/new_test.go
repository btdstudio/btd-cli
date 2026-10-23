package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
	newpkg "github.com/btdstudio/btd-cli/internal/new"
)

type recCall struct {
	dir  string
	name string
	args []string
}

type recRunner struct {
	calls []recCall
}

func (f *recRunner) Run(dir, name string, args ...string) error {
	f.calls = append(f.calls, recCall{dir: dir, name: name, args: append([]string(nil), args...)})
	return nil
}

func TestParseIncludes(t *testing.T) {
	tests := []struct {
		in   string
		want []string
		err  bool
	}{
		{"", nil, false},
		{"tomo", []string{"tomo"}, false},
		{"tomo,hachigo", []string{"tomo", "hachigo"}, false},
		{" tomo , hachigo ", []string{"tomo", "hachigo"}, false},
		{"nope", nil, true},
		{"tomo,nope", nil, true},
	}
	for _, tc := range tests {
		got, err := parseIncludes(tc.in)
		if tc.err && err == nil {
			t.Errorf("parseIncludes(%q): want error, got nil", tc.in)
			continue
		}
		if !tc.err && err != nil {
			t.Errorf("parseIncludes(%q): %v", tc.in, err)
			continue
		}
		if strings.Join(got, ",") != strings.Join(tc.want, ",") {
			t.Errorf("parseIncludes(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestRunNewScaffoldsAndIncludes(t *testing.T) {
	old := install.LookPath
	install.LookPath = func(name string) (string, error) {
		if name == "goph" || name == "hachigo" {
			return "/usr/bin/" + name, nil
		}
		return "", os.ErrNotExist
	}
	defer func() { install.LookPath = old }()

	tmp := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	cache := filepath.Join(tmp, "cache", "tomo")
	r := &recRunner{}
	o := newpkg.Options{
		Template: "goph",
		Name:     "myapp",
		Include:  []string{"tomo", "hachigo"},
		Provider: "all",
		CacheDir: cache,
	}
	if err := newpkg.Run(r, manifest.Default(), o); err != nil {
		t.Fatal(err)
	}

	want := []recCall{
		{name: "goph", args: []string{"new", "myapp"}},
		{name: "git", args: []string{"clone", "--depth", "1", "--branch", "v0.1.0", "https://github.com/btdstudio/tomo.git", cache}},
		{name: "node", args: []string{filepath.Join(cache, "scripts", "install.js"), "--project", "myapp", "--provider", "all"}},
		{name: "hachigo", args: []string{"init", filepath.Join("myapp", "site")}},
	}
	if len(r.calls) != len(want) {
		t.Fatalf("got %d calls, want %d:\n%#v", len(r.calls), len(want), r.calls)
	}
	for i, w := range want {
		g := r.calls[i]
		if g.name != w.name || strings.Join(g.args, " ") != strings.Join(w.args, " ") {
			t.Fatalf("call %d = %q %v, want %q %v", i, g.name, g.args, w.name, w.args)
		}
	}

	// seeded manifest exists
	if _, err := os.Stat(filepath.Join(tmp, "myapp", ".btd", "manifest.json")); err != nil {
		t.Fatalf("seeded manifest missing: %v", err)
	}
	// .gitignore has .tomo/
	gi, err := os.ReadFile(filepath.Join(tmp, "myapp", ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(gi), ".tomo/") {
		t.Fatalf(".gitignore missing .tomo/: %q", string(gi))
	}
}

func TestRunNewInstallsGophWhenMissing(t *testing.T) {
	old := install.LookPath
	install.LookPath = func(string) (string, error) { return "", os.ErrNotExist }
	defer func() { install.LookPath = old }()

	tmp := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	r := &recRunner{}
	if err := newpkg.Run(r, manifest.Default(), newpkg.Options{Template: "goph", Name: "myapp"}); err != nil {
		t.Fatal(err)
	}
	if len(r.calls) != 2 {
		t.Fatalf("got %d calls: %#v", len(r.calls), r.calls)
	}
	if r.calls[0].name != "go" || r.calls[0].args[0] != "install" {
		t.Fatalf("call 0 = %#v, want go install", r.calls[0])
	}
	if r.calls[1].name != "goph" || r.calls[1].args[0] != "new" {
		t.Fatalf("call 1 = %#v, want goph new", r.calls[1])
	}
}

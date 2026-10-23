package new

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

// includeTomo installs tomo into the project and appends .tomo/ to the
// project's .gitignore.
func includeTomo(r install.Runner, m *manifest.Manifest, o Options) error {
	if err := install.Install(r, m, "tomo", install.Options{
		ProjectRoot: o.Name,
		Provider:    o.Provider,
		CacheDir:    o.CacheDir,
	}); err != nil {
		return err
	}
	return appendGitignore(o.Name)
}

// appendGitignore appends ".tomo/" to <root>/.gitignore, creating the file if
// needed and skipping when the entry is already present.
func appendGitignore(root string) error {
	path := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if strings.Contains(string(data), ".tomo/") {
		return nil
	}
	var b strings.Builder
	b.Write(data)
	if len(data) > 0 && data[len(data)-1] != '\n' {
		b.WriteByte('\n')
	}
	b.WriteString(".tomo/\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

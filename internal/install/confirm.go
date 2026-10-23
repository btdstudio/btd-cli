package install

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

// Confirm asks the user to confirm a global install of the named tools,
// printing the install targets. It returns true when the user answers yes.
func Confirm(in io.Reader, out io.Writer, m *manifest.Manifest, names []string) bool {
	goBin := goBinDir()
	fmt.Fprintln(out, "btd will install the following globally:")
	for _, name := range names {
		tool := m.Tools[name]
		switch tool.Kind {
		case manifest.KindGo:
			fmt.Fprintf(out, "  %-8s go install %s@%s  ->  %s\n",
				name, tool.Path, tool.Version, filepath.Join(goBin, name))
		case manifest.KindGit:
			fmt.Fprintf(out, "  %-8s git clone %s @ %s  ->  %s, then node scripts/install.js --global\n",
				name, tool.Repo, tool.Ref, CacheDir())
		}
	}
	fmt.Fprintf(out, "Global installs write to user-level directories (GOPATH/bin, harness config).\n")
	fmt.Fprintf(out, "Continue? [y/N] ")

	line, _ := bufio.NewReader(in).ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes"
}

// goBinDir returns the directory go install writes binaries to, best-effort.
func goBinDir() string {
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		return filepath.Join(strings.TrimSpace(string(out)), "bin")
	}
	if gp := os.Getenv("GOPATH"); gp != "" {
		return filepath.Join(gp, "bin")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "go", "bin")
}

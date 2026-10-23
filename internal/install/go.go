package install

import "github.com/btdstudio/btd-cli/internal/manifest"

// goStrategy installs a go-kind tool with `go install <path>@<version>`.
type goStrategy struct{}

func (goStrategy) Install(r Runner, tool manifest.Tool, _ Options) error {
	return r.Run("", "go", "install", tool.Path+"@"+tool.Version)
}

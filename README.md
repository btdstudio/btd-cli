# btd

`btd` is the orchestrator for the [btdstudio](https://github.com/btdstudio) toolchain. It is a thin
Go CLI (Cobra) that ties three sibling tools into one product:

| Tool | Purpose |
|------|---------|
| [kumi](https://github.com/btdstudio/kumi) | Scaffold Go + HTMX projects |
| [tomo](https://github.com/btdstudio/tomo) | Harness-agnostic AI SDD loop |
| [hachigo](https://github.com/btdstudio/hachigo) | Static site generator |

`btd` is deliberately **not** a reimplementation of any tool — it shells out to each tool's existing
CLI and install flow, holding no templates and no copied tool logic.

## Install

```bash
go install github.com/btdstudio/btd-cli/cmd/btd@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

Requires Go 1.24+. Until the first tagged release, build from source: `go build ./cmd/btd`.

## Usage

```
btd new <template> <name> [module-path] [--include tomo,hachigo]
btd install <kumi|tomo|hachigo|all> [--provider <ids>] [--yes]
btd update [--check]
btd doctor [--install]
btd list | ls
```

### `btd new`

Scaffold a project with kumi and optionally include tomo and/or hachigo:

```bash
btd new kumi my-app                          # kumi only
btd new kumi my-app --include tomo,hachigo   # + SDD loop + static site
```

`btd new` seeds a `.btd/manifest.json` (the project pin) and, with `--include tomo`, appends
`.tomo/` to the project's `.gitignore`. Inclusions install **per-project**, never globally.

### `btd install`

Installs the sibling tools, pinned by the manifest. Install location is **context-aware**:

| Context | `tomo` | `hachigo` | `kumi` |
|---------|--------|-----------|--------|
| Inside a btd project | per-project (`install.js --project <root>`) | `hachigo init <site-dir>` | global |
| Outside a btd project | global (clone + `install.js --global`) | global (`go install`) | global (`go install`) |

- **Inside a project** (a `.btd/manifest.json` is found up the tree), installs touch only the
  project's own directories and never prompt.
- **Outside a project**, installs are global, write to user-level directories (`GOPATH/bin`, harness
  config), and prompt for confirmation first — skip with `--yes`.

`--provider <ids>` selects tomo harnesses (`all | opencode | claude | copilot | antigravity | codex`)
without prompting; it defaults to `all`.

### `btd update`

Moves tool pins to the refs curated for the installed btd release, then re-installs anything that
changed. `--check` prints what would change without applying. Pins never move blindly to `@latest`.

### `btd doctor`

Checks `go`, `git`, `node` (≥20), `kumi`, `hachigo`, and kumi's prerequisites (`templ`, `air`,
`goose`). `--install` installs what it can and delegates the rest to `kumi doctor --install`.

### `btd list`

Prints installed tools, their resolved versions, and whether they match the manifest pins.

## Manifest

Tool pins live in `~/.config/btd/manifest.json` (global) and `.btd/manifest.json` (project), merged
project-wins. The project manifest doubles as the "btd project" marker that `install` and `update`
use to detect context.

## Development

```bash
go build ./...
go vet ./...
go test ./...
```

## License

[MIT](LICENSE)

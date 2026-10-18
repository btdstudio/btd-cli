# btd

`btd` is the orchestrator for the [btdstudio](https://github.com/btdstudio) toolchain. It is a thin
Go CLI (Cobra) that ties three sibling tools together into one product:

| Tool | Purpose | Repo |
|------|---------|------|
| [goph](https://github.com/btdstudio/goph) | Scaffold Go + HTMX projects | `github.com/btdstudio/goph` |
| [tomo](https://github.com/btdstudio/tomo) | Harness-agnostic AI SDD loop | `github.com/btdstudio/tomo` |
| [hachigo](https://github.com/btdstudio/hachigo) | Static site generator | `github.com/btdstudio/hachigo` |

`btd` is deliberately **not** a reimplementation of any tool — it shells out to their existing CLIs
and install flows, holding no templates and no copied tool logic.

## Status

Pre-1.0 and under active development. The CLI is not yet published; install instructions will land
with the first release. Progress is tracked in the repository's SDD workspace (`.tomo/specs/`).

## Planned command surface

- `btd new <template> <name>` — scaffold a project (optionally `--include tomo,hachigo`)
- `btd install <goph|tomo|hachigo|all>` — context-aware install of the sibling tools
- `btd update [--check]` — keep the toolchain mutually compatible via a pinned manifest
- `btd doctor [--install]` — check and repair required tooling
- `btd list` — list installed tools and their versions

## License

[MIT](LICENSE)

# Security Policy

## Supported versions

btd is pre-1.0. Security fixes are applied to the latest commit on `main`;
re-run `go install github.com/btdstudio/btd-cli/cmd/btd@latest` to update.

## Reporting a vulnerability

Please do not open a public issue for security problems.

Report privately through GitHub: go to the
[Security tab](https://github.com/btdstudio/btd-cli/security/advisories/new) and
choose **Report a vulnerability**.

Include what you found, how to reproduce it, and the affected version or
commit. This is a small project maintained in spare time; you can expect an
acknowledgement within a few days. Once a fix is available, the advisory will
be published and you will be credited unless you prefer otherwise.

## Scope

In scope: the `btd` orchestrator CLI, its manifest handling, and its install /
update / new flows.

Out of scope: the tools btd wraps (goph, tomo, hachigo) — report issues in
those tools to their own repositories — and vulnerabilities in third-party
dependencies with no impact on btd.

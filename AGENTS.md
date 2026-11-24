# Agent instructions

## Hard rule: never modify `src/`

The `src/` directory is the fixed Constellation framework. It is already set and is not part of this application's editable code.

- You MAY read anything under `src/` to understand the framework API.
- You MUST NOT create, edit, delete, rename, move, format, or generate any file under `src/`.
- This applies to all tools, shell commands, code generators, formatters, and build steps.
- Do all work outside `src/`: `main.go`, `site/`, `static/`, `data/`, documentation, and root-level configuration files.
- If a requested change appears to require a modification inside `src/`, stop and tell the user instead of changing it.

## Project notes

- `main.go` embeds `static/`, registers the endpoints from `site/`, and starts the framework.
- `site/` contains the novel site: pages, API handlers, JSON store, Markdown renderer, and shared components.
- `data/novels.json` is the local story library. Set `CLIO_DATA_FILE` to use another path.
- Prefer the standard library; this project has no third-party Go dependencies.

## Verification

Run these from the repository root:

```bash
gofmt -l .
go vet ./...
go build ./...
go test ./...
```

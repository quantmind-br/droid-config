# AGENTS.md

## Build/Run Commands
- `make build` - Compile binary to build/droid-config
- `make run` - Build and execute
- `make install` - Install to ~/.local/bin
- `make clean` - Remove build artifacts
- `make test` - Run tests and vet

## Code Style Guidelines
- **Go version**: 1.21+
- **Imports**: Standard library first, then third-party, then local
- **Formatting**: gofmt standard
- **Naming**: camelCase for unexported, PascalCase for exported
- **Error handling**: Return errors, handle at call site
- **TUI framework**: Bubble Tea (TEA pattern) with Lip Gloss styling
- **Config path**: `~/.factory/config.json` (JSON format with indent=2)

## Project Structure
```
cmd/droid-config/main.go      # Entry point
internal/config/              # Config types and file I/O
internal/ui/                  # Bubble Tea model, view, styles
internal/ui/components/       # Reusable UI components
```

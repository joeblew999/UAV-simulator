# Repository Guidelines

## Code Style

- `go fmt ./...` before committing
- Package names lowercase; exported `UpperCamelCase`, unexported `lowerCamelCase`
- Filenames `snake_case.go`; tests `*_test.go`
- Keep functions small; separate rendering from physics

## Testing

- Standard `testing` package; table-driven tests encouraged
- Prioritize math/physics tests (run headless)
- Avoid OpenGL in tests; keep CI-friendly

## Commits

- Conventional Commits: `feat:`, `fix:`, `refactor:`, `chore:`
- Imperative mood, concise scope
- Don't commit binaries or local caches

## Maintainer

Maintained by an experienced game/simulation developer with physics degree. Complex dynamics and control models welcome.

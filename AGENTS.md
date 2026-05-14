# Repository Guidelines

## Project Structure & Module Organization
This repository is a small standalone Go utility library extracted from a larger project. `logger/` contains Zap-based file logging, basic log helpers, and Gin request/response middleware. `jwt/` contains JWT signing, parsing, validation, and refresh helpers built around a local service type. `config/` contains shared configuration structs plus validation helpers for both logging and JWT usage. `README.md` is minimal today, so package comments and exported identifiers should stay clear enough to explain intended usage.

## Build, Test, and Development Commands
Use standard Go tooling from the repository root:

- `go build ./...` builds all packages and catches compile-time issues.
- `go test ./...` runs the full test suite; add tests before changing shared behavior.
- `go fmt ./...` applies canonical Go formatting.
- `go vet ./...` checks for common correctness issues.

Run `go mod tidy` after adding or removing dependencies to keep `go.mod` consistent.

## Coding Style & Naming Conventions
Follow idiomatic Go conventions: tabs for indentation, mixedCaps for exported names, and short lowercase package names (`config`, `logger`, `jwt`). Keep files focused on one concern: config validation, logger setup, middleware behavior, or JWT lifecycle handling. Prefer small helpers over large procedural flows. Use `gofmt` as the source of truth for formatting; do not hand-align code. When adding comments, keep them brief and start exported comments with the identifier name.

Avoid introducing cross-project assumptions. Code in this repository should not depend on application-global config variables or foreign module paths copied from source projects; library packages should be wired through local config structs and explicit constructors instead.

## Testing Guidelines
Add coverage alongside the package you change. Name tests `TestXxx` and keep them in the same package unless black-box testing is more useful. Favor table-driven tests for validators, logger initialization, rotation settings, Gin middleware behavior, and JWT token lifecycle cases such as parse, expiry, and refresh decisions. Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines
Git history is currently minimal (`Initial commit`), so use concise imperative commit messages such as `Add logger config defaults` or `Fix Gin response capture`. Keep each commit scoped to one logical change. PRs should include a short summary, the reason for the change, any config or dependency impact, and sample logs or request traces when middleware behavior changes.

## Configuration Notes
Do not commit secrets, JWT signing keys, or environment-specific log paths. Keep example YAML values sanitized, and document any new logging or JWT fields in code comments or the README when they affect library consumers integrating this package into services.

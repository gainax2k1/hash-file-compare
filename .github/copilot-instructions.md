# Guidance for AI coding agents

Purpose: Help contributors and AI agents quickly understand and modify this small Go CLI project that computes SHA-256 file hashes and discovers duplicates.

- **Big picture:** This repo is a tiny Go module (module path: `github.com/gainax2k1/hash-file-compare`). It contains a single CLI entrypoint (`main.go`) and two focused components:
  - `hashFile` (package): exports `HashFile(filename string) (string, error)` which computes a SHA-256 hex string for a file. See [hashFile/hashFile.go](hashFile/hashFile.go).
  - `walkDir` (intended component): walks directories, hashes files and groups them by hash. The current file is at [walkDir/walkDir.go](walkDir/walkDir.go) but contains implementation issues (package name vs `main`, and incorrect call to the hashing function). Review before running.

- **Why this structure:** The project separates hashing logic into a reusable package (`hashFile`) so CLI code can focus on argument parsing and presentation.

- **Build & run (explicit):**
  - Build the main CLI: `go build -o bin/hash-file-compare .`
  - Run single-file hash (from repo root): `go run . path/to/file`
  - Intended directory walker (requires walkDir to be either `package main` or expose a `WalkDir` function): examples below.

- **Common tasks and examples for edits:**
  - To compute a file hash from another package, call the exported function: `hashfile.HashFile(path)` using the import path `github.com/gainax2k1/hash-file-compare/hashFile` (the code currently aliases it as `hashfile`). Example from `main.go`.
  - If you implement a directory-walking utility, prefer this signature and behavior: `func WalkDir(dir string) (map[string][]string, error)` which returns a map of hash->list-of-files. Then call it from `main.go`.
  - Fixes commonly needed in `walkDir/walkDir.go`:
    - Change `package walkdir` to `package main` if the file should be a separate CLI command, or keep `package walkdir` and export `WalkDir`.
    - Replace the call `hashFile(path)` with `hashfile.HashFile(path)` (exported function).

- **Project-specific conventions & notes:**
  - Module import path is defined in `go.mod`. Use that exact import when referencing local packages.
  - The project keeps package folders at top-level (not under `cmd/`). Small size — prefer single binary approach unless adding subcommands.
  - There are currently no tests; add targeted unit tests for `hashFile.HashFile` and for a `WalkDir` function if you extract one.

- **Integration points & external dependencies:**
  - No external services. Uses standard library crypto packages (`crypto/sha256`, `encoding/hex`) — see `hashFile/hashFile.go`.

- **What to do first (practical steps for an agent):**
  1. Read `main.go`, `hashFile/hashFile.go`, and `walkDir/walkDir.go` to confirm intent.
 2. If tasked to make directory-walking runnable: either convert `walkDir/walkDir.go` to `package main` (then `go run ./walkDir <dir>`) or refactor to expose `WalkDir` and call it from `main.go` as a `-d` flag handler.
 3. Run `go run . <file>` to verify the single-file hashing path works before changing walkDir behavior.

- **Edge cases and safety checks to preserve:**
  - Preserve the error-returning behavior of `HashFile` and propagate errors to `main` for consistent CLI exit codes.
  - When walking directories, skip directories and continue on individual file errors (existing code logs and continues).

If anything here is unclear or you want the instructions tailored (for example: add examples of preferred CLI flags, or convert `walkDir` into an exported package), tell me which direction and I will iterate.

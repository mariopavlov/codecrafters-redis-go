<!--
  Sync Impact Report
  ==================================================
  Version change: N/A → 1.0.0 (initial ratification)
  Modified principles: N/A (initial creation)
  Added sections:
    - Core Principles (5 principles)
    - Development Workflow
    - Quality Standards
    - Governance
  Removed sections: N/A
  Templates requiring updates:
    - specs/templates/plan-template.md ✅ written
    - specs/templates/spec-template.md ✅ written
    - specs/templates/tasks-template.md ✅ written
    - specs/templates/spec-checklist.md ✅ written
  Follow-up TODOs: None
  ==================================================
-->

# CodeCrafters Redis Go Constitution

## Core Principles

### I. RESP Protocol Correctness
All Redis command handling MUST correctly implement the RESP
(REdis Serialization Protocol) specification. Parsing and
serialization MUST handle all RESP data types: Simple Strings,
Errors, Integers, Bulk Strings, and Arrays. Edge cases such as
null bulk strings, empty arrays, and multi-byte characters MUST
be handled. Protocol correctness is verified by CodeCrafters
stage tests and local integration tests.

### II. Concurrent Connection Handling
The server MUST handle multiple simultaneous client connections
without blocking or data corruption. Each connection MUST be
served in its own goroutine. Shared state (e.g., key-value
store) MUST be protected with proper synchronization primitives
(sync.Mutex, sync.RWMutex, or channels). No data races are
permitted — `go test -race` MUST pass.

### III. Standard Library Only
All functionality MUST be implemented using Go's standard
library exclusively. No third-party dependencies are permitted.
This constraint ensures deep understanding of networking, data
structures, and concurrency primitives. The `go.mod` file MUST
NOT contain `require` directives for external modules.

### IV. Incremental Stage Delivery
Each CodeCrafters stage MUST be implemented as an atomic,
working increment. A stage MUST NOT break previously passing
stages. Implementation order MUST follow the CodeCrafters
challenge progression. Each stage commit MUST be independently
submittable via `git push origin master`.

### V. Simplicity and Clarity
Code MUST be straightforward and idiomatic Go. Prefer explicit
logic over clever abstractions. Functions MUST have a single
clear responsibility. Package-level organization MUST remain
flat until complexity genuinely demands separation. YAGNI: do
not implement features ahead of the current stage requirement.

## Development Workflow

- Entry point is `app/main.go`; new files go under `app/`
- Run the server locally via `./your_program.sh`
- Test against CodeCrafters by pushing: `git push origin master`
- Local unit/integration tests use `go test ./...`
- Race detection: `go test -race ./...`
- Formatting: `gofmt -s -w .` before every commit
- Commits follow conventional format: `impl:`, `fix:`, `refactor:`, `test:`

## Quality Standards

- All code MUST compile with zero warnings under `go vet ./...`
- All code MUST pass `gofmt` formatting checks
- Functions exceeding 40 lines SHOULD be refactored unless
  the logic is inherently sequential and splitting would
  reduce clarity
- Error values MUST be checked; ignored errors require an
  explicit comment justifying the decision
- Logging via `fmt.Println` is acceptable for this challenge
  scope; structured logging is not required

## Governance

This constitution governs all development on the
codecrafters-redis-go project. It supersedes ad-hoc practices
and MUST be consulted when making architectural decisions.

**Amendment procedure**: Any principle change requires:
1. A written proposal describing the change and rationale
2. Update to this document with version bump
3. Review of all dependent specs and templates for consistency

**Versioning policy**: MAJOR.MINOR.PATCH semantic versioning.
- MAJOR: Principle removal or incompatible redefinition
- MINOR: New principle or materially expanded guidance
- PATCH: Clarifications, wording, or typo fixes

**Compliance**: All code reviews and spec approvals MUST verify
alignment with these principles. Violations MUST be resolved
before merge.

**Version**: 1.0.0 | **Ratified**: 2026-02-10 | **Last Amended**: 2026-02-10

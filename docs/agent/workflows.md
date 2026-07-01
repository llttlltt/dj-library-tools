# Workflows & Quality

This file is the single source of truth for how work is validated and reported. If any other doc seems to conflict with this one on validation or reporting, this file wins.

## Verification Protocol (MANDATORY)

These rules exist because the most common failure mode is **claiming a check passed without actually running it**. Treat every one as non-negotiable.

1. **Run it, then claim it.** Never report a command, gate, or test as passing unless you have just executed it. Capture and paste the real output and the exit code:

   ```bash
   <command>; echo "exit: $?"
   ```

   "Passes" means you saw `exit: 0`, not that you expected it to.
2. **No claim without evidence.** Statements like "byte-identical", "zero warnings", or "all green" must be backed by the pasted output of the command that proves it. If you did not run it, say so.
3. **Reference artifacts are immutable.** Golden baselines (`tests/golden/`), fixtures, and recorded references represent known-good prior state. **Never edit a baseline to make a diff pass.** If output genuinely must change, stop and get explicit approval, then update the baseline as a deliberate, separately-noted change — never silently.
4. **Do not mark work complete prematurely.** A task or plan is `Completed` only when every check below passes with pasted evidence. Do not tick a task box, set `status: Completed`, or flip a status badge until that is true. When you do complete a task, check its row off (`[x]` or ✅) with the date.
5. **Fix the cause, not the check.** If a gate fails, fix the code or the generator. Do not weaken, skip, or work around the check.
6. **Tests are self-contained.** Go tests must not depend on external files in `tests/fixtures/`. Any required data must be built in-code. `tests/fixtures/` is git-ignored and used only for manual verification and golden CLI baselines. Re-running golden diffs requires the user to provide their own local library fixtures.
7. **No History Blobs; scratch goes to `/tmp`.** Build artifacts (`djlt`, `djlt_base`, `site/`), scratch/analysis output (e.g. `deadcode_*.txt`, comparison dumps), and large fixtures are never tracked. Write any throwaway analysis or temp binaries under `/tmp`, never the repo. **Run `git status` before every commit** and confirm only intended source/docs are staged — no binaries, `.txt` logs, build output, or fixtures.
8. **Commit clean the first time.** Do not commit a mistake and "fix it" in a follow-up commit (e.g. add scratch then remove it). If a wrong file was staged, unstage it before committing. The repo has no remote yet and history rewrites are permitted, but the goal is a clean history that never contained the artifact.

## Validation Gate

Run before every commit and before claiming any task complete. All must pass:

```bash
mise run fmt-check      ; echo "fmt:   $?"   # must be 0; repo is gofmt-clean. Run `mise run fmt` to auto-fix.
go build ./...          ; echo "build: $?"   # must be 0
go vet ./...            ; echo "vet:   $?"   # must be 0
mise run test           ; echo "test:  $?"   # must be 0 (runs: go test ./...)
```

When any file under `cmd/djlt-gui/frontend/src/` changes, also run the frontend checks (requires Node 22+):

```bash
cd cmd/djlt-gui/frontend && pnpm run check ; echo "biome: $?"   # must be 0; Biome lint + format
cd cmd/djlt-gui/frontend && pnpm run build ; echo "frontend: $?" # must be 0
```

`mise run fmt` applies `gofmt -w .`; `mise run fmt-check` fails if any file is not gofmt-clean. Never hand-fight gofmt — run `mise run fmt` and commit the result.

When CLI flags, descriptions, or command structure change, also:

```bash
mise run gen-docs                      # regenerate command docs
mkdocs build --strict ; echo "exit: $?"   # must be 0; treat any link WARNING as failure
```

Notes:

- `mise run gen-docs` must be idempotent: after running it, `git status docs/` should be clean (committed docs already match generator output).
- The Material for MkDocs deprecation **banner** is cosmetic and is not a build warning. A real failure is a `WARNING -` line and/or `Aborted ... in strict mode!` with a non-zero exit.
- Behavior-sensitive CLI changes: diff live output against `tests/golden/cli/*.txt` and paste the result. Report each file as a real `diff` outcome, not a bare "OK".

## Plan Execution (multi-step initiatives)

When working an implementation plan or the master index (`plan/process-architecture-migration-index-1.md`):

- Action member plans strictly in dependency order; do not start a plan until its dependencies are `Completed`.
- Use a dedicated git branch per plan; treat a failed phase as a rollback point (revert to the last green commit rather than patching forward).
- Use `git mv` for file relocations (never delete-and-recreate) to preserve `git blame` / `git log --follow` history.
- Check tasks off with the date as you finish them; sync the plan's front-matter `status` and visible badge only when the Validation Gate passes.

## Conventional Commits

- Use conventional commit messages (`feat:`, `fix:`, `refactor:`, `chore:`, `docs:`, `test:`) to document intent and scope.
- Include code changes and their corresponding documentation regeneration (`mise run gen-docs`) in the same commit so the system stays in sync.
- Commits must be self-contained and gated by a passing Validation Gate.
- The changelog is handled by release-please; do not hand-edit `CHANGELOG.md`.

## Local-Only AI Docs

Per `.gitignore`, these are gitignored and MUST NOT be committed, but MUST be kept accurate: `plan/**/*.md`, `docs/adr/`, and `AGENTS.md`. When a decision changes, update them and confirm via `git status` that they remain unstaged. Note: `docs/agent/` IS tracked — it is the canonical agent-facing documentation (kept out of the public mkdocs nav).

## Beatgrid Maintenance Pipeline

1. **Identification**: Use `ls` with statistical paths to find tracks requiring attention:
   - `beatgrids/bpm-redundancy:>80`: Find clutter/useless markers.
   - `beatgrids/bpm-stability:<100`: Find variable tempo/drift.
   - `beatgrids-count:>1`: Find tracks with multiple anchors.
2. **Segmentation**: Use `sync` to categorize tracks into maintenance playlists:

   ```bash
   djlt sync "rb/tracks" "query" --to "rb/playlists name:Target" --apply
   ```

3. **Audit & Action**: Review categorized tracks in DJ software and perform manual grid simplification.
4. **Pruning**: Re-run the `sync` command to remove tracks from maintenance playlists once their grids have been simplified/fixed.

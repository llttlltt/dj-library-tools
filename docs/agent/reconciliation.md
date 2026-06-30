# Metadata & Reconciliation

- **Reconciliation-First**: The `sync` command handles both membership and metadata.
- **Inherit and Ignore**: Metadata updates must follow an "Inherit and Ignore" policy. Do not wipe target fields if they are missing in the source.
- **Standardized Metadata**: Use neutral types in `models` for structured data like beatgrids.
- **Join Orchestration**: Use the `Join` orchestrator for dataset matching. Skip ambiguous matches unless `--match-force` is used.

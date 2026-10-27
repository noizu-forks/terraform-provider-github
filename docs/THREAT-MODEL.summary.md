# Threat Model (summary)

Crown jewels: GitHub credentials (PAT / app private key) and the managed GitHub estate (repos, teams, org settings, secrets — destroy is real). The provider is a local plugin binary with no listening surface; the interesting boundaries are operator config, the GitHub API, the local filesystem, the `gh` CLI subprocess, and the fork binary ↔ consumer project relationship.

## Register at a glance

| Severity | Count | IDs |
|---|---|---|
| High | 3 | T-001 unsigned dev_overrides binary · T-002 `token` not `Sensitive` · T-006 destructive authority |
| Medium | 6 | T-003 gh CLI credential inheritance · T-004 plaintext response cache · T-005 PEM via env · T-007 no per-run actor · T-009 no CI/dependency alerting · T-011 state holds secrets |
| Low | 3 | T-008 rate-limit amplification · T-010 attacker-controlled `base_url` · T-012 resource-name shadowing |

Coverage: 2 mitigated (T-008, T-012) · 5 partial (T-001, T-003, T-004, T-005, T-006) · 3 open (T-002, T-009, residual gaps) · 3 accepted (T-006, T-007, T-011).

## Open items worth closing

1. Mark the `token` attribute `Sensitive` (one line, `github/provider.go`).
2. Restore some advisory/lint signal on the fork, or state explicitly that dependency review is manual.

## Accepted residual

- Destructive authority is inherent to the tool; consumers must use `prevent_destroy`, state locking and plan review.
- Fork builds have no provenance path: build from a reviewed commit, consume only via `TF_CLI_CONFIG_FILE`, never install globally.
- Credential widening via `gh auth token` disappears if `token`/`app_auth` is pinned explicitly.
- Terraform state encryption is the consumer's backend choice.

Details: [THREAT-MODEL.md](THREAT-MODEL.md), [threats/attack-surface.md](threats/attack-surface.md), [threats/credentials.md](threats/credentials.md), [threats/supply-chain.md](threats/supply-chain.md).

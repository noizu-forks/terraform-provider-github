# Threat — Credentials and Auth Paths

Covers T-002, T-003, T-005, T-007, T-010, T-011.

## The resolution chain (legacy client, default)

```
app_auth block / GITHUB_APP_*  →  installation token (GenerateOAuthTokenFromApp)
        ↓ else
token attribute / GITHUB_TOKEN  →  static token
        ↓ else
gh auth token --hostname <host> →  whatever the operator's gh CLI holds
        ↓ else
anonymous                       →  public read-only
```

## Findings

**T-002 — `token` is not `Sensitive` (Open).** In `github/provider.go` only `app_auth.pem_file` carries `Sensitive: true`. Provider configuration is not written to state, but an unmarked attribute cannot be redacted by Terraform in paths that echo config into diagnostics or validation output. The PEM is masked; the more commonly used PAT is not. Fix is one line; it was left open rather than silently changed here because it alters provider schema behaviour.

**T-003 — `gh` CLI fallback widens credentials (Partial).** `tokenFromGHCLI` executes `gh auth token --hostname <host>` with `GH_PATH` (unvalidated) or `PATH` resolution. Two consequences: an operator with a broad-scope `gh` login may apply with far more authority than the config suggests, and anything that controls `GH_PATH`/`PATH` gets arbitrary execution inside the provider process. Controls: it runs **only** on the legacy path and only when no token/app auth resolved; failures are swallowed and logged at debug. There is no allowlist or opt-out flag.

**T-005 — app private key exposure (Partial).** `pem_file` may come from `GITHUB_APP_PEM_FILE`, i.e. a process-environment secret visible to same-user processes and easily echoed into CI logs. Literal `\n` unescaping means the key is often stored as a single-line env value, which invites pasting into shell history. Controls: `Sensitive: true` on the schema field, key never logged; operator hygiene otherwise.

**T-010 — base URL controls where the token goes (Accepted).** `base_url` accepts any host; the token is sent to it. This is required for GHES/GHEC. The one genuine control present: `insecure` is deprecated and **not wired to anything**, so TLS verification cannot be disabled through provider config. Keep it that way.

**T-011 — state holds secrets (Accepted).** Writing Actions secrets/keys means state contains them in plaintext unless the consumer uses an encrypted backend. `Sensitive` on individual fields only masks CLI output.

**T-007 — repudiation (Accepted).** Every action carries one identity. There is no per-run annotation, so plan/apply attribution relies on audit logs at the GitHub end.

## What a reviewer should check on any auth change

1. Does the change add a new implicit credential source? (Only the documented chain is acceptable.)
2. Is every new secret-bearing schema field marked `Sensitive`?
3. Can the change disable TLS verification?
4. Is the token ever logged (including at debug)?
5. Does app auth still require an explicit `owner`?

#!/usr/bin/env bash
# ---------------------------------------------------------------------------
# Build the noizu-forks GitHub provider fork for dev_overrides use.
#
# Fork:    github.com/noizu-forks/terraform-provider-github
# Branch:  develop (based on upstream release tag v6.13.0)
# Output:  $HOME/.local/share/opentofu/github-fork/terraform-provider-github
#
# The binary is consumed via dev_overrides in a DEDICATED OpenTofu CLI config
# (~/.config/opentofu/github-fork.tfrc), used only through TF_CLI_CONFIG_FILE
# so no other stack is affected. The primary dev_overrides key MUST be the
# OpenTofu registry address (registry.opentofu.org/...); the terraform.io key
# is kept alongside for terraform-CLI compatibility. See NOIZU.md.
#
# Usage:  ./scripts/noizu/build-dev-override.sh [output-dir]
# Default output dir: ~/.local/share/opentofu/github-fork
# ---------------------------------------------------------------------------
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT="${1:-$HOME/.local/share/opentofu/github-fork}"

# Guard: never build into an empty/root path.
[ -n "$OUT" ] || { echo "ERROR: empty output dir" >&2; exit 1; }
[ "$OUT" != "/" ] || { echo "ERROR: refusing output dir /" >&2; exit 1; }

echo "==> building $REPO_DIR -> $OUT/terraform-provider-github"
mkdir -p "$OUT"
( cd "$REPO_DIR" && go build -o "$OUT/terraform-provider-github" . )

TFRC="$HOME/.config/opentofu/github-fork.tfrc"
echo "==> writing $TFRC"
mkdir -p "$(dirname "$TFRC")"
cat > "$TFRC" <<EOF
provider_installation {
  dev_overrides {
    "registry.opentofu.org/integrations/github" = "$OUT" # OpenTofu canonical key
    "registry.terraform.io/integrations/github" = "$OUT" # terraform CLI compatibility
  }
  direct {}
}
EOF

HEAD_SHA="$(git -C "$REPO_DIR" rev-parse HEAD)"
BIN_SHA="$(shasum -a 256 "$OUT/terraform-provider-github" | cut -d' ' -f1)"
echo "==> done"
echo "binary:  $OUT/terraform-provider-github"
echo "commit:  ${HEAD_SHA:0:7}"
echo "sha256:  $BIN_SHA"
echo "tfrc:    $TFRC"
cat <<EOF

Next: run the stack against the fork with:
  export TF_CLI_CONFIG_FILE=$TFRC
EOF

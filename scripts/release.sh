#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/release.sh v1.2.3
# Updates only the version in snapcraft.yaml

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 vX.Y.Z"
  exit 1
fi

RAW_VER="$1"
# normalize without leading v
NO_V_VER="${RAW_VER#v}"

if [[ ! "$RAW_VER" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version format: $RAW_VER"
  echo "Expected vMAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH"
  exit 2
fi

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

echo "Updating snapcraft.yaml version to $NO_V_VER"
sed -E -i.bak "s/^version:[[:space:]]+.*/version: \"${NO_V_VER}\"/" snapcraft.yaml
echo "snapcraft.yaml updated."
rm -f snapcraft.yaml.bak

# Git commit, push, and tag for snapcraft.yaml only
if [ -n "$(git status --porcelain snapcraft.yaml)" ]; then
  git add snapcraft.yaml
  git commit -m "Release $RAW_VER (snapcraft.yaml)"
  git push origin $(git rev-parse --abbrev-ref HEAD)
  git tag -a "$RAW_VER" -m "Release $RAW_VER (snapcraft.yaml)"
  git push origin "$RAW_VER"
  echo "Git commit, push, and tag completed."
else
  echo "No changes to commit."
fi
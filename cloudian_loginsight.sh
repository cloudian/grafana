#!/bin/sh -e

cd "$(dirname "$0")"
SELF=$(basename "$0")

echo "$SELF"

# Files to not change:
# hyperiq_guide.pdf
find "$PWD" -type f \
  -not -name "$SELF" \
  -not -name hyperiq_guide.pdf \
  -print0 |
  xargs -0 sed -i 's/HyperIQ/LogInsight/g'

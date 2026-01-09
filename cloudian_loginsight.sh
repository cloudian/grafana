#!/bin/sh -e

cd "$(dirname "$0")"
SELF=$(basename "$0")

echo "$SELF"

# Files to not change:
# hyperiq_guide.pdf
# navtree.go: Keep 'HyperIQ' tag for dashboard links
find "$PWD" -type f \
  -not -name "$SELF" \
  -not -name hyperiq_guide.pdf \
  -not -name navtree.go \
  -print0 |
  xargs -0 sed -i 's/HyperIQ/LogInsight/g'

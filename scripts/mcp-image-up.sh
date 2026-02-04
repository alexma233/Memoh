#!/usr/bin/env sh
set -e

IMAGE="memoh-mcp:dev"

if [ "$(uname -s)" = "Darwin" ]; then
  limactl shell default -- nerdctl build -f cmd/mcp/Dockerfile -t "$IMAGE" .
  exit $?
fi

if ! command -v nerdctl >/dev/null 2>&1; then
  echo "nerdctl not found. Install nerdctl to build images."
  exit 1
fi

nerdctl build -f cmd/mcp/Dockerfile -t "$IMAGE" .

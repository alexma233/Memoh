#!/usr/bin/env sh
set -e

if [ "$(uname -s)" = "Darwin" ]; then
  limactl start default
  limactl shell default -- sudo containerd --version
  exit $?
fi

if command -v containerd >/dev/null 2>&1 && command -v nerdctl >/dev/null 2>&1; then
  containerd --version
  nerdctl --version
  exit 0
fi

if command -v apt-get >/dev/null 2>&1; then
  sudo apt-get update
  sudo apt-get install -y containerd nerdctl
elif command -v dnf >/dev/null 2>&1; then
  sudo dnf install -y containerd nerdctl
elif command -v yum >/dev/null 2>&1; then
  sudo yum install -y containerd nerdctl
elif command -v apk >/dev/null 2>&1; then
  sudo apk add --no-cache containerd nerdctl
else
  echo "No supported package manager found. Install containerd manually."
  exit 1
fi

containerd --version
nerdctl --version

#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Source common variables
source $(dirname "$0")/common.sh

# Stop and remove containers
if [ "$CONTAINER_RUNTIME" = "docker" ]; then
  docker-compose -f ./test/e2e/docker-compose.yml down
elif [ "$CONTAINER_RUNTIME" = "podman" ]; then
  podman-compose -f ./test/e2e/docker-compose.yml down
else
  echo "Unsupported container runtime: $CONTAINER_RUNTIME"
  exit 1
fi

# Remove the kubeconfig dir if it exists
if [ -d test/e2e/kubeconfig ]; then
    rm -rf test/e2e/kubeconfig
fi

# Remove the images dir if it exists
if [ -d test/e2e/images ]; then
    rm -rf test/e2e/images
fi

# Delete the kind cluster if it exists
if kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "Deleting kind cluster '$CLUSTER_NAME'..."
    kind delete cluster --name "$CLUSTER_NAME"
else
    echo "Kind cluster '$CLUSTER_NAME' not found. Skipping deletion."
fi

# Remove the 'kind' binary
if [ -f "$LOCALBIN/kind" ]; then
    echo "Removing 'kind' binary..."
    rm "$LOCALBIN/kind"
else
    echo "'kind' binary not found. Skipping removal."
fi

echo "=====Successfully Performed cleanups====="

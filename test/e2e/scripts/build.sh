#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Source common variables
source $(dirname "$0")/common.sh

# Ensure the local bin directory exists
mkdir -p "$LOCALBIN"

# Check if 'kind' is installed, download if necessary
if [ ! -f "$KIND" ]; then
  echo "'kind' not found. Installing..."
  GOBIN="$LOCALBIN" go install sigs.k8s.io/kind@$KIND_VERSION
else
  echo "'kind' is already installed."
fi

# Check if the kind cluster exists, create if necessary
if ! "$KIND" get clusters | grep -q "^$CLUSTER_NAME$"; then
  echo "Creating kind cluster '$CLUSTER_NAME'..."
  "$KIND" create cluster --name "$CLUSTER_NAME" --kubeconfig ./test/e2e/kubeconfig/config
else
  echo "Kind cluster '$CLUSTER_NAME' already exists."
fi

# Ensure the kubeconfig file exists and set correct permissions
if [ -f ./test/e2e/kubeconfig/config ]; then
  echo "kubeconfig file exists. Setting correct permissions..."
  chmod 644 ./test/e2e/kubeconfig/config
  chown $(id -u):$(id -g) ./test/e2e/kubeconfig/config
else
  echo "kubeconfig file does not exist. Exiting..."
  exit 1
fi

# Build the pac-service image
$CONTAINER_RUNTIME build -t localhost/pac:e2e-test .

# Load the image to the kind cluster
mkdir -p ./test/e2e/images
$CONTAINER_RUNTIME save -o ./test/e2e/images/pac_e2e_test.tar pac:e2e-test
kind load image-archive ./test/e2e/images/pac_e2e_test.tar --name $CLUSTER_NAME

echo "=====Loaded the image successfully to the kind cluster======"

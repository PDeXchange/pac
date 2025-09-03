#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Source common variables
source $(dirname "$0")/common.sh

# Run make deploy and create a dummy secret and 
KUBECONFIG=./test/e2e/kubeconfig/config \
make test-deploy

# Create a IBM secret (TODO: Replace with actual later to test service endpoint)
kubectl --kubeconfig=./test/e2e/kubeconfig/config create secret generic ibm-credential --from-literal=IBMCLOUD_APIKEY=dummy -n pac-system

# Bring up the necessary containers
if [ "$CONTAINER_RUNTIME" = "docker" ]; then
  KUBERNETES_MASTER=https://host.docker.internal:6443
  KUBERNETES_SERVER=https://host.docker.internal:6443
  docker-compose -f ./test/e2e/docker-compose.yml up -d
elif [ "$CONTAINER_RUNTIME" = "podman" ]; then
  podman-compose -f ./test/e2e/docker-compose.yml up -d
else
  echo "Unsupported container runtime: $CONTAINER_RUNTIME"
  exit 1
fi

# Sleep to make sure pac-go-server is running with all keycloak pre-requisites created
sleep 10

echo "=====Deployed the services successfully======"

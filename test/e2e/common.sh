# Variables
LOCALBIN="$(pwd)/bin"
KIND="$LOCALBIN/kind"
KIND_VERSION="v0.29.0"
CLUSTER_NAME="test-e2e-pac"

# Determine the container runtime to use
if command -v podman &> /dev/null; then
  CONTAINER_RUNTIME="podman"
elif command -v docker &> /dev/null; then
  CONTAINER_RUNTIME="docker"
else
  echo "Neither podman nor docker is installed. Please install one of them to proceed."
  exit 1
fi

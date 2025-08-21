#!/bin/sh

set -e

# Define Keycloak admin credentials and server URL
KEYCLOAK_URL="${KEYCLOAK_HOSTNAME:-http://localhost:8080}"
ADMIN_USER="admin"
ADMIN_PASSWORD="dummypasswd"
REALM="master"  # Use 'master' realm or your custom realm
CLIENT_NAME="my-service-client"  # Name of the service client you want to create

# Wait for Keycloak to be ready
echo "Waiting for Keycloak..."
until curl -s ${KEYCLOAK_URL}/realms/master > /dev/null; do
  echo "Waiting..."
  sleep 5
done
echo "Keycloak is ready!"

# Step 1: Obtain the Keycloak token (admin login)
ACCESS_TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=admin-cli" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASSWORD" \
  -d "grant_type=password" | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')

# Step 2: Create the client in Keycloak
curl -X POST "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "'$CLIENT_NAME'",
    "enabled": true,
    "protocol": "openid-connect",
    "publicClient": false,
    "secret": "my-secret",
    "redirectUris": ["http://localhost:8080/*"],
    "webOrigins": ["*"]
  }'

echo "Service client '$CLIENT_NAME' created in Keycloak."

#!/bin/sh

# Define Keycloak admin credentials and server URL
KEYCLOAK_URL="${KEYCLOAK_HOSTNAME:-http://localhost:8080}"
ADMIN_USER="admin"
ADMIN_PASSWORD="dummypasswd"
REALM="pac"
CLIENT_NAME="my-service-client"
CLIENT_SECRET="my-secret"

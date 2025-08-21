#!/bin/sh

set -e

# Source common variables
source $(dirname "$0")/common.sh

# Wait for Keycloak to be ready
echo "Waiting for Keycloak..."
until curl -s ${KEYCLOAK_URL}/realms/master > /dev/null; do
  echo "Waiting..."
  sleep 5
done
echo "Keycloak is ready!"


# Step 1: Obtain the Keycloak token (admin login) and create a realm
ACCESS_TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=admin-cli" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASSWORD" \
  -d "grant_type=password" | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')


curl -X POST "$KEYCLOAK_URL/admin/realms" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "realm": "'$REALM'",
    "enabled": true
  }'

# Step 2: Create the client in Keycloak
curl -X POST "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "'$CLIENT_NAME'",
    "enabled": true,
    "protocol": "openid-connect",
    "publicClient": false,
    "directAccessGrantsEnabled": true,
    "secret": "'$CLIENT_SECRET'"
  }'

echo "Service client '$CLIENT_NAME' created in Keycloak."


# Step 3: Assign view-groups access to default-roles-pac realm role
REALM_ROLE=default-roles-pac

ACCOUNT_CLIENT_ID=$(curl -s \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM/clients?clientId=account" | \
  grep -o '"id":"[^"]*"' | head -n 1 | cut -d':' -f2 | tr -d '"')

VIEW_GROUPS_ROLE=$(curl -s \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM/clients/$ACCOUNT_CLIENT_ID/roles/view-groups")

curl -X POST "$KEYCLOAK_URL/admin/realms/$REALM/roles/$REALM_ROLE/composites" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "[$VIEW_GROUPS_ROLE]"


# Step 4: Create 2 groups (admin and bronze)
# admin group
ADMIN_LOCATION=$(curl -i -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "admin"}' | grep -Fi Location)

ADMIN_LOCATION_ID=$(echo "$ADMIN_LOCATION" | awk -F '/' '{print $NF}' | tr -d '\r')

# bronze group
BRONZE_LOCATION=$(curl -i -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "bronze"}' | grep -Fi Location)

BRONZE_LOCATION_ID=$(echo "$BRONZE_LOCATION" | awk -F '/' '{print $NF}' | tr -d '\r')


# Step 5: Create 2 users and set the passwords for each (root and test)
# root
ROOT_LOCATION=$(curl -i -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "root", "enabled": true, "email": "root@example.com", "firstName": "root", "lastName": "root"}' | grep -Fi Location)

ROOTUSER_ID=$(echo "$ROOT_LOCATION" | awk -F '/' '{print $NF}' | tr -d '\r')

curl -X PUT "$KEYCLOAK_URL/admin/realms/$REALM/users/$ROOTUSER_ID/reset-password" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "password",
    "temporary": false,
    "value": "passwd"
}'

# test
TEST_LOCATION=$(curl -i -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "enabled": true, "email": "test@example.com", "firstName": "test", "lastName": "test"}' | grep -Fi Location)

TEST_USER_ID=$(echo "$TEST_LOCATION" | awk -F '/' '{print $NF}' | tr -d '\r')

curl -X PUT "$KEYCLOAK_URL/admin/realms/$REALM/users/$TEST_USER_ID/reset-password" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "password",
    "temporary": false,
    "value": "passwd"
}'


# Step 6: Add users to group (root -> admin, test -> bronze)
curl -X PUT "$KEYCLOAK_URL/admin/realms/$REALM/users/$ROOTUSER_ID/groups/$ADMIN_LOCATION_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'

curl -X PUT "$KEYCLOAK_URL/admin/realms/$REALM/users/$TEST_USER_ID/groups/$BRONZE_LOCATION_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'


# Step 7: Create a new realm role "manager" and assign it to the admin group
MANAGER_ROLE="manager"

curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "'$MANAGER_ROLE'"}'

ROLE_JSON=$(curl -s \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM/roles/$MANAGER_ROLE")

curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/groups/$ADMIN_LOCATION_ID/role-mappings/realm" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "[$ROLE_JSON]"


echo "Successfully completed user, group creations and assignments in Keycloak."

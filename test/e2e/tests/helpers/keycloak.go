package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
)

const (
	adminUser = "root"
	adminPass = "passwd"

	testUser = "test"
	testPass = "passwd"
)

func getUserAccessToken(username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		config.Current.KeycloakHost, config.Current.KeycloakRealm)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", config.Current.KeycloakClientID)
	data.Set("client_secret", config.Current.KeycloakClientSecret)
	data.Set("username", username)
	data.Set("password", password)
	data.Set("scope", "openid profile")

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: %s", string(bodyBytes))
	}

	var tokenResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in response")
	}

	return accessToken, nil
}

func GetAdminUserToken() (string, error) {
	return getUserAccessToken(adminUser, adminPass)
}

func GetTestUserToken() (string, error) {
	return getUserAccessToken(testUser, testPass)
}

func GetUserID(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub claim not found")
	}

	return sub, nil
}

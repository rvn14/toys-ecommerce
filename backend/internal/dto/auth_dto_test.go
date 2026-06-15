package dto

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthResponseUsesCamelCaseAccessToken(t *testing.T) {
	payload, err := json.Marshal(AuthResponse{AccessToken: "token-value"})
	if err != nil {
		t.Fatalf("marshal auth response: %v", err)
	}

	jsonBody := string(payload)
	if !strings.Contains(jsonBody, `"accessToken":"token-value"`) {
		t.Fatalf("expected accessToken field, got %s", jsonBody)
	}
	if strings.Contains(jsonBody, `"accesstoken"`) {
		t.Fatalf("unexpected legacy accesstoken field in %s", jsonBody)
	}
}

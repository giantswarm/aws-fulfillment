package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccessPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/success", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Success)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedMessage := `Please wait, you will receive an email from Giant Swarm within 48 hours from hello@giantswarm.io.`
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("Form does not contain the correct message: expected %q, got %q", expectedMessage, rr.Body.String())
	}
}

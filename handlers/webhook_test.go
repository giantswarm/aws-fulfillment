package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWebhookWithoutTokenReturnsError(t *testing.T) {
	req, err := http.NewRequest("POST", "/webhook", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Webhook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestWebhookWithTokenAsHeaderRedirectsToRoot(t *testing.T) {
	req, err := http.NewRequest("POST", "/webhook", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	token := "foo"

	req.Header.Set("x-amzn-marketplace-token", token)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Webhook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedLocation := "/?token=" + token
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Handler returned wrong redirect location: got %v want %v", location, expectedLocation)
	}
}

func TestWebhookWithTokenAsQueryParameterRedirectsToRoot(t *testing.T) {
	req, err := http.NewRequest("POST", "/webhook", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	token := "bar"

	q := req.URL.Query()
	q.Set("x-amzn-marketplace-token", token)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Webhook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedLocation := "/?token=" + token
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Handler returned wrong redirect location: got %v want %v", location, expectedLocation)
	}
}

func TestWebhookWithTokenAsFormRedirectsToRoot(t *testing.T) {
	token := "baz"

	form := url.Values{}
	form.Add("x-amzn-marketplace-token", token)

	req, err := http.NewRequest("POST", "/webhook", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Webhook)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedLocation := "/?token=" + token
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Handler returned wrong redirect location: got %v want %v", location, expectedLocation)
	}

}

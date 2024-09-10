package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/giantswarm/aws-fulfillment/aws"
	"github.com/giantswarm/aws-fulfillment/slack"
)

func getMockHandleRoot() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Root(w, r, &aws.Mock{}, &slack.Mock{})
	}
}

// Test that GETting the root URL without a token returns an error.
func TestRootWithNoTokenError(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(getMockHandleRoot())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Test that GETting the root URL with a token presents a form.
func TestRootWithTokenQueryRenders(t *testing.T) {
	req, err := http.NewRequest("GET", "/?token=example-token", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(getMockHandleRoot())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedMessage := `Please enter your email address below, and Giant Swarm will get in touch within 48 hours to discuss fulfillment of your Giant Swarm AWS SaaS offering.`
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("Form does not contain the correct message: expected %q, got %q", expectedMessage, rr.Body.String())
	}

	expectedField := `<input type="hidden" name="token" value="example-token">`
	if !strings.Contains(rr.Body.String(), expectedField) {
		t.Errorf("Form does not contain the correct token field: expected %q, got %q", expectedField, rr.Body.String())
	}
}

// Test that POSTing the root URL with a form:
// - calls the AWS ResolveCustomer endpoint
// - posts an update to Slack
// - redirects to the success page
func TestRootWithFormRedirectsToSuccess(t *testing.T) {
	var mockAWSService aws.Service = &aws.Mock{}
	var mockSlackPoster slack.Service = &slack.Mock{}

	form := url.Values{}
	form.Add("email", "user@example.com")
	form.Add("token", "example-token")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Root(w, r, mockAWSService, mockSlackPoster)
	})
	handler.ServeHTTP(rr, req)

	if mock, ok := mockAWSService.(*aws.Mock); ok && !mock.Called {
		t.Errorf("Expected the mock ResolveCustomer method to be called, but it was not")
	}
	if mock, ok := mockSlackPoster.(*slack.Mock); ok && !mock.Called {
		t.Errorf("Expected the mock SlackPoster method to be called, but it was not")
	}

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedLocation := "/success"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Handler returned wrong redirect location: got %v want %v", location, expectedLocation)
	}
}

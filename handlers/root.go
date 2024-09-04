package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/giantswarm/aws-fulfillment/aws"
	"github.com/giantswarm/aws-fulfillment/slack"
)

func Root(w http.ResponseWriter, r *http.Request, c aws.Service, s slack.Service) {
	switch r.Method {
	case http.MethodGet:
		rootGet(w, r)
	case http.MethodPost:
		rootPost(w, r, c, s)
	}
}

func rootGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
		return
	}

	escapedToken := url.QueryEscape(token)

	log.Printf("escapedToken: %s", escapedToken)

	err := Template.Execute(w, map[string]string{"Token": escapedToken})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func rootPost(w http.ResponseWriter, r *http.Request, c aws.Service, s slack.Service) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	token := r.FormValue("token")

	unescapedToken, err := url.QueryUnescape(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unescaping token: %s", err), http.StatusBadRequest)
		return
	}
	unescapedToken = strings.Replace(unescapedToken, " ", "+", -1) // TODO: Fix missing plus signs on post properly

	log.Printf("unescapedToken: %s", unescapedToken)

	log.Printf("Form posted with email: %s and token: %s", email, unescapedToken)

	customerData, err := c.FetchCustomerData(unescapedToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error resolving customer: %s", err), http.StatusInternalServerError)
		return
	}

	customerData.Email = email

	log.Printf("Customer data fetched: %+v", customerData)

	if err := s.PostCustomerData(customerData); err != nil {
		http.Error(w, fmt.Sprintf("Error posting to Slack: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

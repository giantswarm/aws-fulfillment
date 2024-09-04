package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("header token: %s", r.Header.Get("x-amzn-marketplace-token"))
	log.Printf("query token: %s", r.URL.Query().Get("x-amzn-marketplace-token"))

	if err := r.ParseForm(); err != nil {
		log.Printf("could not parse form: %s", err)
	}
	formToken := r.FormValue("x-amzn-marketplace-token")
	log.Printf("form token: %s", formToken)

	var token string
	if r.Header.Get("x-amzn-marketplace-token") != "" {
		token = r.Header.Get("x-amzn-marketplace-token")
	} else if r.URL.Query().Get("x-amzn-marketplace-token") != "" {
		token = r.URL.Query().Get("x-amzn-marketplace-token")
	} else if formToken != "" {
		token = formToken
	}

	log.Printf("token: %s", token)

	if token == "" {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
		log.Printf("token missing, erroring")

		return
	}

	url := fmt.Sprintf("/?token=%s", token)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

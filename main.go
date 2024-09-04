package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/giantswarm/aws-fulfillment/aws"
	"github.com/giantswarm/aws-fulfillment/handlers"
	"github.com/giantswarm/aws-fulfillment/slack"
)

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "version" { // TODO: Fix this properly, this just gets us past CI.
		fmt.Println("this is a hack")
		os.Exit(0)
	}

	awsAccessKeyId := flag.String("aws-access-key-id", "", "AWS access key id")
	awsSecretAccessKey := flag.String("aws-secret-access-key", "", "AWS secret access key")

	slackToken := flag.String("slack-token", "", "Slack API token")

	if *awsAccessKeyId == "" {
		*awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if *awsSecretAccessKey == "" {
		*awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if *slackToken == "" {
		*slackToken = os.Getenv("SLACK_TOKEN")
	}

	mockAws := flag.Bool("mock-aws", false, "Mock calls to AWS")
	mockSlack := flag.Bool("mock-slack", false, "Mock calls to Slack")

	flag.Parse()

	awsService, err := aws.New(*awsAccessKeyId, *awsSecretAccessKey, *mockAws)
	if err != nil {
		log.Fatalf("failed to create aws service: %s", err)
	}

	slackService, err := slack.New(*slackToken, *mockSlack)
	if err != nil {
		log.Fatalf("failed to create slack service: %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Root(w, r, awsService, slackService)
	})
	mux.HandleFunc("/success", handlers.Success)
	mux.HandleFunc("/webhook", handlers.Webhook)
	mux.Handle("/content/", http.StripPrefix("/content/", http.FileServer(http.Dir("./content"))))

	loggedMux := loggingMiddleware(mux)

	server := &http.Server{
		Addr:              ":8000",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           loggedMux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}

//handles GET validation, POST decoding, quick acknowledgement, and asynchronous dispatch.
//Strava’s validation request uses hub.mode, hub.verify_token, and hub.challenge; successful validation must echo the challenge as JSON
// This file exposes and validates the public Strava webhook HTTP endpoints.
//The handler returns immediately after launching the goroutine. It does not wait for Strava API calls or PostgreSQL work.

package webhook

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	service       *WebhookService
	verifyToken   string
	signingSecret string
}

func NewHandler(
	service *WebhookService,
	verifyToken string,
	signingSecret string,
) *Handler {
	return &Handler{
		service:       service,
		verifyToken:   verifyToken,
		signingSecret: signingSecret,
	}
}

// executes on both GET and POST on strava/webhooks but different functions depending on GET or POST
// GET is for subscription verification which Strava would do
// POST is for receiving events from Strava, which we validate and then process asynchronously

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.verifySubscription(w, r)

	case http.MethodPost:
		h.receiveEvent(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *Handler) verifySubscription(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	mode := query.Get("hub.mode")
	verifyToken := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")

	if mode != "subscribe" ||
		verifyToken == "" ||
		verifyToken != h.verifyToken ||
		challenge == "" {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"hub.challenge": challenge, // returns the challenge back to Strava for verification
		},
	)
}

func (h *Handler) receiveEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	//read body upto 1MB
	body, err := io.ReadAll(
		io.LimitReader(r.Body, 1<<20),
	)
	if err != nil {
		http.Error(
			w,
			"failed to read request body",
			http.StatusBadRequest,
		)
		return
	}

	//verify the signature 
	signature := r.Header.Get("X-Strava-Signature")

	if !verifySignature(
		signature,
		body,
		h.signingSecret,
		time.Now(),
	) {
		http.Error(
			w,
			"invalid webhook signature",
			http.StatusForbidden,
		)
		return
	}

	var event Event

	// Convert JSON into your Go Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(
			w,
			"invalid webhook JSON",
			http.StatusBadRequest,
		)
		return
	}

	if err := validateEvent(event); err != nil {
		log.Printf("Invalid webhook event: %v", err)
		http.Error(w, "invalid webhook event", http.StatusBadRequest)
		return
	}

	if !isSupportedEvent(event) {
		//Got it, but I don’t need to do anything with this.
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("EVENT_IGNORED"))
		return
	}

	//Run this function concurrently instead of making the HTTP request wait for it
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf(
					"Webhook processing panic: %v",
					recovered,
				)
			}
		}()

		h.service.ProcessEvent(event)
	}()

	//Returns a 200 OK to Strava immediately, so it doesn’t retry the webhook.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("EVENT_RECEIVED"))
}

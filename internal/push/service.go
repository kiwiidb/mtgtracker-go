package push

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mtgtracker/internal/middleware"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type Service struct {
	client *messaging.Client
	repo   *Repository
}

func NewService(app *firebase.App, repo *Repository) (*Service, error) {
	if app == nil {
		// Return a no-op service if Firebase is not configured
		log.Println("Firebase app is nil, push notifications will be disabled")
		return &Service{client: nil, repo: repo}, nil
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize FCM client: %w", err)
	}

	log.Println("Push notification service initialized")
	return &Service{client: client, repo: repo}, nil
}

// SendNotification sends a push notification to all devices of a player
func (s *Service) SendNotification(playerID, title, body string, data map[string]string) error {
	if s.client == nil {
		// Firebase not configured, skip sending
		return nil
	}

	// Get all device tokens for player
	tokens, err := s.repo.GetPlayerTokens(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player tokens: %w", err)
	}

	if len(tokens) == 0 {
		// No devices to notify
		return nil
	}

	// Build multicast message
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokens,
	}

	// Send to FCM
	response, err := s.client.SendMulticast(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	log.Printf("Push notification sent to %d devices (success: %d, failure: %d)",
		len(tokens), response.SuccessCount, response.FailureCount)

	// Remove invalid tokens
	s.handleFailedTokens(response, tokens)

	return nil
}

// handleFailedTokens removes invalid tokens from the database
func (s *Service) handleFailedTokens(response *messaging.BatchResponse, tokens []string) {
	for idx, resp := range response.Responses {
		if !resp.Success {
			// Check if token is invalid or unregistered
			if messaging.IsRegistrationTokenNotRegistered(resp.Error) ||
				messaging.IsInvalidArgument(resp.Error) {
				log.Printf("Removing invalid token: %s", tokens[idx])
				err := s.repo.DeleteToken(tokens[idx])
				if err != nil {
					log.Printf("Failed to delete invalid token: %v", err)
				}
			}
		}
	}
}

// RegisterRoutes registers HTTP endpoints for push token management
func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /push/v1/tokens", s.RegisterToken)
	mux.HandleFunc("DELETE /push/v1/tokens", s.UnregisterToken)
}

// RegisterToken registers a device token for push notifications
func (s *Service) RegisterToken(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	if req.Platform == "" {
		http.Error(w, "Platform is required", http.StatusBadRequest)
		return
	}

	// Validate platform
	if req.Platform != "ios" && req.Platform != "android" && req.Platform != "web" {
		http.Error(w, "Platform must be 'ios', 'android', or 'web'", http.StatusBadRequest)
		return
	}

	err := s.repo.SaveToken(userID, req.Token, req.Platform)
	if err != nil {
		log.Printf("Failed to save token: %v", err)
		http.Error(w, "Failed to register token", http.StatusInternalServerError)
		return
	}

	log.Printf("Registered push token for user %s on platform %s", userID, req.Platform)
	w.WriteHeader(http.StatusNoContent)
}

// UnregisterToken removes a device token
func (s *Service) UnregisterToken(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	err := s.repo.DeletePlayerToken(userID, req.Token)
	if err != nil {
		log.Printf("Failed to delete token: %v", err)
		http.Error(w, "Failed to unregister token", http.StatusInternalServerError)
		return
	}

	log.Printf("Unregistered push token for user %s", userID)
	w.WriteHeader(http.StatusNoContent)
}

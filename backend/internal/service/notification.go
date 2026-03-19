package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type NotificationService struct {
	tokens     repository.DeviceTokenRepository
	fcmKey     string
	httpClient *http.Client
}

func NewNotificationService(tokens repository.DeviceTokenRepository) *NotificationService {
	return &NotificationService{
		tokens:     tokens,
		fcmKey:     os.Getenv("FCM_SERVER_KEY"),
		httpClient: &http.Client{},
	}
}

func (s *NotificationService) RegisterToken(ctx context.Context, membershipID string, req model.RegisterTokenRequest) (*model.DeviceToken, error) {
	dt := &model.DeviceToken{
		MembershipID: membershipID,
		Token:        req.Token,
		Platform:     req.Platform,
	}
	if err := s.tokens.Register(ctx, dt); err != nil {
		return nil, err
	}
	return dt, nil
}

func (s *NotificationService) UnregisterToken(ctx context.Context, membershipID, token string) error {
	return s.tokens.Unregister(ctx, membershipID, token)
}

// SendToMember sends a push notification to all devices of a member.
func (s *NotificationService) SendToMember(ctx context.Context, membershipID string, notif model.Notification) {
	tokens, err := s.tokens.ListByMembership(ctx, membershipID)
	if err != nil {
		log.Printf("error fetching device tokens for member=%s: %v", membershipID, err)
		return
	}
	for _, dt := range tokens {
		go s.sendFCM(dt.Token, notif)
	}
}

// SendToMembers sends a push notification to multiple members.
func (s *NotificationService) SendToMembers(ctx context.Context, membershipIDs []string, notif model.Notification) {
	tokens, err := s.tokens.ListByMemberships(ctx, membershipIDs)
	if err != nil {
		log.Printf("error fetching device tokens: %v", err)
		return
	}
	for _, dt := range tokens {
		go s.sendFCM(dt.Token, notif)
	}
}

func (s *NotificationService) sendFCM(deviceToken string, notif model.Notification) {
	if s.fcmKey == "" {
		log.Printf("FCM_SERVER_KEY not set, skipping push to token=%s title=%s", deviceToken[:8], notif.Title)
		return
	}

	payload := map[string]any{
		"to": deviceToken,
		"notification": map[string]string{
			"title": notif.Title,
			"body":  notif.Body,
		},
		"data": notif.Data,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(body))
	if err != nil {
		log.Printf("error creating FCM request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", s.fcmKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending FCM push: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("FCM returned status=%d for token=%s", resp.StatusCode, deviceToken[:8])
	}
}

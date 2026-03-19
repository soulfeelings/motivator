package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type WebhookService struct {
	webhooks   repository.WebhookRepository
	httpClient *http.Client
}

func NewWebhookService(webhooks repository.WebhookRepository) *WebhookService {
	return &WebhookService{webhooks: webhooks, httpClient: &http.Client{}}
}

func (s *WebhookService) Create(ctx context.Context, companyID string, req model.CreateWebhookRequest) (*model.Webhook, error) {
	w := &model.Webhook{
		CompanyID: companyID,
		Name:      req.Name,
		URL:       req.URL,
		Platform:  req.Platform,
		Events:    req.Events,
	}
	if err := s.webhooks.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WebhookService) ListByCompany(ctx context.Context, companyID string) ([]model.Webhook, error) {
	return s.webhooks.ListByCompany(ctx, companyID)
}

func (s *WebhookService) Delete(ctx context.Context, id string) error {
	return s.webhooks.Delete(ctx, id)
}

// SendEvent sends a webhook notification for a specific event to all matching webhooks.
func (s *WebhookService) SendEvent(ctx context.Context, companyID, event, title, body string) {
	webhooks, err := s.webhooks.ListActiveByEvent(ctx, companyID, event)
	if err != nil {
		log.Printf("error fetching webhooks for event=%s: %v", event, err)
		return
	}
	for _, w := range webhooks {
		go s.send(w, title, body)
	}
}

func (s *WebhookService) send(w model.Webhook, title, body string) {
	var payload []byte

	switch w.Platform {
	case "slack":
		payload, _ = json.Marshal(map[string]any{
			"text": fmt.Sprintf("*%s*\n%s", title, body),
		})
	case "teams":
		payload, _ = json.Marshal(map[string]any{
			"@type":   "MessageCard",
			"summary": title,
			"sections": []map[string]any{
				{"activityTitle": title, "text": body},
			},
		})
	default:
		payload, _ = json.Marshal(map[string]string{"title": title, "body": body})
	}

	resp, err := s.httpClient.Post(w.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("error sending webhook id=%s: %v", w.ID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("webhook id=%s returned status=%d", w.ID, resp.StatusCode)
	}
}

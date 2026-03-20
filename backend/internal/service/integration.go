package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type IntegrationService struct {
	integrations repository.IntegrationRepository
	members      repository.MembershipRepository
	achievements repository.AchievementRepository
}

func NewIntegrationService(integrations repository.IntegrationRepository, members repository.MembershipRepository, achievements repository.AchievementRepository) *IntegrationService {
	return &IntegrationService{integrations: integrations, members: members, achievements: achievements}
}

func (s *IntegrationService) Create(ctx context.Context, companyID string, req model.CreateIntegrationRequest) (*model.Integration, error) {
	secret, err := generateWebhookSecret()
	if err != nil {
		return nil, err
	}
	i := &model.Integration{
		CompanyID:     companyID,
		Provider:      req.Provider,
		Name:          req.Name,
		Config:        req.Config,
		WebhookSecret: secret,
	}
	if i.Config == nil {
		i.Config = map[string]any{}
	}
	if err := s.integrations.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *IntegrationService) ListByCompany(ctx context.Context, companyID string) ([]model.Integration, error) {
	return s.integrations.ListByCompany(ctx, companyID)
}

func (s *IntegrationService) Delete(ctx context.Context, id string) error {
	return s.integrations.Delete(ctx, id)
}

func (s *IntegrationService) CreateMapping(ctx context.Context, integrationID string, req model.CreateMappingRequest) (*model.IntegrationMapping, error) {
	m := &model.IntegrationMapping{
		IntegrationID: integrationID,
		ExternalEvent: req.ExternalEvent,
		Metric:        req.Metric,
		UserField:     req.UserField,
		Transform:     req.Transform,
	}
	if m.UserField == "" {
		m.UserField = "email"
	}
	if m.Transform == nil {
		m.Transform = map[string]any{"value": 1}
	}
	if err := s.integrations.CreateMapping(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *IntegrationService) ListMappings(ctx context.Context, integrationID string) ([]model.IntegrationMapping, error) {
	return s.integrations.ListMappings(ctx, integrationID)
}

func (s *IntegrationService) DeleteMapping(ctx context.Context, id string) error {
	return s.integrations.DeleteMapping(ctx, id)
}

func (s *IntegrationService) ListEvents(ctx context.Context, integrationID string) ([]model.IntegrationEvent, error) {
	return s.integrations.ListEvents(ctx, integrationID, 50)
}

// ProcessWebhook handles an incoming webhook from an external provider.
func (s *IntegrationService) ProcessWebhook(ctx context.Context, secret string, rawBody map[string]any) error {
	integration, err := s.integrations.GetByWebhookSecret(ctx, secret)
	if err != nil {
		return fmt.Errorf("invalid webhook secret")
	}

	// Parse event from provider
	eventName, userEmail := parseProviderEvent(integration.Provider, rawBody)
	if eventName == "" {
		return s.logEvent(ctx, integration.ID, "unknown", rawBody, nil, nil, 0, false, "unrecognized event")
	}

	// Find mapping
	mapping, err := s.integrations.GetMappingByEvent(ctx, integration.ID, eventName)
	if err != nil {
		return err
	}
	if mapping == nil {
		return s.logEvent(ctx, integration.ID, eventName, rawBody, nil, &userEmail, 0, false, "no mapping for event")
	}

	// Determine value from transform
	value := 1
	if v, ok := mapping.Transform["value"]; ok {
		if fv, ok := v.(float64); ok {
			value = int(fv)
		}
	}

	// Log the event
	metric := mapping.Metric
	if err := s.logEvent(ctx, integration.ID, eventName, rawBody, &metric, &userEmail, value, true, ""); err != nil {
		log.Printf("error logging integration event: %v", err)
	}

	// Find membership by email in the company
	// For now, we just log. The achievement evaluation will be triggered
	// when the admin UI or a cron calls the metric evaluation endpoint.
	log.Printf("integration=%s provider=%s event=%s metric=%s email=%s value=%d",
		integration.ID, integration.Provider, eventName, metric, userEmail, value)

	return nil
}

func (s *IntegrationService) logEvent(ctx context.Context, integrationID, event string, raw map[string]any, metric *string, email *string, value int, processed bool, errMsg string) error {
	var errPtr *string
	if errMsg != "" {
		errPtr = &errMsg
	}
	e := &model.IntegrationEvent{
		IntegrationID: integrationID,
		ExternalEvent: event,
		RawData:       raw,
		Metric:        metric,
		UserEmail:     email,
		Value:         value,
		Processed:     processed,
		Error:         errPtr,
	}
	return s.integrations.LogEvent(ctx, e)
}

// parseProviderEvent extracts the event name and user email from provider-specific payloads.
func parseProviderEvent(provider string, body map[string]any) (event string, email string) {
	switch provider {
	case model.ProviderJira:
		return parseJiraEvent(body)
	case model.ProviderGitHub:
		return parseGitHubEvent(body)
	case model.ProviderSalesforce:
		return parseSalesforceEvent(body)
	case model.ProviderZendesk:
		return parseZendeskEvent(body)
	default:
		// Custom: expect {event, email} at top level
		ev, _ := body["event"].(string)
		em, _ := body["email"].(string)
		return ev, em
	}
}

func parseJiraEvent(body map[string]any) (string, string) {
	webhookEvent, _ := body["webhookEvent"].(string)
	user, _ := body["user"].(map[string]any)
	email := ""
	if user != nil {
		email, _ = user["emailAddress"].(string)
	}

	switch webhookEvent {
	case "jira:issue_updated":
		changelog, _ := body["changelog"].(map[string]any)
		if changelog != nil {
			items, _ := changelog["items"].([]any)
			for _, item := range items {
				if m, ok := item.(map[string]any); ok {
					if field, _ := m["field"].(string); field == "status" {
						if to, _ := m["toString"].(string); to == "Done" || to == "Resolved" || to == "Closed" {
							return "issue_resolved", email
						}
					}
				}
			}
		}
		return "issue_updated", email
	case "jira:issue_created":
		return "issue_created", email
	default:
		return webhookEvent, email
	}
}

func parseGitHubEvent(body map[string]any) (string, string) {
	action, _ := body["action"].(string)
	sender, _ := body["sender"].(map[string]any)
	email := ""
	if sender != nil {
		email, _ = sender["email"].(string)
		if email == "" {
			email, _ = sender["login"].(string)
		}
	}

	if _, ok := body["pull_request"]; ok {
		if action == "closed" {
			pr, _ := body["pull_request"].(map[string]any)
			if merged, _ := pr["merged"].(bool); merged {
				return "pr_merged", email
			}
		}
		return "pr_" + action, email
	}
	if _, ok := body["issue"]; ok {
		return "issue_" + action, email
	}
	if _, ok := body["commits"]; ok {
		return "push", email
	}
	return action, email
}

func parseSalesforceEvent(body map[string]any) (string, string) {
	event, _ := body["event"].(string)
	email, _ := body["user_email"].(string)
	return event, email
}

func parseZendeskEvent(body map[string]any) (string, string) {
	event, _ := body["event"].(string)
	email, _ := body["current_user_email"].(string)
	if email == "" {
		if agent, _ := body["agent"].(map[string]any); agent != nil {
			email, _ = agent["email"].(string)
		}
	}
	return event, email
}

func generateWebhookSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

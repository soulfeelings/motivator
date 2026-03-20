package model

import "time"

type Integration struct {
	ID            string         `json:"id"`
	CompanyID     string         `json:"company_id"`
	Provider      string         `json:"provider"`
	Name          string         `json:"name"`
	Config        map[string]any `json:"config"`
	WebhookSecret string         `json:"webhook_secret"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type IntegrationMapping struct {
	ID            string         `json:"id"`
	IntegrationID string         `json:"integration_id"`
	ExternalEvent string         `json:"external_event"`
	Metric        string         `json:"metric"`
	UserField     string         `json:"user_field"`
	Transform     map[string]any `json:"transform"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
}

type IntegrationEvent struct {
	ID            string         `json:"id"`
	IntegrationID string         `json:"integration_id"`
	ExternalEvent string         `json:"external_event"`
	RawData       map[string]any `json:"raw_data,omitempty"`
	Metric        *string        `json:"metric,omitempty"`
	UserEmail     *string        `json:"user_email,omitempty"`
	Value         int            `json:"value"`
	Processed     bool           `json:"processed"`
	Error         *string        `json:"error,omitempty"`
	ReceivedAt    time.Time      `json:"received_at"`
}

type CreateIntegrationRequest struct {
	Provider string         `json:"provider" validate:"required"`
	Name     string         `json:"name" validate:"required"`
	Config   map[string]any `json:"config"`
}

type CreateMappingRequest struct {
	ExternalEvent string         `json:"external_event" validate:"required"`
	Metric        string         `json:"metric" validate:"required"`
	UserField     string         `json:"user_field"`
	Transform     map[string]any `json:"transform"`
}

// Supported providers
const (
	ProviderJira       = "jira"
	ProviderGitHub     = "github"
	ProviderSalesforce = "salesforce"
	ProviderZendesk    = "zendesk"
	ProviderCustom     = "custom"
)

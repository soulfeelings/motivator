package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SupabaseAdmin struct {
	url        string
	serviceKey string
	client     *http.Client
}

func NewSupabaseAdmin(url, serviceKey string) *SupabaseAdmin {
	return &SupabaseAdmin{url: url, serviceKey: serviceKey, client: &http.Client{}}
}

type CreateUserResult struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// CreateUser creates a user in Supabase Auth using the admin API.
func (s *SupabaseAdmin) CreateUser(ctx context.Context, email, password string) (*CreateUserResult, error) {
	if s.serviceKey == "" {
		return nil, fmt.Errorf("SUPABASE_SERVICE_KEY not configured")
	}

	body, _ := json.Marshal(map[string]any{
		"email":            email,
		"password":         password,
		"email_confirm":    true,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.url+"/auth/v1/admin/users", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("apikey", s.serviceKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Supabase: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp map[string]any
		json.Unmarshal(respBody, &errResp)
		msg, _ := errResp["msg"].(string)
		if msg == "" {
			msg, _ = errResp["message"].(string)
		}
		if msg == "" {
			msg = string(respBody)
		}
		return nil, fmt.Errorf("supabase error: %s", msg)
	}

	var result struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &CreateUserResult{ID: result.ID, Email: result.Email}, nil
}

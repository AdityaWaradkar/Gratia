package food

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UserClient defines the contract for internal HTTP communication with the User Service
type UserClient interface {
	IsDonor(ctx context.Context, userID string) (bool, error)
	IsVerifiedNGO(ctx context.Context, userID string) (bool, error)
}

type userClient struct {
	baseURL string
	client  *http.Client
}

// NewUserClient initializes an HTTP client specifically tuned for internal microservice calls
func NewUserClient(baseURL string) UserClient {
	return &userClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// IsDonor checks whether a user has provisioned a valid donor profile
func (u *userClient) IsDonor(ctx context.Context, userID string) (bool, error) {
	url := fmt.Sprintf("%s/internal/users/%s/donor", u.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("user service network request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("user service returned unexpected status: %d", resp.StatusCode)
	}
}

// IsVerifiedNGO checks if a user has an NGO profile that has been administratively verified
func (u *userClient) IsVerifiedNGO(ctx context.Context, userID string) (bool, error) {
	url := fmt.Sprintf("%s/internal/users/%s/ngo", u.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("user service network request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("user service returned unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Verified bool `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode user service response: %w", err)
	}

	return result.Verified, nil
}
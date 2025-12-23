package food

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// UserClient defines calls to user_service
type UserClient interface {
	IsDonor(ctx context.Context, userID string) (bool, error)
	IsVerifiedNGO(ctx context.Context, userID string) (bool, error)
}

type userClient struct {
	baseURL string
	client  *http.Client
}

// NewUserClient creates a new user service client
func NewUserClient(baseURL string) UserClient {
	return &userClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

/* ===================== DONOR ===================== */

// IsDonor checks if user has a donor profile
func (u *userClient) IsDonor(ctx context.Context, userID string) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.baseURL+"/internal/users/"+userID+"/donor",
		nil,
	)
	if err != nil {
		return false, err
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("user service error")
	}

	return true, nil
}

/* ===================== NGO ===================== */

// IsVerifiedNGO checks if user has a verified NGO profile
func (u *userClient) IsVerifiedNGO(ctx context.Context, userID string) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.baseURL+"/internal/users/"+userID+"/ngo",
		nil,
	)
	if err != nil {
		return false, err
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("user service error")
	}

	var result struct {
		Verified bool `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Verified, nil
}

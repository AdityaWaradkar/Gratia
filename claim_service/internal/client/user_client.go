package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidUserResponse = errors.New("invalid response from user service")
)

// UserClient handles communication with the user service
type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUserClient creates a new user service client instance
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// NGOStatusResponse represents the response from user service for NGO verification status
type NGOStatusResponse struct {
	Verified bool `json:"verified"`
}

// IsNGOVerified checks if a user is a verified NGO
func (c *UserClient) IsNGOVerified(
	ctx context.Context,
	userID string,
) (bool, error) {

	url := fmt.Sprintf(
		"%s/internal/users/%s/ngo-status",
		c.baseURL,
		userID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		// Continue processing

	case http.StatusNotFound:
		return false, ErrUserNotFound

	default:
		return false, fmt.Errorf(
			"user service returned unexpected status: %d",
			resp.StatusCode,
		)
	}

	var res NGOStatusResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	if resp.ContentLength == 0 {
		return false, ErrInvalidUserResponse
	}

	return res.Verified, nil
}
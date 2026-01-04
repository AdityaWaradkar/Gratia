package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ngoStatusResponse struct {
	Verified bool `json:"verified"`
}

func (c *UserClient) IsNGOVerified(
	ctx context.Context,
	userID string,
) (bool, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/internal/users/%s/ngo-status", c.baseURL, userID),
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
		// continue
	case http.StatusNotFound:
		return false, ErrUserNotFound
	default:
		return false, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	var res ngoStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	return res.Verified, nil
}

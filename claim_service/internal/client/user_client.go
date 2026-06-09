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

/*
Response DTO
*/

type NGOStatusResponse struct {
	Verified bool `json:"verified"`
}

/*
IsNGOVerified

Expected response from user_service:

{
    "verified": true
}
*/

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
		// continue

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

	// Defensive validation.
	// Currently the response only contains one field,
	// but this gives us a single place to extend later.
	if resp.ContentLength == 0 {
		return false, ErrInvalidUserResponse
	}

	return res.Verified, nil
}
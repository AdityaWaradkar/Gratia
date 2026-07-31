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
	ErrFoodNotFound       = errors.New("food listing not found")
	ErrInvalidFoodResponse = errors.New("invalid response from food service")
)

// FoodClient handles communication with the food service
type FoodClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewFoodClient creates a new food service client instance
func NewFoodClient(baseURL string) *FoodClient {
	return &FoodClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FoodForClaimResponse represents the response from food service for claim info
type FoodForClaimResponse struct {
	DonorUserID string `json:"donorUserId"`
	Status      string `json:"status"`
}

// GetFoodForClaim retrieves food listing information for claim creation
func (c *FoodClient) GetFoodForClaim(
	ctx context.Context,
	foodListingID string,
) (
	donorUserID string,
	status string,
	err error,
) {

	url := fmt.Sprintf(
		"%s/internal/foods/%s/claim-info",
		c.baseURL,
		foodListingID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		// Continue processing

	case http.StatusNotFound:
		return "", "", ErrFoodNotFound

	default:
		return "", "", fmt.Errorf(
			"food service returned unexpected status: %d",
			resp.StatusCode,
		)
	}

	var res FoodForClaimResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	if res.DonorUserID == "" {
		return "", "", ErrInvalidFoodResponse
	}

	if res.Status == "" {
		return "", "", ErrInvalidFoodResponse
	}

	return res.DonorUserID, res.Status, nil
}
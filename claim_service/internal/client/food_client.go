package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FoodClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewFoodClient(baseURL string) *FoodClient {
	return &FoodClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type foodForClaimResponse struct {
	DonorUserID string `json:"donorUserId"`
	Status      string `json:"status"`
}

func (c *FoodClient) GetFoodForClaim(
	ctx context.Context,
	foodListingID string,
) (string, string, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/internal/foods/%s/claim-info", c.baseURL, foodListingID),
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

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("food service returned status %d", resp.StatusCode)
	}

	var res foodForClaimResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	return res.DonorUserID, res.Status, nil
}

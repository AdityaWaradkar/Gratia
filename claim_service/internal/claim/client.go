package claim

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

// IsVerifiedNGO checks if a user is a verified NGO
func (c *UserClient) IsVerifiedNGO(ctx context.Context, userID string) (bool, error) {
	url := fmt.Sprintf("%s/internal/users/%s/ngo", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("user service returned unexpected status: %d", resp.StatusCode)
	}

	var res struct {
		Verified bool `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	return res.Verified, nil
}

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

// GetFoodForClaim retrieves food listing information for claim creation
func (c *FoodClient) GetFoodForClaim(ctx context.Context, foodListingID string) (donorUserID string, status string, err error) {
	url := fmt.Sprintf("%s/internal/foods/%s/validate", c.baseURL, foodListingID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", "", fmt.Errorf("food listing not found")
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("food service returned unexpected status: %d", resp.StatusCode)
	}

	var res struct {
		Claimable bool   `json:"claimable"`
		Reason    string `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	if !res.Claimable {
		return "", "", fmt.Errorf("food listing is not claimable: %s", res.Reason)
	}

	// TODO: Add endpoint to food service that returns donorUserId and status
	return "", "", fmt.Errorf("food listing validation requires donor info - implement separate endpoint")
}

// MarkFoodClaimed marks a food listing as claimed
func (c *FoodClient) MarkFoodClaimed(ctx context.Context, foodListingID string) error {
	url := fmt.Sprintf("%s/internal/foods/%s/claim", c.baseURL, foodListingID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("food service returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}
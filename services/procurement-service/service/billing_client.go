package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type BillingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBillingClient() *BillingClient {
	baseURL := os.Getenv("BILLING_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8010"
	}

	return &BillingClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CreatePaymentIntentRequest struct {
	OrderID     uuid.UUID              `json:"order_id"`
	SupplierID  uuid.UUID              `json:"supplier_id"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	PaymentMode string                 `json:"payment_mode"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CreatePaymentIntentResponse struct {
	PaymentIntentID string  `json:"payment_intent_id"`
	ClientSecret    string  `json:"client_secret"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

func (c *BillingClient) CreatePaymentIntent(token string, req CreatePaymentIntentRequest) (*CreatePaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/api/billing/v1/payments/intent", c.baseURL)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call billing service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("billing service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response CreatePaymentIntentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

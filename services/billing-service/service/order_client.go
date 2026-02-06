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

type OrderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrderClient() *OrderClient {
	baseURL := os.Getenv("PROCUREMENT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8006"
	}

	return &OrderClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type UpdateOrderPaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status"`
}

func (c *OrderClient) UpdateOrderPaymentStatus(orderID uuid.UUID, paymentStatus string, authToken string) error {
	url := fmt.Sprintf("%s/api/v1/purchase-orders/%s/payment-status", c.baseURL, orderID.String())

	reqBody, err := json.Marshal(UpdateOrderPaymentStatusRequest{
		PaymentStatus: paymentStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call procurement service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("procurement service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

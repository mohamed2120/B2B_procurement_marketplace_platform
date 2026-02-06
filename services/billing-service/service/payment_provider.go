package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentProvider interface for payment processing
type PaymentProvider interface {
	CreatePaymentIntent(ctx context.Context, amount float64, currency string, orderID uuid.UUID, metadata map[string]interface{}) (*PaymentIntent, error)
	ConfirmPayment(ctx context.Context, paymentIntentID string) (*PaymentResult, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)
	Refund(ctx context.Context, paymentIntentID string, amount float64, reason string) (*RefundResult, error)
}

// PaymentIntent represents a payment intent
type PaymentIntent struct {
	ID            string
	ClientSecret  string
	Amount        float64
	Currency      string
	Status        string
	PaymentMethod string
}

// PaymentResult represents the result of a payment
type PaymentResult struct {
	Success       bool
	PaymentID     string
	TransactionID string
	Amount        float64
	Currency      string
	Status        string
	Metadata      map[string]interface{}
}

// WebhookEvent represents a webhook event from payment provider
type WebhookEvent struct {
	Type            string
	PaymentIntentID string
	Status          string
	Amount          float64
	Metadata        map[string]interface{}
}

// RefundResult represents the result of a refund
type RefundResult struct {
	Success     bool
	RefundID    string
	Amount      float64
	Status      string
	FailedReason string
}

// MockPaymentProvider is a mock implementation for local development
type MockPaymentProvider struct {
	simulateFailure bool
	webhookDelay    time.Duration
}

func NewMockPaymentProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		simulateFailure: false,
		webhookDelay:    2 * time.Second,
	}
}

func (m *MockPaymentProvider) SetSimulateFailure(fail bool) {
	m.simulateFailure = fail
}

func (m *MockPaymentProvider) CreatePaymentIntent(ctx context.Context, amount float64, currency string, orderID uuid.UUID, metadata map[string]interface{}) (*PaymentIntent, error) {
	intentID := fmt.Sprintf("pi_mock_%s", uuid.New().String()[:8])
	return &PaymentIntent{
		ID:            intentID,
		ClientSecret:  fmt.Sprintf("%s_secret_%s", intentID, uuid.New().String()[:16]),
		Amount:        amount,
		Currency:      currency,
		Status:        "requires_payment_method",
		PaymentMethod: "card",
	}, nil
}

func (m *MockPaymentProvider) ConfirmPayment(ctx context.Context, paymentIntentID string) (*PaymentResult, error) {
	if m.simulateFailure {
		return &PaymentResult{
			Success:   false,
			PaymentID: paymentIntentID,
			Status:    "failed",
			Metadata: map[string]interface{}{
				"error": "Insufficient funds",
			},
		}, nil
	}

	return &PaymentResult{
		Success:       true,
		PaymentID:     paymentIntentID,
		TransactionID: fmt.Sprintf("txn_%s", uuid.New().String()[:8]),
		Status:        "succeeded",
		Metadata: map[string]interface{}{
			"provider": "mock",
		},
	}, nil
}

func (m *MockPaymentProvider) HandleWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error) {
	// In a real implementation, this would verify the signature
	// For mock, we simulate a successful payment webhook
	return &WebhookEvent{
		Type:            "payment_intent.succeeded",
		PaymentIntentID: "pi_mock_webhook",
		Status:          "succeeded",
		Metadata: map[string]interface{}{
			"provider": "mock",
		},
	}, nil
}

func (m *MockPaymentProvider) Refund(ctx context.Context, paymentIntentID string, amount float64, reason string) (*RefundResult, error) {
	if m.simulateFailure {
		return &RefundResult{
			Success:      false,
			RefundID:     "",
			Status:       "failed",
			FailedReason: "Refund not allowed",
		}, nil
	}

	return &RefundResult{
		Success:  true,
		RefundID: fmt.Sprintf("re_%s", uuid.New().String()[:8]),
		Amount:   amount,
		Status:   "succeeded",
	}, nil
}

// SimulateWebhook simulates a webhook callback for testing
func (m *MockPaymentProvider) SimulateWebhook(paymentIntentID string, status string) *WebhookEvent {
	return &WebhookEvent{
		Type:            fmt.Sprintf("payment_intent.%s", status),
		PaymentIntentID: paymentIntentID,
		Status:          status,
		Metadata: map[string]interface{}{
			"provider": "mock",
			"simulated": true,
		},
	}
}

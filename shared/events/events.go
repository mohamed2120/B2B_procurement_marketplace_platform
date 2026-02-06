package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	// Company events
	EventCompanyApproved EventType = "core.company.approved.v1"
	
	// Catalog events
	EventCatalogPartApproved EventType = "catalog.lib_part.approved.v1"
	
	// Procurement events
	EventPRApproved        EventType = "procurement.pr.approved.v1"
	EventRFQCreated       EventType = "procurement.rfq.created.v1"
	EventQuoteSubmitted   EventType = "procurement.quote.submitted.v1"
	EventOrderPlaced      EventType = "procurement.order.placed.v1"
	
	// Logistics events
	EventShipmentLate EventType = "logistics.shipment.late.v1"
	
	// Collaboration events
	EventChatMessageSent EventType = "collab.chat.message_sent.v1"
	
	// Billing events
	EventSubscriptionStarted EventType = "billing.subscription.started.v1"
	EventPaymentSucceeded    EventType = "billing.payment.succeeded.v1"
	EventPaymentFailed       EventType = "billing.payment.failed.v1"
	EventEscrowHeld          EventType = "billing.escrow.held.v1"
	EventSettlementReleased  EventType = "billing.settlement.released.v1"
	EventRefundIssued        EventType = "billing.refund.issued.v1"
)

// EventEnvelope wraps all events with common metadata
type EventEnvelope struct {
	ID        uuid.UUID              `json:"id"`
	Type      EventType              `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  *uuid.UUID             `json:"tenant_id,omitempty"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Payload   map[string]interface{} `json:"payload"`
}

// EventBus interface for publishing and subscribing to events
type EventBus interface {
	Publish(ctx interface{}, event *EventEnvelope) error
	Subscribe(ctx interface{}, eventType EventType, handler func(*EventEnvelope) error) error
	SubscribeAll(ctx interface{}, handler func(*EventEnvelope) error) error
}

// NewEventEnvelope creates a new event envelope
func NewEventEnvelope(eventType EventType, source string, payload map[string]interface{}) *EventEnvelope {
	return &EventEnvelope{
		ID:        uuid.New(),
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// WithTenantID sets the tenant ID on the event
func (e *EventEnvelope) WithTenantID(tenantID uuid.UUID) *EventEnvelope {
	e.TenantID = &tenantID
	return e
}

// WithUserID sets the user ID on the event
func (e *EventEnvelope) WithUserID(userID uuid.UUID) *EventEnvelope {
	e.UserID = &userID
	return e
}

// Serialize converts the event to JSON bytes
func (e *EventEnvelope) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

// DeserializeEventEnvelope creates an event envelope from JSON bytes
func DeserializeEventEnvelope(data []byte) (*EventEnvelope, error) {
	var envelope EventEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to deserialize event: %w", err)
	}
	return &envelope, nil
}

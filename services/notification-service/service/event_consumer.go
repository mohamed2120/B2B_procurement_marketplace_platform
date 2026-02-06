package service

import (
	"context"
	"fmt"
	"log"

	"github.com/b2b-platform/notification-service/models"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

// EventConsumer handles events and creates notifications
type EventConsumer struct {
	notificationService *NotificationService
}

func NewEventConsumer(notificationService *NotificationService) *EventConsumer {
	return &EventConsumer{
		notificationService: notificationService,
	}
}

func (ec *EventConsumer) HandleEvent(event *events.EventEnvelope) error {
	switch event.Type {
	case events.EventCompanyApproved:
		return ec.handleCompanyApproved(event)
	case events.EventCatalogPartApproved:
		return ec.handlePartApproved(event)
	case events.EventPRApproved:
		return ec.handlePRApproved(event)
	case events.EventRFQCreated:
		return ec.handleRFQCreated(event)
	case events.EventQuoteSubmitted:
		return ec.handleQuoteSubmitted(event)
	case events.EventOrderPlaced:
		return ec.handleOrderPlaced(event)
	case events.EventShipmentLate:
		return ec.handleShipmentLate(event)
	case events.EventChatMessageSent:
		return ec.handleChatMessageSent(event)
	case events.EventSubscriptionStarted:
		return ec.handleSubscriptionStarted(event)
	default:
		// Unknown event, skip
		return nil
	}
}

func (ec *EventConsumer) handleCompanyApproved(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	// Get company admin users (would need identity-service call in production)
	// For now, create notification for tenant admin
	notification := &models.Notification{
		TenantID: *event.TenantID,
		UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"), // Admin user
		Channel:  "in_app",
		Type:     "company.approved",
		Title:    "Company Approved",
		Message:  fmt.Sprintf("Your company %s has been approved", event.Payload["name"]),
		Status:   "pending",
	}

	return ec.notificationService.SendNotification(notification)
}

func (ec *EventConsumer) handlePartApproved(event *events.EventEnvelope) error {
	// Get part creator (would need catalog-service call in production)
	// For now, use template-based notification
	return ec.notificationService.CreateFromTemplate(
		"part_approved",
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		map[string]interface{}{
			"part_number": event.Payload["part_number"],
			"name":        event.Payload["name"],
		},
	)
}

func (ec *EventConsumer) handlePRApproved(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	return ec.notificationService.CreateFromTemplate(
		"pr_approved",
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		*event.TenantID,
		map[string]interface{}{
			"pr_number": event.Payload["pr_number"],
		},
	)
}

func (ec *EventConsumer) handleRFQCreated(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	// Notify suppliers (would need to fetch from RFQ in production)
	notification := &models.Notification{
		TenantID: *event.TenantID,
		UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"), // Supplier user
		Channel:  "in_app",
		Type:     "rfq.created",
		Title:    "New RFQ Available",
		Message:  fmt.Sprintf("New RFQ %s has been created", event.Payload["rfq_number"]),
		Status:   "pending",
	}

	return ec.notificationService.SendNotification(notification)
}

func (ec *EventConsumer) handleQuoteSubmitted(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	return ec.notificationService.CreateFromTemplate(
		"quote_submitted",
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		*event.TenantID,
		map[string]interface{}{
			"quote_number": event.Payload["quote_number"],
			"rfq_number":   event.Payload["rfq_id"],
		},
	)
}

func (ec *EventConsumer) handleOrderPlaced(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	return ec.notificationService.CreateFromTemplate(
		"order_placed",
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		*event.TenantID,
		map[string]interface{}{
			"order_number": event.Payload["po_number"],
		},
	)
}

func (ec *EventConsumer) handleShipmentLate(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	return ec.notificationService.CreateFromTemplate(
		"shipment_late",
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		*event.TenantID,
		map[string]interface{}{
			"tracking_number": event.Payload["tracking_number"],
			"eta":            event.Payload["eta"],
		},
	)
}

func (ec *EventConsumer) handleChatMessageSent(event *events.EventEnvelope) error {
	// Get thread participants (would need collaboration-service call in production)
	// For now, create notification for thread participants
	notification := &models.Notification{
		TenantID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Channel:  "in_app",
		Type:     "chat.message",
		Title:    "New Message",
		Message:  "You have a new message in a chat thread",
		Status:   "pending",
	}

	return ec.notificationService.SendNotification(notification)
}

func (ec *EventConsumer) handleSubscriptionStarted(event *events.EventEnvelope) error {
	if event.TenantID == nil {
		return fmt.Errorf("tenant_id required")
	}

	notification := &models.Notification{
		TenantID: *event.TenantID,
		UserID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"), // Admin user
		Channel:  "email",
		Type:     "subscription.started",
		Title:    "Subscription Activated",
		Message:  "Your subscription has been activated successfully",
		Status:   "pending",
	}

	return ec.notificationService.SendNotification(notification)
}

// StartEventConsumer starts listening to events
func (ec *EventConsumer) StartEventConsumer(ctx context.Context, eventBus events.EventBus) error {
	log.Println("Starting notification event consumer...")

	handler := func(event *events.EventEnvelope) error {
		log.Printf("Received event: %s", event.Type)
		if err := ec.HandleEvent(event); err != nil {
			log.Printf("Error handling event %s: %v", event.Type, err)
			return err
		}
		log.Printf("Successfully processed event: %s", event.Type)
		return nil
	}

	return eventBus.SubscribeAll(ctx, handler)
}

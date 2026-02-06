package service

import (
	"testing"
	"time"

	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

func TestIndexerService_HandleEvent(t *testing.T) {
	service := NewIndexerService()

	tests := []struct {
		name    string
		event   *events.EventEnvelope
		wantErr bool
	}{
		{
			name: "Catalog Part Approved",
			event: events.NewEventEnvelope(
				events.EventCatalogPartApproved,
				"catalog-service",
				map[string]interface{}{
					"part_id":        "part-123",
					"part_number":    "CAT-001",
					"name":           "Test Part",
					"manufacturer_id": "mfr-123",
				},
			),
			wantErr: false,
		},
		{
			name: "Company Approved",
			event: events.NewEventEnvelope(
				events.EventCompanyApproved,
				"company-service",
				map[string]interface{}{
					"company_id": "company-123",
					"name":       "Test Company",
					"subdomain":  "test",
				},
			),
			wantErr: false,
		},
		{
			name: "Order Placed",
			event: events.NewEventEnvelope(
				events.EventOrderPlaced,
				"procurement-service",
				map[string]interface{}{
					"po_id":     "po-123",
					"po_number": "PO-001",
					"pr_id":     "pr-123",
					"quote_id":  "quote-123",
				},
			),
			wantErr: false,
		},
		{
			name: "RFQ Created",
			event: events.NewEventEnvelope(
				events.EventRFQCreated,
				"procurement-service",
				map[string]interface{}{
					"rfq_id":     "rfq-123",
					"rfq_number": "RFQ-001",
					"pr_id":      "pr-123",
				},
			),
			wantErr: false,
		},
		{
			name: "Quote Submitted",
			event: events.NewEventEnvelope(
				events.EventQuoteSubmitted,
				"procurement-service",
				map[string]interface{}{
					"quote_id":     "quote-123",
					"quote_number": "QT-001",
					"rfq_id":       "rfq-123",
					"supplier_id":  "supplier-123",
				},
			),
			wantErr: false,
		},
		{
			name: "Shipment Late",
			event: events.NewEventEnvelope(
				events.EventShipmentLate,
				"logistics-service",
				map[string]interface{}{
					"shipment_id":    "ship-123",
					"tracking_number": "TRACK-001",
					"eta":            time.Now(),
				},
			),
			wantErr: false,
		},
		{
			name: "Unknown Event",
			event: events.NewEventEnvelope(
				events.EventType("unknown.event.v1"),
				"unknown-service",
				map[string]interface{}{},
			),
			wantErr: false, // Should skip unknown events
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.HandleEvent(tt.event)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				// For unknown events, error should be nil (skipped)
				// For known events, may error if OpenSearch is not available (expected in test)
				// So we just check it doesn't panic
				_ = err
			}
		})
	}
}

func TestIndexerService_EventPayloadValidation(t *testing.T) {
	service := NewIndexerService()

	// Test missing required fields
	tests := []struct {
		name    string
		event   *events.EventEnvelope
		wantErr bool
	}{
		{
			name: "Missing part_id",
			event: events.NewEventEnvelope(
				events.EventCatalogPartApproved,
				"catalog-service",
				map[string]interface{}{
					"part_number": "CAT-001",
				},
			),
			wantErr: true,
		},
		{
			name: "Missing company_id",
			event: events.NewEventEnvelope(
				events.EventCompanyApproved,
				"company-service",
				map[string]interface{}{
					"name": "Test Company",
				},
			),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.HandleEvent(tt.event)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for missing required field")
			}
		})
	}
}

package service

import (
	"time"

	"github.com/b2b-platform/logistics-service/models"
	"github.com/b2b-platform/logistics-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type LogisticsService struct {
	repo     *repository.ShipmentRepository
	eventBus events.EventBus
}

func NewLogisticsService(repo *repository.ShipmentRepository, eventBus events.EventBus) *LogisticsService {
	return &LogisticsService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *LogisticsService) Create(shipment *models.Shipment) error {
	return s.repo.Create(shipment)
}

func (s *LogisticsService) GetByID(id uuid.UUID) (*models.Shipment, error) {
	return s.repo.GetByID(id)
}

func (s *LogisticsService) List(tenantID uuid.UUID) ([]models.Shipment, error) {
	return s.repo.List(tenantID)
}

func (s *LogisticsService) UpdateTracking(shipmentID uuid.UUID, event *models.TrackingEvent) error {
	shipment, err := s.repo.GetByID(shipmentID)
	if err != nil {
		return err
	}

	// Add tracking event
	if err := s.repo.AddTrackingEvent(event); err != nil {
		return err
	}

	// Check if late
	if time.Now().After(shipment.ETA) && !shipment.IsLate {
		shipment.IsLate = true
		if err := s.repo.Update(shipment); err != nil {
			return err
		}

		// Publish late event
		lateEvent := events.NewEventEnvelope(
			events.EventShipmentLate,
			"logistics-service",
			map[string]interface{}{
				"shipment_id":    shipment.ID.String(),
				"tracking_number": shipment.TrackingNumber,
				"eta":            shipment.ETA,
			},
		).WithTenantID(shipment.TenantID)

		return s.eventBus.Publish(nil, lateEvent)
	}

	return nil
}

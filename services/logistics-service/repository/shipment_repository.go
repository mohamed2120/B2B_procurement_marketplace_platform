package repository

import (
	"github.com/b2b-platform/logistics-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(shipment *models.Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *ShipmentRepository) GetByID(id uuid.UUID) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.db.Preload("TrackingEvents").Preload("POD").Where("id = ?", id).First(&shipment).Error
	return &shipment, err
}

func (r *ShipmentRepository) List(tenantID uuid.UUID) ([]models.Shipment, error) {
	var shipments []models.Shipment
	err := r.db.Where("tenant_id = ?", tenantID).Find(&shipments).Error
	return shipments, err
}

func (r *ShipmentRepository) Update(shipment *models.Shipment) error {
	return r.db.Save(shipment).Error
}

func (r *ShipmentRepository) GetLateShipments() ([]models.Shipment, error) {
	var shipments []models.Shipment
	err := r.db.Where("is_late = ? AND late_alert_sent = ?", true, false).Find(&shipments).Error
	return shipments, err
}

func (r *ShipmentRepository) AddTrackingEvent(event *models.TrackingEvent) error {
	return r.db.Create(event).Error
}

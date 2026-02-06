package repository

import (
	"time"

	"github.com/b2b-platform/diagnostics-service/models"
	"gorm.io/gorm"
)

type DiagnosticsRepository struct {
	db *gorm.DB
}

func NewDiagnosticsRepository(db *gorm.DB) *DiagnosticsRepository {
	return &DiagnosticsRepository{db: db}
}

// Service Heartbeats
func (r *DiagnosticsRepository) ListHeartbeats() ([]models.ServiceHeartbeat, error) {
	var heartbeats []models.ServiceHeartbeat
	err := r.db.Order("last_seen_at DESC").Find(&heartbeats).Error
	return heartbeats, err
}

func (r *DiagnosticsRepository) GetHeartbeat(serviceName, instanceID string) (*models.ServiceHeartbeat, error) {
	var heartbeat models.ServiceHeartbeat
	err := r.db.Where("service_name = ? AND instance_id = ?", serviceName, instanceID).First(&heartbeat).Error
	if err != nil {
		return nil, err
	}
	return &heartbeat, nil
}

// Incidents
func (r *DiagnosticsRepository) ListIncidents(filters map[string]interface{}) ([]models.Incident, error) {
	var incidents []models.Incident
	query := r.db.Model(&models.Incident{})

	if severity, ok := filters["severity"].(string); ok && severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if serviceName, ok := filters["service_name"].(string); ok && serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("occurred_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("occurred_at <= ?", endDate)
	}
	if resolved, ok := filters["resolved"].(bool); ok {
		if resolved {
			query = query.Where("resolved_at IS NOT NULL")
		} else {
			query = query.Where("resolved_at IS NULL")
		}
	}

	err := query.Order("occurred_at DESC").Limit(100).Find(&incidents).Error
	return incidents, err
}

func (r *DiagnosticsRepository) GetIncident(id string) (*models.Incident, error) {
	var incident models.Incident
	err := r.db.Where("id = ?", id).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *DiagnosticsRepository) ResolveIncident(id string, notes string) error {
	now := time.Now()
	return r.db.Model(&models.Incident{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"resolved_at":      &now,
			"resolution_notes": notes,
			"updated_at":       now,
		}).Error
}

// Event Failures
func (r *DiagnosticsRepository) ListEventFailures(filters map[string]interface{}) ([]models.EventFailure, error) {
	var failures []models.EventFailure
	query := r.db.Model(&models.EventFailure{})

	if eventName, ok := filters["event_name"].(string); ok && eventName != "" {
		query = query.Where("event_name = ?", eventName)
	}
	if direction, ok := filters["direction"].(string); ok && direction != "" {
		query = query.Where("direction = ?", direction)
	}
	if serviceName, ok := filters["service_name"].(string); ok && serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("last_attempt_at DESC").Limit(100).Find(&failures).Error
	return failures, err
}

func (r *DiagnosticsRepository) GetEventFailure(id string) (*models.EventFailure, error) {
	var failure models.EventFailure
	err := r.db.Where("id = ?", id).First(&failure).Error
	if err != nil {
		return nil, err
	}
	return &failure, nil
}

func (r *DiagnosticsRepository) RetryEventFailure(id string) error {
	return r.db.Model(&models.EventFailure{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         "RETRYING",
			"retry_count":    r.db.Raw("retry_count + 1"),
			"last_attempt_at": time.Now(),
			"updated_at":     time.Now(),
		}).Error
}

// Jobs
func (r *DiagnosticsRepository) ListJobs(filters map[string]interface{}) ([]models.Job, error) {
	var jobs []models.Job
	query := r.db.Model(&models.Job{})

	if jobName, ok := filters["job_name"].(string); ok && jobName != "" {
		query = query.Where("job_name = ?", jobName)
	}
	if serviceName, ok := filters["service_name"].(string); ok && serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("started_at DESC").Limit(100).Find(&jobs).Error
	return jobs, err
}

// Metrics
func (r *DiagnosticsRepository) GetMetrics(serviceName string, startTime, endTime time.Time) ([]models.APIMetricMinute, error) {
	var metrics []models.APIMetricMinute
	query := r.db.Model(&models.APIMetricMinute{}).
		Where("minute_ts >= ? AND minute_ts <= ?", startTime, endTime)

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	err := query.Order("minute_ts DESC").Find(&metrics).Error
	return metrics, err
}

// Summary stats
func (r *DiagnosticsRepository) GetSummary() (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Unhealthy services (no heartbeat in last 5 minutes)
	var unhealthyCount int64
	r.db.Model(&models.ServiceHeartbeat{}).
		Where("last_seen_at < ?", time.Now().Add(-5*time.Minute)).
		Count(&unhealthyCount)
	summary["unhealthy_services"] = unhealthyCount

	// Incidents last 24h
	var incidents24h int64
	r.db.Model(&models.Incident{}).
		Where("occurred_at >= ?", time.Now().Add(-24*time.Hour)).
		Count(&incidents24h)
	summary["incidents_24h"] = incidents24h

	// Event failures
	var eventFailures int64
	r.db.Model(&models.EventFailure{}).
		Where("status = ?", "FAILED").
		Count(&eventFailures)
	summary["event_failures"] = eventFailures

	// Critical incidents
	var criticalIncidents int64
	r.db.Model(&models.Incident{}).
		Where("severity = ? AND resolved_at IS NULL", "CRITICAL").
		Count(&criticalIncidents)
	summary["critical_incidents"] = criticalIncidents

	return summary, nil
}

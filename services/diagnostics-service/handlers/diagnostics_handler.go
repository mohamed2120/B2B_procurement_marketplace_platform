package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/b2b-platform/diagnostics-service/repository"
)

type DiagnosticsHandler struct {
	repo *repository.DiagnosticsRepository
}

func NewDiagnosticsHandler(repo *repository.DiagnosticsRepository) *DiagnosticsHandler {
	return &DiagnosticsHandler{repo: repo}
}

func (h *DiagnosticsHandler) GetServices(c *gin.Context) {
	heartbeats, err := h.repo.ListHeartbeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, heartbeats)
}

func (h *DiagnosticsHandler) GetIncidents(c *gin.Context) {
	filters := make(map[string]interface{})
	
	if severity := c.Query("severity"); severity != "" {
		filters["severity"] = severity
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		filters["service_name"] = serviceName
	}
	if resolved := c.Query("resolved"); resolved != "" {
		filters["resolved"] = resolved == "true"
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filters["end_date"] = t
		}
	}

	incidents, err := h.repo.ListIncidents(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func (h *DiagnosticsHandler) GetIncident(c *gin.Context) {
	id := c.Param("id")
	incident, err := h.repo.GetIncident(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}
	c.JSON(http.StatusOK, incident)
}

func (h *DiagnosticsHandler) ResolveIncident(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.ResolveIncident(id, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Incident resolved"})
}

func (h *DiagnosticsHandler) GetEventFailures(c *gin.Context) {
	filters := make(map[string]interface{})
	
	if eventName := c.Query("event_name"); eventName != "" {
		filters["event_name"] = eventName
	}
	if direction := c.Query("direction"); direction != "" {
		filters["direction"] = direction
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		filters["service_name"] = serviceName
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	failures, err := h.repo.ListEventFailures(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, failures)
}

func (h *DiagnosticsHandler) RetryEventFailure(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.RetryEventFailure(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event failure marked for retry"})
}

func (h *DiagnosticsHandler) GetMetrics(c *gin.Context) {
	serviceName := c.Query("service")
	rangeParam := c.DefaultQuery("range", "1h")
	
	var startTime time.Time
	switch rangeParam {
	case "1h":
		startTime = time.Now().Add(-1 * time.Hour)
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-1 * time.Hour)
	}
	endTime := time.Now()

	metrics, err := h.repo.GetMetrics(serviceName, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func (h *DiagnosticsHandler) GetSummary(c *gin.Context) {
	summary, err := h.repo.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

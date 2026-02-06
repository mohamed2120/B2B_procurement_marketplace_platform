package observability

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Service     string                 `json:"service"`
	Env         string                 `json:"env"`
	RequestID   string                 `json:"request_id,omitempty"`
	TenantID    string                 `json:"tenant_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Route       string                 `json:"route,omitempty"`
	Method      string                 `json:"method,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	LatencyMS   int64                  `json:"latency_ms,omitempty"`
	ErrorCode   string                 `json:"error_code,omitempty"`
	ErrorMessage string                `json:"error_message,omitempty"`
	Message     string                 `json:"message"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

type Logger struct {
	service string
	env     string
}

func NewLogger(service string) *Logger {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	return &Logger{
		service: service,
		env:     env,
	}
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Service:   l.service,
		Env:       l.env,
		Message:   message,
		Fields:    fields,
	}

	// Extract common fields
	if fields != nil {
		if reqID, ok := fields["request_id"].(string); ok {
			entry.RequestID = reqID
		}
		if tenantID, ok := fields["tenant_id"].(string); ok {
			entry.TenantID = tenantID
		}
		if userID, ok := fields["user_id"].(string); ok {
			entry.UserID = userID
		}
		if route, ok := fields["route"].(string); ok {
			entry.Route = route
		}
		if method, ok := fields["method"].(string); ok {
			entry.Method = method
		}
		if statusCode, ok := fields["status_code"].(int); ok {
			entry.StatusCode = statusCode
		}
		if latency, ok := fields["latency_ms"].(int64); ok {
			entry.LatencyMS = latency
		}
		if errorCode, ok := fields["error_code"].(string); ok {
			entry.ErrorCode = errorCode
		}
		if errorMsg, ok := fields["error_message"].(string); ok {
			entry.ErrorMessage = errorMsg
		}
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log(LogLevelDebug, message, fields)
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(LogLevelInfo, message, fields)
}

func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(LogLevelWarn, message, fields)
}

func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(LogLevelError, message, fields)
}

func (l *Logger) WithRequest(requestID, route, method string) *RequestLogger {
	return &RequestLogger{
		logger:    l,
		requestID: requestID,
		route:     route,
		method:    method,
		startTime: time.Now(),
	}
}

type RequestLogger struct {
	logger    *Logger
	requestID string
	route     string
	method    string
	startTime time.Time
	tenantID  string
	userID    string
}

func (rl *RequestLogger) SetTenantID(tenantID string) {
	rl.tenantID = tenantID
}

func (rl *RequestLogger) SetUserID(userID string) {
	rl.userID = userID
}

func (rl *RequestLogger) LogStart() {
	rl.logger.Info("Request started", map[string]interface{}{
		"request_id": rl.requestID,
		"route":      rl.route,
		"method":     rl.method,
		"tenant_id":  rl.tenantID,
		"user_id":    rl.userID,
	})
}

func (rl *RequestLogger) LogEnd(statusCode int) {
	latency := time.Since(rl.startTime).Milliseconds()
	rl.logger.Info("Request completed", map[string]interface{}{
		"request_id": rl.requestID,
		"route":      rl.route,
		"method":     rl.method,
		"status_code": statusCode,
		"latency_ms": latency,
		"tenant_id":  rl.tenantID,
		"user_id":    rl.userID,
	})
}

func (rl *RequestLogger) LogError(statusCode int, errorCode, errorMessage string) {
	latency := time.Since(rl.startTime).Milliseconds()
	rl.logger.Error("Request failed", map[string]interface{}{
		"request_id":   rl.requestID,
		"route":        rl.route,
		"method":       rl.method,
		"status_code":  statusCode,
		"latency_ms":   latency,
		"error_code":   errorCode,
		"error_message": errorMessage,
		"tenant_id":    rl.tenantID,
		"user_id":      rl.userID,
	})
}

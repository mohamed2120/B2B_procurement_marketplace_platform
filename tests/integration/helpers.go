package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	IdentityServiceURL    = "http://localhost:8001"
	CompanyServiceURL     = "http://localhost:8002"
	CatalogServiceURL     = "http://localhost:8003"
	ProcurementServiceURL = "http://localhost:8006"
	LogisticsServiceURL   = "http://localhost:8007"
	NotificationServiceURL = "http://localhost:8009"
	OpenSearchURL         = "http://localhost:9200"
)

type TestClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	TenantID   uuid.UUID
	UserID     uuid.UUID
}

func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *TestClient) SetToken(token string) {
	c.Token = token
}

func (c *TestClient) SetTenantID(tenantID uuid.UUID) {
	c.TenantID = tenantID
}

func (c *TestClient) SetUserID(userID uuid.UUID) {
	c.UserID = userID
}

func (c *TestClient) Do(method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	if c.TenantID != uuid.Nil {
		req.Header.Set("X-Tenant-ID", c.TenantID.String())
	}

	return c.HTTPClient.Do(req)
}

func (c *TestClient) Get(path string) (*http.Response, error) {
	return c.Do("GET", path, nil)
}

func (c *TestClient) Post(path string, body interface{}) (*http.Response, error) {
	return c.Do("POST", path, body)
}

func (c *TestClient) Put(path string, body interface{}) (*http.Response, error) {
	return c.Do("PUT", path, body)
}

func (c *TestClient) Delete(path string) (*http.Response, error) {
	return c.Do("DELETE", path, nil)
}

func ParseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	TenantID  uuid.UUID `json:"tenant_id"`
}

func Login(tenantID uuid.UUID, email, password string) (*TestClient, error) {
	client := NewTestClient(IdentityServiceURL)
	
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
		TenantID: tenantID.String(),
	}

	resp, err := client.Post("/api/v1/auth/login", loginReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(body))
	}

	var loginResp LoginResponse
	if err := ParseResponse(resp, &loginResp); err != nil {
		return nil, err
	}

	client.SetToken(loginResp.Token)
	client.SetTenantID(loginResp.User.TenantID)
	client.SetUserID(loginResp.User.ID)

	return client, nil
}

func WaitForService(url string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("service at %s not ready after %d retries", url, maxRetries)
}

func WaitForServices() error {
	services := []struct {
		name string
		url  string
	}{
		{"identity", IdentityServiceURL},
		{"company", CompanyServiceURL},
		{"catalog", CatalogServiceURL},
		{"procurement", ProcurementServiceURL},
		{"logistics", LogisticsServiceURL},
		{"notification", NotificationServiceURL},
	}

	for _, svc := range services {
		if err := WaitForService(svc.url, 30); err != nil {
			return fmt.Errorf("%s service: %w", svc.name, err)
		}
	}

	return nil
}

func GenerateUniqueID() string {
	return uuid.New().String()
}

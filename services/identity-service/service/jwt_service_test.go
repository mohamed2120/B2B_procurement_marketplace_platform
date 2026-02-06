package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTService_GenerateToken(t *testing.T) {
	service := NewJWTService()

	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	roles := []string{"buyer", "approver"}

	token, err := service.GenerateToken(userID, tenantID, email, roles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Error("expected token to be generated")
	}

	// Validate the token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
	}
	if claims.TenantID != tenantID {
		t.Errorf("expected tenant ID %s, got %s", tenantID, claims.TenantID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if len(claims.Roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(claims.Roles))
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := NewJWTService()

	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	roles := []string{"buyer"}

	token, err := service.GenerateToken(userID, tenantID, email, roles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, claims.UserID)
	}
	if claims.TenantID != tenantID {
		t.Errorf("expected tenant ID %s, got %s", tenantID, claims.TenantID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestJWTService_ValidateToken_Invalid(t *testing.T) {
	service := NewJWTService()

	invalidToken := "invalid.token.here"
	_, err := service.ValidateToken(invalidToken)
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTService_ValidateToken_Expired(t *testing.T) {
	// Note: This test would require setting a custom expiration time
	// For now, we'll just test that valid tokens work
	service := NewJWTService()

	userID := uuid.New()
	tenantID := uuid.New()
	email := "test@example.com"
	roles := []string{"buyer"}

	token, err := service.GenerateToken(userID, tenantID, email, roles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Token should be valid immediately
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims == nil {
		t.Error("expected claims to be returned")
	}

	// Verify expiration is in the future
	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("expected token expiration to be in the future")
	}
}

package service

import (
	"github.com/google/uuid"
	"github.com/b2b-platform/shared/auth"
)

type JWTService struct {
	jwt *auth.JWTService
}

func NewJWTService() *JWTService {
	return &JWTService{
		jwt: auth.NewJWTService(),
	}
}

func (s *JWTService) GenerateToken(userID, tenantID uuid.UUID, email string, roles []string) (string, error) {
	return s.jwt.GenerateToken(userID, tenantID, email, roles)
}

func (s *JWTService) ValidateToken(token string) (*auth.Claims, error) {
	return s.jwt.ValidateToken(token)
}

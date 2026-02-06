package service

import (
	"time"

	"github.com/b2b-platform/billing-service/models"
	"github.com/b2b-platform/billing-service/repository"
	"github.com/b2b-platform/shared/events"
	"github.com/google/uuid"
)

type BillingService struct {
	planRepo         *repository.PlanRepository
	subscriptionRepo *repository.SubscriptionRepository
	eventBus         events.EventBus
}

func NewBillingService(
	planRepo *repository.PlanRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	eventBus events.EventBus,
) *BillingService {
	return &BillingService{
		planRepo:         planRepo,
		subscriptionRepo: subscriptionRepo,
		eventBus:         eventBus,
	}
}

func (s *BillingService) CreatePlan(plan *models.Plan) error {
	return s.planRepo.Create(plan)
}

func (s *BillingService) GetPlan(id uuid.UUID) (*models.Plan, error) {
	return s.planRepo.GetByID(id)
}

func (s *BillingService) GetPlanByCode(code string) (*models.Plan, error) {
	return s.planRepo.GetByCode(code)
}

func (s *BillingService) ListPlans() ([]models.Plan, error) {
	return s.planRepo.List()
}

func (s *BillingService) CreateSubscription(subscription *models.Subscription) error {
	subscription.Status = "active"
	subscription.StartedAt = time.Now()

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return err
	}

	// Publish event
	event := events.NewEventEnvelope(
		events.EventSubscriptionStarted,
		"billing-service",
		map[string]interface{}{
			"subscription_id": subscription.ID.String(),
			"tenant_id":       subscription.TenantID.String(),
			"plan_id":         subscription.PlanID.String(),
		},
	).WithTenantID(subscription.TenantID)

	return s.eventBus.Publish(nil, event)
}

func (s *BillingService) GetSubscription(id uuid.UUID) (*models.Subscription, error) {
	return s.subscriptionRepo.GetByID(id)
}

func (s *BillingService) GetTenantSubscription(tenantID uuid.UUID) (*models.Subscription, error) {
	return s.subscriptionRepo.GetByTenant(tenantID)
}

func (s *BillingService) CancelSubscription(subscriptionID uuid.UUID) error {
	return s.subscriptionRepo.Cancel(subscriptionID)
}

func (s *BillingService) CheckEntitlement(tenantID uuid.UUID, feature string) (bool, int, error) {
	subscription, err := s.subscriptionRepo.GetByTenant(tenantID)
	if err != nil {
		return false, 0, err
	}

	plan, err := s.planRepo.GetByID(subscription.PlanID)
	if err != nil {
		return false, 0, err
	}

	for _, entitlement := range plan.Entitlements {
		if entitlement.Feature == feature {
			return true, entitlement.Limit, nil
		}
	}

	return false, 0, nil
}

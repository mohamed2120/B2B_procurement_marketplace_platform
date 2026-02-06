package service

import (
	"github.com/b2b-platform/marketplace-service/models"
	"github.com/b2b-platform/marketplace-service/repository"
	"github.com/google/uuid"
)

type MarketplaceService struct {
	storeRepo   *repository.StoreRepository
	listingRepo *repository.ListingRepository
	mediaRepo   *repository.MediaRepository
}

func NewMarketplaceService(
	storeRepo *repository.StoreRepository,
	listingRepo *repository.ListingRepository,
	mediaRepo *repository.MediaRepository,
) *MarketplaceService {
	return &MarketplaceService{
		storeRepo:   storeRepo,
		listingRepo: listingRepo,
		mediaRepo:   mediaRepo,
	}
}

func (s *MarketplaceService) CreateStore(store *models.Store) error {
	return s.storeRepo.Create(store)
}

func (s *MarketplaceService) GetStore(id uuid.UUID) (*models.Store, error) {
	return s.storeRepo.GetByID(id)
}

func (s *MarketplaceService) ListStores(tenantID uuid.UUID, limit, offset int) ([]models.Store, error) {
	return s.storeRepo.List(tenantID, limit, offset)
}

func (s *MarketplaceService) UpdateStore(store *models.Store) error {
	return s.storeRepo.Update(store)
}

func (s *MarketplaceService) CreateListing(listing *models.Listing) error {
	return s.listingRepo.Create(listing)
}

func (s *MarketplaceService) GetListing(id uuid.UUID) (*models.Listing, error) {
	return s.listingRepo.GetByID(id)
}

func (s *MarketplaceService) ListListings(tenantID uuid.UUID, limit, offset int, listingType, status string) ([]models.Listing, error) {
	return s.listingRepo.List(tenantID, limit, offset, listingType, status)
}

func (s *MarketplaceService) UpdateListing(listing *models.Listing) error {
	return s.listingRepo.Update(listing)
}

func (s *MarketplaceService) UpdateStock(listingID uuid.UUID, quantity float64) error {
	return s.listingRepo.UpdateStock(listingID, quantity)
}

func (s *MarketplaceService) AddMedia(media *models.ListingMedia) error {
	return s.mediaRepo.Create(media)
}

func (s *MarketplaceService) GetListingMedia(listingID uuid.UUID) ([]models.ListingMedia, error) {
	return s.mediaRepo.GetByListing(listingID)
}

func (s *MarketplaceService) SetPrimaryMedia(listingID, mediaID uuid.UUID) error {
	return s.mediaRepo.SetPrimary(listingID, mediaID)
}

func (s *MarketplaceService) GetStoreListings(storeID uuid.UUID) ([]models.Listing, error) {
	return s.listingRepo.GetByStore(storeID)
}

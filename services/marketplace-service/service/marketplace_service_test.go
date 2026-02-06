package service

import (
	"testing"

	"github.com/b2b-platform/marketplace-service/models"
	"github.com/google/uuid"
)

// MockStoreRepository for testing
type MockStoreRepository struct {
	stores map[uuid.UUID]*models.Store
}

func (m *MockStoreRepository) Create(store *models.Store) error {
	if m.stores == nil {
		m.stores = make(map[uuid.UUID]*models.Store)
	}
	m.stores[store.ID] = store
	return nil
}

func (m *MockStoreRepository) GetByID(id uuid.UUID) (*models.Store, error) {
	if store, ok := m.stores[id]; ok {
		return store, nil
	}
	return nil, nil
}

func (m *MockStoreRepository) List(tenantID uuid.UUID, limit, offset int) ([]models.Store, error) {
	var result []models.Store
	for _, store := range m.stores {
		if store.TenantID == tenantID {
			result = append(result, *store)
		}
	}
	return result, nil
}

func (m *MockStoreRepository) Update(store *models.Store) error {
	m.stores[store.ID] = store
	return nil
}

// MockListingRepository for testing
type MockListingRepository struct {
	listings map[uuid.UUID]*models.Listing
}

func (m *MockListingRepository) Create(listing *models.Listing) error {
	if m.listings == nil {
		m.listings = make(map[uuid.UUID]*models.Listing)
	}
	m.listings[listing.ID] = listing
	return nil
}

func (m *MockListingRepository) GetByID(id uuid.UUID) (*models.Listing, error) {
	if listing, ok := m.listings[id]; ok {
		return listing, nil
	}
	return nil, nil
}

func (m *MockListingRepository) List(tenantID uuid.UUID, limit, offset int, listingType, status string) ([]models.Listing, error) {
	var result []models.Listing
	for _, listing := range m.listings {
		if listing.TenantID == tenantID {
			if listingType == "" || listing.Type == listingType {
				if status == "" || listing.Status == status {
					result = append(result, *listing)
				}
			}
		}
	}
	return result, nil
}

func (m *MockListingRepository) Update(listing *models.Listing) error {
	m.listings[listing.ID] = listing
	return nil
}

func (m *MockListingRepository) UpdateStock(listingID uuid.UUID, quantity float64) error {
	if listing, ok := m.listings[listingID]; ok {
		listing.StockQuantity = quantity
	}
	return nil
}

func (m *MockListingRepository) GetByStore(storeID uuid.UUID) ([]models.Listing, error) {
	var result []models.Listing
	for _, listing := range m.listings {
		if listing.StoreID == storeID {
			result = append(result, *listing)
		}
	}
	return result, nil
}

// MockMediaRepository for testing
type MockMediaRepository struct {
	media map[uuid.UUID]*models.ListingMedia
}

func (m *MockMediaRepository) Create(media *models.ListingMedia) error {
	if m.media == nil {
		m.media = make(map[uuid.UUID]*models.ListingMedia)
	}
	m.media[media.ID] = media
	return nil
}

func (m *MockMediaRepository) GetByListing(listingID uuid.UUID) ([]models.ListingMedia, error) {
	var result []models.ListingMedia
	for _, m := range m.media {
		if m.ListingID == listingID {
			result = append(result, *m)
		}
	}
	return result, nil
}

func (m *MockMediaRepository) SetPrimary(listingID, mediaID uuid.UUID) error {
	for _, m := range m.media {
		if m.ListingID == listingID {
			m.IsPrimary = (m.ID == mediaID)
		}
	}
	return nil
}

func TestMarketplaceService_CreateStore(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	tenantID := uuid.New()
	store := &models.Store{
		ID:       uuid.New(),
		TenantID: tenantID,
		Name:     "Test Store",
	}

	err := service.CreateStore(store)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify store was created
	created, _ := mockStoreRepo.GetByID(store.ID)
	if created == nil {
		t.Errorf("expected store to be created")
	}
	if created.Name != "Test Store" {
		t.Errorf("expected name 'Test Store', got %s", created.Name)
	}
}

func TestMarketplaceService_CreateListing(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	tenantID := uuid.New()
	storeID := uuid.New()
	listing := &models.Listing{
		ID:          uuid.New(),
		TenantID:    tenantID,
		StoreID:     storeID,
		Type:        "product",
		Name:        "Test Product",
		Price:       99.99,
		StockQuantity: 100.0,
		Status:      "active",
	}

	err := service.CreateListing(listing)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify listing was created
	created, _ := mockListingRepo.GetByID(listing.ID)
	if created == nil {
		t.Errorf("expected listing to be created")
	}
	if created.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got %s", created.Name)
	}
}

func TestMarketplaceService_UpdateStock(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	listingID := uuid.New()
	listing := &models.Listing{
		ID:           listingID,
		StockQuantity: 100.0,
	}
	mockListingRepo.Create(listing)

	newQuantity := 50.0
	err := service.UpdateStock(listingID, newQuantity)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify stock was updated
	updated, _ := mockListingRepo.GetByID(listingID)
	if updated.StockQuantity != newQuantity {
		t.Errorf("expected stock quantity %.2f, got %.2f", newQuantity, updated.StockQuantity)
	}
}

func TestMarketplaceService_AddMedia(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	listingID := uuid.New()
	media := &models.ListingMedia{
		ID:        uuid.New(),
		ListingID: listingID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
	}

	err := service.AddMedia(media)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify media was added
	mediaList, _ := service.GetListingMedia(listingID)
	if len(mediaList) != 1 {
		t.Errorf("expected 1 media item, got %d", len(mediaList))
	}
	if mediaList[0].URL != "https://example.com/image.jpg" {
		t.Errorf("expected URL 'https://example.com/image.jpg', got %s", mediaList[0].URL)
	}
}

func TestMarketplaceService_SetPrimaryMedia(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	listingID := uuid.New()
	media1 := &models.ListingMedia{
		ID:        uuid.New(),
		ListingID: listingID,
		IsPrimary: false,
	}
	media2 := &models.ListingMedia{
		ID:        uuid.New(),
		ListingID: listingID,
		IsPrimary: false,
	}
	mockMediaRepo.Create(media1)
	mockMediaRepo.Create(media2)

	err := service.SetPrimaryMedia(listingID, media2.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify primary media was set
	mediaList, _ := service.GetListingMedia(listingID)
	for _, m := range mediaList {
		if m.ID == media2.ID && !m.IsPrimary {
			t.Errorf("expected media %s to be primary", media2.ID)
		}
		if m.ID == media1.ID && m.IsPrimary {
			t.Errorf("expected media %s not to be primary", media1.ID)
		}
	}
}

func TestMarketplaceService_GetStoreListings(t *testing.T) {
	mockStoreRepo := &MockStoreRepository{}
	mockListingRepo := &MockListingRepository{}
	mockMediaRepo := &MockMediaRepository{}

	service := NewMarketplaceService(mockStoreRepo, mockListingRepo, mockMediaRepo)

	storeID := uuid.New()
	listing1 := &models.Listing{
		ID:      uuid.New(),
		StoreID: storeID,
		Name:    "Product 1",
	}
	listing2 := &models.Listing{
		ID:      uuid.New(),
		StoreID: storeID,
		Name:    "Product 2",
	}
	mockListingRepo.Create(listing1)
	mockListingRepo.Create(listing2)

	listings, err := service.GetStoreListings(storeID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(listings) != 2 {
		t.Errorf("expected 2 listings, got %d", len(listings))
	}
}

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/b2b-platform/shared/events"
)

type IndexerService struct {
	opensearchURL string
	httpClient    *http.Client
}

func NewIndexerService() *IndexerService {
	opensearchURL := os.Getenv("OPENSEARCH_URL")
	if opensearchURL == "" {
		opensearchURL = "http://localhost:9200"
	}

	return &IndexerService{
		opensearchURL: opensearchURL,
		httpClient:    &http.Client{},
	}
}

func (s *IndexerService) HandleEvent(event *events.EventEnvelope) error {
	switch event.Type {
	case events.EventCatalogPartApproved:
		return s.indexPart(event)
	case events.EventCompanyApproved:
		return s.indexCompany(event)
	case events.EventOrderPlaced:
		return s.indexOrder(event)
	case events.EventRFQCreated:
		return s.indexRFQ(event)
	case events.EventQuoteSubmitted:
		return s.indexQuote(event)
	case events.EventShipmentLate:
		return s.indexShipment(event)
	default:
		// Unknown event type, skip
		return nil
	}
}

func (s *IndexerService) indexPart(event *events.EventEnvelope) error {
	partID, ok := event.Payload["part_id"].(string)
	if !ok {
		return fmt.Errorf("part_id not found in event payload")
	}

	indexName := "parts"
	document := map[string]interface{}{
		"id":          partID,
		"part_number": event.Payload["part_number"],
		"name":        event.Payload["name"],
		"manufacturer_id": event.Payload["manufacturer_id"],
		"event_type":  event.Type,
		"timestamp":   event.Timestamp,
	}

	return s.indexDocument(indexName, partID, document)
}

func (s *IndexerService) indexCompany(event *events.EventEnvelope) error {
	companyID, ok := event.Payload["company_id"].(string)
	if !ok {
		return fmt.Errorf("company_id not found in event payload")
	}

	indexName := "companies"
	document := map[string]interface{}{
		"id":        companyID,
		"name":      event.Payload["name"],
		"subdomain": event.Payload["subdomain"],
		"event_type": event.Type,
		"timestamp": event.Timestamp,
	}

	return s.indexDocument(indexName, companyID, document)
}

func (s *IndexerService) indexOrder(event *events.EventEnvelope) error {
	orderID, ok := event.Payload["po_id"].(string)
	if !ok {
		return fmt.Errorf("po_id not found in event payload")
	}

	indexName := "orders"
	document := map[string]interface{}{
		"id":         orderID,
		"po_number":  event.Payload["po_number"],
		"pr_id":      event.Payload["pr_id"],
		"quote_id":   event.Payload["quote_id"],
		"event_type": event.Type,
		"timestamp":  event.Timestamp,
	}

	return s.indexDocument(indexName, orderID, document)
}

func (s *IndexerService) indexRFQ(event *events.EventEnvelope) error {
	rfqID, ok := event.Payload["rfq_id"].(string)
	if !ok {
		return fmt.Errorf("rfq_id not found in event payload")
	}

	indexName := "rfqs"
	document := map[string]interface{}{
		"id":         rfqID,
		"rfq_number": event.Payload["rfq_number"],
		"pr_id":      event.Payload["pr_id"],
		"event_type": event.Type,
		"timestamp":  event.Timestamp,
	}

	return s.indexDocument(indexName, rfqID, document)
}

func (s *IndexerService) indexQuote(event *events.EventEnvelope) error {
	quoteID, ok := event.Payload["quote_id"].(string)
	if !ok {
		return fmt.Errorf("quote_id not found in event payload")
	}

	indexName := "quotes"
	document := map[string]interface{}{
		"id":          quoteID,
		"quote_number": event.Payload["quote_number"],
		"rfq_id":      event.Payload["rfq_id"],
		"supplier_id": event.Payload["supplier_id"],
		"event_type":  event.Type,
		"timestamp":   event.Timestamp,
	}

	return s.indexDocument(indexName, quoteID, document)
}

func (s *IndexerService) indexShipment(event *events.EventEnvelope) error {
	shipmentID, ok := event.Payload["shipment_id"].(string)
	if !ok {
		return fmt.Errorf("shipment_id not found in event payload")
	}

	indexName := "shipments"
	document := map[string]interface{}{
		"id":             shipmentID,
		"tracking_number": event.Payload["tracking_number"],
		"eta":            event.Payload["eta"],
		"event_type":     event.Type,
		"timestamp":      event.Timestamp,
	}

	return s.indexDocument(indexName, shipmentID, document)
}

func (s *IndexerService) indexDocument(indexName, documentID string, document map[string]interface{}) error {
	// Ensure index exists
	if err := s.createIndexIfNotExists(indexName); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Index document
	url := fmt.Sprintf("%s/%s/_doc/%s", s.opensearchURL, indexName, documentID)
	
	jsonData, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch error: %s - %s", resp.Status, string(body))
	}

	return nil
}

func (s *IndexerService) createIndexIfNotExists(indexName string) error {
	url := fmt.Sprintf("%s/%s", s.opensearchURL, indexName)
	
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Index exists
	if resp.StatusCode == 200 {
		return nil
	}

	// Create index
	createReq, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(`{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"name": {"type": "text"},
				"timestamp": {"type": "date"}
			}
		}
	}`)))
	if err != nil {
		return err
	}

	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := s.httpClient.Do(createReq)
	if err != nil {
		return err
	}
	defer createResp.Body.Close()

	if createResp.StatusCode >= 400 {
		body, _ := io.ReadAll(createResp.Body)
		return fmt.Errorf("failed to create index: %s - %s", createResp.Status, string(body))
	}

	return nil
}

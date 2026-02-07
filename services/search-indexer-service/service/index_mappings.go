package service

// GetIndexMapping returns the OpenSearch mapping for a given index type
func GetIndexMapping(indexType string) string {
	switch indexType {
	case "parts":
		return `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"autocomplete": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": ["lowercase", "autocomplete_filter"]
						}
					},
					"filter": {
						"autocomplete_filter": {
							"type": "edge_ngram",
							"min_gram": 2,
							"max_gram": 20
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"id": {"type": "keyword"},
					"type": {"type": "keyword"},
					"part_number": {
						"type": "text",
						"fields": {
							"keyword": {"type": "keyword"},
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"manufacturer_code": {
						"type": "text",
						"fields": {
							"keyword": {"type": "keyword"},
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"name": {
						"type": "text",
						"fields": {
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"description": {"type": "text"},
					"manufacturer": {"type": "keyword"},
					"manufacturer_id": {"type": "keyword"},
					"category": {"type": "keyword"},
					"visibility": {"type": "keyword"},
					"status": {"type": "keyword"},
					"company_status": {"type": "keyword"},
					"price": {"type": "float"},
					"currency": {"type": "keyword"},
					"stock": {"type": "integer"},
					"rating": {"type": "float"},
					"timestamp": {"type": "date"}
				}
			}
		}`
	case "equipment":
		return `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"autocomplete": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": ["lowercase", "autocomplete_filter"]
						}
					},
					"filter": {
						"autocomplete_filter": {
							"type": "edge_ngram",
							"min_gram": 2,
							"max_gram": 20
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"id": {"type": "keyword"},
					"type": {"type": "keyword"},
					"model": {
						"type": "text",
						"fields": {
							"keyword": {"type": "keyword"},
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"series": {"type": "text"},
					"name": {
						"type": "text",
						"fields": {
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"description": {"type": "text"},
					"manufacturer": {"type": "keyword"},
					"manufacturer_id": {"type": "keyword"},
					"category": {"type": "keyword"},
					"visibility": {"type": "keyword"},
					"status": {"type": "keyword"},
					"company_status": {"type": "keyword"},
					"price": {"type": "float"},
					"currency": {"type": "keyword"},
					"rating": {"type": "float"},
					"eta": {"type": "integer"},
					"timestamp": {"type": "date"}
				}
			}
		}`
	case "companies":
		return `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"autocomplete": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": ["lowercase", "autocomplete_filter"]
						}
					},
					"filter": {
						"autocomplete_filter": {
							"type": "edge_ngram",
							"min_gram": 2,
							"max_gram": 20
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"id": {"type": "keyword"},
					"type": {"type": "keyword"},
					"name": {
						"type": "text",
						"fields": {
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"subdomain": {"type": "keyword"},
					"visibility": {"type": "keyword"},
					"status": {"type": "keyword"},
					"company_status": {"type": "keyword"},
					"rating": {"type": "float"},
					"timestamp": {"type": "date"}
				}
			}
		}`
	case "listings":
		return `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"autocomplete": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": ["lowercase", "autocomplete_filter"]
						}
					},
					"filter": {
						"autocomplete_filter": {
							"type": "edge_ngram",
							"min_gram": 2,
							"max_gram": 20
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"id": {"type": "keyword"},
					"type": {"type": "keyword"},
					"title": {
						"type": "text",
						"fields": {
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"name": {
						"type": "text",
						"fields": {
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"sku": {
						"type": "text",
						"fields": {
							"keyword": {"type": "keyword"},
							"autocomplete": {"type": "text", "analyzer": "autocomplete"}
						}
					},
					"description": {"type": "text"},
					"brand": {"type": "keyword"},
					"category": {"type": "keyword"},
					"supplier_id": {"type": "keyword"},
					"visibility": {"type": "keyword"},
					"status": {"type": "keyword"},
					"company_status": {"type": "keyword"},
					"price": {"type": "float"},
					"currency": {"type": "keyword"},
					"price_restricted": {"type": "boolean"},
					"stock": {"type": "integer"},
					"rating": {"type": "float"},
					"eta": {"type": "integer"},
					"timestamp": {"type": "date"}
				}
			}
		}`
	default:
		return `{
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
		}`
	}
}

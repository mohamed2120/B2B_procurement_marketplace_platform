package integration

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACEnforcement(t *testing.T) {
	// Wait for services to be ready
	require.NoError(t, WaitForServices())

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Test 1: Buyer should be able to create PR but not approve
	t.Run("Buyer can create PR but cannot approve", func(t *testing.T) {
		buyerClient, err := Login(tenantID, "buyer@demo.com", "demo123456")
		require.NoError(t, err, "Failed to login as buyer")

		// Create PR as buyer (should succeed)
		prReq := map[string]interface{}{
			"title":       "Buyer PR " + GenerateUniqueID(),
			"description": "Test PR from buyer",
			"status":      "draft",
			"items": []map[string]interface{}{
				{
					"description": "Test item",
					"quantity":    5,
					"unit_price":  50.0,
				},
			},
		}

		resp, err := buyerClient.Post("/api/v1/purchase-requests", prReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode, "Buyer should be able to create PR")

		var pr map[string]interface{}
		require.NoError(t, ParseResponse(resp, &pr))
		prID := pr["id"].(string)

		// Try to approve PR as buyer (should fail)
		approveResp, err := buyerClient.Post("/api/v1/purchase-requests/"+prID+"/approve", map[string]interface{}{})
		if err == nil {
			// If request succeeds, it should return 403
			assert.NotEqual(t, http.StatusOK, approveResp.StatusCode, "Buyer should not be able to approve PR")
			if approveResp.StatusCode != http.StatusOK {
				assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, approveResp.StatusCode)
			}
		}
	})

	// Test 2: Approver should be able to approve PR
	t.Run("Approver can approve PR", func(t *testing.T) {
		// First create PR as buyer
		buyerClient, err := Login(tenantID, "buyer@demo.com", "demo123456")
		require.NoError(t, err)

		prReq := map[string]interface{}{
			"title":       "PR for approval " + GenerateUniqueID(),
			"description": "Test PR for approval",
			"status":      "draft",
			"items": []map[string]interface{}{
				{
					"description": "Test item",
					"quantity":    10,
					"unit_price":  100.0,
				},
			},
		}

		resp, err := buyerClient.Post("/api/v1/purchase-requests", prReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var pr map[string]interface{}
		require.NoError(t, ParseResponse(resp, &pr))
		prID := pr["id"].(string)

		// Login as admin (who has approver role)
		adminClient, err := Login(tenantID, "admin@demo.com", "demo123456")
		require.NoError(t, err, "Failed to login as admin")

		// Approve PR as admin (should succeed)
		approveResp, err := adminClient.Post("/api/v1/purchase-requests/"+prID+"/approve", map[string]interface{}{})
		if err == nil {
			// Approval should succeed (200 or 201)
			assert.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, approveResp.StatusCode,
				"Admin should be able to approve PR")
		}
	})

	// Test 3: Catalog admin can manage catalog but buyer cannot
	t.Run("Catalog admin can manage catalog, buyer cannot", func(t *testing.T) {
		// Try to create a manufacturer as buyer (should fail)
		buyerClient, err := Login(tenantID, "buyer@demo.com", "demo123456")
		require.NoError(t, err)

		manufacturerReq := map[string]interface{}{
			"name":        "Test Manufacturer " + GenerateUniqueID(),
			"description": "Test manufacturer",
		}

		resp, err := buyerClient.Post("/api/v1/manufacturers", manufacturerReq)
		if err == nil {
			// Should fail with 403
			assert.NotEqual(t, http.StatusCreated, resp.StatusCode,
				"Buyer should not be able to create manufacturer")
			if resp.StatusCode != http.StatusCreated {
				assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, resp.StatusCode)
			}
		}

		// Admin should be able to create manufacturer
		adminClient, err := Login(tenantID, "admin@demo.com", "demo123456")
		require.NoError(t, err)

		resp, err = adminClient.Post("/api/v1/manufacturers", manufacturerReq)
		if err == nil {
			// Admin should succeed
			assert.Contains(t, []int{http.StatusCreated, http.StatusOK}, resp.StatusCode,
				"Admin should be able to create manufacturer")
		}
	})

	// Test 4: Supplier cannot create PR
	t.Run("Supplier cannot create PR", func(t *testing.T) {
		supplierClient, err := Login(tenantID, "supplier@demo.com", "demo123456")
		require.NoError(t, err, "Failed to login as supplier")

		prReq := map[string]interface{}{
			"title":       "Supplier PR " + GenerateUniqueID(),
			"description": "Test PR from supplier",
			"status":      "draft",
			"items": []map[string]interface{}{
				{
					"description": "Test item",
					"quantity":    5,
					"unit_price":  50.0,
				},
			},
		}

		resp, err := supplierClient.Post("/api/v1/purchase-requests", prReq)
		if err == nil {
			// Should fail with 403
			assert.NotEqual(t, http.StatusCreated, resp.StatusCode,
				"Supplier should not be able to create PR")
			if resp.StatusCode != http.StatusCreated {
				assert.Contains(t, []int{http.StatusForbidden, http.StatusUnauthorized}, resp.StatusCode)
			}
		}
	})

	t.Logf("RBAC enforcement tests passed")
}

# RBAC Navigation Testing Guide

This guide helps you test all user types and verify they can only access their allowed pages.

## Demo Accounts

All accounts use password: `demo123456`

| Email | Role | Expected Dashboard |
|-------|------|-------------------|
| `buyer.requester@demo.com` | Requester | `/app/customer/dashboard` |
| `buyer.procurement@demo.com` | Procurement Manager | `/app/customer/dashboard` |
| `supplier@demo.com` | Supplier | `/app/supplier/dashboard` |
| `admin@demo.com` | Platform Admin | `/app/admin/dashboard` |

## Manual Testing Steps

### 1. Test Requester (buyer.requester@demo.com)

**Login:**
1. Go to http://localhost:3002/login
2. Email: `buyer.requester@demo.com`
3. Password: `demo123456`
4. Click "Login"

**Expected Access:**
- ✅ `/app/customer/dashboard` - Should see dashboard
- ✅ `/app/customer/pr` - Should see PR list
- ✅ `/app/customer/pr/create` - Should see create PR form
- ✅ `/app/customer/rfq` - Should see RFQ list
- ✅ `/app/customer/orders` - Should see orders
- ✅ `/app/customer/shipments` - Should see shipments
- ✅ `/app/my-plan` - Should see plan page
- ✅ `/app/notifications` - Should see notifications
- ✅ `/app/chat` - Should see chat
- ✅ `/app/profile` - Should see profile

**Expected Sidebar:**
- ✅ Should see: Dashboard, Purchase Requests, RFQs, Orders, Shipments
- ❌ Should NOT see: Admin links, Supplier links

**Expected Restrictions:**
- ❌ `/app/admin/dashboard` - Should redirect to login or customer dashboard
- ❌ `/app/admin/tenants` - Should redirect
- ❌ `/app/admin/users` - Should redirect
- ❌ `/app/supplier/dashboard` - Should redirect
- ❌ `/app/supplier/rfq` - Should redirect

**PR Detail Page:**
- ✅ Can view PR details
- ❌ Should NOT see "Approve" or "Reject" buttons (procurement only)

---

### 2. Test Procurement Manager (buyer.procurement@demo.com)

**Login:**
1. Go to http://localhost:3002/login
2. Email: `buyer.procurement@demo.com`
3. Password: `demo123456`
4. Click "Login"

**Expected Access:**
- ✅ `/app/customer/dashboard` - Should see dashboard
- ✅ `/app/customer/pr` - Should see PR list
- ✅ `/app/customer/pr/[id]` - Should see PR detail WITH approve/reject buttons
- ✅ `/app/customer/rfq` - Should see RFQ list
- ✅ `/app/customer/rfq/[id]` - Should see RFQ detail WITH award quote button
- ✅ `/app/customer/orders` - Should see orders
- ✅ `/app/customer/shipments` - Should see shipments

**Expected Sidebar:**
- ✅ Should see: Dashboard, Purchase Requests, RFQs, Orders, Shipments
- ❌ Should NOT see: Admin links, Supplier links

**Expected Restrictions:**
- ❌ `/app/customer/pr/create` - Should redirect or not be accessible
- ❌ `/app/admin/dashboard` - Should redirect
- ❌ `/app/admin/tenants` - Should redirect
- ❌ `/app/supplier/dashboard` - Should redirect

**Special Features:**
- ✅ PR Detail: Should see "Approve" and "Reject" buttons
- ✅ RFQ Detail: Should see "Award Quote" button for submitted quotes

---

### 3. Test Supplier (supplier@demo.com)

**Login:**
1. Go to http://localhost:3002/login
2. Email: `supplier@demo.com`
3. Password: `demo123456`
4. Click "Login"

**Expected Access:**
- ✅ `/app/supplier/dashboard` - Should see dashboard
- ✅ `/app/supplier/rfq` - Should see RFQ inbox
- ✅ `/app/supplier/rfq/[id]` - Should see RFQ detail WITH quote submission form
- ✅ `/app/supplier/quotes` - Should see quotes list
- ✅ `/app/supplier/listings` - Should see listings list
- ✅ `/app/supplier/listings/create` - Should see create listing form
- ✅ `/app/supplier/orders` - Should see orders
- ✅ `/app/supplier/shipments` - Should see shipments

**Expected Sidebar:**
- ✅ Should see: Dashboard, RFQ Inbox, My Quotes, Listings, Orders, Shipments
- ❌ Should NOT see: Admin links, Buyer links

**Expected Restrictions:**
- ❌ `/app/admin/dashboard` - Should redirect
- ❌ `/app/admin/tenants` - Should redirect
- ❌ `/app/customer/dashboard` - Should redirect
- ❌ `/app/customer/pr` - Should redirect
- ❌ `/app/customer/rfq` - Should redirect

**Special Features:**
- ✅ RFQ Detail: Should see quote submission form with line items
- ✅ Create Listing: Should see form to create new product listing

---

### 4. Test Admin (admin@demo.com)

**Login:**
1. Go to http://localhost:3002/login
2. Email: `admin@demo.com`
3. Password: `demo123456`
4. Click "Login"

**Expected Access:**
- ✅ `/app/admin/dashboard` - Should see admin dashboard
- ✅ `/app/admin/tenants` - Should see tenant management table
- ✅ `/app/admin/tenants/[id]` - Should see tenant detail page
- ✅ `/app/admin/users` - Should see user management table
- ✅ `/app/admin/users/[id]` - Should see user detail page
- ✅ `/app/admin/roles-permissions` - Should see RBAC matrix
- ✅ `/app/admin/company-verification` - Should see company verification
- ✅ `/app/admin/catalog-approvals` - Should see catalog approvals
- ✅ `/app/admin/disputes` - Should see disputes
- ✅ `/app/admin/subscriptions` - Should see subscriptions
- ✅ `/app/admin/audit-logs` - Should see audit logs
- ✅ `/app/admin/diagnostics` - Should see diagnostics

**Expected Sidebar:**
- ✅ Should see: Dashboard, Tenants, Users, Roles & Permissions, Company Verification, Catalog Approvals, Disputes, Subscriptions, Audit Logs, Diagnostics
- ❌ Should NOT see: Buyer links, Supplier links

**Expected Restrictions:**
- ❌ `/app/customer/dashboard` - Should redirect
- ❌ `/app/customer/pr` - Should redirect
- ❌ `/app/supplier/dashboard` - Should redirect
- ❌ `/app/supplier/rfq` - Should redirect

**Special Features:**
- ✅ Tenant Management: Should see table with search/filter
- ✅ User Management: Should see table with user actions
- ✅ Roles & Permissions: Should see RBAC matrix

---

## Automated Testing

### Run Playwright RBAC Tests

```bash
cd frontend
npx playwright install  # First time only
npm run test:e2e -- tests/rbac-navigation.spec.ts
```

### Run All E2E Tests

```bash
cd frontend
npm run test:e2e
```

### Run with UI (Interactive)

```bash
cd frontend
npm run test:e2e:ui
```

---

## Quick Test Checklist

### ✅ Requester Checklist
- [ ] Can access customer dashboard
- [ ] Can create PR
- [ ] Can view PR list
- [ ] Cannot see approve/reject buttons on PR detail
- [ ] Cannot access admin pages
- [ ] Cannot access supplier pages
- [ ] Sidebar shows only customer links

### ✅ Procurement Checklist
- [ ] Can access customer dashboard
- [ ] Can view PR list
- [ ] Can see approve/reject buttons on PR detail
- [ ] Can view RFQ list
- [ ] Can see award quote button on RFQ detail
- [ ] Cannot create PR (requester only)
- [ ] Cannot access admin pages
- [ ] Cannot access supplier pages

### ✅ Supplier Checklist
- [ ] Can access supplier dashboard
- [ ] Can view RFQ inbox
- [ ] Can see quote submission form on RFQ detail
- [ ] Can create listing
- [ ] Cannot access admin pages
- [ ] Cannot access buyer pages
- [ ] Sidebar shows only supplier links

### ✅ Admin Checklist
- [ ] Can access admin dashboard
- [ ] Can view tenant management
- [ ] Can view user management
- [ ] Can view roles & permissions
- [ ] Can view audit logs
- [ ] Cannot access buyer pages
- [ ] Cannot access supplier pages
- [ ] Sidebar shows only admin links

---

## Troubleshooting

### If pages redirect unexpectedly:
1. Check browser console for errors
2. Verify user is logged in (check localStorage/auth cookies)
3. Check network tab for API errors
4. Verify role is correctly assigned in database

### If sidebar shows wrong links:
1. Check `frontend/components/layout/Sidebar.tsx`
2. Verify `hasRole()` function works correctly
3. Check user's roles in JWT token

### If buttons are missing:
1. Check page component for role checks
2. Verify `hasRole()` is called correctly
3. Check browser console for errors

---

## Test Results Summary

After testing, you should verify:

✅ **All users can access their allowed pages**
✅ **All users are blocked from unauthorized pages**
✅ **Sidebar navigation is filtered by role**
✅ **Action buttons (approve, award, etc.) are role-specific**
✅ **Redirects work correctly for unauthorized access**

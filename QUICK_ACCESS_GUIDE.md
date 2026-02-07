# Quick Access Guide - All Pages

## Frontend URL: http://localhost:3002

## Demo Accounts (Password: demo123456)

### 1. Platform Admin
- **Email:** `admin@demo.com`
- **Password:** `demo123456`
- **Dashboard:** http://localhost:3002/app/admin/dashboard

**Accessible Pages:**
- http://localhost:3002/app/admin/tenants (Tenant Management)
- http://localhost:3002/app/admin/tenants/[id] (Tenant Detail)
- http://localhost:3002/app/admin/users (User Management)
- http://localhost:3002/app/admin/users/[id] (User Detail)
- http://localhost:3002/app/admin/roles-permissions (RBAC Matrix)
- http://localhost:3002/app/admin/audit-logs (Audit Logs)
- http://localhost:3002/app/admin/company-verification
- http://localhost:3002/app/admin/catalog-approvals
- http://localhost:3002/app/admin/disputes
- http://localhost:3002/app/admin/subscriptions
- http://localhost:3002/app/admin/diagnostics

---

### 2. Requester (Buyer)
- **Email:** `buyer.requester@demo.com`
- **Password:** `demo123456`
- **Dashboard:** http://localhost:3002/app/customer/dashboard

**Accessible Pages:**
- http://localhost:3002/app/customer/pr (PR List)
- http://localhost:3002/app/customer/pr/create (Create PR)
- http://localhost:3002/app/customer/pr/[id] (PR Detail - NO approve button)
- http://localhost:3002/app/customer/rfq (RFQ List)
- http://localhost:3002/app/customer/rfq/[id] (RFQ Detail)
- http://localhost:3002/app/customer/orders
- http://localhost:3002/app/customer/shipments

---

### 3. Procurement Manager (Buyer)
- **Email:** `buyer.procurement@demo.com`
- **Password:** `demo123456`
- **Dashboard:** http://localhost:3002/app/customer/dashboard

**Accessible Pages:**
- http://localhost:3002/app/customer/pr (PR List)
- http://localhost:3002/app/customer/pr/[id] (PR Detail - WITH Approve/Reject buttons)
- http://localhost:3002/app/customer/rfq (RFQ List)
- http://localhost:3002/app/customer/rfq/[id] (RFQ Detail - WITH Award Quote button)
- http://localhost:3002/app/customer/orders
- http://localhost:3002/app/customer/shipments

**Note:** Cannot access `/app/customer/pr/create` (requester only)

---

### 4. Supplier
- **Email:** `supplier@demo.com`
- **Password:** `demo123456`
- **Dashboard:** http://localhost:3002/app/supplier/dashboard

**Accessible Pages:**
- http://localhost:3002/app/supplier/rfq (RFQ Inbox)
- http://localhost:3002/app/supplier/rfq/[id] (RFQ Detail - WITH Quote Submission Form)
- http://localhost:3002/app/supplier/quotes (My Quotes)
- http://localhost:3002/app/supplier/listings (Listings)
- http://localhost:3002/app/supplier/listings/create (Create Listing)
- http://localhost:3002/app/supplier/orders
- http://localhost:3002/app/supplier/shipments

---

## Public Pages (No Login Required)

- http://localhost:3002/ (Home)
- http://localhost:3002/login
- http://localhost:3002/register
- http://localhost:3002/register/buyer
- http://localhost:3002/register/supplier
- http://localhost:3002/how-it-works
- http://localhost:3002/pricing
- http://localhost:3002/contact
- http://localhost:3002/terms
- http://localhost:3002/privacy

---

## Shared Pages (After Login)

- http://localhost:3002/app/my-plan
- http://localhost:3002/app/notifications
- http://localhost:3002/app/chat
- http://localhost:3002/app/profile

---

## Quick Test Steps

1. **Open browser:** http://localhost:3002
2. **Click "Login"** or go to http://localhost:3002/login
3. **Login with any demo account** above
4. **Check sidebar** - should show only relevant links for your role
5. **Navigate to pages** - all should be accessible based on your role
6. **Try accessing other role's pages** - should redirect/block

---

## Verification Checklist

✅ All public pages load without login
✅ Login works for all 4 user types
✅ Each user type redirects to correct dashboard
✅ Sidebar shows only relevant links per role
✅ Admin can access all admin pages
✅ Buyer can access all buyer pages
✅ Supplier can access all supplier pages
✅ Cross-role access is blocked
✅ Action buttons appear based on role (approve, award, etc.)

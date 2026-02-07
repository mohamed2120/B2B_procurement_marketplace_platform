# New Pages Access Guide

## ✅ All New Pages Are Created and Available

The following pages have been added to the frontend:

### 1. Customer/Buyer Pages

#### `/app/customer/pr/[id]` - PR Detail Page
**How to access:**
1. Login as `buyer.requester@demo.com` or `buyer.procurement@demo.com` (password: `demo123456`)
2. Go to: http://localhost:3002/app/customer/pr
3. Click the **"View →"** button on any Purchase Request
4. Or directly visit: http://localhost:3002/app/customer/pr/[PR_ID]

**Features:**
- View PR details
- View PR items
- **Approve/Reject buttons** (visible only for `procurement_manager` role)

---

### 2. Supplier Pages

#### `/app/supplier/listings/create` - Create Listing Page
**How to access:**
1. Login as `supplier@demo.com` (password: `demo123456`)
2. Go to: http://localhost:3002/app/supplier/listings
3. Click the **"Create Listing"** button
4. Or directly visit: http://localhost:3002/app/supplier/listings/create

**Features:**
- Create new product listings
- Set pricing, stock, lead time
- Choose category and status

---

### 3. Admin Pages

#### `/app/admin/users` - User Management
**How to access:**
1. Login as `admin@demo.com` (password: `demo123456`)
2. In the sidebar, click **"Users"** (👥 icon)
3. Or directly visit: http://localhost:3002/app/admin/users

**Features:**
- View all users across tenants
- Search and filter users
- View user details
- Deactivate/Activate users
- Reset passwords

#### `/app/admin/users/[id]` - User Detail Page
**How to access:**
1. From the Users page, click **"View"** on any user
2. Or directly visit: http://localhost:3002/app/admin/users/[USER_ID]

**Features:**
- View user information
- View assigned roles
- View audit logs
- Toggle user active status
- Reset password

#### `/app/admin/roles-permissions` - RBAC Matrix
**How to access:**
1. Login as `admin@demo.com`
2. In the sidebar, click **"Roles & Permissions"** (🔐 icon)
3. Or directly visit: http://localhost:3002/app/admin/roles-permissions

**Features:**
- View all roles and their permissions
- See RBAC matrix (which roles have which permissions)
- Read-only view (no editing yet)

#### `/app/admin/audit-logs` - Audit Log Viewer
**How to access:**
1. Login as `admin@demo.com`
2. In the sidebar, click **"Audit Logs"** (📋 icon)
3. Or directly visit: http://localhost:3002/app/admin/audit-logs

**Features:**
- View system activity logs
- Filter by action, resource, user, date
- See user actions and system events

---

## Navigation Sidebar

After logging in, you'll see different navigation items based on your role:

### Admin Sidebar (when logged in as admin@demo.com):
- 🏠 Dashboard
- 🏢 Company Verification
- **👥 Users** ← NEW
- **🔐 Roles & Permissions** ← NEW
- 📚 Catalog Approvals
- ⚖️ Disputes
- 💳 Subscriptions
- **📋 Audit Logs** ← NEW
- 🔧 Diagnostics

### Customer Sidebar (when logged in as buyer):
- 🏠 Dashboard
- 📋 Purchase Requests (click "View →" to see detail page)
- 📝 RFQs
- 📦 Orders
- 🚚 Shipments

### Supplier Sidebar (when logged in as supplier):
- 🏠 Dashboard
- 📥 RFQ Inbox
- 💵 My Quotes
- 🏪 Listings (click "Create Listing" button)
- 📦 Orders
- 🚚 Shipments

---

## Quick Test

1. **Test Admin Pages:**
   ```bash
   # Login at http://localhost:3002/login
   # Email: admin@demo.com
   # Password: demo123456
   # Then click "Users" in sidebar
   ```

2. **Test Customer PR Detail:**
   ```bash
   # Login as buyer.requester@demo.com
   # Go to Purchase Requests
   # Click "View →" on any PR
   ```

3. **Test Supplier Create Listing:**
   ```bash
   # Login as supplier@demo.com
   # Go to Listings
   # Click "Create Listing" button
   ```

---

## Troubleshooting

If you can't see the pages:

1. **Make sure you're logged in** - All `/app/*` routes require authentication
2. **Check your role** - Some pages are role-specific:
   - Admin pages require `admin` or `super_admin` role
   - PR detail approve/reject requires `procurement_manager` role
3. **Check the sidebar** - Pages appear in the sidebar based on your role
4. **Try direct URL** - If sidebar doesn't show it, try the direct URL
5. **Clear browser cache** - Sometimes browser cache can cause issues

---

## All New Pages Summary

| Page | Route | Role Required | Access Method |
|------|-------|---------------|---------------|
| PR Detail | `/app/customer/pr/[id]` | requester/procurement_manager | Click "View →" from PR list |
| Create Listing | `/app/supplier/listings/create` | supplier | Click "Create Listing" button |
| User Management | `/app/admin/users` | admin | Sidebar → "Users" |
| User Detail | `/app/admin/users/[id]` | admin | Click "View" from users list |
| Roles & Permissions | `/app/admin/roles-permissions` | admin | Sidebar → "Roles & Permissions" |
| Audit Logs | `/app/admin/audit-logs` | admin | Sidebar → "Audit Logs" |

---

**Note:** All pages are protected routes. You must be logged in to access them.

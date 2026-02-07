# B2B Procurement Marketplace - Frontend

Next.js frontend for the B2B Procurement Marketplace Platform.

## Features

- **Tenant Resolution**: Automatic tenant detection via subdomain
- **Authentication**: JWT-based authentication with local dev mode
- **Role-Based Routing**: Protected routes based on user roles
- **UI Components**: Simple, reusable UI kit with Tailwind CSS
- **E2E Tests**: Playwright tests for critical flows

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn
- Backend services running (see main README)

### Installation

```bash
cd frontend
npm install
```

### Environment Setup

Copy the example environment file:

```bash
cp .env.local.example .env.local
```

Edit `.env.local` to configure service URLs if needed.

### Development

```bash
npm run dev
```

Open [http://localhost:3002](http://localhost:3002) in your browser.

### Build

```bash
npm run build
npm start
```

## Demo Accounts

- **Admin**: `admin@demo.com` / `demo123456`
- **Buyer**: `buyer@demo.com` / `demo123456`
- **Supplier**: `supplier@demo.com` / `demo123456`
- **Procurement Manager**: `procurement@demo.com` / `demo123456`

All accounts use tenant ID: `00000000-0000-0000-0000-000000000001`

## Testing

### E2E Tests (Playwright)

```bash
# Run all E2E tests
npm run test:e2e

# Run with UI
npm run test:e2e:ui
```

Tests include:
- Login flow
- Create Purchase Request
- Send Chat Message

## Project Structure

```
frontend/
├── app/                    # Next.js app directory
│   ├── customer/          # Customer-facing pages
│   ├── supplier/          # Supplier-facing pages
│   ├── admin/             # Admin pages
│   ├── chat/              # Chat page
│   ├── login/             # Login page
│   └── layout.tsx         # Root layout
├── components/
│   ├── layout/            # Layout components (Header, Sidebar)
│   └── ui/                # UI components (Button, Input, Card)
├── lib/                   # Utilities
│   ├── auth.ts           # Authentication logic
│   ├── api.ts            # API client setup
│   └── tenant.ts         # Tenant resolution
├── tests/
│   └── e2e/              # Playwright E2E tests
└── public/               # Static assets
```

## Key Features

### Tenant Resolution

The app automatically detects the tenant from the subdomain:
- `tenant1.localhost:3002` → tenant1
- `tenant2.example.com` → tenant2
- Falls back to default tenant for localhost development

### Role-Based Routing

Routes are protected based on user roles:
- `/customer/*` - Requires buyer or procurement_manager role
- `/supplier/*` - Requires supplier role
- `/admin/*` - Requires admin or super_admin role

### Authentication

- JWT tokens stored in cookies
- User data stored in localStorage
- Automatic token validation
- Redirect to login on 401 errors

## API Integration

The frontend connects to backend services via environment-configured URLs:
- Identity Service (authentication)
- Company Service
- Catalog Service
- Procurement Service
- Logistics Service
- Collaboration Service
- Notification Service

All API calls include:
- JWT token in Authorization header
- Tenant ID in X-Tenant-ID header

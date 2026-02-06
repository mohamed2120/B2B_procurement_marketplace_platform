/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8001',
  },
  // Suppress router mounting errors in development (known Next.js 14.0.4 issue)
  onDemandEntries: {
    maxInactiveAge: 25 * 1000,
    pagesBufferLength: 2,
  },
  // Disable React strict mode warnings for router
  experimental: {
    appDir: true,
  },
}

module.exports = nextConfig

FROM node:20-alpine

# Build argument for version tracking
ARG SERVICE_VERSION=unknown

# Install wget for health checks
RUN apk add --no-cache wget

# Label for version tracking
LABEL service.version="${SERVICE_VERSION}"

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3002

# Set environment for dev mode
ENV NODE_ENV=development

# Start Next.js in dev mode
CMD ["npm", "run", "dev"]

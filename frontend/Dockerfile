FROM node:20-alpine

# Install wget for health checks
RUN apk add --no-cache wget

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Set environment for dev mode
ENV NODE_ENV=development

# Start Next.js in dev mode
CMD ["npm", "run", "dev"]

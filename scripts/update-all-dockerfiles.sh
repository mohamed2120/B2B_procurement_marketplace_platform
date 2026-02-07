#!/bin/bash

# Update all Dockerfiles to include version labels
# This adds ARG SERVICE_VERSION and LABEL service.version to each Dockerfile

SERVICES=(
    "identity-service"
    "company-service"
    "catalog-service"
    "equipment-service"
    "marketplace-service"
    "procurement-service"
    "logistics-service"
    "collaboration-service"
    "notification-service"
    "billing-service"
    "virtual-warehouse-service"
    "diagnostics-service"
)

for SERVICE in "${SERVICES[@]}"; do
    DOCKERFILE="services/${SERVICE}/Dockerfile"
    
    if [ ! -f "$DOCKERFILE" ]; then
        echo "⚠️  Skipping $SERVICE: Dockerfile not found"
        continue
    fi
    
    echo "Updating $DOCKERFILE..."
    
    # Check if already updated
    if grep -q "ARG SERVICE_VERSION" "$DOCKERFILE" && grep -q "LABEL service.version" "$DOCKERFILE"; then
        echo "  ✅ Already has version labels"
        continue
    fi
    
    # Create backup
    cp "$DOCKERFILE" "${DOCKERFILE}.bak"
    
    # Add ARG in builder stage (if multi-stage)
    if grep -q "FROM.*AS builder" "$DOCKERFILE"; then
        # Add ARG after FROM line in builder stage
        sed -i.bak2 '/FROM.*AS builder/a\
ARG SERVICE_VERSION=unknown
' "$DOCKERFILE"
        rm -f "${DOCKERFILE}.bak2"
    fi
    
    # Add ARG and LABEL in final stage
    # Find the final FROM line
    FINAL_FROM_LINE=$(grep -n "^FROM" "$DOCKERFILE" | tail -1 | cut -d: -f1)
    
    if [ -n "$FINAL_FROM_LINE" ]; then
        # Add ARG after final FROM
        sed -i.bak3 "${FINAL_FROM_LINE}a\\
ARG SERVICE_VERSION=unknown\\
\\
# Add label with service version\\
LABEL service.version=\"\${SERVICE_VERSION}\"
" "$DOCKERFILE"
        rm -f "${DOCKERFILE}.bak3"
    fi
    
    # Clean up backup if successful
    if grep -q "ARG SERVICE_VERSION" "$DOCKERFILE" && grep -q "LABEL service.version" "$DOCKERFILE"; then
        rm -f "${DOCKERFILE}.bak"
        echo "  ✅ Updated successfully"
    else
        echo "  ❌ Update failed, restoring backup"
        mv "${DOCKERFILE}.bak" "$DOCKERFILE"
    fi
done

echo ""
echo "✅ Dockerfile update complete"

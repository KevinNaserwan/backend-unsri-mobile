#!/bin/bash

# Docker Compose Build Partial Script
# Usage: ./scripts/docker-build-partial.sh <service1> [service2] [service3] ...
# Example: ./scripts/docker-build-partial.sh auth-service user-service
# Example: ./scripts/docker-build-partial.sh auth-service --no-cache

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Docker compose file path
COMPOSE_FILE="deployments/docker-compose/docker-compose.yml"

# Check if docker-compose file exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo -e "${RED}Error: Docker compose file not found at $COMPOSE_FILE${NC}"
    exit 1
fi

# Check if services are provided
if [ $# -eq 0 ]; then
    echo -e "${RED}Error: No services specified${NC}"
    echo ""
    echo "Usage: $0 <service1> [service2] [service3] ... [--no-cache]"
    echo ""
    echo "Examples:"
    echo "  $0 auth-service"
    echo "  $0 auth-service user-service api-gateway"
    echo "  $0 auth-service --no-cache"
    echo ""
    echo "Available services:"
    echo "  Infrastructure: postgres, redis, rabbitmq"
    echo "  Core: auth-service, user-service, api-gateway"
    echo "  Academic: course-service, schedule-service, attendance-service, calendar-service"
    echo "  Communication: broadcast-service, notification-service"
    echo "  Location: location-service, access-service, qr-service"
    echo "  Additional: quick-actions-service, file-storage-service, search-service, report-service, master-data-service, leave-service"
    exit 1
fi

# Extract services and build args
SERVICES=()
BUILD_ARGS=()

for arg in "$@"; do
    if [[ "$arg" == "--"* ]]; then
        BUILD_ARGS+=("$arg")
    else
        SERVICES+=("$arg")
    fi
done

# Check if services array is empty (only build args provided)
if [ ${#SERVICES[@]} -eq 0 ]; then
    echo -e "${RED}Error: No services specified (only build arguments provided)${NC}"
    exit 1
fi

# Display what will be built
echo -e "${GREEN}Building services: ${SERVICES[*]}${NC}"
if [ ${#BUILD_ARGS[@]} -gt 0 ]; then
    echo -e "${YELLOW}Build arguments: ${BUILD_ARGS[*]}${NC}"
fi
echo ""

# Build services
docker-compose -f "$COMPOSE_FILE" build "${BUILD_ARGS[@]}" "${SERVICES[@]}"

echo ""
echo -e "${GREEN}✓ Build complete!${NC}"
echo ""
echo "To start the services:"
echo "  docker-compose -f $COMPOSE_FILE up -d ${SERVICES[*]}"
echo ""
echo "Or use make:"
echo "  make docker-up"

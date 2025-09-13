#!/bin/bash
set -e

echo "=== Starting deployment at $(date) ==="

PROJECT_DIR="/opt/sos-kotopes"
BACKUP_DIR="/opt/backups/$(date +%Y%m%d_%H%M%S)"

cd "$PROJECT_DIR"

echo "Building containers..."
docker-compose build

echo "Stopping old services..."
docker-compose down --timeout 30 || true

echo "Starting new services..."
docker-compose up -d --remove-orphans

echo "Waiting for services to start..."
sleep 50

echo "Checking services status..."
if ! docker-compose ps | grep -q "Up"; then
    echo "ERROR: Some services failed to start!"
    docker-compose logs --tail=20
    exit 1
fi

echo "Cleaning up old images..."
docker image prune -f

echo "=== Deployment completed successfully at $(date) ==="
echo "Status:"
docker-compose ps

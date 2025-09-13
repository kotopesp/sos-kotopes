#!/bin/bash
set -e

echo "=== Database restore utility ==="

PROJECT_DIR="/opt/sos-kotopes"
BACKUP_DIR="/opt/backups"

if [ ! -d "$PROJECT_DIR" ]; then
    echo "ERROR: Project directory not found: $PROJECT_DIR"
    exit 1
fi

cd "$PROJECT_DIR" || { echo "ERROR: Cannot cd to $PROJECT_DIR"; exit 1; }

if [ ! -f .env ]; then
    echo "ERROR: .env file not found!"
    exit 1
fi

POSTGRES_USER=$(grep -E '^POSTGRES_USER=' .env | cut -d '=' -f2- | tr -d '"'"'")
POSTGRES_DB=$(grep -E '^POSTGRES_DB=' .env | cut -d '=' -f2- | tr -d '"'"'")
POSTGRES_PASSWORD=$(grep -E '^POSTGRES_PASSWORD=' .env | cut -d '=' -f2- | tr -d '"'"'")


if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_DB" ]; then
    echo "ERROR: PostgreSQL credentials not found in .env"
    exit 1
fi


echo "Available backups:"
ls -la "$BACKUP_DIR"/*.sql.gz 2>/dev/null | head -10 || echo "No backups found"

if [ $# -eq 0 ]; then
    echo ""
    echo "Usage: $0 backup_file.sql.gz"
    echo "Example: $0 soskot_production_backup_20250101_120000.sql.gz"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    if [ ! -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
        echo "ERROR: Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILE"
fi


if [ ! -s "$BACKUP_FILE" ]; then
    echo "ERROR: Backup file is empty: $BACKUP_FILE"
    exit 1
fi

echo "Restoring from: $BACKUP_FILE"
echo "WARNING: This will OVERWRITE the current database!"
echo "Target database: $POSTGRES_DB"
read -p "Are you sure? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Restore cancelled"
    exit 0
fi


if ! docker-compose ps postgres >/dev/null 2>&1; then
    echo "ERROR: PostgreSQL container is not running"
    exit 1
fi

echo "Stopping application..."
docker-compose stop backend || echo "Warning: Could not stop backend"

echo "Restoring database..."
export PGPASSWORD="$POSTGRES_PASSWORD"
if gunzip -c "$BACKUP_FILE" | docker-compose exec -T postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1; then
    unset PGPASSWORD
    echo "Database restore successful"
else
    unset PGPASSWORD
    echo "ERROR: Database restore failed!"
    echo "Starting application back..."
    docker-compose start backend || true
    exit 1
fi

echo "Starting application..."
docker-compose start backend

echo "=== Database restore completed successfully ==="

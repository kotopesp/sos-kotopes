#!/bin/bash
set -e

echo "=== Starting database backup at $(date) ==="

PROJECT_DIR="/opt/sos-kotopes"
BACKUP_DIR="/opt/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

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

if ! mkdir -p "$BACKUP_DIR"; then
    echo "ERROR: Cannot create backup directory: $BACKUP_DIR"
    exit 1
fi

if [ ! -w "$BACKUP_DIR" ]; then
    echo "ERROR: No write permission to backup directory: $BACKUP_DIR"
    exit 1
fi

BACKUP_FILE="$BACKUP_DIR/${POSTGRES_DB}_backup_${TIMESTAMP}.sql.gz"

if ! docker-compose ps postgres | grep -q "Up"; then
    echo "ERROR: PostgreSQL container is not running"
    exit 1
fi

echo "Backing up database: $POSTGRES_DB"

export PGPASSWORD="$POSTGRES_PASSWORD"
if docker-compose exec -T postgres pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" | gzip > "$BACKUP_FILE"; then
    unset PGPASSWORD
    
    if [ ! -s "$BACKUP_FILE" ]; then
        echo "ERROR: Backup file is empty!"
        rm -f "$BACKUP_FILE"
        exit 1
    fi
    
    echo "Backup successful: $BACKUP_FILE"
    echo "File size: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    unset PGPASSWORD
    echo "ERROR: Backup failed!"
    rm -f "$BACKUP_FILE"
    exit 1
fi

echo "=== Database backup completed ==="

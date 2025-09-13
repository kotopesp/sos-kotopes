#!/bin/bash
set -e

echo "=== Cleaning old backups at $(date) ==="

BACKUP_DIR="/opt/backups"
RETENTION_DAYS=7

if [ ! -d "$BACKUP_DIR" ]; then
    echo "Backup directory not found, nothing to clean"
    exit 0
fi

echo "Deleting backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Cleanup completed"

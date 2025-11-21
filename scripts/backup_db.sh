#!/bin/bash

# Database backup script for Tropa Nartov

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/tropa_nartov"
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql"
DB_NAME="tropa_nartov"
DB_USER="postgres"

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Perform backup
echo "Starting backup at $(date)"
pg_dump -U $DB_USER -d $DB_NAME > $BACKUP_FILE

if [ $? -eq 0 ]; then
    echo "✅ Backup successful: $BACKUP_FILE"
    
    # Compress backup
    gzip $BACKUP_FILE
    echo "✅ Backup compressed: $BACKUP_FILE.gz"
    
    # Keep only last 30 days of backups
    find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
    echo "✅ Old backups cleaned up"
    
    # Optional: Upload to S3 or cloud storage
    # aws s3 cp $BACKUP_FILE.gz s3://tropa-nartov-backups/
else
    echo "❌ Backup failed!"
    exit 1
fi

echo "Backup completed at $(date)"


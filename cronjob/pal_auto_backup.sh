#!/bin/bash

# Ref : https://indiefaq.com/guides/palworld-how-to-backup-your-last-5-world-saves-automatically-linux-server.html

# Define source directories and base destination directory
SRC="/home/jason/docker-palworld-dedicated-server/game/Pal/Saved/SaveGames/0/"
BASE_DEST="/home/jason/palworld_backup"

# Create a timestamp
TIMESTAMP=$(date +"%Y%m%d")

# Create new backup directories for this run
DEST="$BASE_DEST/palworld_savegames_$TIMESTAMP"

# if the backup exists, exit
if [ -d "$DEST" ]; then
    echo "Backup already exists for today"
    exit 1
fi

mkdir -p "$DEST"

# Copy files with rsync
rsync -av "$SRC" "$DEST/"

# Function to delete oldest backups if more than 5 exist
function delete_oldest_backup {
    local backup_base_path=$1
    local max_backups=5

    # Create an array of backup directories sorted by date (newest to oldest)
    local backups=($(ls -1d $backup_base_path* | sort -r))

    # Keep only the newest $max_backups directories
    local count=${#backups[@]}
    if [ $count -gt $max_backups ]; then
        local backups_to_delete=(${backups[@]:$max_backups})
        for backup in "${backups_to_delete[@]}"; do
            rm -rf "$backup"
        done
    fi
}

# Apply the backup deletion to each type of backup
delete_oldest_backup "$BASE_DEST/palworld_savegames"

echo "Backup completed on $(date)"

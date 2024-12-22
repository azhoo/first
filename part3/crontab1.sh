#!/bin/bash

# Variables
SCRIPT_PATH="/part2/auto_test.sh"    # Path to your testing script
CRON_TIME="30 14 25 12 *"               # Cron schedule: 2:30 PM on Dec 25th
USER_CRON=$(crontab -l )     # Fetch existing crontab (if any)



# Create the new cron entry
NEW_CRON_ENTRY="$CRON_TIME $SCRIPT_PATH"

# Check if the entry already exists
if echo "$USER_CRON" | grep -Fq "$SCRIPT_PATH"; then
    echo "Crontab entry for $SCRIPT_PATH already exists."
else
  
    (echo "$USER_CRON"; echo "$NEW_CRON_ENTRY") | crontab -
    
    
    echo "Crontab entry added: $NEW_CRON_ENTRY"
 
fi

# Verify crontab
echo "Updated crontab:"
crontab -l

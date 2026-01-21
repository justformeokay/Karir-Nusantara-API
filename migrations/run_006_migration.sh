#!/bin/bash

# Script untuk menjalankan migration password reset tokens
# Untuk XAMPP MySQL

echo "Running migration: 006_password_reset_tokens.sql"
echo "----------------------------------------"

# Path to migration file
MIGRATION_FILE="/Users/putramac/Desktop/Loker/karir-nusantara-api/migrations/006_password_reset_tokens.sql"

# Database details
DB_NAME="karir_nusantara"

# Check if XAMPP MySQL is running
if ! /Applications/XAMPP/xamppfiles/bin/mysql -V &> /dev/null; then
    echo "Error: MySQL (XAMPP) is not accessible"
    exit 1
fi

# Run migration
/Applications/XAMPP/xamppfiles/bin/mysql -u root "$DB_NAME" < "$MIGRATION_FILE"

if [ $? -eq 0 ]; then
    echo "✓ Migration completed successfully!"
    echo "Table 'password_reset_tokens' has been created."
else
    echo "✗ Migration failed!"
    exit 1
fi

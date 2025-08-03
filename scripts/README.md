# Wenxi Network Disk - Database Management Scripts

Author: Wenxi

This directory contains database management scripts for Wenxi Network Disk.

## Script Descriptions

### `init_db.bat`
**Function:** Initialize a new database and clear all existing user data and stored files.

**Warning:** This script will:
- Delete all existing user accounts
- Delete all stored files
- Reset the entire database to a clean state

**Usage:**
1. Run the script: `scripts\init_db.bat`
2. When prompted, type `YES` to confirm the operation
3. Any other input will cancel the operation

### `force_init_db.bat`
**Function:** Force complete reset of the database system, deleting ALL data including user accounts and stored files.

**Warning:** This script is more destructive than `init_db.bat` and will:
- **PERMANENTLY DELETE** all user accounts
- **PERMANENTLY DELETE** all stored files
- **COMPLETELY RESET** the entire system
- **IRREVERSIBLE** operation - data cannot be recovered

**Usage:**
1. Run the script: `scripts\force_init_db.bat`
2. When prompted, type `YES` to confirm the operation
3. Any other input will cancel the operation

## Safety Features

Both scripts include mandatory confirmation:
- **YES** (case-insensitive) - Proceed with operation
- **Any other input** - Cancel operation safely

## Requirements

- Python 3.8+ installed and in PATH
- All dependencies listed in `requirements.txt`
- Windows PowerShell or Command Prompt

## Emergency Recovery

If you accidentally run these scripts:
1. **STOP** the operation immediately by closing the terminal
2. **DO NOT** restart the application
3. **CHECK** if any backup files exist in the project directory
4. **CONTACT** system administrator for data recovery assistance

## Testing

Run the confirmation script tests:
```bash
python -m pytest tests/test_confirmation_scripts.py -v
```
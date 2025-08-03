@echo off
REM Wenxi NetDisk - Force Database Initialization Script
REM This script will completely reset the database and clear all stored data

echo [Wenxi] ================================================
echo [Wenxi] WARNING: FORCE DATABASE INITIALIZATION
echo [Wenxi] ================================================
echo [Wenxi] This will DELETE all user data and stored files!
echo [Wenxi] This action is IRREVERSIBLE and will result in:
echo [Wenxi] - Complete loss of all user accounts
echo [Wenxi] - Complete loss of all stored files
echo [Wenxi] - Complete reset of the entire system
set /p confirm="[Wenxi] Are you absolutely sure? Type YES to confirm: "
if /i "%confirm%" neq "YES" (
    echo [Wenxi] Operation cancelled by user.
    pause
    exit /b 1
)
echo [Wenxi] Continuing with force initialization...
echo.

REM Check if Python is available
python --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Python is not installed or not in PATH
    pause
    exit /b 1
)

REM Change to project root directory
cd /d "%~dp0\.."
echo [Wenxi] Working directory: %cd%

REM Remove existing database file if it exists (backend directory)
if exist "backend\wenxi_netdisk.db" (
    echo [Wenxi] Removing existing database file from backend directory...
    del /f /q "backend\wenxi_netdisk.db"
    echo [Wenxi] Database file removed successfully
) else (
    echo [Wenxi] No existing database file found in backend directory
)

REM Also check project root for compatibility
if exist "wenxi_netdisk.db" (
    echo [Wenxi] Removing old database file from project root...
    del /f /q "wenxi_netdisk.db"
    echo [Wenxi] Old database file removed successfully
)

REM Get storage path dynamically
echo [Wenxi] Getting storage path configuration...
for /f "tokens=2 delims==" %%i in ('python scripts\get_storage_paths.py ^| find "STORAGE_RELATIVE="') do set STORAGE_PATH=%%i
echo [Wenxi] Storage path: %STORAGE_PATH%

REM Remove uploads directory with all user files
echo [Wenxi] Removing uploads directory with all user files...
if exist "%STORAGE_PATH%" (
    rmdir /s /q "%STORAGE_PATH%"
    echo [Wenxi] Storage directory removed successfully
) else (
    echo [Wenxi] No storage directory found at %STORAGE_PATH%
)

REM Run database initialization
echo [Wenxi] Initializing new database...
python -c "import os, sys; sys.path.insert(0, 'backend'); from database import init_db; init_db(); print('[Wenxi] Database initialization completed successfully')"

if errorlevel 1 (
    echo [ERROR] Database initialization failed
    pause
    exit /b 1
)

echo.
echo [SUCCESS] Database has been completely reset and initialized
echo [SUCCESS] All stored data has been cleared
echo [SUCCESS] Uploads directory has been removed
echo.
pause
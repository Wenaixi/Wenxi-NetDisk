@echo off
echo ================================================
echo Wenxi Database Initialization - Files Reset
echo ================================================
echo.
echo [WARNING] This will DELETE ALL USER FILES including:
echo [WARNING] - All uploaded files and folders
echo [WARNING] - All file metadata and shared links
echo [WARNING] - User accounts will be PRESERVED
echo [WARNING] This action is IRREVERSIBLE!
echo.
set /p confirm="Are you absolutely sure? Type YES to confirm: "
if /i "%confirm%" neq "YES" (
    echo Operation cancelled by user.
    pause
    exit /b 1
)
echo.

REM Set working directory to project root
cd /d "%~dp0\.."
echo [Wenxi] Working directory: %cd%

REM Check Python environment
echo [1/4] Checking Python environment...
python --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Python is not installed or not in PATH
    echo Please install Python 3.8+ and add it to your PATH
    pause
    exit /b 1
)

REM Clear files table in database
echo [2/4] Clearing files table...
python -c "import sys; sys.path.insert(0, 'backend'); from database import get_db; from sqlalchemy.orm import Session; from sqlalchemy import text; db = next(get_db()); db.execute(text('DELETE FROM files')); db.commit(); db.close(); print('[SUCCESS] Files table cleared successfully')"

REM Remove uploads directory with all user files
echo [3/4] Removing all user uploaded files...

REM Get storage path dynamically and remove uploads directory
echo [3.1/4] Getting storage path configuration...
for /f "tokens=2 delims==" %%i in ('python scripts\get_storage_paths.py ^| find "STORAGE_RELATIVE="') do set STORAGE_PATH=%%i
echo [INFO] Storage path: %STORAGE_PATH%

if exist "%STORAGE_PATH%" (
    rmdir /s /q "%STORAGE_PATH%"
    echo [SUCCESS] Storage directory removed with all user files
) else (
    echo [INFO] No storage directory found at %STORAGE_PATH%
)

REM Also check project root uploads directory for compatibility
echo [3.2/4] Checking project root uploads directory...
if exist "uploads" (
    rmdir /s /q "uploads"
    echo [SUCCESS] Project root uploads directory removed
) else (
    echo [INFO] No project root uploads directory found
)

REM Install dependencies
echo [4/4] Installing Python dependencies...
if exist "requirements.txt" (
    pip install -r requirements.txt
    if errorlevel 1 (
        echo [ERROR] Failed to install dependencies
        pause
        exit /b 1
    )
) else (
    echo [WARNING] requirements.txt not found
)

REM Initialize clean state
echo [Wenxi] Initializing clean file storage...
python -c "import sys; sys.path.insert(0, 'backend'); from database import init_db; init_db()"

if errorlevel 1 (
    echo [ERROR] Database initialization failed
    echo Please check the error message above
    pause
    exit /b 1
) else (
    echo.
    echo ================================================
    echo [SUCCESS] Files reset completed successfully!
    echo [SUCCESS] All user files have been cleared
    echo [SUCCESS] User accounts are preserved
    echo [SUCCESS] System is ready for new uploads
    echo ================================================
)

pause
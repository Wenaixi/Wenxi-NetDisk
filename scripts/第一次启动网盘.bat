@echo off
REM Wenxi NetDisk - First Run Setup Script
REM This script checks and installs all required dependencies for Wenxi NetDisk

echo ================================================
echo Wenxi NetDisk - First Time Setup
echo ================================================
echo.

REM Check environment variables configuration first
set "PROJECT_ROOT=%~dp0.."
cd /d "%PROJECT_ROOT%"

if not exist ".env" (
    echo [ERROR] Environment configuration file not found!
    echo.
    echo Please complete the following steps:
    echo 1. Copy .env.example to .env
    echo 2. Configure the environment variables in .env
    echo 3. Run this script again
    echo.
    echo Location: %CD%\.env.example
    pause
    exit /b 1
)
echo [SUCCESS] Environment configuration file found

REM Set console to UTF-8 encoding
chcp 65001 > nul
echo Console encoding set to UTF-8

REM Save current directory
set "CURRENT_DIR=%CD%"
echo Current directory: %CURRENT_DIR%

REM Check Python installation
echo.
echo [1/5] Checking Python environment...
python --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Python is not installed or not in PATH
    echo Please install Python 3.8+ from https://python.org
    pause
    exit /b 1
)

REM Get Python version
for /f "tokens=2" %%i in ('python --version') do set PYTHON_VERSION=%%i
echo [SUCCESS] Python version: %PYTHON_VERSION%

REM Check Node.js installation
echo.
echo [2/5] Checking Node.js environment...
node --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Node.js is not installed or not in PATH
    echo Please install Node.js 16+ from https://nodejs.org
    pause
    exit /b 1
)

REM Get Node.js version
for /f %%i in ('node --version') do set NODE_VERSION=%%i
echo [SUCCESS] Node.js version: %NODE_VERSION%

REM Check pip installation
echo.
echo [3/5] Checking pip environment...
python -m pip --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] pip is not installed
    echo Installing pip...
    python -m ensurepip --upgrade
)
echo [SUCCESS] pip is available

REM Choose pip source
echo.
echo [4/5] Configuring pip source...
echo Please select a pip mirror for faster downloads:
echo 1. Alibaba Cloud (https://mirrors.aliyun.com/pypi/simple/) [Recommended]
echo 2. Tsinghua University (https://pypi.tuna.tsinghua.edu.cn/simple)
echo 3. Douban (https://pypi.douban.com/simple)
echo 4. Official (https://pypi.org/simple)
set /p PIP_SOURCE="Enter your choice [1-4, default 1]: "

if "%PIP_SOURCE%"=="" set PIP_SOURCE=1
if "%PIP_SOURCE%"=="1" set PIP_INDEX_URL=https://mirrors.aliyun.com/pypi/simple/
if "%PIP_SOURCE%"=="2" set PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple
if "%PIP_SOURCE%"=="3" set PIP_INDEX_URL=https://pypi.douban.com/simple
if "%PIP_SOURCE%"=="4" set PIP_INDEX_URL=https://pypi.org/simple

echo Selected pip source: %PIP_INDEX_URL%

REM Install Python dependencies
echo.
echo [5/5] Installing Python dependencies...
cd /d "%~dp0\..\backend"
if exist "requirements.txt" (
    echo Installing backend dependencies...
    python -m pip install -i %PIP_INDEX_URL% -r requirements.txt --upgrade
    if errorlevel 1 (
        echo [ERROR] Failed to install backend dependencies
        pause
        exit /b 1
    )
    echo [SUCCESS] Backend dependencies installed
) else (
    echo [WARNING] backend/requirements.txt not found
)

REM Install Node.js dependencies
echo.
echo Installing Node.js dependencies...
cd /d "%~dp0\..\frontend"
if exist "package.json" (
    echo Installing frontend dependencies...
    call npm install
    if errorlevel 1 (
        echo [ERROR] Failed to install frontend dependencies
        pause
        exit /b 1
    )
    echo [SUCCESS] Frontend dependencies installed
) else (
    echo [WARNING] frontend/package.json not found
)

REM Initialize database
echo.
echo Initializing database...
cd /d "%~dp0\..\backend"
python -c "import sys; sys.path.insert(0, '.'); from database import init_db; init_db(); print('[SUCCESS] Database initialized')"

REM Final success message
echo.
echo ================================================
echo [SUCCESS] Wenxi NetDisk setup completed!
echo.
echo You can now run: scripts\start_netdisk.bat
echo to start the application.
echo.
pause
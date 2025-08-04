@echo off
REM Wenxi NetDisk - Start Services Script
REM This script checks environment and starts Wenxi NetDisk services

echo ================================================
echo Wenxi NetDisk - Starting Services
echo ================================================
echo.

REM Set console to UTF-8 encoding
chcp 65001 > nul
echo Console encoding set to UTF-8

REM Set project root directory (parent of scripts)
set "PROJECT_ROOT=%~dp0.."
echo Project root: %PROJECT_ROOT%

REM Check Python installation
echo.
echo [1/4] Checking Python environment...
python --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Python is not installed or not in PATH
    echo Please run scripts\first_run_setup.bat first
    pause
    exit /b 1
)

REM Check Node.js installation
echo.
echo [2/4] Checking Node.js environment...
node --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Node.js is not installed or not in PATH
    echo Please run scripts\first_run_setup.bat first
    pause
    exit /b 1
)

REM Check backend dependencies
echo.
echo [3/4] Checking backend dependencies...
cd /d "%PROJECT_ROOT%\backend"
python -c "import fastapi; import sqlalchemy; print('[SUCCESS] Backend dependencies available')" >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Backend dependencies not found
    echo Please run scripts\first_run_setup.bat first
    pause
    exit /b 1
)

REM Check frontend dependencies
echo.
echo [4/4] Checking frontend dependencies...
cd /d "%PROJECT_ROOT%\frontend"
if not exist "node_modules" (
    echo [ERROR] Frontend dependencies not found
    echo Please run scripts\first_run_setup.bat first
    pause
    exit /b 1
)

REM Kill existing processes on ports
echo.
echo Stopping existing services...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3008') do taskkill /F /PID %%a 2>nul
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do taskkill /F /PID %%a 2>nul
timeout /t 2 /nobreak >nul

REM Start backend service
echo.
echo Starting backend service...
cd /d "%PROJECT_ROOT%\backend"
start "Wenxi Backend" cmd /k "cd /d %PROJECT_ROOT%\backend && echo Starting backend on port 3008... && python main.py"

REM Wait for backend to initialize
echo Waiting for backend to start...
timeout /t 3 /nobreak >nul

REM Start frontend service
echo.
echo Starting frontend service...
cd /d "%PROJECT_ROOT%\frontend"
start "Wenxi Frontend" cmd /k "cd /d %PROJECT_ROOT%\frontend && echo Starting frontend on port 5173... && npm run dev"

REM Display startup information
echo.
echo ================================================
echo Wenxi NetDisk services are starting...
echo.
echo Backend: http://localhost:3008
    echo Frontend: http://localhost:5173
    echo API Documentation: http://localhost:3008/docs
    echo Health Check: http://localhost:3008/health
    echo.
    echo Press any key to continue...
    pause > nul
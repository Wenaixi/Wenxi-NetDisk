@echo off
echo ================================================
echo Wenxi Quick Start - One Click Service Launcher
echo ================================================
echo.

REM Kill processes on ports 3008 and 5173
echo Stopping services on ports 3008 and 5173...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3008') do taskkill /F /PID %%a 2>nul
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do taskkill /F /PID %%a 2>nul

REM Wait a moment for cleanup
timeout /t 2 /nobreak >nul

echo.
echo Starting Wenxi NetDisk services...
echo.

REM Check and install root dependencies
echo Checking root Python dependencies...
cd /d %~dp0\..
if exist "requirements.txt" (
    python -c "import sys; sys.path.append('.'); import pkg_resources; pkg_resources.require(open('requirements.txt').read())" 2>nul
    if errorlevel 1 (
        echo Installing root Python dependencies...
        pip install -r requirements.txt
    ) else (
        echo Root Python dependencies already installed, skipping...
    )
) else (
    echo No root requirements.txt found, skipping root Python dependencies...
)

REM Check and install backend dependencies
echo Checking backend dependencies...
cd /d %~dp0\..\backend
if exist "requirements.txt" (
    python -c "import sys; sys.path.append('.'); import pkg_resources; pkg_resources.require(open('requirements.txt').read())" 2>nul
    if errorlevel 1 (
        echo Installing backend Python dependencies...
        pip install -r requirements.txt
    ) else (
        echo Backend Python dependencies already installed, skipping...
    )
) else (
    echo No backend requirements.txt found, skipping backend Python dependencies...
)

REM Initialize database
echo Initializing database...
cd /d %~dp0\..\backend
python -c "from database import init_db; init_db()" 2>nul
if errorlevel 1 (
    echo Warning: Database initialization failed, continuing anyway...
)

REM Start backend in background
start "Wenxi Backend" cmd /k "cd /d %~dp0\..\backend && echo Starting backend... && python main.py"

REM Wait for backend to initialize
echo Waiting for backend to start...
timeout /t 5 /nobreak >nul

echo.
echo Starting frontend...

REM Check and install frontend dependencies if needed
cd /d %~dp0\..\frontend
if exist "package.json" (
    if exist "node_modules" (
        echo Node.js dependencies already installed, skipping...
    ) else (
        echo Installing Node.js dependencies...
        npm install
    )
) else (
    echo No package.json found, skipping Node.js dependencies...
)

REM Start frontend in background
start "Wenxi Frontend" cmd /k "cd /d %~dp0\..\frontend && echo Starting dev server... && npm run dev"

echo.
echo ================================================
echo Services starting...
echo Backend: http://localhost:3008 (starting...)
echo Frontend: http://localhost:5173 (starting...)
echo ================================================
echo.
echo If login still hangs, run: scripts\diagnose.bat
echo Press any key to close this window (services will continue)
pause >nul
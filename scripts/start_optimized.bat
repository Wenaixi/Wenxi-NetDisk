@echo off
echo Wenxi Cloud Storage Quick Optimization Startup Script
echo.
echo Checking Redis service...
tasklist /FI "IMAGENAME eq redis-server.exe" 2>NUL | find /I /N "redis-server.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo Redis service is already running
) else (
    echo Starting Redis service...
    start "" "C:\Program Files\Redis\redis-server.exe" --maxmemory 512mb
    timeout /t 3 /nobreak >NUL
)

echo.
echo Starting backend service...
cd backend
start cmd /k "python -m uvicorn main:app --reload --host 0.0.0.0 --port 8000"

timeout /t 3 /nobreak >NUL

echo.
echo Starting frontend service...
cd ..\frontend
start cmd /k "npm run dev"

echo.
echo Wenxi Cloud Storage has been started!
echo Access URL: http://localhost:5173
echo Backend API: http://localhost:8000
echo.
echo Press any key to exit...
pause >NUL
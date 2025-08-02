@echo off
echo ================================================
echo Wenxi NetDiagnosis - Network & Service Checker
echo ================================================
echo.

echo [1] Checking backend service on port 3008...
curl -s -o nul -w "%%{http_code}" http://localhost:3008/health > temp.txt
set /p status=<temp.txt
del temp.txt

if "%status%"=="200" (
    echo ✅ Backend service is running properly
) else (
    echo ❌ Backend service is not responding (port 3008)
    echo    Please check if backend is running: cd backend && python main.py
    echo.
)

echo.
echo [2] Checking frontend service on port 5173...
curl -s -o nul -w "%%{http_code}" http://localhost:5173 > temp.txt
set /p frontend=<temp.txt
del temp.txt

if "%frontend%"=="200" (
    echo ✅ Frontend service is running properly
) else (
    echo ❌ Frontend service is not responding (port 5173)
    echo    Please check if frontend is running: cd frontend && npm run dev
    echo.
)

echo.
echo [3] Testing API endpoints...
echo Testing /api/auth/login endpoint...
curl -s -X POST http://localhost:3008/api/auth/login -F "username=test" -F "password=test" > temp.txt
find "401" temp.txt > nul
if errorlevel 1 (
    echo ❌ API endpoint not accessible
) else (
    echo ✅ API endpoint is accessible
)
del temp.txt

echo.
echo [4] Checking Node.js and Python installations...
node --version
python --version

echo.
echo ================================================
echo Diagnosis complete. Check above results.
echo ================================================
pause
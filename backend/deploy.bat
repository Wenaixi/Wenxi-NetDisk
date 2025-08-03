@echo off
echo =========================================
echo Wenxi Cloud Storage - One-click Deployment Script
echo Author: Wenxi
echo =========================================

echo 1. Checking Python environment...
python --version
if errorlevel 1 (
    echo ❌ Python is not installed or not in PATH
    pause
    exit /b 1
)

echo 2. Installing dependencies...
pip install -r requirements.txt
if errorlevel 1 (
    echo ❌ Dependency installation failed
    pause
    exit /b 1
)

echo 3. Running encryption module tests...
python -m pytest tests/test_encryption.py -v
if errorlevel 1 (
    echo ❌ Encryption module tests failed
    pause
    exit /b 1
)

echo 4. Running end-to-end tests...
python test_end_to_end.py
if errorlevel 1 (
    echo ❌ End-to-end tests failed
    pause
    exit /b 1
)

echo 5. Cleaning up temporary files...
del /f /q *.tmp 2>nul
del /f /q uploads\*.tmp 2>nul

echo =========================================
echo ✅ Deployment completed! Wenxi Cloud Storage is ready
echo ✅ New Features:
echo    - AES CTR mode encryption with unchanged file size
echo    - Encrypted files without extensions, unified format
echo    - MP4 files work normally after decryption
echo =========================================
pause
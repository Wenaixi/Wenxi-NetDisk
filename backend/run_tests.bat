@echo off
echo Starting Wenxi Network Disk Test Suite...
echo =========================================

echo Running encryption tests...
cd /d %~dp0
python -m pytest tests/test_encryption.py -v

echo Running file upload tests...
python -m pytest tests/test_file_upload.py -v

echo Running database tests...
python -m pytest tests/test_database_config.py -v

echo Running uploads directory tests...
python -m pytest tests/test_uploads_dir.py -v

echo All tests completed!
pause
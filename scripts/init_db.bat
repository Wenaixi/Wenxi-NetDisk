@echo off
echo Initializing Wenxi NetDisk Database...
cd /d %~dp0\..
python -c "from backend.database import init_db; init_db()"
echo Database initialized successfully!
pause
@echo off
chcp 65001 >nul
echo === Wenxi NetDisk Git Helper ===
echo.

if "%~1"=="" goto :menu
if "%~1"=="status" goto :status
if "%~1"=="diff" goto :diff
if "%~1"=="log" goto :log
if "%~1"=="review" goto :review
goto :menu

:menu
echo 用法: git-helper.bat [命令]
echo.
echo 可用命令:
echo   status  - 查看 git 状态
echo   diff    - 查看变更
echo   log     - 查看提交历史
echo   review  - 代码审查准备
echo.
pause
goto :eof

:status
echo [Git Status]
git status --short
goto :eof

:diff
echo [Git Diff Statistics]
git diff --stat
echo.
echo [Detailed Diff]
git diff
goto :eof

:log
echo [Recent Commits]
git log --oneline -20
goto :eof

:review
echo [Preparing for Code Review...]
echo 未提交的变更:
git diff --stat
echo.
echo 已暂存的变更:
git diff --cached --stat
goto :eof

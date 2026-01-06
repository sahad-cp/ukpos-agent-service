@echo off
SET SCRIPT_DIR=%~dp0
SET NSSM=%SCRIPT_DIR%nssm.exe
SET SERVICE_NAME=UKPOSPrintAgent

net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Please run as Administrator.
    pause
    exit /b 1
)

"%NSSM%" stop %SERVICE_NAME%
"%NSSM%" remove %SERVICE_NAME% confirm

echo Service removed.
pause
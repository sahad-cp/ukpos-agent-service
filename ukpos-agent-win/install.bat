@echo off
SETLOCAL ENABLEEXTENSIONS

echo ========================================
echo   UKPOS Print Agent - Windows Installer
echo ========================================
echo.

:: Require Administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ERROR: Please right-click and Run as Administrator.
    pause
    exit /b 1
)

:: Resolve directory of this script
SET SCRIPT_DIR=%~dp0
SET NSSM=%SCRIPT_DIR%nssm.exe

SET INSTALL_DIR=%USERPROFILE%\ukpos-agent-runtime
SET SERVICE_NAME=UKPOSPrintAgent

echo [1/5] Creating directories...
mkdir "%INSTALL_DIR%\logs" >nul 2>&1

echo [2/5] Copying files...
copy "%SCRIPT_DIR%ukpos-agent.exe" "%INSTALL_DIR%\" >nul
copy "%SCRIPT_DIR%config.json" "%INSTALL_DIR%\" >nul

echo [3/5] Installing Windows Service...
"%NSSM%" install %SERVICE_NAME% "%INSTALL_DIR%\ukpos-agent.exe"
"%NSSM%" set %SERVICE_NAME% AppDirectory "%INSTALL_DIR%"
"%NSSM%" set %SERVICE_NAME% Start SERVICE_AUTO_START
"%NSSM%" set %SERVICE_NAME% AppStdout "%INSTALL_DIR%\logs\stdout.log"
"%NSSM%" set %SERVICE_NAME% AppStderr "%INSTALL_DIR%\logs\stderr.log"

echo [4/5] Starting service...
"%NSSM%" start %SERVICE_NAME%

echo [5/5] Installation complete.
echo Logs: %INSTALL_DIR%\logs
pause
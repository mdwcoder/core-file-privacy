@echo off
setlocal EnableDelayedExpansion

set "REPO_URL=https://github.com/mdwcoder/core-file-privacy"
set "INSTALL_DIR=%USERPROFILE%\bin"
set "BINARY_NAME=cfp.exe"

for /f "usebackq tokens=*" %%a in (`powershell -NoProfile -Command "[System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture.ToString().ToLower()"`) do (
    set "ARCH=%%a"
)

if "%ARCH%"=="x64" set "ARCH=amd64"
if "%ARCH%"=="arm64" set "ARCH=arm64"

set "ASSET_NAME=cfp_windows_%ARCH%.zip"
set "URL=%REPO_URL%/releases/latest/download/%ASSET_NAME%"

echo Downloading precompiled binary: %URL%
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

set "TMP_DIR=%TEMP%\cfp_install_%RANDOM%"
mkdir "%TMP_DIR%"

powershell -NoProfile -Command "try { Invoke-WebRequest -Uri '%URL%' -OutFile '%TMP_DIR%\%ASSET_NAME%' -UseBasicParsing } catch { exit 1 }"
if errorlevel 1 (
    echo Failed to download binary. Please install Go or download manually from %REPO_URL%/releases
    rmdir /s /q "%TMP_DIR%"
    exit /b 1
)

powershell -NoProfile -Command "Expand-Archive -Path '%TMP_DIR%\%ASSET_NAME%' -DestinationPath '%TMP_DIR%' -Force"

if exist "%TMP_DIR%\cfp.exe" (
    copy /Y "%TMP_DIR%\cfp.exe" "%INSTALL_DIR%\cfp.exe" >nul
) else (
    echo Binary not found in downloaded archive.
    rmdir /s /q "%TMP_DIR%"
    exit /b 1
)

rmdir /s /q "%TMP_DIR%"

echo Installed: %INSTALL_DIR%\cfp.exe

if not exist "%INSTALL_DIR%\cfp.exe" (
    echo ERROR: Installation failed.
    exit /b 1
)

echo.
echo core-file-privacy installed successfully!
echo Run 'cfp --help' to get started.
echo.
echo NOTE: Make sure %INSTALL_DIR% is in your PATH.

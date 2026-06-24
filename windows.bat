@echo off
setlocal enabledelayedexpansion

set "APP_NAME=tex.exe"
set "SRC_DIR=.\src"
set "BUILD_DIR=.\build"
set "INSTALL_PATH=C:\tex"

echo 📦 Building %APP_NAME%...

if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
)

cd "%SRC_DIR%"
call go mod tidy
go build -o "..\%BUILD_DIR%\%APP_NAME%" .
cd ..

if %ERRORLEVEL% neq 0 (
    echo ❌ Build failed!
    exit /b %ERRORLEVEL%
)

if not exist "%INSTALL_PATH%" (
    echo 📂 Creating directory %INSTALL_PATH%...
    mkdir "%INSTALL_PATH%"
)

echo 🚚 Installing to %INSTALL_PATH%...
copy /y "%BUILD_DIR%\%APP_NAME%" "%INSTALL_PATH%\%APP_NAME%" >nul

:: Safe path injection using embedded PowerShell
echo ⚙️ Adding %INSTALL_PATH% to User PATH safely...
powershell -NoProfile -Command ^
    "$oldPath = [Environment]::GetEnvironmentVariable('Path', 'User');" ^
    "if ($oldPath -notlike '*%INSTALL_PATH%*') {" ^
    "    $newPath = if ($oldPath) { \"$oldPath;%INSTALL_PATH%\" } else { '%INSTALL_PATH%' };" ^
    "    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User');" ^
    "}"

echo.
echo ✅ Installed successfully!
echo.
echo 👉 IMPORTANT: Close this terminal and open a BRAND NEW terminal window.
echo    Then type: tex

endlocal
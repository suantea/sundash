@echo off
chcp 65001 >nul
echo ================================
echo  SunDash Build ^& Export
echo ================================

cd /d "%~dp0"

set IMAGE_NAME=sundash
set TAG=latest
set EXPORT_FILE=sundash-image.tar.gz

echo.
echo [1/4] Stopping old container...
docker stop sundash 2>nul
docker rm sundash 2>nul

echo.
echo [2/4] Building new image...
docker build -f docker/Dockerfile -t %IMAGE_NAME%:%TAG% .

if errorlevel 1 (
    echo.
    echo [ERROR] Build failed!
    pause
    exit /b 1
)

echo.
echo [3/4] Exporting image to %EXPORT_FILE%...
docker save %IMAGE_NAME%:%TAG% | gzip > %EXPORT_FILE%

if errorlevel 1 (
    echo.
    echo [ERROR] Export failed!
    pause
    exit /b 1
)

echo.
echo [4/4] Starting container...
docker run -d --name sundash -p 3000:3000 -v sundash-data:/app/data -e SUNDASH_PORT=3000 -e SUNDASH_DATA_DIR=/app/data -e TZ=Asia/Shanghai %IMAGE_NAME%:%TAG%

echo.
echo ================================
echo  Build ^& Export Complete!
echo ================================
echo.
echo  Local access: http://localhost:3000
echo  Image file: %EXPORT_FILE%
echo.
echo  NAS Deployment:
echo   1. Copy %EXPORT_FILE% to NAS
echo   2. Run start.bat or start.sh
echo.
pause
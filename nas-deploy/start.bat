@echo off
chcp 65001 >nul
echo ================================
echo  SunDash NAS Deploy
echo ================================

cd /d "%~dp0"

if not exist data mkdir data

echo.
echo [1/3] Loading image...
docker load ^< sundash-image.tar.gz

if errorlevel 1 (
    echo [ERROR] Load failed!
    pause
    exit /b 1
)

echo.
echo [2/3] Stopping old container...
docker stop sundash 2>nul
docker rm sundash 2>nul

echo.
echo [3/3] Starting container...
docker run -d --name sundash --restart unless-stopped -p 3000:3000 -v %cd%\data:/app/data -e SUNDASH_PORT=3000 -e SUNDASH_DATA_DIR=/app/data -e TZ=Asia/Shanghai sundash:latest

echo.
echo ================================
echo  Deploy Complete!
echo ================================
echo.
echo  Access: http://YOUR_NAS_IP:3000
echo  Default account: admin / admin
echo.
pause
@echo off
chcp 65001 >nul
echo ================================
echo  SunDash NAS Deploy
echo ================================

cd /d "%~dp0"

if not exist data mkdir data

rem Ensure JWT secret (required by the server; auto-generate if missing)
if exist .env (
    findstr /b "SUNDASH_JWT_SECRET=" .env >nul 2>&1
    if errorlevel 1 goto gen_secret
    for /f "tokens=1,* delims==" %%a in ('findstr /b "SUNDASH_JWT_SECRET=" .env') do if "%%b"=="" goto gen_secret
    echo [OK] SUNDASH_JWT_SECRET found in .env
    goto load_image
)
:gen_secret
set "SECRET="
for /f "delims=" %%i in ('powershell -NoProfile -Command "[Convert]::ToBase64String((1..48 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))"') do set "SECRET=%%i"
echo SUNDASH_JWT_SECRET=%SECRET%> .env
echo [INFO] Generated SUNDASH_JWT_SECRET into .env

:load_image
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
docker run -d --name sundash --restart unless-stopped -p 3000:3000 -v %cd%\data:/app/data --env-file .env -e SUNDASH_PORT=3000 -e SUNDASH_DATA_DIR=/app/data -e TZ=Asia/Shanghai sundash:latest

echo.
echo ================================
echo  Deploy Complete!
echo ================================
echo.
echo  Access: http://YOUR_NAS_IP:3000
echo  Default account: admin / admin
echo.
pause
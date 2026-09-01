#!/bin/bash
echo "================================"
echo "  SunDash NAS Deploy"
echo "================================"

# Create data directory
mkdir -p ./data

# Ensure JWT secret (required by the server; auto-generate if missing)
if [ -f .env ] && grep -q "^SUNDASH_JWT_SECRET=" .env && [ -n "$(grep '^SUNDASH_JWT_SECRET=' .env | cut -d= -f2-)" ]; then
    echo "[OK] SUNDASH_JWT_SECRET found in .env"
else
    SECRET="$(openssl rand -base64 48 2>/dev/null || head -c 48 /dev/urandom | base64)"
    echo "SUNDASH_JWT_SECRET=${SECRET}" > .env
    echo "[INFO] Generated SUNDASH_JWT_SECRET into .env"
fi

# Load image
echo ""
echo "[1/3] Loading image..."
docker load < sundash-image.tar.gz

if [ $? -ne 0 ]; then
    echo "[ERROR] Load failed!"
    exit 1
fi

# Stop old container
echo ""
echo "[2/3] Stopping old container..."
docker stop sundash 2>/dev/null
docker rm sundash 2>/dev/null

# Start new container (env-file passes SUNDASH_JWT_SECRET from .env)
echo ""
echo "[3/3] Starting container..."
docker run -d --name sundash --restart unless-stopped -p 3000:3000 \
  -v "$(pwd)/data:/app/data" \
  --env-file .env \
  -e SUNDASH_PORT=3000 -e SUNDASH_DATA_DIR=/app/data -e TZ=Asia/Shanghai \
  sundash:latest

echo ""
echo "================================"
echo "  Deploy Complete!"
echo "================================"
echo ""
echo "  Access: http://YOUR_NAS_IP:3000"
echo "  Default account: admin / admin"
echo ""
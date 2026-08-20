#!/bin/bash
echo "================================"
echo "  SunDash NAS Deploy"
echo "================================"

# Create data directory
mkdir -p ./data

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

# Start new container
echo ""
echo "[3/3] Starting container..."
docker run -d --name sundash --restart unless-stopped -p 3000:3000 -v "$(pwd)/data:/app/data" -e SUNDASH_PORT=3000 -e SUNDASH_DATA_DIR=/app/data -e TZ=Asia/Shanghai sundash:latest

echo ""
echo "================================"
echo "  Deploy Complete!"
echo "================================"
echo ""
echo "  Access: http://YOUR_NAS_IP:3000"
echo "  Default account: admin / admin"
echo ""
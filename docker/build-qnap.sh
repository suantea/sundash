#!/bin/bash
# ============================================================
# 构建 sundash 的 QNAP (linux/amd64) 可运行 Docker 镜像
# 用法: ./docker/build-qnap.sh  [镜像标签]
# 默认标签: sundash:latest
#
# 说明:
#  - 宿主机 arm64 (M 系 Mac) 需用 buildx 交叉编译 linux/amd64
#  - 产物: qnap/sundash-image.tar.gz (供 NAS 上 docker load)
#  - QNAP 部署: 见 nas-deploy/start.sh / start.bat
# ============================================================
set -euo pipefail

# 项目根目录
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

TAG="${1:-sundash:latest}"
OUT_DIR="${ROOT}/qnap"
OUT_TAR="${OUT_DIR}/sundash-image.tar.gz"

echo "======================"
echo "  构建平台: linux/amd64"
echo "  目标镜像: ${TAG}"
echo "  输出 tar: ${OUT_TAR}"
echo "======================"

# 1) buildx 构建 amd64 镜像
echo "[1/3] docker buildx build (linux/amd64) ..."
docker buildx build \
  --platform linux/amd64 \
  -f docker/Dockerfile \
  -t "${TAG}" \
  --load \
  .

# 2) 用 buildx 输出为 docker-archive 再压缩
echo "[2/3] 导出镜像文件 ..."
mkdir -p "${OUT_DIR}"
docker buildx build \
  --platform linux/amd64 \
  -f docker/Dockerfile \
  -t "${TAG}" \
  --output "type=docker,dest=${OUT_TAR}" \
  --progress=plain \
  .

echo "[3/3] 完成: ${OUT_TAR}"
ls -lh "${OUT_TAR}"
echo ""
echo "  部署方式: 将 tar 传到 NAS 后执行 nas-deploy/start.sh"

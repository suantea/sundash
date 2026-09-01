#!/bin/sh
# SunDash 入口脚本：
#  1. 确保数据目录存在
#  2. 以 root 运行时：修正数据目录属主并降权到目标用户（默认 sundash）
#  3. 非 root（用户通过 compose 指定了 user）时：校验可写性，给清晰报错
set -e

DATA_DIR="${SUNDASH_DATA_DIR:-/app/data}"
RUN_USER="${SUNDASH_RUN_USER:-sundash}"

mkdir -p "$DATA_DIR"

if [ "$(id -u)" = "0" ]; then
    # root：修正属主后降权运行（默认 sundash；可经 SUNDASH_RUN_USER 覆盖）
    chown -R "$RUN_USER:$RUN_USER" "$DATA_DIR" 2>/dev/null || true
    exec su-exec "$RUN_USER" "$@"
fi

# 非 root：校验可写性，失败时给出可操作的提示（而不是晦涩的 SQLite 14 错误码）
if [ ! -w "$DATA_DIR" ]; then
    echo "ERROR: data directory $DATA_DIR is not writable by uid $(id -u)." >&2
    echo "       Fix: chown -R $(id -u):$(id -g) $DATA_DIR  (on the host), then recreate the container." >&2
    exit 1
fi

exec "$@"

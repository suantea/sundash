# SunDash 部署说明

本文档覆盖三种部署方式：**Docker**（推荐）、**裸二进制 + systemd**、**NAS 脚本**。

> ⚠️ **重要**：`SUNDASH_JWT_SECRET` 为必填环境变量，未设置时服务**拒绝启动**（`SUNDASH_DEBUG=true` 仅限开发）。生产环境务必使用随机长密钥：`openssl rand -base64 48`。

## 构件说明

| 构件 | 来源 | 说明 |
|---|---|---|
| `sundash`（二进制） | `server/` 编译 | Go 1.26+，纯 Go 无 CGO |
| `static/` | `web/` 构建产物 | `npm run build` 输出 `web/dist/` 重命名 |
| `data/` | 运行时创建 | SQLite 数据库 `sundash.db` 所在目录 |

---

## 方式 A：Docker（推荐）

### 1. 准备密钥

```bash
cd docker
echo "SUNDASH_JWT_SECRET=$(openssl rand -base64 48)" > .env
```

> `docker-compose.yml` 使用 `${SUNDASH_JWT_SECRET:?...}` 强制检查，未设置会直接报错。

### 2. 构建并启动

```bash
docker compose up -d --build
```

- 镜像：多阶段构建（node 构建前端 → golang 编译后端 → alpine 运行时）
- 容器以**非 root** 用户运行，带 **healthcheck**（`/api/site-config` 探活）
- 数据持久化到 Docker volume `sundash-data`（`/app/data`）
- 访问：`http://<主机IP>:3000`

### 3. 更新与日志

```bash
docker compose pull && docker compose up -d   # 更新
docker compose logs -f sundash                  # 查看日志
docker compose down                           # 停止
```

---

## 方式 B：裸二进制 + systemd

### 1. 构建（本机执行）

```bash
# 前端
cd web
npm install
npm run build            # 产物 web/dist/

# 后端交叉编译（Linux/amd64；ARM 机器改 GOARCH=arm64）
cd ../server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sundash .
```

### 2. 上传与目录布局（服务器执行）

```bash
mkdir -p /opt/sundash
# 上传 sundash 二进制与 web/dist 内容（重命名为 static/）到 /opt/sundash/
# 最终结构：
#   /opt/sundash/sundash
#   /opt/sundash/static/
cd /opt/sundash
chmod +x sundash
```

### 3. 前台验证

```bash
export SUNDASH_PORT=3000
export SUNDASH_DATA_DIR=/opt/sundash/data
export SUNDASH_JWT_SECRET='随机长字符串'
./sundash

# 另开终端自测
curl -s http://127.0.0.1:3000/api/site-config   # 期望 {"site_title":"",...}
```

### 4. 注册 systemd 服务

`/etc/systemd/system/sundash.service`：

```ini
[Unit]
Description=SunDash Panel
After=network.target

[Service]
WorkingDirectory=/opt/sundash
ExecStart=/opt/sundash/sundash
Restart=always
RestartSec=3
Environment=SUNDASH_PORT=3000
Environment=SUNDASH_DATA_DIR=/opt/sundash/data
Environment=SUNDASH_JWT_SECRET=随机长字符串

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable --now sundash
systemctl status sundash
```

---

## 方式 C：NAS 部署（QNAP / 群晖 / TrueNAS）

`nas-deploy/` 目录提供 NAS 场景脚本：

- `docker-compose.nas.yml`：NAS Docker 部署编排（卷挂载到 NAS 存储路径）
- `start.sh` / `start.bat`：一键启动脚本（自动检查 JWT 密钥并生成 `.env`）

```bash
cd nas-deploy
./start.sh        # Linux/NAS（群晖等）
# 或 start.bat（Windows）
```

### NAS 部署注意事项

1. **数据持久化**：确保 `docker-compose.nas.yml` 中配置的宿主机路径（如 `/share/Container/sundash/data`）指向 NAS 存储卷，避免容器重建时数据丢失。

2. **权限**：NAS Docker 通常以特定 UID/GUID 运行，请确保容器对数据目录有读写权限。如遇权限问题，可在 `docker-compose.nas.yml` 中配置 `user: "UID:GID"`。

3. **网络**：
   - 默认端口 `3000`，如被占用可修改 `SUNDASH_PORT` 环境变量。
   - 如需通过 80/443 端口访问，建议在前端加 nginx 反代（示例见下方）。

4. **RSS 订阅后台抓取**：默认每 1 分钟检查所有订阅源是否需要更新（基于 `update_interval` 字段）。首次添加订阅时会立即抓取一次，后续按间隔自动更新。抓取的文章保留最新 50 条。

5. **系统监控**：基于 `gopsutil` 读取系统指标，在容器内运行时依赖 `/proc` 和 `/sys` 文件系统。如容器未挂载这些目录，监控数据可能为空。建议在 `docker-compose` 中添加：
   ```yaml
   volumes:
     - /proc:/host/proc:ro
     - /sys:/host/sys:ro
   ```
   并相应调整服务代码读取 `host.proc` 与 `host.sys`。

6. **天气 Widget**：默认使用固定坐标（北京）。如需自动定位，需在前端实现 IP 地理定位，并通过 `lat` / `lon` 参数调用 `/api/weather`。

---

## nginx 反向代理（80/443 端口访问）

```nginx
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

如需 HTTPS，请为 `server_name` 配置证书并添加 443 监听；若 nginx 与后端不同机，记得把 `SUNDASH_ALLOWED_ORIGINS` 配成你的域名。

---

## 环境变量一览

| 变量 | 默认 | 说明 |
|---|---|---|
| `SUNDASH_PORT` | `3000` | 监听端口 |
| `SUNDASH_DATA_DIR` | `./data` | 数据库目录 |
| `SUNDASH_JWT_SECRET` | **必填** | 签名密钥，生产用 `openssl rand -base64 48` |
| `SUNDASH_DEBUG` | 空 | `true` = 开发模式（Gin 调试 + 开发密钥） |
| `SUNDASH_ALLOWED_ORIGINS` | 空 | CORS 白名单（逗号分隔），空 = 仅同源 |
| `SUNDASH_MCP_TOKEN` | 空 | MCP 端点静态 token（可选），设置后 AI 客户端可用它管理书签 |
| `SUNDASH_MCP_USERNAME` | `admin` | MCP 静态 token 绑定的用户 |

## MCP（AI 接入）

服务启动后 MCP 端点为 `POST /mcp`（Streamable HTTP，与 Web 同端口）。启用步骤：

1. 设置 `SUNDASH_MCP_TOKEN`（如 `openssl rand -base64 32`）与 `SUNDASH_MCP_USERNAME`（默认 admin）
2. 在 AI 客户端（Claude Desktop / Cursor 等）配置 MCP server：

```json
{
  "mcpServers": {
    "sundash": {
      "type": "http",
      "url": "http://<主机>:3000/mcp",
      "headers": { "Authorization": "Bearer <SUNDASH_MCP_TOKEN>" }
    }
  }
}
```

3. 可用工具（共 18 个，详见 README「MCP」章节）：书签/分组增删改查与排序（`sundash_list_groups` / `sundash_create_group` / `sundash_rename_group` / `sundash_delete_group` / `sundash_reorganize` / `sundash_create_card` / `sundash_update_card` / `sundash_move_card` / `sundash_delete_card`）、图标批量管理（`sundash_set_icons` / `sundash_suggest_icons` / `sundash_auto_iconify`）、系统状态（`sundash_system_status`）、全局搜索（`sundash_search`）、便签（`sundash_list_memo` / `sundash_add_memo` / `sundash_archive_memo` / `sundash_delete_memo`）

## 验收清单

- [ ] 页面打开正常（静态资源 + SPA 路由）
- [ ] 登录成功（默认 `admin / admin`，登录后立即修改密码）
- [ ] 服务重启后自动拉起（systemd / Docker restart）
- [ ] `data/sundash.db` 持久化且可备份（试一次 `GET /api/admin/backup` 下载快照）
- [ ] `SUNDASH_JWT_SECRET` 已设置且非默认值
- [ ] Docker 方式下 healthcheck 状态为 healthy
- [ ] 系统监控数据正常显示（CPU/内存/磁盘/网络）
- [ ] 天气 Widget 正常获取数据（默认北京坐标）
- [ ] 搜索功能返回匹配书签（FTS5 全文索引）
- [ ] RSS 订阅后台抓取正常（首次添加后等待 1-2 分钟）
- [ ] 便签创建/归档/删除功能正常
- [ ] （可选）`SUNDASH_MCP_TOKEN` 已设置，AI 客户端可调用 `/mcp` 管理书签

## 故障排查

| 现象 | 原因与处理 |
|---|---|
| 启动即退出，日志提示 `SUNDASH_JWT_SECRET ... required` | 未设置密钥；生产环境设置随机密钥，开发用 `SUNDASH_DEBUG=true` |
| 打开页面只有 API 没有界面 | `static/` 目录缺失或为空，重新执行前端构建并上传 |
| 登录后接口返回 401 | 更换过 `SUNDASH_JWT_SECRET` 导致旧 token 失效，重新登录即可 |
| 端口被占用 | 检查 `SUNDASH_PORT` 是否被其他进程占用（`ss -lntp` / `netstat -ano`） |
| 数据库文件被锁 | 确保只有一个实例在运行；SQLite 使用 WAL 模式，勿直接拷贝运行中的 db 文件——用**管理员后台备份端点** `GET /api/admin/backup`（运行中生成一致性完整快照下载，内部即 `VACUUM INTO`） |
| 系统监控无数据 | 容器内缺少 `/proc`、`/sys` 挂载；NAS 部署时建议只读挂载这两个目录 |
| 天气 Widget 显示异常 | 检查网络是否能访问 `api.open-meteo.com`；NAS 需确保容器有外网访问权限 |
| 搜索无结果 | FTS5 索引在数据库迁移时自动创建；如数据是旧版本导入，需手动执行 `INSERT INTO cards_fts(cards_fts) VALUES('rebuild')` 重建索引 |
| RSS 订阅无文章 | 检查订阅 URL 是否有效（需支持 RSS/Atom）；首次添加后等待后台抓取（约 1 分钟）；查看容器日志确认抓取是否报错 |

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
- 数据持久化到 Docker volume `asuan-data`（`/app/data`）
- 访问：`http://<主机IP>:3000`

### 3. 更新与日志

```bash
docker compose pull && docker compose up -d   # 更新
docker compose logs -f asuan                  # 查看日志
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
mkdir -p /opt/asuan
# 上传 sundash 二进制与 web/dist 内容（重命名为 static/）到 /opt/asuan/
# 最终结构：
#   /opt/asuan/sundash
#   /opt/asuan/static/
cd /opt/asuan
chmod +x sundash
```

### 3. 前台验证

```bash
export SUNDASH_PORT=3000
export SUNDASH_DATA_DIR=/opt/asuan/data
export SUNDASH_JWT_SECRET='随机长字符串'
./sundash

# 另开终端自测
curl -s http://127.0.0.1:3000/api/site-config   # 期望 {"site_title":"",...}
```

### 4. 注册 systemd 服务

`/etc/systemd/system/asuan.service`：

```ini
[Unit]
Description=SunDash Panel
After=network.target

[Service]
WorkingDirectory=/opt/asuan
ExecStart=/opt/asuan/sundash
Restart=always
RestartSec=3
Environment=SUNDASH_PORT=3000
Environment=SUNDASH_DATA_DIR=/opt/asuan/data
Environment=SUNDASH_JWT_SECRET=随机长字符串

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable --now asuan
systemctl status asuan
```

---

## 方式 C：NAS 部署

`nas-deploy/` 目录提供 NAS 场景脚本：

- `docker-compose.yml`：NAS Docker 部署编排（卷挂载到 NAS 存储路径）
- `start.sh` / `start.bat`：一键启动脚本（自动检查 JWT 密钥并生成 `.env`）

```bash
cd nas-deploy
./start.sh        # Linux/NAS（群晖等）
# 或 start.bat（Windows）
```

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

3. 可用工具：`sundash_list_groups` / `sundash_create_group` / `sundash_rename_group` / `sundash_delete_group` / `sundash_create_card` / `sundash_update_card` / `sundash_move_card` / `sundash_delete_card`（详见 README「MCP」章节）

## 验收清单

- [ ] 页面打开正常（静态资源 + SPA 路由）
- [ ] 登录成功（默认 `admin / admin`，登录后立即修改密码）
- [ ] 服务重启后自动拉起（systemd / Docker restart）
- [ ] `data/sundash.db` 持久化且可备份
- [ ] `SUNDASH_JWT_SECRET` 已设置且非默认值
- [ ] Docker 方式下 healthcheck 状态为 healthy
- [ ] （可选）`SUNDASH_MCP_TOKEN` 已设置，AI 客户端可调用 `/mcp` 管理书签

## 故障排查

| 现象 | 原因与处理 |
|---|---|
| 启动即退出，日志提示 `SUNDASH_JWT_SECRET ... required` | 未设置密钥；生产环境设置随机密钥，开发用 `SUNDASH_DEBUG=true` |
| 打开页面只有 API 没有界面 | `static/` 目录缺失或为空，重新执行前端构建并上传 |
| 登录后接口返回 401 | 更换过 `SUNDASH_JWT_SECRET` 导致旧 token 失效，重新登录即可 |
| 端口被占用 | 检查 `SUNDASH_PORT` 是否被其他进程占用（`ss -lntp` / `netstat -ano`） |
| 数据库文件被锁 | 确保只有一个实例在运行；SQLite 使用 WAL 模式，勿直接拷贝运行中的 db 文件（应使用备份或 `VACUUM INTO`） |

<div align="center">

<img src="https://img.icons8.com/fluency/128/star.png" alt="SunDash" width="80"/>

# SunDash

**自托管个人导航面板** — 多用户书签管理 · 系统监控 · 天气 · RSS · 便签 · MCP

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.5-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vitejs.dev)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://docker.com)
[![MCP](https://img.shields.io/badge/MCP-AI%20Ready-000000?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xMiAyTDQgN3YxMGw4IDUgOC01VjdsLTgtNXptMCAyLjE4bDUuNzYgMy40NkwxMiAxMS44MmwtNS43Ni00LjE4TDEyIDQuMTh6Ii8+PC9zdmc+)](https://modelcontextprotocol.io)

[English](README.md) · [中文](README.zh.md)

</div>

---

**SunDash** 是为 NAS、云服务器和个人电脑设计的开源导航面板。作为浏览器主页 / 新标签页使用，提供书签管理、系统监控、天气、RSS、便签等一站式功能。

## ✨ 功能特性

### 📑 面板管理
- **多用户系统**：注册开关、审批流（pending/approve/reject）、管理员后台、JWT 认证
- **书签分组与卡片**：拖拽排序、折叠/隐藏、批量应用图标色/背景色、内网地址切换
- **站点配置**：标题、图标、CDN 前缀、统计代码、自定义 head/footer 注入
- **数据管理**：JSON 导出 / 导入（追加合并）、内置**常用导航模板**一键导入

### 🎨 个性化
- **多语言**：中文 / 英文切换（设置抽屉，基于 vue-i18n）
- **深浅主题** + 自定义主题色（Naive UI 主题覆盖）
- **壁纸**：必应每日壁纸（服务端缓存）、渐变、自定义图片，支持模糊/透明度/版权标注

### 📊 小组件（均可在设置抽屉开关）
- **系统监控**：CPU / 内存 / 磁盘 / 网络实时状态（NAS 场景刚需）
- **天气 Widget**：Open-Meteo 数据（温度 / 风速 / 天气状况），10 分钟缓存
- **RSS 订阅**：后台定时抓取（按订阅源配置的间隔）、支持多订阅源、文章展开阅读
- **快捷便签**：随手记备忘，支持归档管理
- **全局搜索**：SQLite FTS5 全文索引，搜索书签标题 / 链接 / 描述，含搜索建议

### 🚀 性能优化
- 前端资源哈希化（Vite）：JS/CSS 文件名包含内容哈希，支持长期浏览器缓存
- 多阶段 Docker 镜像：最终镜像仅含可执行文件和静态资源，体积更小、启动更快
- SQLite WAL 模式：读写并发更友好，适合读多写少的导航场景
- 面板服务缓存：用户面板数据内存缓存 5 分钟，所有变更操作即时失效，改动立刻可见
- NAS 友好部署：提供 `docker-compose.nas.yml` 示例，支持本地持久化卷、时区设置、健康检查与自动重启

### 🔄 跨浏览器书签同步
- **`/bookmarks` 子页面**（不在主页展示）：展示从外部 **bookmark-sync** 服务器同步来的完整书签树——与 Chrome MV3 / Safari 扩展共享的文件夹与书签
- **完整读写**：在仪表盘中直接新增、编辑、移动、删除书签；变更推送到同步服务器后传播到所有其他设备
- **LWW + 墓碑语义**：删除采用软删除墓碑（服务端保留、绝不物理删除），但会传播到其他电脑——在这里删掉一个书签，其他设备的 Chrome/Safari 也会移除它；`updatedAt` 仅在真实变更时刷新
- **配置**：在「管理面板 → 全局设置」填写同步服务器地址与 Token（`bmsync_server_url`、`bmsync_token`）；本地镜像存于 SQLite（`bmsync_nodes`），通过 pull / push 刷新

### 🔧 工程化
- 后端三层分层架构（handler → service → repository）+ 手动依赖注入，可单元测试
- 版本化数据库迁移（`schema_migrations`）、SQLite WAL 并发读、FTS5 全文索引
- `/api/bootstrap` 首屏聚合接口（1 个请求拿齐设置/资料/面板）、gzip 压缩、静态资源 immutable 缓存

### 🤖 Chrome 扩展（SunDash Monitor）
- 页面内容变化监控（哈希比对 + MutationObserver），变更桌面通知
- 弹窗 / 侧边栏管理监控项，30s 轮询（带超时与大小限制）

---

## 🛠 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.26 · Gin · modernc.org/sqlite · golang-jwt/v5 · bcrypt |
| 前端 | Vue 3.5 · Vite 8 · TypeScript · Pinia · Vue Router · Naive UI · Iconify · vuedraggable |
| 扩展 | Chrome MV3（service worker + popup + sidepanel + content script） |
| 部署 | Docker · docker-compose · systemd · nginx |

---

## 🏗 架构

```
server/
├── main.go          入口：依赖注入、路由注册、SPA fallback、优雅关闭
├── config/          环境变量配置（JWT 密钥强制校验）
├── database/        SQLite 初始化 + 版本化迁移 + 默认管理员
├── repository/      数据访问层（user / panel / settings / memo / rss）
├── service/         业务逻辑层（认证、面板、搜索、便签、RSS、壁纸、favicon、系统监控）
├── handlers/        HTTP 层（参数绑定 + 调用 service + 响应）
├── middleware/      JWT 鉴权（HS256 固定）/ 管理员鉴权 / CORS 白名单 / gzip / 登录限流
├── mcp/             MCP 服务端（书签/搜索/便签/系统工具）+ 鉴权（静态 token 或 JWT）
└── models/          数据模型与请求/响应结构
```

前端（`web/src/`）：`views/` 页面 · `components/` 组件（卡片 / 时钟 / 便签 / RSS / 搜索 / 设置 / 状态 / 天气）· `stores/` Pinia 状态 · `composables/` 组合式逻辑（壁纸、分组可见性）· `api/` axios 封装（baseURL `/api`）。

---

## 🚀 快速开始

前置：Go 1.26+、Node.js 20+。

```bash
# 1. 启动后端（开发模式，自动使用开发 JWT 密钥）
cd server
SUNDASH_DEBUG=true go run .

# 2. 启动前端（Vite dev，代理 /api → localhost:3000）
cd ../web
npm install
npm run dev
```

打开 http://localhost:5173 ，默认账号 **admin / admin**（首次部署自动创建，登录后请立即修改）。

---

## 📦 部署

### 方式 A：Docker — 预构建镜像（推荐，免编译）

```bash
# 1. 复制环境变量模板，并填入随机 JWT 密钥
cp .env.example .env
# 或直接生成: echo "SUNDASH_JWT_SECRET=$(openssl rand -base64 48)" > .env

# 2. 启动（自动拉取 ghcr.io/suantea/sundash，amd64 + arm64 双架构）
docker compose up -d --pull always
# 访问 http://<主机>:3000
```

根目录 `docker-compose.yml` 直接使用已发布的镜像，无需本地编译。如需从源码构建，使用 `docker/docker-compose.yml`：

```bash
cd docker
cp ../.env.example .env
docker compose up -d --build
```

常用的 `docker compose` 操作：

```bash
docker compose logs -f sundash         # 查看日志
docker compose down                    # 停止（数据保留在 volume 中）
docker compose up -d --pull always     # 升级到最新镜像
docker compose exec sundash sh         # 进入容器
```

如需修改宿主机端口（例如改用 8080），编辑 `docker-compose.yml` 的 `ports:`（改成 `"8080:3000"`），再重新执行 `docker compose up -d` 即可。

#### NAS / QNAP 部署

NAS（QNAP / 群晖 / TrueNAS）上的两种部署方式：

- **Container Station 直接拉取镜像**（QNAP，最简单，无需离线包）：
  1. Container Station → **映像（Images）** → **拉取（Pull）** → 输入 `ghcr.io/suantea/sundash:latest`
  2. **创建（Create）**：端口映射 `3000:3000`、卷挂载 `/share/Container/sundash/data:/app/data`、环境变量 `SUNDASH_JWT_SECRET`（必填）+ `TZ=Asia/Shanghai`
  3. 访问 `http://<NAS_IP>:3000`，默认账号 `admin / admin`

- **离线包**（NAS 无外网时）：用 `./docker/build-qnap.sh` 构建 `qnap/sundash-image.tar.gz`，拷到 NAS 后执行 `cd nas-deploy && ./start.sh`（Windows 用 `start.bat`）。脚本会自动加载镜像、生成 JWT 密钥并启动容器。

完整步骤见 [DEPLOY.md → NAS 部署](DEPLOY.md#方式-cnas-部署qnap--群晖--truenas)。

### 方式 B：裸二进制 + systemd

```bash
# 前端构建
cd web && npm install && npm run build   # 产物 dist/

# 后端编译（Linux/amd64）
cd ../server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sundash .

# 服务器上：目录布局 /opt/sundash/{sundash, static/, data/}
# static/ = web/dist 的内容；data/ 运行时自动创建
```

```bash
export SUNDASH_JWT_SECRET='随机长字符串'   # 必填，否则拒绝启动
export SUNDASH_DATA_DIR=/opt/sundash/data
./sundash
```

systemd / nginx 反代示例见 [DEPLOY.md](DEPLOY.md)。

---

## 🔧 环境变量

| 变量 | 默认 | 说明 |
|---|---|---|
| `SUNDASH_PORT` | `3000` | 监听端口 |
| `SUNDASH_DATA_DIR` | `./data` | SQLite 数据库目录 |
| `SUNDASH_JWT_SECRET` | **无** | **必填**；生产环境请用 `openssl rand -base64 48` 生成；`SUNDASH_DEBUG=true` 时可用开发密钥 |
| `SUNDASH_DEBUG` | 空 | `true` 开启 debug（Gin 调试 + 开发 JWT 密钥） |
| `SUNDASH_ALLOWED_ORIGINS` | 空 | 逗号分隔的 CORS 白名单；空 = 仅同源 |
| `SUNDASH_MCP_TOKEN` | 空 | **MCP 端点专用静态 token**（可选）。设置后 AI 客户端以 `Authorization: Bearer <token>` 访问 `/mcp` 管理书签 |
| `SUNDASH_MCP_USERNAME` | `admin` | MCP 静态 token 绑定的用户（默认管理员） |

---

## 🤖 MCP（AI 接入）

SunDash 内置 [MCP](https://modelcontextprotocol.io)（Model Context Protocol）服务端，AI 助手（如 Claude Desktop / Cursor / 本地 agent）可直接接入来**录入、整理书签、记便签、查系统状态**（分组 = 标签）。

- **端点**：`POST /mcp`（Streamable HTTP），挂载在 Web 服务同端口
- **认证**：`Authorization: Bearer <SUNDASH_MCP_TOKEN>`（推荐），或直接使用登录 JWT
- **工具**（均绑定到 token 对应用户）：

| 工具 | 说明 |
|---|---|
| `sundash_list_groups` | 列出全部分组及组内卡片（只读） |
| `sundash_create_group` | 新建分组（标签分类） |
| `sundash_rename_group` | 重命名分组 |
| `sundash_delete_group` | 删除分组（组内卡片一并删除） |
| `sundash_reorganize` | AI 辅助批量整理：自然语言指令 → 一系列分组/卡片操作 |
| `sundash_create_card` | 组内新增书签（title/url 必填） |
| `sundash_update_card` | 更新卡片标题/链接/描述 |
| `sundash_move_card` | 把卡片移动到其他分组（整理归类） |
| `sundash_delete_card` | 删除卡片 |
| `sundash_set_icons` | 批量设置卡片图标色/背景色 |
| `sundash_suggest_icons` | 为卡片建议图标 |
| `sundash_auto_iconify` | 自动为全部分组配图标配色 |
| `sundash_system_status` | 获取系统监控状态（CPU/内存/磁盘/网络/主机信息） |
| `sundash_search` | 全局搜索书签（FTS5 全文索引，按相关度排序） |
| `sundash_list_memo` | 获取便签列表 |
| `sundash_add_memo` | 新增便签 |
| `sundash_archive_memo` | 归档/取消归档便签 |
| `sundash_delete_memo` | 删除便签 |

### 配置示例（Claude Desktop）

```json
{
  "mcpServers": {
    "sundash": {
      "type": "http",
      "url": "http://localhost:3000/mcp",
      "headers": { "Authorization": "Bearer <你的 SUNDASH_MCP_TOKEN>" }
    }
  }
}
```

---

## 📋 主要 API

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/auth/login` `/api/auth/register` | 登录 / 注册（登录接口限流） |
| GET | `/api/bootstrap` | **首屏聚合**：设置 + 资料 + 面板（1 个请求） |
| GET/PUT | `/api/profile` `/api/profile/password` | 个人资料 / 改密码 |
| GET/POST | `/api/panels` `/api/panels/groups` | 面板 / 分组 |
| PUT | `/api/panels/reorder` | 拖拽排序（校验归属） |
| GET/PUT | `/api/settings` `/api/settings/batch` | 用户设置 |
| GET | `/api/wallpaper/bing[:date]` | 必应壁纸（服务端缓存） |
| GET | `/api/favicon?url=` | favicon 识别（含 SSRF 防护） |
| GET/PUT | `/api/admin/settings` | 全局站点设置（管理员） |
| GET/PUT/DELETE | `/api/users` `/api/users/:id` | 用户管理（管理员） |
| GET | `/api/system/stats` | 系统监控（CPU/内存/磁盘/网络） |
| GET | `/api/weather` | 天气数据（Open-Meteo） |
| GET/POST | `/api/search` `/api/search/suggestions` | 全局搜索 / 搜索建议 |
| GET/POST/PUT/DELETE | `/api/memo` `/api/memo/:id` | 便签列表 / 新建 / 归档 / 更新 / 删除 |
| GET/POST/PUT/DELETE | `/api/rss` `/api/rss/:id` | RSS 订阅管理 |
| GET | `/api/rss/:id/items` | RSS 文章列表（支持 `limit`，校验订阅归属） |
| GET | `/api/bmsync/status` | 书签同步连接状态（是否配置 / 是否已同步 / rev） |
| GET | `/api/bmsync/tree` | 本地书签镜像（只读，不触发网络请求） |
| POST | `/api/bmsync/pull` | 从同步服务器拉取全量规范书签树 |
| POST | `/api/bmsync/push` | 推送新增/编辑/删除变更（LWW + 墓碑），并存储返回的规范状态 |

---

## 📥 数据导入

- 设置页 / 设置抽屉 → **导入数据**：从 JSON 备份文件追加导入（`{settings, groups}` 格式）
- **模板**：内置「常用导航」模板（5 分组 28 书签）一键导入

---

## 🌐 浏览器扩展

`extension/` 为 Chrome MV3 扩展（SunDash Monitor）：加载方式为浏览器「开发者模式 → 加载已解压的扩展程序」选择该目录。功能：监控任意页面内容变化、变更通知、弹窗与侧边栏管理。

---

## 🤝 贡献

欢迎提交 Issue 和 PR！请遵循以下规范：

- 提交前请确保 `go build ./...`、`go test ./...` 和 `npm run build` 通过
- 新功能请附带相应测试
- 提交信息请遵循 [Conventional Commits](https://www.conventionalcommits.org/)

---

## 📜 License

[MIT](LICENSE) © 2026 SunDash Contributors

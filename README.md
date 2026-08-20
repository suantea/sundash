# SunDash（原 Asuan）

自托管个人导航面板：**多用户书签管理 + 必应壁纸 + 页面监控扩展**。适合部署在 NAS、云服务器或个人电脑上，作为浏览器主页 / 新标签页使用。

- 后端：Go 1.26 + Gin + SQLite（纯 Go 驱动，零 CGO）
- 前端：Vue 3 + Vite + TypeScript + Pinia + Naive UI
- 扩展：Chrome Manifest V3（页面内容监控）
- 部署：Docker / docker-compose / 裸二进制 + systemd / NAS 脚本

## 功能特性

**面板管理**
- 多用户系统：注册开关、审批流（pending/approve/reject）、管理员后台、JWT 认证
- 书签分组与卡片：拖拽排序、折叠/隐藏、批量应用图标色/背景色、内网地址切换
- 站点配置：标题、图标、CDN 前缀、统计代码、自定义 head/footer 注入
- 数据管理：JSON 导出 / 导入（追加合并）、内置**常用导航模板**一键导入

**个性化**
- 深浅主题 + 自定义主题色（Naive UI 主题覆盖）
- 壁纸：必应每日壁纸（服务端缓存）、渐变、自定义图片，支持模糊/透明度/版权标注
- 时钟（12/24 小时制）、Logo（图片/文字）、搜索栏（多引擎）、系统状态显示

**工程化**
- 后端三层分层架构（handler → service → repository）+ 手动依赖注入，可单元测试
- 版本化数据库迁移（`schema_migrations`）、SQLite WAL 并发读
- 性能：`/api/bootstrap` 首屏聚合接口（1 个请求拿齐设置/资料/面板）、gzip 压缩、静态资源 immutable 缓存、站点配置短 TTL 缓存

**Chrome 扩展（SunDash Monitor）**
- 页面内容变化监控（哈希比对 + MutationObserver），变更桌面通知
- 弹窗 / 侧边栏管理监控项，30s 轮询（带超时与大小限制）

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.26 · Gin · modernc.org/sqlite · golang-jwt/v5 · bcrypt |
| 前端 | Vue 3.5 · Vite 8 · TypeScript · Pinia · Vue Router · Naive UI · Iconify · vuedraggable |
| 扩展 | Chrome MV3（service worker + popup + sidepanel + content script） |
| 部署 | Docker · docker-compose · systemd · nginx |

## 架构

```
server/
├── main.go          入口：依赖注入、路由注册、SPA fallback、优雅关闭
├── config/          环境变量配置（JWT 密钥强制校验）
├── database/        SQLite 初始化 + 版本化迁移 + 默认管理员
├── repository/      数据访问层（user / panel / settings）
├── service/         业务逻辑层（认证、面板、设置、壁纸、favicon、统一错误）
├── handlers/        HTTP 层（参数绑定 + 调用 service + 响应）
├── middleware/      JWT 鉴权（HS256 固定）/ 管理员鉴权 / CORS 白名单 / gzip / 登录限流
└── models/          数据模型与请求/响应结构
```

前端（`web/src/`）：`views/` 页面 · `components/` 组件 · `stores/` Pinia 状态 · `composables/` 组合式逻辑（壁纸、分组可见性）· `api/` axios 封装（baseURL `/api`）。

## 快速开始（本地开发）

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

## 构建与部署

### 方式 A：Docker（推荐）

```bash
cd docker
# 先设置密钥（docker/.env，compose 强制要求）
echo "SUNDASH_JWT_SECRET=$(openssl rand -base64 48)" > .env
docker compose up -d --build
# 访问 http://<主机>:3000
```

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

## 环境变量

| 变量 | 默认 | 说明 |
|---|---|---|
| `SUNDASH_PORT` | `3000` | 监听端口 |
| `SUNDASH_DATA_DIR` | `./data` | SQLite 数据库目录 |
| `SUNDASH_JWT_SECRET` | **无** | **必填**；生产环境请用 `openssl rand -base64 48` 生成；`SUNDASH_DEBUG=true` 时可用开发密钥 |
| `SUNDASH_DEBUG` | 空 | `true` 开启 debug（Gin 调试 + 开发 JWT 密钥） |
| `SUNDASH_ALLOWED_ORIGINS` | 空 | 逗号分隔的 CORS 白名单；空 = 仅同源 |
| `SUNDASH_MCP_TOKEN` | 空 | **MCP 端点专用静态 token**（可选）。设置后 AI 客户端以 `Authorization: Bearer <token>` 访问 `/mcp` 管理书签 |
| `SUNDASH_MCP_USERNAME` | `admin` | MCP 静态 token 绑定的用户（默认管理员） |

## MCP（AI 接入）

SunDash 内置 [MCP](https://modelcontextprotocol.io)（Model Context Protocol）服务端，AI 助手（如 Claude Desktop / Cursor / 本地 agent）可直接接入来**录入、整理书签**（分组 = 标签）。

- **端点**：`POST /mcp`（Streamable HTTP），挂载在 Web 服务同端口
- **认证**：`Authorization: Bearer <SUNDASH_MCP_TOKEN>`（推荐），或直接使用登录 JWT
- **工具**（均绑定到 token 对应用户）：

| 工具 | 说明 |
|---|---|
| `sundash_list_groups` | 列出全部分组及组内卡片（只读） |
| `sundash_create_group` | 新建分组（标签分类） |
| `sundash_rename_group` | 重命名分组 |
| `sundash_delete_group` | 删除分组（组内卡片一并删除） |
| `sundash_create_card` | 组内新增书签（title/url 必填） |
| `sundash_update_card` | 更新卡片标题/链接/描述 |
| `sundash_move_card` | 把卡片移动到其他分组（整理归类） |
| `sundash_delete_card` | 删除卡片 |

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

## 主要 API

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/auth/login` `/api/auth/register` | 登录 / 注册（登录接口限流） |
| GET | `/api/bootstrap` | **首屏聚合**：设置 + 资料 + 面板（1 个请求） |
| GET/PUT | `/api/profile` `/api/profile/password` | 个人资料 / 改密码 |
| GET/POST | `/api/panels` `/api/panels/groups` | 面板 / 分组 |
| PUT | `/api/panels/reorder` | 拖拽排序（校验归属） |
| GET/PUT | `/api/settings` `/api/settings/batch` | 用户设置 |
| GET | `/api/wallpaper/bing[:date]` | 必应壁纸（1 小时服务端缓存） |
| GET | `/api/favicon?url=` | favicon 识别（含 SSRF 防护） |
| GET/PUT | `/api/admin/settings` | 全局站点设置（管理员） |
| GET/PUT/DELETE | `/api/users` `/api/users/:id` | 用户管理（管理员） |

## 数据导入

- 设置页 / 设置抽屉 → **导入数据**：从 JSON 备份文件追加导入（`{settings, groups}` 格式）
- **模板**：内置「常用导航」模板（5 分组 28 书签）一键导入
- `temp/` 目录含 SunPanel 数据导入脚本（`import-sunpanel.js`，Node 脚本，供迁移参考）

## 浏览器扩展

`extension/` 为 Chrome MV3 扩展（SunDash Monitor）：加载方式为浏览器「开发者模式 → 加载已解压的扩展程序」选择该目录。功能：监控任意页面内容变化、变更通知、弹窗与侧边栏管理。

## License

MIT

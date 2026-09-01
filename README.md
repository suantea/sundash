<div align="center">

<img src="https://img.icons8.com/fluency/128/star.png" alt="SunDash" width="80"/>

# SunDash

**Self-hosted personal dashboard** — Multi-user bookmarks · System monitor · Weather · RSS · Memos · MCP

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.5-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vitejs.dev)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://docker.com)
[![MCP](https://img.shields.io/badge/MCP-AI%20Ready-000000?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xMiAyTDQgN3YxMGw4IDUgOC01VjdsLTgtNXptMCAyLjE4bDUuNzYgMy40NkwxMiAxMS44MmwtNS43Ni00LjE4TDEyIDQuMTh6Ii8+PC9zdmc+)](https://modelcontextprotocol.io)

[English](README.md) · [中文](README.zh.md)

</div>

---

**SunDash** is an open-source dashboard for NAS devices, cloud servers, and personal computers. Designed as your browser home page / new-tab page, it brings bookmarks, system monitoring, weather, RSS, and quick memos together in one place.

## ✨ Features

### 📑 Panel Management
- **Multi-user system**: registration toggle, approval flow (pending/approve/reject), admin console, JWT auth
- **Bookmark groups & cards**: drag-and-drop ordering, collapse/hide, batch icon/background colors, internal-URL switching
- **Site configuration**: title, icon, CDN prefix, analytics code, custom head/footer injection
- **Data management**: JSON export / import (append-merge), one-click import of a built-in starter template

### 🎨 Personalization
- **i18n**: Chinese / English switcher (settings drawer, vue-i18n)
- **Light/dark themes** + custom accent colors (Naive UI theme overrides)
- **Wallpapers**: Bing daily wallpaper (server-cached), gradients, custom images, with blur/opacity/attribution

### 📊 Widgets (each toggleable in the settings drawer)
- **System monitor**: live CPU / memory / disk / network status (a NAS essential)
- **Weather widget**: Open-Meteo data (temperature / wind / conditions), 10-minute cache
- **RSS feeds**: background fetching on a per-feed interval, multiple sources, expandable article reading
- **Quick memos**: jot notes from the home page, with archive support
- **Global search**: SQLite FTS5 full-text index over bookmark titles / URLs / descriptions, with suggestions

### 🚀 Performance
- Hashed frontend assets (Vite): content-hashed JS/CSS filenames for long-lived browser caching
- Multi-stage Docker image: final image contains only the binary and static assets
- SQLite WAL mode: friendlier read/write concurrency for read-heavy dashboard workloads
- Panel cache: per-user panel data cached in memory for 5 minutes; every mutation invalidates it instantly so changes are visible immediately
- NAS-friendly deploys: `docker-compose.nas.yml` example with local volumes, timezone, healthcheck, and auto-restart

### 🔄 Cross-Browser Bookmark Sync
- **`/bookmarks` sub-page** (not the home page) shows the full bookmark tree synced from an external **bookmark-sync** server — folders and bookmarks shared with the Chrome MV3 / Safari extensions
- **Full read/write**: create, edit, move and delete bookmarks right from the dashboard; changes are pushed to the sync server and propagate to every other device
- **LWW + tombstone semantics**: deletions are soft-deleted tombstones kept on the server (never physically erased) but propagated to other computers, so removing a bookmark here removes it from Chrome/Safari on other devices; `updatedAt` is only stamped on real changes
- **Config**: set the sync server URL + bearer token in **Admin → global settings** (`bmsync_server_url`, `bmsync_token`); the local mirror lives in SQLite (`bmsync_nodes`) and is refreshed via pull/push

### 🔧 Engineering
- Three-layer backend (handler → service → repository) with manual dependency injection, unit-testable
- Versioned database migrations (`schema_migrations`), SQLite WAL, FTS5 full-text index
- `/api/bootstrap` first-screen aggregate (settings + profile + panels in one request), gzip, immutable static caching

### 🤖 Chrome Extension (SunDash Monitor)
- Page content change monitoring (hash comparison + MutationObserver) with desktop notifications
- Popup / sidepanel management, 30s polling (with timeout and size limits)

---

## 🛠 Tech Stack

| Layer | Technologies |
|---|---|
| Backend | Go 1.26 · Gin · modernc.org/sqlite · golang-jwt/v5 · bcrypt |
| Frontend | Vue 3.5 · Vite 8 · TypeScript · Pinia · Vue Router · Naive UI · Iconify · vuedraggable |
| Extension | Chrome MV3 (service worker + popup + sidepanel + content script) |
| Deployment | Docker · docker-compose · systemd · nginx |

---

## 🏗 Architecture

```
server/
├── main.go          entrypoint: dependency injection, routes, SPA fallback, graceful shutdown
├── config/          environment configuration (JWT secret enforced)
├── database/        SQLite init + versioned migrations + default admin
├── repository/      data access (user / panel / settings / memo / rss)
├── service/         business logic (auth, panels, search, memos, RSS, wallpaper, favicon, system)
├── handlers/        HTTP layer (binding + service calls + responses)
├── middleware/      JWT auth (HS256 pinned) / admin auth / CORS allowlist / gzip / login rate limit
├── mcp/             MCP server (bookmark / search / memo / system tools) + auth (static token or JWT)
└── models/          data models and request/response structs
```

Frontend (`web/src/`): `views/` pages · `components/` components (cards / clock / memo / rss / search / settings / status / weather) · `stores/` Pinia state · `composables/` composables (wallpaper, group visibility) · `api/` axios wrapper (baseURL `/api`).

---

## 🚀 Quick Start

Prerequisites: Go 1.26+, Node.js 20+.

```bash
# 1. Start the backend (dev mode uses a development JWT secret)
cd server
SUNDASH_DEBUG=true go run .

# 2. Start the frontend (Vite dev, proxies /api → localhost:3000)
cd ../web
npm install
npm run dev
```

Open http://localhost:5173 and sign in with the default account **admin / admin** (created automatically on first start — change the password immediately).

---

## 📦 Deployment

### Option A: Docker — prebuilt image (recommended, no build)

```bash
# 1. Copy the env template and set a random JWT secret
cp .env.example .env
# or: echo "SUNDASH_JWT_SECRET=$(openssl rand -base64 48)" > .env

# 2. Start (pulls ghcr.io/suantea/sundash, amd64 + arm64)
docker compose up -d --pull always
# Visit http://<host>:3000
```

The root `docker-compose.yml` uses the published image directly — no local build needed. Prefer the in-repo build? Use `docker/docker-compose.yml` (build from source):

```bash
cd docker
cp ../.env.example .env
docker compose up -d --build
```

### Option B: Bare binary + systemd

```bash
# Build the frontend
cd web && npm install && npm run build   # output: dist/

# Build the backend (Linux/amd64)
cd ../server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sundash .

# On the server: layout /opt/sundash/{sundash, static/, data/}
# static/ = contents of web/dist; data/ is created at runtime
```

```bash
export SUNDASH_JWT_SECRET='a-long-random-string'   # required, refuses to start otherwise
export SUNDASH_DATA_DIR=/opt/sundash/data
./sundash
```

See [DEPLOY.md](DEPLOY.md) for systemd / nginx reverse-proxy examples.

---

## 🔧 Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SUNDASH_PORT` | `3000` | Listen port |
| `SUNDASH_DATA_DIR` | `./data` | SQLite database directory |
| `SUNDASH_JWT_SECRET` | — | **Required**; generate with `openssl rand -base64 48`; a development key is allowed when `SUNDASH_DEBUG=true` |
| `SUNDASH_DEBUG` | — | `true` enables debug mode (Gin debug + development JWT secret) |
| `SUNDASH_ALLOWED_ORIGINS` | — | Comma-separated CORS allowlist; empty = same-origin only |
| `SUNDASH_MCP_TOKEN` | — | **Static token for the MCP endpoint** (optional). AI clients then call `/mcp` with `Authorization: Bearer <token>` |
| `SUNDASH_MCP_USERNAME` | `admin` | User bound to the MCP static token (default admin) |

---

## 🤖 MCP (AI Integration)

SunDash ships with a built-in [MCP](https://modelcontextprotocol.io) (Model Context Protocol) server, so AI assistants (Claude Desktop / Cursor / local agents) can **add and organize bookmarks, manage memos, and check system status** on your behalf (groups double as tags).

- **Endpoint**: `POST /mcp` (Streamable HTTP), mounted on the same port as the web server
- **Auth**: `Authorization: Bearer <SUNDASH_MCP_TOKEN>` (recommended), or a regular login JWT
- **Tools** (all bound to the token's user):

| Tool | Description |
|---|---|
| `sundash_list_groups` | List all groups with their cards (read-only) |
| `sundash_create_group` | Create a group (tag category) |
| `sundash_rename_group` | Rename a group |
| `sundash_delete_group` | Delete a group (cards are deleted with it) |
| `sundash_reorganize` | AI-assisted bulk reorganization: natural-language instruction → a series of group/card operations |
| `sundash_create_card` | Add a bookmark to a group (title/url required) |
| `sundash_update_card` | Update a card's title/link/description |
| `sundash_move_card` | Move a card to another group |
| `sundash_delete_card` | Delete a card |
| `sundash_set_icons` | Batch-set icon/background colors on cards |
| `sundash_suggest_icons` | Suggest icons for cards |
| `sundash_auto_iconify` | Auto-assign icon colors across groups |
| `sundash_system_status` | System monitor status (CPU/memory/disk/network/host) |
| `sundash_search` | Full-text bookmark search (FTS5, relevance-ranked) |
| `sundash_list_memo` | List memos |
| `sundash_add_memo` | Create a memo |
| `sundash_archive_memo` | Archive / unarchive a memo |
| `sundash_delete_memo` | Delete a memo |

### Example configuration (Claude Desktop)

```json
{
  "mcpServers": {
    "sundash": {
      "type": "http",
      "url": "http://localhost:3000/mcp",
      "headers": { "Authorization": "Bearer <your SUNDASH_MCP_TOKEN>" }
    }
  }
}
```

---

## 📋 API Overview

| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/login` `/api/auth/register` | Login / register (login is rate-limited) |
| GET | `/api/bootstrap` | **First-screen aggregate**: settings + profile + panels in one request |
| GET/PUT | `/api/profile` `/api/profile/password` | Profile / change password |
| GET/POST | `/api/panels` `/api/panels/groups` | Panel / groups |
| PUT | `/api/panels/reorder` | Drag-and-drop ordering (ownership checked) |
| GET/PUT | `/api/settings` `/api/settings/batch` | User settings |
| GET | `/api/wallpaper/bing[:date]` | Bing wallpaper (server-cached) |
| GET | `/api/favicon?url=` | Favicon discovery (with SSRF protection) |
| GET/PUT | `/api/admin/settings` | Global site settings (admin) |
| GET/PUT/DELETE | `/api/users` `/api/users/:id` | User management (admin) |
| GET | `/api/system/stats` | System monitor (CPU/memory/disk/network) |
| GET | `/api/weather` | Weather (Open-Meteo) |
| GET/POST | `/api/search` `/api/search/suggestions` | Global search / suggestions |
| GET/POST/PUT/DELETE | `/api/memo` `/api/memo/:id` | Memos: list / create / archive / update / delete |
| GET/POST/PUT/DELETE | `/api/rss` `/api/rss/:id` | RSS feed management |
| GET | `/api/rss/:id/items` | RSS article list (supports `limit`, ownership checked) |
| GET | `/api/bmsync/status` | Bookmark-sync connection status (configured / synced / rev) |
| GET | `/api/bmsync/tree` | Local bookmark mirror (read-only, no network) |
| POST | `/api/bmsync/pull` | Pull the full canonical bookmark tree from the sync server |
| POST | `/api/bmsync/push` | Push create/update/delete changes (LWW + tombstones) and store the returned canonical state |

---

## 📥 Data Import

- Settings drawer → **Import data**: append-merge from a JSON backup file (`{settings, groups}` format)
- **Template**: built-in "Starter Navigation" template (5 groups, 28 bookmarks), one-click import

---

## 🌐 Browser Extension

`extension/` is a Chrome MV3 extension (SunDash Monitor). Load it via the browser's "Developer mode → Load unpacked" pointing at the directory. Features: monitor any page for content changes, change notifications, popup and sidepanel management.

---

## 🤝 Contributing

Issues and PRs are welcome! Please:

- Make sure `go build ./...`, `go test ./...`, and `npm run build` pass before submitting
- Include tests for new features
- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages

---

## 📜 License

[MIT](LICENSE) © 2026 SunDash Contributors

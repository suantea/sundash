# Asuan 部署交接说明（给 workbuddy 执行）

本文件由 opencode 生成，因本机无法连通腾讯云服务器，部署任务转交 workbuddy 执行。

## 一、目标服务器

- 公网 IP：`43.133.45.67`
- SSH 用户：`root`
- SSH 私钥：`D:\下载\temp.pem`（RSA 私钥，OpenSSH 格式，可直接 `-i` 使用）
- 系统：Linux（架构未知，需先探测，大概率 x86_64 / amd64）
- 项目类型：腾讯云轻量应用服务器（Lighthouse）

## 二、连接注意事项（重要）

> 本机（opencode 所在 Windows）测试 `43.133.45.67:22` 一直 **Connection timed out**，
> 排查过密钥、ssh 客户端均正常，判断是服务器侧网络/防火墙/实例状态问题。
> workbuddy 执行前务必先确认：

1. 实例处于「运行中」（轻量服务器关机时端口全不通）。
2. 轻量服务器控制台「防火墙」入站规则放行 `TCP:22`，来源 `0.0.0.0/0`。
3. 公网 IP 确实为 `43.133.45.67`（可能已变更）。
4. 若本机到服务器的 22 仍不通，可让用户从自己电脑连一次：
   `ssh -i D:\下载\temp.pem root@43.133.45.67`，确认服务器本身可用。

连接测试命令（Windows PowerShell）：

```powershell
ssh -i "D:\下载\temp.pem" -o StrictHostKeyChecking=no -o ConnectTimeout=15 root@43.133.45.67 "uname -a; nproc; free -h"
```

## 三、部署内容

项目名 **Asuan**：个人导航/面板应用，前端 Vue 3 + Vite，后端 Go + Gin + SQLite。
Docker 方案见 `docker/` 与 `nas-deploy/`，本次目标是**不用 Docker、直接跑二进制**。

部署后访问：`http://43.133.45.67:3000`，默认账号 `admin / admin`。

### 构件清单（3 样东西放同一目录，例如 `/opt/asuan/`）

| 构件 | 来源 | 说明 |
|---|---|---|
| `sundash` | 交叉编译 | Go 源码 `server/`，Linux/amd64 二进制 |
| `static/` | 前端构建 | `web/` 下 `npm run build` 产物，重命名为 `static/` |
| `data/` | 运行时创建 | SQLite 数据库 `sundash.db` 所在目录 |

### 环境变量（`server/config/config.go`）

| 变量 | 默认值 | 说明 |
|---|---|---|
| `SUNDASH_PORT` | `3000` | 监听端口 |
| `SUNDASH_DATA_DIR` | `./data` | 数据库目录 |
| `SUNDASH_JWT_SECRET` | `sundash-default-secret-change-me` | **务必改成随机长字符串**（默认值已公开） |
| `SUNDASH_DEBUG` | 空 | 设 `true` 开启 debug |

## 四、构建步骤（在本机执行）

前端（在项目根目录 `D:\WEB 开发\index`）：

```powershell
cd "D:\WEB 开发\index\web"
npm install
npm run build
```

后端交叉编译（Go 1.26.4 已装）：

```powershell
cd "D:\WEB 开发\index\server"
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -o sundash .
# 若服务器是 arm64，把 GOARCH 改为 arm64
```

## 五、部署步骤（在服务器上执行）

```bash
mkdir -p /opt/asuan
# 上传 sundash、static/（前端 dist）、data/ 到 /opt/asuan/
cd /opt/asuan
chmod +x sundash
```

启动（前台验证）：

```bash
export SUNDASH_PORT=3000
export SUNDASH_DATA_DIR=/opt/asuan/data
export SUNDASH_JWT_SECRET='换成随机长字符串'
./sundash
```

curl 自测：

```bash
curl -s http://127.0.0.1:3000/api/site-config
```

### 建议注册为 systemd 服务 `/etc/systemd/system/asuan.service`

```ini
[Unit]
Description=Asuan Panel
After=network.target

[Service]
WorkingDirectory=/opt/asuan
ExecStart=/opt/asuan/sundash
Restart=always
RestartSec=3
Environment=SUNDASH_PORT=3000
Environment=SUNDASH_DATA_DIR=/opt/asuan/data
Environment=SUNDASH_JWT_SECRET=换成随机长字符串

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable --now asuan
```

## 六、可选：nginx 反代（80/443 端口访问）

腾讯云轻量服务器防火墙需额外放行 `TCP:80`（和 443）。

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

## 七、验收清单

- [ ] `http://43.133.45.67:3000` 打开页面（静态 + SPA 路由正常）
- [ ] 登录成功（admin / admin）
- [ ] 重启后服务自动拉起（systemd）
- [ ] `data/sundash.db` 已持久化到 `/opt/asuan/data/`
- [ ] JWT 密钥已改，非默认值

## 八、源码说明（给执行者快速定位）

- `server/main.go`：入口，Gin 路由 + 静态托管 + SPA fallback + 站点配置注入
- `server/config/config.go`：环境变量加载
- `server/database/`：SQLite 初始化
- `server/handlers/`：业务接口
- `server/middleware/`：JWT 鉴权 / 管理鉴权 / 登录限流
- `web/`：Vue 3 前端源码（Vite 构建，产物 `web/dist/`）
- `docker/Dockerfile`：官方 Docker 构建方式（交叉编译参考）

---
name: api-documenter
description: 为 Go + Gin 后端（server/handlers/）的 REST API 生成 OpenAPI 文档。用于记录 API 端点、参数、请求/响应结构。
---

为 $ARGUMENTS 指定的 API 端点生成 OpenAPI 3.0 文档。

1. 读取对应 handler 源码（如 server/handlers/auth.go、panel.go、wallpaper.go）
2. 提取：路径、HTTP 方法、认证要求（JWT Bearer）、请求参数、请求/响应 JSON 结构
3. 参考 server/models/ 中的结构体定义 JSON tag 生成 schema
4. 输出 OpenAPI YAML，包含：
   - paths: 每个端点完整定义
   - components.schemas: 复用模型结构
   - securitySchemes: bearerAuth（JWT）

检查清单：
- 每个端点都有 summary 和 description
- 参数标注 required 与类型
- 错误响应（400/401/403/429）有描述

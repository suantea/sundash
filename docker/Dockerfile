# Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build
# Pre-compress static assets for better performance (node:20-alpine lacks brotli)
RUN apk add --no-cache brotli && \
    find dist -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" -o -name "*.svg" -o -name "*.json" \) | while read f; do \
      gzip -k -f "$f" && brotli -k -f "$f"; \
    done

# Build backend (match the Go version required by server/go.mod)
FROM golang:1.26-alpine AS backend
# 容器内默后默认 proxy.golang.org 易超时，改用国内镜像
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app/server
COPY server/go.* ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o sundash .

# Runtime: unprivileged user + healthcheck
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S sundash && adduser -S sundash -G sundash
WORKDIR /app
COPY --from=backend /app/server/sundash .
COPY --from=frontend /app/web/dist ./static
RUN mkdir -p /app/data && chown -R sundash:sundash /app

USER sundash
EXPOSE 3000
ENV SUNDASH_PORT=3000
ENV SUNDASH_DATA_DIR=/app/data

VOLUME ["/app/data"]
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:3000/api/site-config || exit 1

CMD ["./sundash"]

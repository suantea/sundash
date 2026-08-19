package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	DataDir        string
	Debug          bool
	AllowedOrigins []string
	// MCPToken 是 MCP 端点专用静态 token（可选）。设置后 AI 客户端用
	// Authorization: Bearer <MCPToken> 访问 /mcp，绑定到 MCPUsername 指定的用户。
	MCPToken    string
	MCPUsername string
}

func Load() *Config {
	debug := os.Getenv("SUNDASH_DEBUG") == "true"
	port := getEnv("SUNDASH_PORT", "3000")
	dataDir := getEnv("SUNDASH_DATA_DIR", "./data")

	jwtSecret := os.Getenv("SUNDASH_JWT_SECRET")
	if jwtSecret == "" {
		if debug {
			jwtSecret = "dev-insecure-secret"
			log.Println("WARNING: SUNDASH_DEBUG=true, using insecure development JWT secret")
		} else {
			log.Fatal("SUNDASH_JWT_SECRET environment variable is required. Set a strong secret, or run with SUNDASH_DEBUG=true for development only.")
		}
	}

	// Comma-separated list of allowed CORS origins. Empty = same-origin only (secure default).
	var origins []string
	if raw := os.Getenv("SUNDASH_ALLOWED_ORIGINS"); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if o = strings.TrimSpace(o); o != "" {
				origins = append(origins, o)
			}
		}
	}

	return &Config{
		Port:           port,
		DBPath:         filepath.Join(dataDir, "sundash.db"),
		JWTSecret:      jwtSecret,
		DataDir:        dataDir,
		Debug:          debug,
		AllowedOrigins: origins,
		MCPToken:       os.Getenv("SUNDASH_MCP_TOKEN"),
		MCPUsername:    getEnv("SUNDASH_MCP_USERNAME", "admin"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"sundash/config"
	"sundash/database"
	"sundash/handlers"
	"sundash/mcp"
	"sundash/middleware"
	"sundash/repository"
	"sundash/service"

	"github.com/gin-gonic/gin"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Load()

	db, err := database.Init(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Wire the dependency graph: repository → service → handler.
	userRepo := repository.NewUserRepo(db)
	panelRepo := repository.NewPanelRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)

	userSvc := service.NewUserService(userRepo, settingsRepo, []byte(cfg.JWTSecret))
	panelSvc := service.NewPanelService(panelRepo, userRepo)
	settingsSvc := service.NewSettingsService(settingsRepo)
	wallpaperSvc := service.NewWallpaperService()
	faviconSvc := service.NewFaviconService()
	systemSvc := service.NewSystemService()
	searchSvc := service.NewSearchService(db)
	weatherSvc := service.NewWeatherService()
	memoSvc := service.NewMemoService(repository.NewMemoRepo(db))
	rssSvc := service.NewRSSService(repository.NewRSSFeedRepo(db), repository.NewRSSItemRepo(db))
	bmsyncSvc := service.NewBmsyncService(repository.NewBmsyncRepo(db), settingsRepo)
	dockerSvc := service.NewDockerService()

	authH := handlers.NewAuthHandler(userSvc, panelSvc, settingsSvc)
	panelH := handlers.NewPanelHandler(panelSvc, settingsSvc)
	wallpaperH := handlers.NewWallpaperHandler(wallpaperSvc)
	faviconH := handlers.NewFaviconHandler(faviconSvc)
	systemH := handlers.NewSystemHandler(systemSvc)
	searchH := handlers.NewSearchHandler(searchSvc)
	weatherH := handlers.NewWeatherHandler(weatherSvc)
	memoH := handlers.NewMemoHandler(memoSvc)
	rssH := handlers.NewRSSHandler(rssSvc)
	bmsyncH := handlers.NewBmsyncHandler(bmsyncSvc)
	bootstrapH := handlers.NewBootstrapHandler(userSvc, panelSvc, settingsSvc)
	backupH := handlers.NewBackupHandler(db)

	// MCP server: AI agents can list / create / organize bookmarks, search
	// cards and manage memos.
	// Endpoint: POST /mcp (Streamable HTTP). Auth: Bearer <SUNDASH_MCP_TOKEN>
	// or a regular sundash JWT; the resolved user is bound to the session.
	mcpSrv := mcp.New(panelSvc, faviconSvc, systemSvc, searchSvc, memoSvc, dockerSvc)
	mcpHTTP := mcpserver.NewStreamableHTTPServer(mcpSrv.MCPServer())

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.Gzip())

	// Long-lived cache for hashed build assets (immutable filenames).
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	})

	staticPath := filepath.Join(".", "static")
	if _, err := os.Stat(staticPath); err == nil {
		r.Static("/assets", filepath.Join(staticPath, "assets"))
		r.StaticFile("/vite.svg", filepath.Join(staticPath, "vite.svg"))
		// web/public/favicon.svg ships with the build; without this route it
		// falls through to the SPA handler and the tab icon gets HTML.
		if _, err := os.Stat(filepath.Join(staticPath, "favicon.svg")); err == nil {
			r.StaticFile("/favicon.svg", filepath.Join(staticPath, "favicon.svg"))
		}
	}

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.LoginRateLimit(), authH.Login)
			auth.POST("/register", authH.Register)
			auth.GET("/settings", authH.GetAuthSettings)
			auth.GET("/setup-status", authH.SetupStatus)
			auth.POST("/setup", authH.Setup)
		}

		api.GET("/site-config", authH.GetSiteConfig)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret)))
		{
			protected.GET("/profile", authH.GetProfile)
			protected.PUT("/profile", authH.UpdateProfile)
			protected.PUT("/profile/password", authH.ChangePassword)

			// Home page bootstrap: settings + profile + panels in one round-trip.
			protected.GET("/bootstrap", bootstrapH.GetBootstrap)

			protected.GET("/panels", panelH.GetPanel)
			protected.POST("/panels/groups", panelH.CreateGroup)
			protected.PUT("/panels/groups/:id", panelH.UpdateGroup)
			protected.DELETE("/panels/groups/:id", panelH.DeleteGroup)
			protected.POST("/panels/cards", panelH.CreateCard)
			protected.PUT("/panels/cards/:id", panelH.UpdateCard)
			protected.DELETE("/panels/cards/:id", panelH.DeleteCard)
			protected.PUT("/panels/reorder", panelH.Reorder)

			protected.GET("/settings", panelH.GetSettings)
			protected.PUT("/settings", panelH.UpdateSetting)
			protected.PUT("/settings/batch", panelH.BatchUpdateSettings)

			protected.GET("/wallpaper/bing", wallpaperH.GetBingWallpaper)
			protected.GET("/wallpaper/bing/:date", wallpaperH.GetWallpaperByDate)

			protected.GET("/favicon", faviconH.FetchFavicon)
			protected.GET("/system/stats", systemH.GetStats)

			// Search API
			protected.GET("/search", searchH.Search)
			protected.GET("/search/suggestions", searchH.Suggestions)

			// Weather API
			protected.GET("/weather", weatherH.GetWeather)

			// Memo API
			protected.GET("/memo", memoH.ListMemos)
			protected.POST("/memo", memoH.CreateMemo)
			protected.PUT("/memo/:id/archive", memoH.ArchiveMemo)
			protected.DELETE("/memo/:id", memoH.DeleteMemo)
			protected.PUT("/memo/:id", memoH.UpdateMemo)

			// RSS API
			protected.GET("/rss", rssH.ListFeeds)
			protected.POST("/rss", rssH.AddFeed)
			protected.PUT("/rss/:id", rssH.UpdateFeed)
			protected.DELETE("/rss/:id", rssH.DeleteFeed)
			protected.GET("/rss/:id/items", rssH.GetFeedItems)

			// Bookmark-sync: full tree mirror driven by the canonical
			// bookmark-sync server. Pull refreshes the mirror from the
			// server; Push uploads local changes and stores the returned
			// canonical state. All endpoints require login.
			protected.GET("/bmsync/status", bmsyncH.Status)
			protected.GET("/bmsync/tree", bmsyncH.Tree)
			protected.POST("/bmsync/pull", bmsyncH.Pull)
			protected.POST("/bmsync/push", bmsyncH.Push)

			admin := protected.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", authH.ListUsers)
				admin.PUT("/users/:id", authH.UpdateUser)
				admin.DELETE("/users/:id", authH.DeleteUser)
				admin.POST("/users/:id/reset-password", authH.AdminResetPassword)
				admin.POST("/users/:id/approve", authH.ApproveUser)
				admin.POST("/users/:id/reject", authH.RejectUser)
				admin.GET("/users/:id/panel", authH.GetUserPanel)
				admin.GET("/admin/settings", authH.GetGlobalSettings)
				admin.PUT("/admin/settings", authH.UpdateGlobalSettings)
				// Full SQLite snapshot download (consistent while live, #备份).
				admin.GET("/admin/backup", backupH.Download)
			}
		}
	}

	// MCP endpoint (Streamable HTTP). Authenticate via static token or JWT,
	// then inject the resolved user id into the request context.
	r.Any("/mcp", func(c *gin.Context) {
		uid := mcp.Auth(c.Request, userRepo, []byte(cfg.JWTSecret), cfg.MCPToken, cfg.MCPUsername)
		if uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing MCP token"})
			return
		}
		req := c.Request.WithContext(mcp.WithUserID(c.Request.Context(), uid))
		mcpHTTP.ServeHTTP(c.Writer, req)
	})

	// SPA fallback with cached index.html + per-request site config injection.
	spa := newSPAHandler(staticPath, settingsSvc)
	r.NoRoute(spa.ServeHTTP)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("SunDash server starting", "addr", "http://0.0.0.0:"+cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	// Stop background workers after the HTTP server drained.
	rssSvc.Stop()
	slog.Info("Server exited")
}

// spaHandler serves index.html with site settings injected, caching the file.
type spaHandler struct {
	staticPath string
	settings   *service.SettingsService

	mu      sync.Mutex
	html    []byte
	loaded  bool
	loadErr error
}

func newSPAHandler(staticPath string, settings *service.SettingsService) *spaHandler {
	return &spaHandler{staticPath: staticPath, settings: settings}
}

func (s *spaHandler) ServeHTTP(c *gin.Context) {
	indexPath := filepath.Join(s.staticPath, "index.html")

	s.mu.Lock()
	if !s.loaded {
		content, err := os.ReadFile(indexPath)
		s.html, s.loadErr = content, err
		s.loaded = true
	}
	html := s.html
	loadErr := s.loadErr
	s.mu.Unlock()

	if loadErr != nil {
		slog.Error("failed to read index.html", "error", loadErr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	settings, err := s.settings.GetSiteConfig()
	if err != nil {
		slog.Error("failed to load site config", "error", err)
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(injectSiteConfig(string(html), settings)))
}

// injectSiteConfig replaces placeholders in the built HTML with site settings.
func injectSiteConfig(html string, settings map[string]string) string {
	if title := settings["site_title"]; title != "" {
		html = strings.Replace(html, "<title>SunDash</title>", fmt.Sprintf("<title>%s</title>", title), 1)
	}
	if icon := settings["site_icon_url"]; icon != "" {
		html = strings.Replace(html, "/favicon.svg", icon, 1)
	}
	if cdn := settings["site_cdn_url"]; cdn != "" {
		html = strings.ReplaceAll(html, "/assets/", cdn+"/assets/")
	}

	var headInjection []string
	if customHead := settings["site_custom_head"]; customHead != "" {
		headInjection = append(headInjection, customHead)
	}
	if analytics := settings["site_analytics_code"]; analytics != "" {
		headInjection = append(headInjection, analytics)
	}
	if len(headInjection) > 0 {
		html = strings.Replace(html, "</head>", strings.Join(headInjection, "\n")+"\n</head>", 1)
	}
	if footer := settings["site_custom_footer"]; footer != "" {
		html = strings.Replace(html, "</body>", footer+"\n</body>", 1)
	}
	return html
}

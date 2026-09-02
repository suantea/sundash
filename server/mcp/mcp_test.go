package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"sundash/database"
	"sundash/models"
	"sundash/repository"
	"sundash/service"

	"github.com/mark3labs/mcp-go/mcp"
)

func setup(t *testing.T) *Server {
	t.Helper()
	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	userRepo := repository.NewUserRepo(db)
	// database.Init no longer auto-creates an admin; create one explicitly in
	// the SAME db so FK constraints (groups.user_id → users.id) hold.
	ensureAdmin(t, userRepo)
	panelSvc := service.NewPanelService(repository.NewPanelRepo(db), userRepo)
	faviconSvc := service.NewFaviconService()
	systemSvc := service.NewSystemService()
	searchSvc := service.NewSearchService(db)
	memoSvc := service.NewMemoService(repository.NewMemoRepo(db))
	return New(panelSvc, faviconSvc, systemSvc, searchSvc, memoSvc)
}

// ensureAdmin inserts the fixed admin user if missing.
func ensureAdmin(t *testing.T, repo *repository.UserRepo) {
	t.Helper()
	u, err := repo.FindByUsername("admin")
	if err != nil {
		t.Fatalf("find admin: %v", err)
	}
	if u == nil {
		admin := &models.User{
			ID:           "admin",
			Username:     "admin",
			PasswordHash: "",
			DisplayName:  "Administrator",
			Role:         "admin",
			Status:       "approved",
		}
		if err := repo.Create(admin); err != nil {
			t.Fatalf("create admin: %v", err)
		}
	}
}

// adminID returns the id of the admin user created by setup().
func adminID(t *testing.T) string {
	t.Helper()
	return "admin"
}

// withUser wraps the authenticated user id into ctx, same as the HTTP path.
func withUser(t *testing.T, uid string) context.Context {
	t.Helper()
	return WithUserID(context.Background(), uid)
}

// invoke calls a tool handler directly (in-process, no HTTP / no token).
func invoke(t *testing.T, s *Server, ctx context.Context, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	var res *mcp.CallToolResult
	var err error
	switch name {
	case "sundash_list_groups":
		res, err = s.handleListGroups(ctx, req)
	case "sundash_create_group":
		res, err = s.handleCreateGroup(ctx, req)
	case "sundash_rename_group":
		res, err = s.handleRenameGroup(ctx, req)
	case "sundash_delete_group":
		res, err = s.handleDeleteGroup(ctx, req)
	case "sundash_create_card":
		res, err = s.handleCreateCard(ctx, req)
	case "sundash_update_card":
		res, err = s.handleUpdateCard(ctx, req)
	case "sundash_move_card":
		res, err = s.handleMoveCard(ctx, req)
	case "sundash_delete_card":
		res, err = s.handleDeleteCard(ctx, req)
	default:
		t.Fatalf("unknown tool %q", name)
	}
	if err != nil {
		t.Fatalf("tool %s error: %v", name, err)
	}
	return res
}

// resultText extracts the text payload from a CallToolResult.
func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		return ""
	}
	tc, ok := res.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return tc.Text
}

func parseGroupID(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	var g models.PanelGroup
	if err := json.Unmarshal([]byte(resultText(t, res)), &g); err != nil {
		t.Fatalf("parse group: %v", err)
	}
	return g.ID
}

func parseCardID(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	var c models.Card
	if err := json.Unmarshal([]byte(resultText(t, res)), &c); err != nil {
		t.Fatalf("parse card: %v", err)
	}
	return c.ID
}

func mustList(t *testing.T, s *Server, ctx context.Context) *models.PanelData {
	t.Helper()
	res := invoke(t, s, ctx, "sundash_list_groups", map[string]any{})
	if res.IsError {
		t.Fatalf("list groups error: %s", resultText(t, res))
	}
	var data models.PanelData
	if err := json.Unmarshal([]byte(resultText(t, res)), &data); err != nil {
		t.Fatalf("parse panel data: %v", err)
	}
	return &data
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }

func TestMCPAuthRequired(t *testing.T) {
	s := setup(t)
	// 未注入 user_id 的 ctx → 工具应返回错误信息而非执行。
	res := invoke(t, s, context.Background(), "sundash_list_groups", map[string]any{})
	if !res.IsError {
		t.Fatalf("expected auth error, got %+v", res)
	}
}

func TestMCPGroupAndCardLifecycle(t *testing.T) {
	s := setup(t)
	ctx := withUser(t, adminID(t))

	// 1) 初始为空
	res := invoke(t, s, ctx, "sundash_list_groups", map[string]any{})
	if res.IsError {
		t.Fatalf("list groups error: %s", resultText(t, res))
	}

	// 2) 创建分组
	res = invoke(t, s, ctx, "sundash_create_group", map[string]any{"name": "开发工具"})
	if res.IsError {
		t.Fatalf("create group error: %s", resultText(t, res))
	}
	groupID := parseGroupID(t, res)

	// 3) 重命名分组
	res = invoke(t, s, ctx, "sundash_rename_group", map[string]any{"group_id": groupID, "name": "Dev Tools"})
	if res.IsError {
		t.Fatalf("rename group error: %s", resultText(t, res))
	}

	// 4) 组内新建卡片
	res = invoke(t, s, ctx, "sundash_create_card", map[string]any{
		"group_id":    groupID,
		"title":       "GitHub",
		"url":         "https://github.com",
		"description": "代码托管",
	})
	if res.IsError {
		t.Fatalf("create card error: %s", resultText(t, res))
	}
	cardID := parseCardID(t, res)

	// 5) 更新卡片
	res = invoke(t, s, ctx, "sundash_update_card", map[string]any{
		"card_id":     cardID,
		"description": "代码托管 + CI",
	})
	if res.IsError {
		t.Fatalf("update card error: %s", resultText(t, res))
	}

	// 6) 列表确认数据落库
	res = invoke(t, s, ctx, "sundash_list_groups", map[string]any{})
	if res.IsError {
		t.Fatalf("list groups error: %s", resultText(t, res))
	}
	text := resultText(t, res)
	if !contains(text, "Dev Tools") || !contains(text, "GitHub") || !contains(text, "代码托管 + CI") {
		t.Fatalf("list result missing expected data: %s", text)
	}

	// 7) 删除卡片
	res = invoke(t, s, ctx, "sundash_delete_card", map[string]any{"card_id": cardID})
	if res.IsError {
		t.Fatalf("delete card error: %s", resultText(t, res))
	}

	// 8) 删除分组
	res = invoke(t, s, ctx, "sundash_delete_group", map[string]any{"group_id": groupID})
	if res.IsError {
		t.Fatalf("delete group error: %s", resultText(t, res))
	}
}

func TestMCPMoveCard(t *testing.T) {
	s := setup(t)
	ctx := withUser(t, adminID(t))

	g1 := parseGroupID(t, invoke(t, s, ctx, "sundash_create_group", map[string]any{"name": "A"}))
	g2 := parseGroupID(t, invoke(t, s, ctx, "sundash_create_group", map[string]any{"name": "B"}))
	cardID := parseCardID(t, invoke(t, s, ctx, "sundash_create_card", map[string]any{
		"group_id": g1, "title": "C1", "url": "https://x.com",
	}))

	res := invoke(t, s, ctx, "sundash_move_card", map[string]any{"card_id": cardID, "group_id": g2})
	if res.IsError {
		t.Fatalf("move card error: %s", resultText(t, res))
	}

	// 确认卡片已落在 g2 下
	data := mustList(t, s, ctx)
	for _, g := range data.Groups {
		if g.ID == g2 {
			if len(g.Cards) != 1 || g.Cards[0].ID != cardID {
				t.Fatalf("card not moved into g2: %+v", g.Cards)
			}
		}
		if g.ID == g1 && len(g.Cards) != 0 {
			t.Fatalf("g1 should be empty, got %+v", g.Cards)
		}
	}
}

func TestMCPMissingArgs(t *testing.T) {
	s := setup(t)
	ctx := withUser(t, adminID(t))

	// create_card 缺 url
	res := invoke(t, s, ctx, "sundash_create_card", map[string]any{"group_id": "g", "title": "T"})
	if !res.IsError {
		t.Fatalf("expected error for missing url, got %+v", res)
	}
	// create_group 缺 name
	res = invoke(t, s, ctx, "sundash_create_group", map[string]any{})
	if !res.IsError {
		t.Fatalf("expected error for missing name, got %+v", res)
	}
}

// Package mcp exposes sundash's bookmark management as an MCP (Model Context
// Protocol) server, so AI agents can list / create / update / organize groups
// and cards on behalf of a user.
//
// Wire-up: main.go mounts the Streamable HTTP handler behind JWT auth and
// injects the user id into the request context via [UserIDContextKey].
package mcp

import (
	"context"
	"fmt"

	"sundash/models"
	"sundash/service"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// UserIDContextKey is the context key under which the authenticated user id is
// stored (injected by the auth middleware in main.go).
type ctxKey string

const UserIDContextKey ctxKey = "sundash_user_id"

// Server wraps the mcp-go server with access to the panel service.
type Server struct {
	server *server.MCPServer
	panels *service.PanelService
}

// New creates an MCPServer with all sundash bookmark-management tools.
func New(panels *service.PanelService) *Server {
	s := &Server{panels: panels}
	mcps := server.NewMCPServer(
		"sundash",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerGroupTools(mcps, s)
	registerCardTools(mcps, s)

	s.server = mcps
	return s
}

// MCPServer returns the underlying mcp-go server (for StreamableHTTP).
func (s *Server) MCPServer() *server.MCPServer { return s.server }

// --- helpers ---------------------------------------------------------------

func (s *Server) userID(ctx context.Context) (string, error) {
	uid, ok := ctx.Value(UserIDContextKey).(string)
	if !ok || uid == "" {
		return "", fmt.Errorf("未认证: 缺少 user_id（请通过 Sundash MCP token 认证）")
	}
	return uid, nil
}

func strArg(args any, key string) string {
	m, ok := args.(map[string]any)
	if !ok {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func strPtrArg(args any, key string) *string {
	v := strArg(args, key)
	if v == "" {
		return nil
	}
	return &v
}

func intArg(args any, key string) *int {
	m, ok := args.(map[string]any)
	if !ok {
		return nil
	}
	switch v := m[key].(type) {
	case float64:
		n := int(v)
		return &n
	case int:
		return &v
	default:
		return nil
	}
}

// --- group tools -----------------------------------------------------------

func registerGroupTools(mcps *server.MCPServer, s *Server) {
	// sundash_list_groups: list all groups (with their cards) for the user.
	mcps.AddTool(mcp.NewTool(
		"sundash_list_groups",
		mcp.WithDescription("列出当前用户的全部分组及组内卡片。无参数。返回 groups[]（含 id/name/sort_order/cards[]）。"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleListGroups)

	// sundash_create_group
	mcps.AddTool(mcp.NewTool(
		"sundash_create_group",
		mcp.WithDescription("新建一个分组（标签分类）。返回新分组的 id/name。"),
		mcp.WithString("name", mcp.Required(), mcp.Description("分组名称，如「开发工具」「影音娱乐」")),
	), s.handleCreateGroup)

	// sundash_rename_group
	mcps.AddTool(mcp.NewTool(
		"sundash_rename_group",
		mcp.WithDescription("重命名分组。"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("分组 id")),
		mcp.WithString("name", mcp.Required(), mcp.Description("新的分组名称")),
	), s.handleRenameGroup)

	// sundash_delete_group
	mcps.AddTool(mcp.NewTool(
		"sundash_delete_group",
		mcp.WithDescription("删除分组（组内卡片会一并删除）。"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("分组 id")),
		mcp.WithDestructiveHintAnnotation(true),
	), s.handleDeleteGroup)
}

func (s *Server) handleListGroups(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	data, err := s.panels.GetPanel(uid)
	if err != nil {
		return mcp.NewToolResultError("加载面板失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(data)
}

func (s *Server) handleCreateGroup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	name := strArg(req.Params.Arguments, "name")
	if name == "" {
		return mcp.NewToolResultError("缺少参数 name"), nil
	}
	group, err := s.panels.CreateGroup(uid, name)
	if err != nil {
		return mcp.NewToolResultError("创建分组失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(group)
}

func (s *Server) handleRenameGroup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	groupID := strArg(req.Params.Arguments, "group_id")
	name := strArg(req.Params.Arguments, "name")
	if groupID == "" || name == "" {
		return mcp.NewToolResultError("缺少参数 group_id 或 name"), nil
	}
	if err := s.panels.UpdateGroup(uid, groupID, &name, nil); err != nil {
		return mcp.NewToolResultError("重命名失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已重命名分组 %s 为 %s", groupID, name)), nil
}

func (s *Server) handleDeleteGroup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	groupID := strArg(req.Params.Arguments, "group_id")
	if groupID == "" {
		return mcp.NewToolResultError("缺少参数 group_id"), nil
	}
	if err := s.panels.DeleteGroup(uid, groupID); err != nil {
		return mcp.NewToolResultError("删除分组失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已删除分组 %s", groupID)), nil
}

// --- card tools ------------------------------------------------------------

func registerCardTools(mcps *server.MCPServer, s *Server) {
	// sundash_create_card
	mcps.AddTool(mcp.NewTool(
		"sundash_create_card",
		mcp.WithDescription("在指定分组内新增一个书签卡片（录入）。title/url 必填，其余可选。"),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("目标分组 id（可用 sundash_list_groups 查询）")),
		mcp.WithString("title", mcp.Required(), mcp.Description("书签标题")),
		mcp.WithString("url", mcp.Required(), mcp.Description("书签链接")),
		mcp.WithString("description", mcp.Description("描述/备注")),
		mcp.WithString("url_internal", mcp.Description("内网链接（可选）")),
		mcp.WithString("open_type", mcp.Description("打开方式: new_tab/current/panel，默认 new_tab")),
	), s.handleCreateCard)

	// sundash_update_card
	mcps.AddTool(mcp.NewTool(
		"sundash_update_card",
		mcp.WithDescription("更新书签卡片的字段（只更新传入的字段）。"),
		mcp.WithString("card_id", mcp.Required(), mcp.Description("卡片 id")),
		mcp.WithString("title", mcp.Description("标题")),
		mcp.WithString("url", mcp.Description("链接")),
		mcp.WithString("description", mcp.Description("描述")),
		mcp.WithString("url_internal", mcp.Description("内网链接")),
	), s.handleUpdateCard)

	// sundash_move_card
	mcps.AddTool(mcp.NewTool(
		"sundash_move_card",
		mcp.WithDescription("把卡片移动到另一个分组（整理归类）。"),
		mcp.WithString("card_id", mcp.Required(), mcp.Description("卡片 id")),
		mcp.WithString("group_id", mcp.Required(), mcp.Description("目标分组 id")),
	), s.handleMoveCard)

	// sundash_delete_card
	mcps.AddTool(mcp.NewTool(
		"sundash_delete_card",
		mcp.WithDescription("删除书签卡片。"),
		mcp.WithString("card_id", mcp.Required(), mcp.Description("卡片 id")),
		mcp.WithDestructiveHintAnnotation(true),
	), s.handleDeleteCard)
}

func (s *Server) handleCreateCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	args := req.Params.Arguments
	groupID, title, url := strArg(args, "group_id"), strArg(args, "title"), strArg(args, "url")
	if groupID == "" || title == "" || url == "" {
		return mcp.NewToolResultError("缺少参数 group_id/title/url"), nil
	}
	openType := strArg(args, "open_type")
	if openType == "" {
		openType = "new_tab"
	}
	card, err := s.panels.CreateCard(uid, models.CreateCardRequest{
		GroupID:     groupID,
		Title:       title,
		URL:         url,
		URLInternal: strArg(args, "url_internal"),
		Description: strArg(args, "description"),
		OpenType:    openType,
	})
	if err != nil {
		return mcp.NewToolResultError("创建卡片失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(card)
}

func (s *Server) handleUpdateCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	args := req.Params.Arguments
	cardID := strArg(args, "card_id")
	if cardID == "" {
		return mcp.NewToolResultError("缺少参数 card_id"), nil
	}
	updateReq := models.UpdateCardRequest{
		Title:       strPtrArg(args, "title"),
		URL:         strPtrArg(args, "url"),
		URLInternal: strPtrArg(args, "url_internal"),
		Description: strPtrArg(args, "description"),
	}
	if err := s.panels.UpdateCard(uid, cardID, updateReq); err != nil {
		return mcp.NewToolResultError("更新卡片失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已更新卡片 %s", cardID)), nil
}

func (s *Server) handleMoveCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	args := req.Params.Arguments
	cardID := strArg(args, "card_id")
	groupID := strArg(args, "group_id")
	if cardID == "" || groupID == "" {
		return mcp.NewToolResultError("缺少参数 card_id 或 group_id"), nil
	}
	// Move = update the group via Reorder with a single card order entry.
	if err := s.panels.Reorder(uid, models.ReorderRequest{
		CardOrders: []models.CardOrder{{ID: cardID, GroupID: groupID, SortOrder: 0}},
	}); err != nil {
		return mcp.NewToolResultError("移动卡片失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已把卡片 %s 移动到分组 %s", cardID, groupID)), nil
}

func (s *Server) handleDeleteCard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cardID := strArg(req.Params.Arguments, "card_id")
	if cardID == "" {
		return mcp.NewToolResultError("缺少参数 card_id"), nil
	}
	if err := s.panels.DeleteCard(uid, cardID); err != nil {
		return mcp.NewToolResultError("删除卡片失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已删除卡片 %s", cardID)), nil
}

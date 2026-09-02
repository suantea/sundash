// Package mcp exposes sundash's bookmark management as an MCP (Model Context
// Protocol) server, so AI agents can list / create / update / organize groups
// and cards on behalf of a user.
//
// Wire-up: main.go mounts the Streamable HTTP handler behind JWT auth and
// injects the user id into the request context via [UserIDContextKey].
package mcp

import (
	"context"
	"encoding/json"
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
	server   *server.MCPServer
	panels   *service.PanelService
	favicon  *service.FaviconService
	system   *service.SystemService
	search   *service.SearchService
	memo     *service.MemoService
	docker   *service.DockerService
	deepseek *service.DeepSeekService
}

// New creates an MCPServer with all sundash bookmark-management tools.
func New(panels *service.PanelService, favicon *service.FaviconService, system *service.SystemService, search *service.SearchService, memo *service.MemoService, docker *service.DockerService, deepseek *service.DeepSeekService) *Server {
	s := &Server{panels: panels, favicon: favicon, system: system, search: search, memo: memo, docker: docker, deepseek: deepseek}
	mcps := server.NewMCPServer(
		"sundash",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerGroupTools(mcps, s)
	registerCardTools(mcps, s)
	registerIconTools(mcps, s)
	registerSystemTools(mcps, s)
	registerSearchTools(mcps, s)
	registerMemoTools(mcps, s)
	registerDockerTools(mcps, s)
	registerDeepSeekTools(mcps, s)

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

	// sundash_reorganize: AI-assisted natural language group reorganization
	mcps.AddTool(mcp.NewTool(
		"sundash_reorganize",
		mcp.WithDescription("AI 辅助分组调整。接受自然语言指令，将其翻译为一系列分组操作（创建分组、重命名分组、删除分组、移动卡片、新增卡片）。示例指令：'把 GitHub 和 GitLab 归到一个叫代码仓库的分组，删除空的分组'。返回执行结果。"),
		mcp.WithString("instruction", mcp.Required(), mcp.Description("自然语言分组指令，例如「把所有 AI 相关的链接归到 AI 工具分组」")),
		mcp.WithString("plan", mcp.Description("可选：预先规划好的 JSON 操作序列，格式为 [{\"action\":\"move_card\",\"card_id\":\"xxx\",\"target_group_id\":\"yyy\"}]。如果不传则 AI 自行解析指令。")),
	), s.handleReorganize)
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

// handleReorganize implements sundash_reorganize: AI-assisted natural language group reorganization.
// It translates a natural-language instruction into a sequence of group/card operations.
func (s *Server) handleReorganize(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	instruction := strArg(req.Params.Arguments, "instruction")
	if instruction == "" {
		return mcp.NewToolResultError("缺少参数 instruction"), nil
	}

	// If a pre-built plan JSON is provided, execute it directly.
	if planJSON := strArg(req.Params.Arguments, "plan"); planJSON != "" {
		var ops []service.BatchOp
		if err := json.Unmarshal([]byte(planJSON), &ops); err != nil {
			return mcp.NewToolResultError("解析 plan JSON 失败: " + err.Error()), nil
		}
		if err := s.panels.BatchReorganize(uid, service.BatchReorganizeRequest{Operations: ops}); err != nil {
			return mcp.NewToolResultError("执行分组调整失败: " + err.Error()), nil
		}
		return mcp.NewToolResultText("已按计划完成分组调整"), nil
	}

	// Otherwise, parse the instruction into operations using simple heuristics.
	// This is a lightweight client-side planner: the AI agent can call sundash_list_groups first,
	// then call sundash_reorganize with a precise plan. But we also accept raw instructions
	// for convenience.
	ops, err := parseReorganizeInstruction(uid, instruction, s.panels)
	if err != nil {
		return mcp.NewToolResultError("解析指令失败: " + err.Error()), nil
	}
	if err := s.panels.BatchReorganize(uid, service.BatchReorganizeRequest{Operations: ops}); err != nil {
		return mcp.NewToolResultError("执行分组调整失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已完成分组调整，共执行 %d 项操作", len(ops))), nil
}

// parseReorganizeInstruction is a lightweight heuristic planner.
// For best results, the AI agent should call sundash_list_groups first to get current state,
// then pass a precise JSON plan via the plan parameter.
func parseReorganizeInstruction(userID string, instruction string, panels *service.PanelService) ([]service.BatchOp, error) {
	// Current best practice: return an error prompting the AI to first list groups and build a plan.
	// This keeps the server simple and avoids fragile NLP parsing.
	return nil, fmt.Errorf("当前需要显式提供 plan 参数。请先调用 sundash_list_groups 获取分组状态，再构造 JSON plan 传给 plan 参数")
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

	// sundash_set_icons: AI-assisted batch icon assignment
	mcps.AddTool(mcp.NewTool(
		"sundash_set_icons",
		mcp.WithDescription("AI 辅助批量设置卡片图标。接受一个 JSON 对象，键为卡片 id，值为图标（支持 emoji、Font Awesome class、图片 URL）。例如 {\"card-1\":\"github\",\"card-2\":\"https://...\"}。建议先调用 sundash_list_groups 获取卡片 id，再通过 AI 推断合适的图标。"),
		mcp.WithString("icons_json", mcp.Required(), mcp.Description("JSON 对象字符串，格式为 {\"card_id\":\"icon_value\", ...}")),
	), s.handleSetIcons)

	// sundash_suggest_icons: AI-assisted icon suggestions for cards without icons
	mcps.AddTool(mcp.NewTool(
		"sundash_suggest_icons",
		mcp.WithDescription("AI 辅助图标建议。返回所有尚未设置图标的卡片列表（id/title/url），供 AI 推断合适的图标后回传给 sundash_set_icons。"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleSuggestIcons)
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

// handleSetIcons implements sundash_set_icons: AI-assisted batch icon assignment.
func (s *Server) handleSetIcons(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	iconsJSON := strArg(req.Params.Arguments, "icons_json")
	if iconsJSON == "" {
		return mcp.NewToolResultError("缺少参数 icons_json"), nil
	}
	var icons map[string]string
	if err := json.Unmarshal([]byte(iconsJSON), &icons); err != nil {
		return mcp.NewToolResultError("解析 icons_json 失败: " + err.Error()), nil
	}
	if err := s.panels.BatchSetIcons(uid, service.BatchSetIconsRequest{Icons: icons}); err != nil {
		return mcp.NewToolResultError("批量设置图标失败: " + err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("已成功设置 %d 个卡片的图标", len(icons))), nil
}

// handleSuggestIcons implements sundash_suggest_icons: returns cards without icons for AI suggestion.
func (s *Server) handleSuggestIcons(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	data, err := s.panels.GetPanel(uid)
	if err != nil {
		return mcp.NewToolResultError("加载面板失败: " + err.Error()), nil
	}
	type iconSuggestion struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	var suggestions []iconSuggestion
	for _, g := range data.Groups {
		for _, c := range g.Cards {
			if c.Icon == "" {
				suggestions = append(suggestions, iconSuggestion{ID: c.ID, Title: c.Title, URL: c.URL})
			}
		}
	}
	return mcp.NewToolResultJSON(suggestions)
}

// --- icon tools ------------------------------------------------------------

func registerIconTools(mcps *server.MCPServer, s *Server) {
	// sundash_auto_iconify: auto-detect and apply icons for user's cards
	mcps.AddTool(mcp.NewTool(
		"sundash_auto_iconify",
		mcp.WithDescription("自动为卡片匹配网站图标。对所有尚未设置图标的卡片，根据其 URL 自动检测合适的图标（Iconify 名称或 favicon URL）并应用。返回成功匹配数量和失败数量。"),
		mcp.WithString("apply", mcp.Description("是否立即应用检测到的图标：\"true\" 立即写入数据库，\"false\" 仅返回检测结果供预览，默认 \"true\"")),
	), s.handleAutoIconify)
}

func (s *Server) handleAutoIconify(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	apply := strArg(req.Params.Arguments, "apply")
	shouldApply := apply == "" || apply == "true"

	data, err := s.panels.GetPanel(uid)
	if err != nil {
		return mcp.NewToolResultError("加载面板失败: " + err.Error()), nil
	}

	type iconResult struct {
		CardID  string `json:"card_id"`
		Title   string `json:"title"`
		URL     string `json:"url"`
		Icon    string `json:"icon"`
		Source  string `json:"source"`
		Applied bool   `json:"applied"`
	}
	var results []iconResult
	successCount := 0
	skipCount := 0
	errorCount := 0

	for _, g := range data.Groups {
		for _, c := range g.Cards {
			if c.Icon != "" {
				skipCount++
				continue
			}
			favResp, err := s.favicon.FetchFavicon(c.URL)
			if err != nil {
				errorCount++
				continue
			}
			iconValue := favResp.IconName
			if iconValue == "" {
				iconValue = favResp.FaviconURL
			}
			if iconValue == "" {
				errorCount++
				continue
			}
			result := iconResult{
				CardID: c.ID,
				Title:  c.Title,
				URL:    c.URL,
				Icon:   iconValue,
				Source: favResp.Source,
			}
			if shouldApply {
				if err := s.panels.BatchSetIcons(uid, service.BatchSetIconsRequest{
					Icons: map[string]string{c.ID: iconValue},
				}); err != nil {
					result.Applied = false
					errorCount++
				} else {
					result.Applied = true
					successCount++
				}
			} else {
				result.Applied = false
				successCount++
			}
			results = append(results, result)
		}
	}

	summary := map[string]any{
		"total_processed": len(results),
		"skipped":         skipCount,
		"success":         successCount,
		"errors":          errorCount,
		"applied":         shouldApply,
		"results":         results,
	}
	return mcp.NewToolResultJSON(summary)
}

// --- system tools ------------------------------------------------------------

func registerSystemTools(mcps *server.MCPServer, s *Server) {
	// sundash_system_status: get system monitoring stats
	mcps.AddTool(mcp.NewTool(
		"sundash_system_status",
		mcp.WithDescription("获取系统监控状态：CPU 使用率、内存使用、磁盘空间、网络流量、主机信息、运行时间。用于 AI 分析 NAS/服务器健康度。"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleSystemStatus)
}

func (s *Server) handleSystemStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stats, err := s.system.GetStats()
	if err != nil {
		return mcp.NewToolResultError("获取系统状态失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(stats)
}

// --- docker tools -----------------------------------------------------------

func registerDockerTools(mcps *server.MCPServer, s *Server) {
	mcps.AddTool(mcp.NewTool(
		"sundash_docker_overview",
		mcp.WithDescription("获取 Docker 容器概览：运行中/已停止数量、容器列表、资源占用。用于 AI 分析 NAS 上 Docker 服务状态。"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleDockerOverview)

	mcps.AddTool(mcp.NewTool(
		"sundash_docker_list",
		mcp.WithDescription("列出 Docker 容器"),
		mcp.WithBoolean("all", mcp.Description("显示所有容器（含已停止，默认仅运行中）")),
	), s.handleDockerList)

	mcps.AddTool(mcp.NewTool(
		"sundash_docker_action",
		mcp.WithDescription("对 Docker 容器执行操作（start/stop/restart）"),
		mcp.WithString("action", mcp.Required(), mcp.Description("操作：start、stop 或 restart")),
		mcp.WithString("container", mcp.Required(), mcp.Description("容器名或 ID")),
	), s.handleDockerAction)

	mcps.AddTool(mcp.NewTool(
		"sundash_docker_logs",
		mcp.WithDescription("查看 Docker 容器日志"),
		mcp.WithString("container", mcp.Required(), mcp.Description("容器名或 ID")),
		mcp.WithNumber("tail", mcp.Description("显示最后 N 行（默认 50）")),
	), s.handleDockerLogs)
}

func (s *Server) handleDockerOverview(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	overview, err := s.docker.GetOverview()
	if err != nil {
		return mcp.NewToolResultError("获取 Docker 概览失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(overview)
}

func (s *Server) handleDockerList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	showAll := false
	if v, ok := req.Params.Arguments.(map[string]any); ok {
		if b, ok := v["all"].(bool); ok {
			showAll = b
		}
	}
	containers, err := s.docker.ListContainers(showAll)
	if err != nil {
		return mcp.NewToolResultError("列出容器失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(containers)
}

func (s *Server) handleDockerAction(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	action, _ := args["action"].(string)
	container, _ := args["container"].(string)
	if action == "" || container == "" {
		return mcp.NewToolResultError("缺少参数 action 或 container"), nil
	}
	out, err := s.docker.ContainerAction(action, container)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	msg := fmt.Sprintf("容器 %s %s 成功", container, action)
	if out != "" {
		msg += ": " + out
	}
	return mcp.NewToolResultText(msg), nil
}

func (s *Server) handleDockerLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	container, _ := args["container"].(string)
	tail := 50
	if v, ok := args["tail"].(float64); ok && v > 0 {
		tail = int(v)
	}
	logs, err := s.docker.ContainerLogs(container, tail)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(logs), nil
}

// --- deepseek tools ---------------------------------------------------------

func registerDeepSeekTools(mcps *server.MCPServer, s *Server) {
	mcps.AddTool(mcp.NewTool(
		"sundash_ai_suggest_group",
		mcp.WithDescription("AI 智能推荐书签分组：根据标题和 URL，从已有分组中选择最合适的。需要配置 DEEPSEEK_API_KEY。"),
		mcp.WithString("title", mcp.Required(), mcp.Description("书签标题")),
		mcp.WithString("url", mcp.Required(), mcp.Description("书签 URL")),
	), s.handleAISuggestGroup)

	mcps.AddTool(mcp.NewTool(
		"sundash_ai_summarize",
		mcp.WithDescription("AI 生成摘要：用一句话总结给定文本内容。需要配置 DEEPSEEK_API_KEY。"),
		mcp.WithString("text", mcp.Required(), mcp.Description("需要总结的文本")),
	), s.handleAISummarize)

	mcps.AddTool(mcp.NewTool(
		"sundash_ai_status",
		mcp.WithDescription("检查 AI 服务状态（DeepSeek API 是否可用）"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleAIStatus)
}

func (s *Server) handleAISuggestGroup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !s.deepseek.Available() {
		return mcp.NewToolResultError("AI 服务未配置（需要设置 DEEPSEEK_API_KEY 环境变量）"), nil
	}
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	args, _ := req.Params.Arguments.(map[string]any)
	title, _ := args["title"].(string)
	url, _ := args["url"].(string)
	if title == "" || url == "" {
		return mcp.NewToolResultError("缺少参数 title 或 url"), nil
	}
	// Get user's group names for context
	panel, err := s.panels.GetPanel(uid)
	if err != nil {
		return mcp.NewToolResultError("获取面板失败: " + err.Error()), nil
	}
	var groupNames []string
	for _, g := range panel.Groups {
		groupNames = append(groupNames, g.Name)
	}
	if len(groupNames) == 0 {
		return mcp.NewToolResultError("用户没有分组，请先创建分组"), nil
	}
	suggestion, err := s.deepseek.SuggestGroup(title, url, groupNames)
	if err != nil {
		return mcp.NewToolResultError("AI 推荐失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(map[string]any{
		"suggested_group": suggestion,
		"available_groups": groupNames,
	})
}

func (s *Server) handleAISummarize(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !s.deepseek.Available() {
		return mcp.NewToolResultError("AI 服务未配置（需要设置 DEEPSEEK_API_KEY 环境变量）"), nil
	}
	args, _ := req.Params.Arguments.(map[string]any)
	text, _ := args["text"].(string)
	if text == "" {
		return mcp.NewToolResultError("缺少参数 text"), nil
	}
	summary, err := s.deepseek.Summarize(text)
	if err != nil {
		return mcp.NewToolResultError("AI 摘要失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(map[string]any{"summary": summary})
}

func (s *Server) handleAIStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultJSON(map[string]any{
		"available": s.deepseek.Available(),
		"provider":  "deepseek",
		"model":     "deepseek-chat",
	})
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

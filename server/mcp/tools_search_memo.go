package mcp

import (
	"context"

	"sundash/models"
	"sundash/service"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// --- search & memo tools ----------------------------------------------------
// The tools the README advertises alongside the bookmark tools: card search
// over FTS5 and memo CRUD.

func registerSearchTools(mcps *server.MCPServer, s *Server) {
	// sundash_search: full-text search across the user's bookmark cards.
	mcps.AddTool(mcp.NewTool(
		"sundash_search",
		mcp.WithDescription("全文搜索当前用户的书签卡片（FTS5，按相关度排序）。query 为搜索词，limit 可选（默认 10，最大 50）。"),
		mcp.WithString("query", mcp.Description("搜索关键词")),
		mcp.WithNumber("limit", mcp.Description("返回条数上限，默认 10")),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleSearch)
}

func registerMemoTools(mcps *server.MCPServer, s *Server) {
	// sundash_list_memo: list the user's memos.
	mcps.AddTool(mcp.NewTool(
		"sundash_list_memo",
		mcp.WithDescription("列出当前用户的全部便签（含已归档，is_archived 标记）。无参数。"),
		mcp.WithReadOnlyHintAnnotation(true),
	), s.handleListMemos)

	// sundash_add_memo: create a memo.
	mcps.AddTool(mcp.NewTool(
		"sundash_add_memo",
		mcp.WithDescription("新建一条便签。content 为便签文本（必填）。"),
		mcp.WithString("content", mcp.Required(), mcp.Description("便签内容")),
	), s.handleAddMemo)

	// sundash_archive_memo: archive or unarchive a memo.
	mcps.AddTool(mcp.NewTool(
		"sundash_archive_memo",
		mcp.WithDescription("归档或取消归档一条便签。id 必填；archived 可选，默认 true。"),
		mcp.WithString("id", mcp.Required(), mcp.Description("便签 ID")),
		mcp.WithBoolean("archived", mcp.Description("true 归档 / false 取消归档，默认 true")),
	), s.handleArchiveMemo)

	// sundash_delete_memo: delete a memo permanently.
	mcps.AddTool(mcp.NewTool(
		"sundash_delete_memo",
		mcp.WithDescription("永久删除一条便签。id 必填。该操作不可恢复。"),
		mcp.WithString("id", mcp.Required(), mcp.Description("便签 ID")),
	), s.handleDeleteMemo)
}

func (s *Server) handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	query := strArg(request.Params.Arguments, "query")
	if query == "" {
		return mcp.NewToolResultError("缺少参数 query"), nil
	}
	limit := 10
	if l := intArg(request.Params.Arguments, "limit"); l != nil && *l > 0 && *l <= 50 {
		limit = *l
	}
	results, err := s.search.SearchCards(ctx, uid, query, limit)
	if err != nil {
		return mcp.NewToolResultError("搜索失败: " + err.Error()), nil
	}
	if results == nil {
		results = []service.SearchResult{}
	}
	return mcp.NewToolResultJSON(results)
}

func (s *Server) handleListMemos(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	memos, err := s.memo.GetMemos(ctx, uid)
	if err != nil {
		return mcp.NewToolResultError("读取便签失败: " + err.Error()), nil
	}
	if memos == nil {
		memos = []*models.Memo{}
	}
	return mcp.NewToolResultJSON(memos)
}

func (s *Server) handleAddMemo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	content := strArg(request.Params.Arguments, "content")
	if content == "" {
		return mcp.NewToolResultError("缺少参数 content"), nil
	}
	memo, err := s.memo.CreateMemo(ctx, uid, content)
	if err != nil {
		return mcp.NewToolResultError("创建便签失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(memo)
}

func (s *Server) handleArchiveMemo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id := strArg(request.Params.Arguments, "id")
	if id == "" {
		return mcp.NewToolResultError("缺少参数 id"), nil
	}
	archived := true
	if v, ok := request.Params.Arguments.(map[string]any); ok {
		if b, ok := v["archived"].(bool); ok {
			archived = b
		}
	}
	if err := s.memo.ArchiveMemo(ctx, uid, id, archived); err != nil {
		return mcp.NewToolResultError("更新便签失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(map[string]any{"id": id, "is_archived": archived})
}

func (s *Server) handleDeleteMemo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	id := strArg(request.Params.Arguments, "id")
	if id == "" {
		return mcp.NewToolResultError("缺少参数 id"), nil
	}
	if err := s.memo.DeleteMemo(ctx, uid, id); err != nil {
		return mcp.NewToolResultError("删除便签失败: " + err.Error()), nil
	}
	return mcp.NewToolResultJSON(map[string]any{"id": id, "deleted": true})
}

package models

// SyncNode mirrors the canonical bookmark node from the bookmark-sync server.
// Field names match the wire protocol exactly (camelCase JSON).
type SyncNode struct {
	SyncID       string `json:"syncId"`
	Type         string `json:"type"` // "bookmark" | "folder"
	Title        string `json:"title"`
	URL          string `json:"url,omitempty"`
	ParentSyncID string `json:"parentSyncId,omitempty"`
	Index        int    `json:"index"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	DeletedAt    string `json:"deletedAt,omitempty"` // tombstone
	Rev          int    `json:"rev,omitempty"`
}

// SyncNodeResponse is the wire shape returned by the bookmark-sync server
// (POST /api/sync and GET /api/state): full canonical state + monotonic rev.
type SyncNodeResponse struct {
	Rev   int         `json:"rev"`
	Nodes []*SyncNode `json:"nodes"`
}

// ChangeIntent is a user mutation coming from the sundash UI.
// Op is one of "create" | "update" | "delete".
type ChangeIntent struct {
	Op           string `json:"op" binding:"required"`
	SyncID       string `json:"syncId"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	URL          string `json:"url,omitempty"`
	ParentSyncID string `json:"parentSyncId,omitempty"`
	Index        int    `json:"index"`
}

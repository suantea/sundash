package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// BackupHandler serves a consistent snapshot of the SQLite database.
//
// DEPLOY.md tells operators to back up data/sundash.db regularly, but until
// now the only export was the dashboard's client-side JSON (bookmarks only —
// no users, settings, memos, feeds). This endpoint produces a complete,
// transactionally consistent snapshot via `VACUUM INTO`, which is safe to run
// while the server is live: SQLite writes the snapshot to a brand-new file
// and never touches the original database.
type BackupHandler struct {
	db *sql.DB
}

func NewBackupHandler(db *sql.DB) *BackupHandler {
	return &BackupHandler{db: db}
}

// Download handles GET /api/admin/backup (admin only, enforced by routing).
func (h *BackupHandler) Download(c *gin.Context) {
	tmp, err := os.CreateTemp("", "sundash-backup-*.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建备份临时文件失败: " + err.Error()})
		return
	}
	tmp.Close()
	// VACUUM INTO requires the target file to NOT exist.
	if err := os.Remove(tmp.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(tmp.Name())

	if _, err := h.db.ExecContext(c.Request.Context(), "VACUUM INTO ?", tmp.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "备份失败: " + err.Error()})
		return
	}

	name := "sundash-backup-" + time.Now().Format("20060102-150405") + ".db"
	c.Header("Content-Disposition", `attachment; filename="`+name+`"`)
	// FileAttachment streams the snapshot; the deferred Remove runs afterwards.
	c.FileAttachment(tmp.Name(), filepath.Base(name))
}

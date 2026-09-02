package handlers

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	db      *sql.DB
	dataDir string
}

func NewBackupHandler(db *sql.DB, dataDir string) *BackupHandler {
	return &BackupHandler{db: db, dataDir: dataDir}
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

// Restore handles POST /api/admin/restore (admin only).
// Accepts a multipart file upload containing a SQLite database.
// It replaces the current database file and requires a server restart to take effect.
func (h *BackupHandler) Restore(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传数据库文件"})
		return
	}
	defer file.Close()

	// Validate file size (max 100MB)
	if header.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，最大支持 100MB"})
		return
	}

	// Write uploaded file to a temp location
	tmp, err := os.CreateTemp("", "sundash-restore-*.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时文件失败"})
		return
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入临时文件失败"})
		return
	}
	tmp.Close()

	// Validate it's a valid SQLite database (magic header)
	f, err := os.Open(tmp.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传文件"})
		return
	}
	magic := make([]byte, 16)
	if _, err := f.Read(magic); err != nil || string(magic[:15]) != "SQLite format 3" {
		f.Close()
		c.JSON(http.StatusBadRequest, gin.H{"error": "不是有效的 SQLite 数据库文件"})
		return
	}
	f.Close()

	// Backup current database before replacing
	dbPath := filepath.Join(h.dataDir, "sundash.db")
	backupPath := dbPath + ".bak-" + time.Now().Format("20060102-150405")
	if err := copyFile(dbPath, backupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "备份当前数据库失败: " + err.Error()})
		return
	}

	// Replace current database
	if err := copyFile(tmp.Name(), dbPath); err != nil {
		// Restore from backup on failure
		copyFile(backupPath, dbPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复数据库失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "恢复成功，请重启服务以加载新数据",
		"backup":   filepath.Base(backupPath),
		"file":     header.Filename,
		"size":     header.Size,
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

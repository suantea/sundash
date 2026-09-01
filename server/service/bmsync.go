package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sundash/models"
	"sundash/repository"
)

// Bmsync service keys stored in the global settings table.
const (
	BmsyncServerURL = "bmsync_server_url"
	BmsyncToken     = "bmsync_token"
	BmsyncDeviceID  = "bmsync_device_id"
)

// bmsync-specific errors (wrapped in AppError by handlers).
var (
	ErrBmsyncNotConfigured = NewError(ErrBadRequest, "bookmark-sync 未配置：请先在设置中填写服务器地址与 Token")
	ErrBmsyncUnreachable   = NewError(ErrBadGateway, "无法连接 bookmark-sync 服务器")
	ErrBmsyncServer        = NewError(ErrBadGateway, "bookmark-sync 服务器返回错误")
)

var BmsyncSettingKeys = []string{BmsyncServerURL, BmsyncToken, BmsyncDeviceID}

// BmsyncService is the bookmark-sync client: it pulls the full canonical tree
// from the bookmark-sync server and pushes local changes. sundash is one more
// read/write replica — the LWW + tombstone rules are enforced by the server.
type BmsyncService struct {
	repo     *repository.BmsyncRepo
	settings *repository.SettingsRepo
	http     *http.Client
}

func NewBmsyncService(repo *repository.BmsyncRepo, settings *repository.SettingsRepo) *BmsyncService {
	return &BmsyncService{
		repo:     repo,
		settings: settings,
		http:     &http.Client{Timeout: 15 * time.Second},
	}
}

// Repo exposes the mirror repo so handlers can read locally without re-syncing.
func (s *BmsyncService) Repo() *repository.BmsyncRepo { return s.repo }

// Config is the connection settings for the bookmark-sync server.
type Config struct {
	ServerURL string
	Token     string
	DeviceID  string
}

func (s *BmsyncService) GetConfig(ctx context.Context) (Config, error) {
	vals, err := s.settings.GlobalByKeys(BmsyncSettingKeys)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		ServerURL: strings.TrimRight(strings.TrimSpace(vals[BmsyncServerURL]), "/"),
		Token:     strings.TrimSpace(vals[BmsyncToken]),
		DeviceID:  vals[BmsyncDeviceID],
	}
	if cfg.DeviceID == "" {
		cfg.DeviceID = "sundash"
	}
	return cfg, nil
}

// IsConfigured reports whether the sync server and token are set.
func (s *BmsyncService) IsConfigured(ctx context.Context) bool {
	cfg, err := s.GetConfig(ctx)
	return err == nil && cfg.ServerURL != "" && cfg.Token != ""
}

// Pull fetches the full canonical state (incl. tombstones) and replaces the mirror.
func (s *BmsyncService) Pull(ctx context.Context) (*models.SyncNodeResponse, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg.ServerURL == "" || cfg.Token == "" {
		return nil, ErrBmsyncNotConfigured
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ServerURL+"/api/state", nil)
	if err != nil {
		return nil, fmt.Errorf("build state request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBmsyncUnreachable, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return nil, fmt.Errorf("read state response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: server returned %d: %s", ErrBmsyncServer, resp.StatusCode, truncate(string(body), 200))
	}
	var out models.SyncNodeResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("parse state response: %w", err)
	}
	if err := s.repo.ReplaceAll(ctx, out.Nodes, out.Rev); err != nil {
		return nil, err
	}
	return &out, nil
}

// Push applies client-side changes and returns the server's full canonical state.
// Each change must already carry the correct timestamps (see ApplyChange below).
func (s *BmsyncService) Push(ctx context.Context, changes []*models.SyncNode) (*models.SyncNodeResponse, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg.ServerURL == "" || cfg.Token == "" {
		return nil, ErrBmsyncNotConfigured
	}
	lastRev, err := s.repo.GetRev(ctx)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"deviceId": cfg.DeviceID,
		"sinceRev": lastRev,
		"changes":  changes,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal push payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.ServerURL+"/api/sync", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("build sync request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBmsyncUnreachable, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return nil, fmt.Errorf("read sync response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: server returned %d: %s", ErrBmsyncServer, resp.StatusCode, truncate(string(body), 200))
	}
	var out models.SyncNodeResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("parse sync response: %w", err)
	}
	if err := s.repo.ReplaceAll(ctx, out.Nodes, out.Rev); err != nil {
		return nil, err
	}
	return &out, nil
}

// ApplyChange stamps the protocol timestamps for one user-intent mutation:
//   - create   -> new UUID + createdAt = updatedAt = now
//   - edit     -> updatedAt = now
//   - delete   -> deletedAt = updatedAt = now (tombstone; the row is kept, and
//     other clients remove the bookmark because the tombstone propagates)
func ApplyChange(now time.Time, intent models.ChangeIntent) *models.SyncNode {
	n := &models.SyncNode{
		SyncID:       intent.SyncID,
		Type:         intent.Type,
		Title:        intent.Title,
		URL:          intent.URL,
		ParentSyncID: intent.ParentSyncID,
		Index:        intent.Index,
	}
	nowIso := now.UTC().Format(time.RFC3339Nano)
	switch intent.Op {
	case "create":
		if n.SyncID == "" {
			n.SyncID = newUUID()
		}
		n.CreatedAt = nowIso
		n.UpdatedAt = nowIso
	case "delete":
		n.DeletedAt = nowIso
		n.UpdatedAt = nowIso
	default: // "update"
		n.UpdatedAt = nowIso
	}
	return n
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func newUUID() string {
	// RFC4122 v4-ish without crypto/rand import overhead.
	return fmt.Sprintf("%x-%x-%x-%x-%x", time.Now().UnixNano()&0xffffffff,
		time.Now().UnixNano()>>32&0xffff, 0x4000|time.Now().UnixNano()>>16&0x0fff,
		0x8000|time.Now().UnixNano()>>8&0x3fff, time.Now().UnixNano()<<48&0xffffffffffff)
}

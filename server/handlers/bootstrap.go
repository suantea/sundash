package handlers

import (
	"sync"

	"sundash/models"
	"sundash/service"

	"github.com/gin-gonic/gin"
)

// BootstrapHandler aggregates everything the home page needs into a single
// request, cutting the initial page load from several round-trips to one.
type BootstrapHandler struct {
	users    *service.UserService
	panels   *service.PanelService
	settings *service.SettingsService
}

func NewBootstrapHandler(users *service.UserService, panels *service.PanelService, settings *service.SettingsService) *BootstrapHandler {
	return &BootstrapHandler{users: users, panels: panels, settings: settings}
}

type bootstrapResult struct {
	Settings map[string]string    `json:"settings"`
	Profile  *models.User         `json:"profile"`
	Panels   *models.PanelData    `json:"panels"`
}

// GetBootstrap fetches settings/profile/panel concurrently in one response.
func (h *BootstrapHandler) GetBootstrap(c *gin.Context) {
	userID := c.GetString("user_id")

	var (
		settings map[string]string
		profile  *models.User
		panel    *models.PanelData
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	record := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		s, err := h.settings.GetForUser(userID)
		mu.Lock()
		settings = s
		mu.Unlock()
		record(err)
	}()
	go func() {
		defer wg.Done()
		p, err := h.users.GetProfile(userID)
		mu.Lock()
		profile = p
		mu.Unlock()
		record(err)
	}()
	go func() {
		defer wg.Done()
		p, err := h.panels.GetPanel(userID)
		mu.Lock()
		panel = p
		mu.Unlock()
		record(err)
	}()
	wg.Wait()

	if firstErr != nil {
		respondError(c, firstErr)
		return
	}
	c.JSON(200, bootstrapResult{Settings: settings, Profile: profile, Panels: panel})
}

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"sundash/models"
	"sundash/repository"
)

// RSSFeed represents an RSS feed subscription.
type RSSFeed struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	Description    string    `json:"description,omitempty"`
	ImageURL       string    `json:"image_url,omitempty"`
	LastFetched    time.Time `json:"last_fetched,omitempty"`
	UpdateInterval int       `json:"update_interval"` // How often to fetch (in minutes)
	Items          []RSSItem `json:"items,omitempty"`
}

// RSSItem represents a single item in an RSS feed.
type RSSItem struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description,omitempty"`
	PubDate     time.Time `json:"pub_date,omitempty"`
	Author      string    `json:"author,omitempty"`
	Guid        string    `json:"guid,omitempty"`
}

// RSSService provides RSS feed fetching and management.
type RSSService struct {
	feedRepo *repository.RSSFeedRepo
	itemRepo *repository.RSSItemRepo
	mu       sync.RWMutex
	feeds    map[string]*RSSFeed // cache of feeds by ID
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewRSSService creates a new RSSService.
func NewRSSService(feedRepo *repository.RSSFeedRepo, itemRepo *repository.RSSItemRepo) *RSSService {
	s := &RSSService{
		feedRepo: feedRepo,
		itemRepo: itemRepo,
		feeds:    make(map[string]*RSSFeed),
		stopCh:   make(chan struct{}),
	}
	// Start background fetcher
	s.wg.Add(1)
	go s.backgroundFetcher()
	return s
}

// Stop stops the background fetcher.
func (s *RSSService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

// AddFeed adds a new RSS feed subscription.
func (s *RSSService) AddFeed(ctx context.Context, userID, url string) (*RSSFeed, error) {
	// Parse the feed to get title and description
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Create feed record
	rssFeed := &models.RSSFeed{
		ID:             generateID(),
		UserID:         userID,
		Title:          feed.Title,
		URL:            url,
		Description:    feed.Description,
		ImageURL:       extractImageURL(feed),
		UpdateInterval: 60, // default 60 minutes
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	if err := s.feedRepo.Create(ctx, rssFeed); err != nil {
		return nil, fmt.Errorf("failed to save RSS feed: %w", err)
	}

	// Fetch initial items
	if err := s.fetchAndStoreItems(ctx, rssFeed); err != nil {
		// Don't fail the addition if initial fetch fails
		// Items will be fetched in background
	}

	// Add to cache
	s.mu.Lock()
	s.feeds[rssFeed.ID] = &RSSFeed{
		ID:             rssFeed.ID,
		UserID:         rssFeed.UserID,
		Title:          rssFeed.Title,
		URL:            rssFeed.URL,
		Description:    rssFeed.Description,
		ImageURL:       rssFeed.ImageURL,
		LastFetched:    rssFeed.CreatedAt,
		UpdateInterval: rssFeed.UpdateInterval,
	}
	s.mu.Unlock()

	return &RSSFeed{
		ID:             rssFeed.ID,
		UserID:         rssFeed.UserID,
		Title:          rssFeed.Title,
		URL:            rssFeed.URL,
		Description:    rssFeed.Description,
		ImageURL:       rssFeed.ImageURL,
		LastFetched:    rssFeed.CreatedAt,
		UpdateInterval: rssFeed.UpdateInterval,
		Items:          []RSSItem{}, // Will be populated by background fetcher
	}, nil
}

// GetFeeds returns all RSS feeds for a user.
func (s *RSSService) GetFeeds(ctx context.Context, userID string) ([]*RSSFeed, error) {
	dbFeeds, err := s.feedRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get RSS feeds: %w", err)
	}

	// Convert to service format and include recent items
	feeds := make([]*RSSFeed, 0)
	for _, dbFeed := range dbFeeds {
		// Get recent items for this feed
		items, err := s.itemRepo.ListByFeed(ctx, dbFeed.ID, 10) // latest 10 items
		if err != nil {
			items = []*models.RSSItem{} // continue without items if fetch fails
		}

		// Convert items
		var rssItems []RSSItem
		for _, item := range items {
			rssItems = append(rssItems, RSSItem{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     item.PubDate,
				Author:      item.Author,
				Guid:        item.Guid,
			})
		}

		feeds = append(feeds, &RSSFeed{
			ID:             dbFeed.ID,
			UserID:         dbFeed.UserID,
			Title:          dbFeed.Title,
			URL:            dbFeed.URL,
			Description:    dbFeed.Description,
			ImageURL:       dbFeed.ImageURL,
			LastFetched:    dbFeed.LastFetched,
			UpdateInterval: dbFeed.UpdateInterval,
			Items:          rssItems,
		})
	}
	return feeds, nil
}

// GetFeedItems returns recent items of one feed, after checking the feed
// belongs to the requesting user.
func (s *RSSService) GetFeedItems(ctx context.Context, userID, feedID string, limit int) ([]*models.RSSItem, error) {
	if _, err := s.feedRepo.GetByIDAndUser(ctx, feedID, userID); err != nil {
		return nil, fmt.Errorf("feed not found or unauthorized: %w", err)
	}
	return s.itemRepo.ListByFeed(ctx, feedID, limit)
}

// RemoveFeed deletes an RSS feed subscription.
func (s *RSSService) RemoveFeed(ctx context.Context, userID, feedID string) error {
	// Delete items first (cascade)
	if err := s.itemRepo.DeleteByFeed(ctx, feedID); err != nil {
		return fmt.Errorf("failed to delete RSS feed items: %w", err)
	}
	// Delete feed
	if err := s.feedRepo.Delete(ctx, userID, feedID); err != nil {
		return fmt.Errorf("failed to delete RSS feed: %w", err)
	}
	// Remove from cache
	s.mu.Lock()
	delete(s.feeds, feedID)
	s.mu.Unlock()
	return nil
}

// UpdateFeed updates an RSS feed subscription (e.g., change URL).
func (s *RSSService) UpdateFeed(ctx context.Context, userID, feedID, url string) (*RSSFeed, error) {
	// Check ownership
	existing, err := s.feedRepo.GetByIDAndUser(ctx, feedID, userID)
	if err != nil {
		return nil, fmt.Errorf("feed not found or unauthorized: %w", err)
	}

	// Parse new feed to validate
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Update feed record
	existing.Title = feed.Title
	existing.Description = feed.Description
	existing.ImageURL = extractImageURL(feed)
	existing.URL = url
	existing.UpdatedAt = time.Now()

	if err := s.feedRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update RSS feed: %w", err)
	}

	// Clear items for this feed (will be refetched)
	if err := s.itemRepo.DeleteByFeed(ctx, feedID); err != nil {
		return nil, fmt.Errorf("failed to clear RSS feed items: %w", err)
	}

	// Fetch initial items
	if err := s.fetchAndStoreItems(ctx, existing); err != nil {
		// Continue even if initial fetch fails
	}

	// Update cache
	s.mu.Lock()
	if f, ok := s.feeds[feedID]; ok {
		f.Title = existing.Title
		f.Description = existing.Description
		f.ImageURL = existing.ImageURL
		f.URL = existing.URL
		f.LastFetched = existing.UpdatedAt
	}
	s.mu.Unlock()

	return &RSSFeed{
		ID:             existing.ID,
		UserID:         existing.UserID,
		Title:          existing.Title,
		URL:            existing.URL,
		Description:    existing.Description,
		ImageURL:       existing.ImageURL,
		LastFetched:    existing.LastFetched,
		UpdateInterval: existing.UpdateInterval,
		Items:          []RSSItem{},
	}, nil
}

// backgroundFetcher runs in a goroutine to periodically fetch feeds.
func (s *RSSService) backgroundFetcher() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Minute) // check every minute
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.fetchAllFeeds()
		}
	}
}

// fetchAllFeeds fetches all feeds that are due for update.
func (s *RSSService) fetchAllFeeds() {
	// We need to get all feeds from the database to check update intervals
	// For simplicity, we'll fetch all feeds for all users (could be optimized)
	// In a real app, we might want to track which feeds need updating
	ctx := context.Background()

	// Get all feeds (this could be expensive; consider adding a "next_fetch_at" column)
	// For now, we'll implement a simple approach: fetch all feeds every time
	// but respect the update interval by checking LastFetched
	feeds, err := s.feedRepo.ListAll(ctx) // Need to add this method
	if err != nil {
		// Log error but continue
		return
	}

	now := time.Now()
	for _, dbFeed := range feeds {
		// Check if it's time to fetch
		if dbFeed.LastFetched.IsZero() ||
			now.Sub(dbFeed.LastFetched) >= time.Duration(dbFeed.UpdateInterval)*time.Minute {
			// Fetch in background to avoid blocking the ticker
			go func(f *models.RSSFeed) {
				_ = s.fetchAndStoreItems(ctx, f)
			}(dbFeed)
		}
	}
}

// fetchAndStoreItems fetches a feed and stores its items.
func (s *RSSService) fetchAndStoreItems(ctx context.Context, feed *models.RSSFeed) error {
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURLWithContext(feed.URL, context.Background())
	if err != nil {
		return fmt.Errorf("failed to parse RSS feed %s: %w", feed.URL, err)
	}

	// Update feed metadata
	feed.Title = parsedFeed.Title
	feed.Description = parsedFeed.Description
	feed.ImageURL = extractImageURL(parsedFeed)
	feed.LastFetched = time.Now()
	feed.UpdatedAt = time.Now()

	// Save updated feed info
	if err := s.feedRepo.Update(ctx, feed); err != nil {
		return fmt.Errorf("failed to update feed metadata: %w", err)
	}

	// Delete old items (we keep only recent items, e.g., last 50)
	if err := s.itemRepo.DeleteByFeed(ctx, feed.ID); err != nil {
		return fmt.Errorf("failed to delete old items: %w", err)
	}

	// Store new items (limit to 50 most recent)
	maxItems := 50
	if len(parsedFeed.Items) < maxItems {
		maxItems = len(parsedFeed.Items)
	}
	for i := 0; i < maxItems; i++ {
		item := parsedFeed.Items[i]
		// Not every entry carries a publish date; gofeed leaves
		// PublishedParsed nil then, and dereferencing it would panic.
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubDate = *item.UpdatedParsed
		}
		rssItem := &models.RSSItem{
			ID:          generateID(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     pubDate,
			Author:      extractAuthor(item),
			Guid:        item.GUID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.itemRepo.Create(ctx, rssItem); err != nil {
			// Continue with other items if one fails
			continue
		}
	}

	// Update cache
	s.mu.Lock()
	if cached, ok := s.feeds[feed.ID]; ok {
		cached.Title = feed.Title
		cached.Description = feed.Description
		cached.ImageURL = feed.ImageURL
		cached.LastFetched = feed.LastFetched
	}
	s.mu.Unlock()

	return nil
}

// extractImageURL tries to extract an image URL from the feed.
func extractImageURL(feed *gofeed.Feed) string {
	if feed.Image != nil && feed.Image.URL != "" {
		return feed.Image.URL
	}
	// TODO: also look for icon or other images
	return ""
}

// extractAuthor extracts author from an item.
func extractAuthor(item *gofeed.Item) string {
	if item.Author != nil {
		if item.Author.Name != "" {
			return item.Author.Name
		}
		if item.Author.Email != "" {
			return item.Author.Email
		}
	}
	// Feeds often carry the author only in the item's authors list (Atom).
	for _, a := range item.Authors {
		if a != nil && a.Name != "" {
			return a.Name
		}
	}
	return ""
}

// generateID is defined in memo.go

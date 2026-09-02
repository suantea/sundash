package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// FaviconResponse represents the result of favicon detection.
type FaviconResponse struct {
	IconName   string `json:"icon_name"`   // Iconify icon name like "mdi:github"
	FaviconURL string `json:"favicon_url"` // Direct favicon URL if no iconify match
	Source     string `json:"source"`      // "mapping", "iconify", "favicon", "none"
}

// siteIconMap maps known domains to Iconify icons.
var siteIconMap = map[string]string{
	"github.com": "mdi:github", "gitlab.com": "mdi:gitlab", "github.io": "mdi:github",
	"google.com": "mdi:google", "google.cn": "mdi:google", "youtube.com": "mdi:youtube",
	"youtu.be": "mdi:youtube", "bilibili.com": "mdi:play-box-outline", "twitter.com": "mdi:twitter",
	"x.com": "mdi:twitter", "facebook.com": "mdi:facebook", "instagram.com": "mdi:instagram",
	"reddit.com": "mdi:reddit", "stackoverflow.com": "mdi:stack-overflow", "stackexchange.com": "mdi:stack-overflow",
	"npmjs.com": "mdi:npm", "docker.com": "mdi:docker", "docker.io": "mdi:docker",
	"hub.docker.com": "mdi:docker", "vercel.com": "mdi:triangle-outline", "netlify.com": "mdi:triangles-diagonal",
	"apple.com": "mdi:apple", "microsoft.com": "mdi:microsoft", "live.com": "mdi:microsoft",
	"office.com": "mdi:microsoft-office", "amazon.com": "mdi:amazon", "aws.amazon.com": "mdi:aws",
	"wikipedia.org": "mdi:wikipedia", "notion.so": "mdi:notebook-outline", "figma.com": "mdi:figma",
	"spotify.com": "mdi:spotify", "netflix.com": "mdi:netflix", "taobao.com": "mdi:shopping",
	"tmall.com": "mdi:shopping", "jd.com": "mdi:shopping-outline", "baidu.com": "mdi:alpha-b-circle-outline",
	"zhihu.com": "mdi:alpha-z-circle-outline", "weibo.com": "mdi:alpha-w-circle-outline",
	"douyin.com": "mdi:music-note", "v2ex.com": "mdi:alpha-v-circle-outline", "juejin.cn": "mdi:diamond-stone",
	"gitee.com": "mdi:code-tags", "aliyun.com": "mdi:cloud-outline", "cloudflare.com": "mdi:cloud-outline",
	"telegram.org": "mdi:telegram", "t.me": "mdi:telegram", "whatsapp.com": "mdi:whatsapp",
	"slack.com": "mdi:slack", "discord.com": "mdi:discord", "twitch.tv": "mdi:twitch",
	"medium.com": "mdi:medium", "dev.to": "mdi:dev-to", "linkedin.com": "mdi:linkedin",
	"pinterest.com": "mdi:pinterest", "tiktok.com": "mdi:music-note", "steam": "mdi:steam",
	"store.steampowered.com": "mdi:steam", "play.google.com": "mdi:google-play", "apps.apple.com": "mdi:apple",
	"chromewebstore.google.com": "mdi:chrome", "maps.google.com": "mdi:google-maps", "docs.google.com": "mdi:google-docs",
	"drive.google.com": "mdi:google-drive", "gmail.com": "mdi:gmail", "mail.google.com": "mdi:gmail",
	"calendar.google.com": "mdi:google-calendar", "photos.google.com": "mdi:google-photos",
	"translate.google.com": "mdi:google-translate", "analytics.google.com": "mdi:google-analytics",
	"console.cloud.google.com": "mdi:google-cloud", "console.aws.amazon.com": "mdi:aws",
	"gitbook.io": "mdi:book-open-page-variant", "readthedocs.io": "mdi:book-open-variant",
	"pkg.go.dev": "mdi:language-go", "crates.io": "mdi:language-rust", "pypi.org": "mdi:language-python",
	"rubygems.org": "mdi:language-ruby", "nuget.org": "mdi:package-variant", "homebrew.sh": "mdi:beer-outline",
	"archlinux.org": "mdi:penguin", "ubuntu.com": "mdi:ubuntu", "debian.org": "mdi:debian",
	"centos.org": "mdi:linux", "fedoraproject.org": "mdi:fedora", "raspberrypi.org": "mdi:raspberry-pi",
	"arduino.cc": "mdi:arduino", "raspberrypi.com": "mdi:raspberry-pi",
}

type FaviconService struct {
	client *http.Client
}

func NewFaviconService() *FaviconService {
	return &FaviconService{client: &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}}
}

// FetchFavicon detects icon info for a URL, refusing internal/private targets (SSRF guard).
func (s *FaviconService) FetchFavicon(rawURL string) (FaviconResponse, error) {
	resp := FaviconResponse{Source: "none"}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	domain := extractDomain(rawURL)
	if domain == "" {
		return resp, NewError(ErrBadRequest, "invalid URL")
	}

	// SSRF guard: reject targets resolving to private/internal addresses.
	if err := s.validatePublicTarget(rawURL); err != nil {
		return resp, NewError(ErrBadRequest, "URL is not allowed")
	}

	// 1. Known site mappings (including parent domain fallback).
	if iconName, ok := lookupIcon(domain); ok {
		resp.IconName = iconName
		resp.Source = "mapping"
		return resp, nil
	}

	// 2. Fetch page HTML and look for icon references.
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		resp.FaviconURL = googleFavicon(domain)
		resp.Source = "favicon"
		return resp, nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	res, err := s.client.Do(req)
	if err != nil {
		resp.FaviconURL = googleFavicon(domain)
		resp.Source = "favicon"
		return resp, nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 50*1024))
	if err != nil {
		resp.FaviconURL = googleFavicon(domain)
		resp.Source = "favicon"
		return resp, nil
	}
	html := string(body)

	// Iconify references: icon-sets.iconify.design/PREFIX/NAME or api.iconify.design/PREFIX/NAME.svg
	iconifyRegex := regexp.MustCompile(`(?:icon-sets\.iconify\.design|api\.iconify\.design)[/"]+([\w-]+)/([\w-]+)`)
	if matches := iconifyRegex.FindStringSubmatch(html); len(matches) >= 3 {
		resp.IconName = matches[1] + ":" + matches[2]
		resp.Source = "iconify"
		return resp, nil
	}

	iconLinkRegex := regexp.MustCompile(`<link[^>]+rel=["'](?:shortcut )?icon["'][^>]+href=["']([^"']+)["']`)
	iconLinkRegex2 := regexp.MustCompile(`<link[^>]+href=["']([^"']+)["'][^>]+rel=["'](?:shortcut )?icon["']`)
	if faviconHref := findFaviconHref(html, iconLinkRegex, iconLinkRegex2); faviconHref != "" {
		resp.FaviconURL = absolutize(faviconHref, domain)
		resp.Source = "favicon"
		return resp, nil
	}

	resp.FaviconURL = googleFavicon(domain)
	resp.Source = "favicon"
	return resp, nil
}

// validatePublicTarget blocks requests to loopback, private, link-local and
// other reserved addresses (also covers cloud metadata 169.254.169.254).
func (s *FaviconService) validatePublicTarget(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	host := strings.Trim(u.Hostname(), "[]")

	if ip := net.ParseIP(host); ip != nil {
		if isBlockedIP(ip) {
			return fmt.Errorf("target address is blocked")
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return fmt.Errorf("target address is blocked")
		}
	}
	return nil
}

func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified()
}

func lookupIcon(domain string) (string, bool) {
	if name, ok := siteIconMap[domain]; ok {
		return name, true
	}
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts)-1; i++ {
		if name, ok := siteIconMap[strings.Join(parts[i:], ".")]; ok {
			return name, true
		}
	}
	return "", false
}

func extractDomain(rawURL string) string {
	domain := strings.TrimPrefix(rawURL, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	if idx := strings.IndexAny(domain, "/:?#"); idx != -1 {
		domain = domain[:idx]
	}
	return domain
}

func googleFavicon(domain string) string {
	return fmt.Sprintf("https://favicon.im/%s", domain)
}

func findFaviconHref(html string, re1, re2 *regexp.Regexp) string {
	if m := re1.FindStringSubmatch(html); len(m) >= 2 {
		return m[1]
	}
	if m := re2.FindStringSubmatch(html); len(m) >= 2 {
		return m[1]
	}
	return ""
}

func absolutize(href, domain string) string {
	switch {
	case strings.HasPrefix(href, "//"):
		return "https:" + href
	case strings.HasPrefix(href, "/"):
		return "https://" + domain + href
	case !strings.HasPrefix(href, "http"):
		return "https://" + domain + "/" + href
	default:
		return href
	}
}

// BatchFetchFavicons 批量获取多个 URL 的 favicon
func (s *FaviconService) BatchFetchFavicons(urls []string) []FaviconResponse {
	results := make([]FaviconResponse, len(urls))
	for i, u := range urls {
		results[i], _ = s.FetchFavicon(u)
	}
	return results
}

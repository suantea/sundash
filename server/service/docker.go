package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// DockerContainer represents a Docker container.
type DockerContainer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Ports   string `json:"ports"`
	Created string `json:"created"`
}

// DockerService provides Docker container management.
type DockerService struct{}

func NewDockerService() *DockerService {
	return &DockerService{}
}

// Available checks if Docker is installed and accessible.
func (d *DockerService) Available() bool {
	err := exec.Command("docker", "info").Run()
	return err == nil
}

// ListContainers returns running (or all) containers.
func (d *DockerService) ListContainers(showAll bool) ([]DockerContainer, error) {
	args := []string{"ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.State}}\t{{.Ports}}\t{{.CreatedAt}}"}
	if showAll {
		args = append(args, "-a")
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker ps: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var containers []DockerContainer
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 7)
		if len(parts) >= 5 {
			c := DockerContainer{
				ID:     parts[0],
				Name:   parts[1],
				Image:  parts[2],
				Status: parts[3],
				State:  parts[4],
			}
			if len(parts) >= 6 {
				c.Ports = parts[5]
			}
			if len(parts) >= 7 {
				c.Created = parts[6]
			}
			containers = append(containers, c)
		}
	}
	return containers, nil
}

// ContainerAction performs start/stop/restart on a container.
func (d *DockerService) ContainerAction(action, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("容器名不能为空")
	}
	out, err := exec.Command("docker", action, name).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker %s %s: %w (%s)", action, name, err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// ContainerLogs returns the last N lines of a container's logs.
func (d *DockerService) ContainerLogs(name string, tail int) (string, error) {
	if name == "" {
		return "", fmt.Errorf("容器名不能为空")
	}
	if tail <= 0 {
		tail = 50
	}
	out, err := exec.Command("docker", "logs", "--tail", fmt.Sprintf("%d", tail), name).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker logs %s: %w", name, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// DockerStats represents per-container resource usage from `docker stats --no-stream`.
type DockerStats struct {
	Name           string `json:"name"`
	CPUPercent     string `json:"cpu_percent"`
	MemUsage       string `json:"mem_usage"`
	MemPercent     string `json:"mem_percent"`
	NetIO          string `json:"net_io"`
	BlockIO        string `json:"block_io"`
	PIDs           string `json:"pids"`
}

// GetStats returns resource usage for running containers.
func (d *DockerService) GetStats() ([]DockerStats, error) {
	out, err := exec.Command("docker", "stats", "--no-stream", "--format",
		"{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.PIDs}}").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker stats: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var stats []DockerStats
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 7)
		if len(parts) >= 6 {
			s := DockerStats{
				Name:       parts[0],
				CPUPercent: parts[1],
				MemUsage:   parts[2],
				MemPercent: parts[3],
				NetIO:      parts[4],
				BlockIO:    parts[5],
			}
			if len(parts) >= 7 {
				s.PIDs = parts[6]
			}
			stats = append(stats, s)
		}
	}
	return stats, nil
}

// DockerOverview is a summary for the dashboard widget.
type DockerOverview struct {
	Available   bool              `json:"available"`
	Running     int               `json:"running"`
	Stopped     int               `json:"stopped"`
	Total       int               `json:"total"`
	Containers  []DockerContainer `json:"containers,omitempty"`
	Stats       []DockerStats     `json:"stats,omitempty"`
}

// GetOverview returns a dashboard-friendly Docker summary.
func (d *DockerService) GetOverview() (*DockerOverview, error) {
	if !d.Available() {
		return &DockerOverview{Available: false}, nil
	}

	all, err := d.ListContainers(true)
	if err != nil {
		return nil, err
	}

	overview := &DockerOverview{
		Available:  true,
		Total:      len(all),
		Containers: all,
	}
	for _, c := range all {
		if c.State == "running" {
			overview.Running++
		} else {
			overview.Stopped++
		}
	}

	// Stats only for running containers
	if overview.Running > 0 {
		stats, err := d.GetStats()
		if err == nil {
			overview.Stats = stats
		}
	}

	return overview, nil
}

// toJSON helper for logging/debugging.
func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// runCmd runs a shell command and returns stdout.
func runCmd(name string, args ...string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

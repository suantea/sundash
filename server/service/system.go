package service

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// SystemStats holds system monitoring metrics.
type SystemStats struct {
	Timestamp   int64       `json:"timestamp"`
	CPU         CPUStats    `json:"cpu"`
	Memory      MemoryStats `json:"memory"`
	Disk        DiskStats   `json:"disk"`
	Network     NetworkStats `json:"network"`
	Host        HostInfo    `json:"host"`
	Processes   ProcessStats `json:"processes"`
}

type CPUStats struct {
	UsagePercent float64   `json:"usage_percent"`
	CoreCount    int       `json:"core_count"`
	ModelName    string    `json:"model_name"`
	LoadAvg      [3]uint64 `json:"load_avg_1_5_15"` // 百分比格式 (load*100)
}

type MemoryStats struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapPercent float64 `json:"swap_percent"`
}

type DiskStats struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	Partitions  []PartitionStat `json:"partitions"`
}

type PartitionStat struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	Fstype      string  `json:"fstype"`
}

type NetworkStats struct {
	BytesSent     uint64           `json:"bytes_sent"`
	BytesRecv     uint64           `json:"bytes_recv"`
	PacketsSent   uint64           `json:"packets_sent"`
	PacketsRecv   uint64           `json:"packets_recv"`
	Interfaces    []InterfaceStat  `json:"interfaces"`
}

type InterfaceStat struct {
	Name      string `json:"name"`
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
}

type HostInfo struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Platform      string `json:"platform"`
	KernelVersion string `json:"kernel_version"`
	Uptime        uint64 `json:"uptime_seconds"`
	BootTime      uint64 `json:"boot_time"`
}

type ProcessStats struct {
	NumGoroutines int `json:"num_goroutines"`
	NumCPU        int `json:"num_cpu"`
}

// SystemService provides system monitoring capabilities.
type SystemService struct {
	lastNetStats *net.IOCountersStat
	lastNetTime  time.Time
}

// NewSystemService creates a new SystemService.
func NewSystemService() *SystemService {
	// Initialize baseline for network delta calculation
	netStats, _ := net.IOCounters(false)
	if len(netStats) > 0 {
		return &SystemService{
			lastNetStats: &netStats[0],
			lastNetTime:  time.Now(),
		}
	}
	return &SystemService{}
}

// GetStats collects and returns current system statistics.
func (s *SystemService) GetStats() (SystemStats, error) {
	stats := SystemStats{Timestamp: time.Now().Unix()}

	if err := s.collectCPU(&stats); err != nil {
		return stats, fmt.Errorf("collect cpu: %w", err)
	}
	if err := s.collectMemory(&stats); err != nil {
		return stats, fmt.Errorf("collect memory: %w", err)
	}
	if err := s.collectDisk(&stats); err != nil {
		return stats, fmt.Errorf("collect disk: %w", err)
	}
	if err := s.collectNetwork(&stats); err != nil {
		return stats, fmt.Errorf("collect network: %w", err)
	}
	s.collectHost(&stats)
	s.collectProcesses(&stats)

	return stats, nil
}

func (s *SystemService) collectCPU(stats *SystemStats) error {
	percent, err := cpu.Percent(0, false)
	if err == nil && len(percent) > 0 {
		stats.CPU.UsagePercent = percent[0]
	}
	counts, err := cpu.Counts(true)
	if err == nil {
		stats.CPU.CoreCount = counts
	}
	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		stats.CPU.ModelName = info[0].ModelName
	}
	return nil
}

func (s *SystemService) collectMemory(stats *SystemStats) error {
	vm, err := mem.VirtualMemory()
	if err == nil {
		stats.Memory.Total = vm.Total
		stats.Memory.Used = vm.Used
		stats.Memory.Free = vm.Free
		stats.Memory.UsedPercent = vm.UsedPercent
	}
	swap, err := mem.SwapMemory()
	if err == nil {
		stats.Memory.SwapTotal = swap.Total
		stats.Memory.SwapUsed = swap.Used
		stats.Memory.SwapPercent = swap.UsedPercent
	}
	return nil
}

func (s *SystemService) collectDisk(stats *SystemStats) error {
	usage, err := disk.Usage("/")
	if err == nil {
		stats.Disk.Total = usage.Total
		stats.Disk.Used = usage.Used
		stats.Disk.Free = usage.Free
		stats.Disk.UsedPercent = usage.UsedPercent
	}

	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, p := range partitions {
			pu, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			stats.Disk.Partitions = append(stats.Disk.Partitions, PartitionStat{
				Device:      p.Device,
				MountPoint:  p.Mountpoint,
				Total:       pu.Total,
				Used:        pu.Used,
				Free:        pu.Free,
				UsedPercent: pu.UsedPercent,
				Fstype:      p.Fstype,
			})
		}
	}
	return nil
}

func (s *SystemService) collectNetwork(stats *SystemStats) error {
	netStats, err := net.IOCounters(false)
	if err == nil && len(netStats) > 0 {
		current := netStats[0]
		now := time.Now()

		stats.Network.BytesSent = current.BytesSent
		stats.Network.BytesRecv = current.BytesRecv
		stats.Network.PacketsSent = current.PacketsSent
		stats.Network.PacketsRecv = current.PacketsRecv

		// Per-interface stats
		perIface, err := net.IOCounters(true)
		if err == nil {
			for _, iface := range perIface {
				if iface.Name == "lo" {
					continue
				}
				stats.Network.Interfaces = append(stats.Network.Interfaces, InterfaceStat{
					Name:      iface.Name,
					BytesSent: iface.BytesSent,
					BytesRecv: iface.BytesRecv,
				})
			}
		}

		_ = now
		_ = s.lastNetStats
	}
	return nil
}

func (s *SystemService) collectHost(stats *SystemStats) {
	info, err := host.Info()
	if err == nil {
		stats.Host.Hostname = info.Hostname
		stats.Host.OS = info.OS
		stats.Host.Platform = info.Platform
		stats.Host.KernelVersion = info.KernelVersion
		stats.Host.Uptime = info.Uptime
		stats.Host.BootTime = info.BootTime
	}
}

func (s *SystemService) collectProcesses(stats *SystemStats) {
	stats.Processes.NumGoroutines = runtime.NumGoroutine()
	stats.Processes.NumCPU = runtime.NumCPU()
}

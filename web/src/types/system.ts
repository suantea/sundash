export interface CPUStats {
  usage_percent: number
  core_count: number
  model_name: string
  load_avg_1_5_15: [number, number, number]
}

export interface MemoryStats {
  total: number
  used: number
  free: number
  used_percent: number
  swap_total: number
  swap_used: number
  swap_percent: number
}

export interface PartitionStat {
  device: string
  mount_point: string
  total: number
  used: number
  free: number
  used_percent: number
  fstype: string
}

export interface DiskStats {
  total: number
  used: number
  free: number
  used_percent: number
  partitions: PartitionStat[]
}

export interface InterfaceStat {
  name: string
  bytes_sent: number
  bytes_recv: number
}

export interface NetworkStats {
  bytes_sent: number
  bytes_recv: number
  packets_sent: number
  packets_recv: number
  interfaces: InterfaceStat[]
}

export interface HostInfo {
  hostname: string
  os: string
  platform: string
  kernel_version: string
  uptime_seconds: number
  boot_time: number
}

export interface ProcessStats {
  num_goroutines: number
  num_cpu: number
}

export interface SystemStats {
  timestamp: number
  cpu: CPUStats
  memory: MemoryStats
  disk: DiskStats
  network: NetworkStats
  host: HostInfo
  processes: ProcessStats
}

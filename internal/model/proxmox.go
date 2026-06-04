package model

import (
	"net"
	"time"
)

// Cluster maps to proxmox.clusters
type Cluster struct {
	ClusterID      string    `json:"cluster_id"`
	Name           string    `json:"name"`
	APIURL         string    `json:"api_url"`
	Port           int       `json:"port"`
	Username       string    `json:"username"`
	APITokenID     *string   `json:"api_token_id,omitempty"`
	APITokenSecret *string   `json:"-"` // never expose secret in JSON
	TLSVerify      bool      `json:"tls_verify"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// Node maps to proxmox.nodes
type Node struct {
	NodeID              string     `json:"node_id"`
	ClusterID           string     `json:"cluster_id"`
	ProxmoxNodeName     string     `json:"proxmox_node_name"`
	Status              *string    `json:"status,omitempty"`
	CPUUsage            *float64   `json:"cpu_usage,omitempty"`
	MemoryTotal         *int64     `json:"memory_total,omitempty"`
	MemoryUsed          *int64     `json:"memory_used,omitempty"`
	LastSyncedAt        *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

// Instance maps to proxmox.instances
type Instance struct {
	InstanceID            string     `json:"instance_id"`
	ClusterID             string     `json:"cluster_id"`
	NodeID                *string    `json:"node_id,omitempty"`
	OrganizationID        string     `json:"organization_id"`
	OwnerAgentID          *string    `json:"owner_agent_id,omitempty"`
	ProxmoxVMID           int        `json:"proxmox_vmid"`
	InstanceType          string     `json:"instance_type"` // qemu or lxc
	Name                  string     `json:"name"`
	Status                *string    `json:"status,omitempty"`
	IPAddress             *net.IP    `json:"ip_address,omitempty"`
	CPUCores              *int       `json:"cpu_cores,omitempty"`
	MemoryMB              *int       `json:"memory_mb,omitempty"`
	DiskGB                *int       `json:"disk_gb,omitempty"`
	GuacamoleConnectionID *int       `json:"guacamole_connection_id,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}

// ProxmoxSession maps to proxmox.sessions
type ProxmoxSession struct {
	SessionID          string     `json:"session_id"`
	ClusterID          *string    `json:"cluster_id,omitempty"`
	InstanceID         *string    `json:"instance_id,omitempty"`
	CreatedByAgentID   *string    `json:"created_by_agent_id,omitempty"`
	GuacConnectionID   *int       `json:"guac_connection_id,omitempty"`
	ProxmoxUPID        *string    `json:"proxmox_upid,omitempty"`
	Protocol           *string    `json:"protocol,omitempty"`
	TaskType           string     `json:"task_type"`
	Status             string     `json:"status"`
	Message            *string    `json:"message,omitempty"`
	StartedAt          time.Time  `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
}

// Template maps to proxmox.templates
type Template struct {
	TemplateID        string    `json:"template_id"`
	ClusterID         *string   `json:"cluster_id,omitempty"`
	Name              string    `json:"name"`
	ProxmoxTemplateID int       `json:"proxmox_template_id"`
	OSType            *string   `json:"os_type,omitempty"`
	DefaultCPU        *int      `json:"default_cpu,omitempty"`
	DefaultMemoryMB   *int      `json:"default_memory_mb,omitempty"`
	DefaultDiskGB     *int      `json:"default_disk_gb,omitempty"`
	NetworkBridge     *string   `json:"network_bridge,omitempty"`
	StorageName       *string   `json:"storage_name,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}


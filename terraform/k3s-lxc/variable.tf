variable "proxmox_node" {
  description = "Nama target node di Proxmox VE"
  type        = string
}

variable "template_name" {
  description = "Template OS LXC yang digunakan"
  type        = string
}

variable "storage_pool" {
  description = "Storage pool untuk disk container"
  type        = string
}

variable "ssh_public_keys" {
  description = "Public key SSH untuk akses ke dalam LXC"
  type        = string
}
variable "network_bridge" {
  description = "Bridge interface di Proxmox"
  type        = string
  default     = "vmbr0"
}

variable "default_gateway" {
  description = "Default gateway untuk akses internet node"
  type        = string
}
variable "proxmox_api_url" {
  type        = string
  description = "URL API Server Proxmox"
}
variable "k3s_nodes" {
  description = "Konfigurasi node k3s cluster"
  type = map(object({
    vmid    = number
    cores   = number
    memory  = number
    ip_addr = string
    role    = string
  }))
}

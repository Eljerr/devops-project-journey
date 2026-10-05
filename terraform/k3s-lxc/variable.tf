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

variable "ssh_public_key" {
  description = "Public key SSH untuk akses ke dalam LXC"
  type        = string
}


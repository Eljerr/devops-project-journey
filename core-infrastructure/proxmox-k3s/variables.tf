variable "proxmox_api_id" {
  type = string
}

variable "proxmox_api_secret" {
  type      = string
  sensitive = true
}

variable "proxmox_api_url" {
  type        = string
  description = "URL API Proxmox (https://<ip>:8006/api2/json)"
}

variable "xlc_password" {
  type = string
}

variable "xlc_template" {
  type = string
}

variable "ssh_public_keys" {
  type = string
}

variable "lxc_ip_base" {
  type        = string
  description = "Base IP address untuk LXC (misal 192.168.18.20, nantinya akan ditambah count.index)"
  default     = "192.168.18.20"
}

variable "lxc_gateway" {
  type        = string
  description = "Gateway IP address untuk LXC"
  default     = "192.168.18.1"
}

variable "ssh_private_key_path" {
  type        = string
  description = "Path ke private key SSH untuk remote-exec provisioner"
  default     = "~/.ssh/id_ed25519"
}

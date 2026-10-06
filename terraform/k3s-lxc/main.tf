resource "proxmox_lxc" "k3s_nodes" {
  for_each     = var.k3s_nodes
  target_node  = var.proxmox_node
  vmid         = each.value.vmid
  hostname     = each.key
  ostemplate   = var.template_name
  unprivileged = true
  start        = true

  cores  = each.value.cores
  memory = each.value.memory
  swap   = 512

  features {
    nesting = true
    keyctl  = true
  }

  network {
    name   = "eth0"
    bridge = var.network_bridge
    ip     = each.value.ip_addr
    gw     = var.default_gateway
  }

  ssh_public_keys = var.ssh_public_keys

  rootfs {
    storage = var.storage_pool
    size    = "7G"
  }
}

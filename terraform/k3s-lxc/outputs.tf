output "node_details" {
  description = "Informasi mengenai node k3s yang telah dibuat"
  value = {
    for k, v in proxmox_lxc.k3s_nodes : k => {
      vmid = v.vmid
      ip   = v.network[0].ip
    }
  }
}

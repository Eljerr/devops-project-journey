terraform {

  required_version = ">= 0.13.0"
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = "2.9.3"
    }
  }
}

provider "proxmox" {
  pm_api_url          = var.proxmox_api_url
  pm_api_token_id     = var.proxmox_api_id
  pm_api_token_secret = var.proxmox_api_secret

  pm_tls_insecure = true
}

resource "proxmox_lxc" "k8s_nodes" {
  target_node     = "desktop"
  hostname        = "k8s-master"
  ostemplate      = var.xlc_template
  password        = var.xlc_password
  unprivileged    = true
  vmid            = 200
  ssh_public_keys = var.ssh_public_keys

  cores  = 2
  memory = 2048
  swap   = 0
  start  = true

  rootfs {
    storage = "local-lvm"
    size    = "10G"
  }

  network {
    name   = "eth0"
    bridge = "vmbr0"
    ip     = "${var.lxc_ip_base}0/24"
    gw     = var.lxc_gateway
  }

  connection {
    type        = "ssh"
    user        = "root"
    host        = "${var.lxc_ip_base}0"
    private_key = file(var.ssh_private_key_path)
  }

  provisioner "remote-exec" {
    inline = [
      "apt-get update",
      "apt-get install -y curl",
      "curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=\"--kubelet-arg=feature-gates=KubeletInUserNamespace=true\" sh -"
    ]
  }

  features {
    nesting = true
  }

}


resource "proxmox_lxc" "k8s_workers" {
  count = 2

  target_node     = "desktop"
  hostname        = "k8s-workers-${count.index + 1}"
  ostemplate      = var.xlc_template
  password        = var.xlc_password
  unprivileged    = true
  vmid            = 201 + count.index
  ssh_public_keys = var.ssh_public_keys

  cores  = 2
  memory = 2048
  swap   = 0
  start  = true

  rootfs {
    storage = "local-lvm"
    size    = "10G"
  }

  network {
    name   = "eth0"
    bridge = "vmbr0"
    ip     = "${var.lxc_ip_base}${count.index + 1}/24"
    gw     = var.lxc_gateway
  }

  connection {
    type        = "ssh"
    user        = "root"
    host        = "${var.lxc_ip_base}${count.index + 1}"
    private_key = file(var.ssh_private_key_path)
  }

  provisioner "remote-exec" {
    inline = [
      "apt-get update",
      "apt-get install -y curl"
    ]
  }

  features {
    nesting = true
  }

}


resource "proxmox_lxc" "monitoring" {

  target_node     = "desktop"
  hostname        = "monitoring-secops"
  ostemplate      = var.xlc_template
  password        = var.xlc_password
  unprivileged    = true
  vmid            = 203
  ssh_public_keys = var.ssh_public_keys

  cores  = 1
  memory = 1024
  swap   = 512
  start  = true

  rootfs {
    storage = "local-lvm"
    size    = "10G"
  }

  network {
    name   = "eth0"
    bridge = "vmbr0"
    ip     = "${var.lxc_ip_base}3/24"
    gw     = var.lxc_gateway
  }

  connection {
    type        = "ssh"
    user        = "root"
    host        = "${var.lxc_ip_base}3"
    private_key = file(var.ssh_private_key_path)
  }

  provisioner "remote-exec" {
    inline = [
      "apt-get update",
      "apt-get install -y vim"
    ]
  }

  features {
    nesting = true
  }

}

output "lxc_ips" {
  description = "Daftar IP untuk K8s Master dan Workers"
  value = merge(
    { (proxmox_lxc.k8s_nodes.hostname) = proxmox_lxc.k8s_nodes.network[0].ip },
    { for node in proxmox_lxc.k8s_workers : node.hostname => node.network[0].ip }
  )
}

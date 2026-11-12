---
runbook-id: RB-XXXX
title: "{Judul Prosedur}"
date: YYYY-MM-DD
last-tested: YYYY-MM-DD
author: Asylon
target-system: "{INNER CachyOS | OUTER Proxmox | OUTER STB | K3s Cluster}"
tested-versions: "{misal: Ubuntu 24.04 ARM64, OpenSSH 9.x}"
automation-status: "{Manual CLI | Bash Scripted | Terraform/Ansible}"
tags: [runbook, sop, operations]
---

# Runbook: {Judul}

## 1. Tujuan & Kapan Dipakai
<!-- OPERATOR VOICE -->
* **Tujuan**: {…}
* **Pemicu (gejala/alert yang membuat runbook ini dibuka)**: {…}

## 2. Prasyarat & Pre-Flight Check
* **Akses**: {…}
* **Sesi cadangan** (wajib untuk perubahan jaringan/SSH): {…}
* **Cek kondisi awal**:
  ```bash
  {…}
  ```

## 3. Backup Sebelum Perubahan
```bash
{Perintah backup — misal: sudo cp -a /etc/ssh/sshd_config /etc/ssh/sshd_config.bak.$(date +%F)}
```

## 4. Langkah Eksekusi
<!-- Perintah idempoten + alasan teknis. Sertakan validasi sintaks sebelum reload (nginx -t, sshd -t, dll.) -->

### Langkah 4.1: {…}
* **Alasan teknis**: <!-- OPERATOR VOICE --> {…}
* **Perintah**:
  ```bash
  {…}
  ```
* **Validasi langkah ini**: `{…}`

### Langkah 4.2: {…}
* **Alasan teknis**: <!-- OPERATOR VOICE --> {…}
* **Perintah**:
  ```bash
  {…}
  ```

## 5. Verifikasi Pasca-Eksekusi
```bash
{Perintah verifikasi}
```
Hasil yang diharapkan (diambil dari uji nyata pada tanggal `last-tested`):
```text
{Output asli}
```

## 6. Rollback
1. {Langkah rollback + perintah}
2. {Verifikasi kembali ke kondisi awal}

## 7. Eskalasi
* Jika rollback gagal: {akses konsol fisik / serial / Proxmox console / dll.}

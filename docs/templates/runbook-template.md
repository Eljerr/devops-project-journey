---
runbook-id: RB-XXXX
title: "{Judul Prosedur Operasional}"
date: YYYY-MM-DD
target-system: "{Inner Loop WSL2 | Proxmox Host | STB Gateway | K3s Cluster}"
automation-status: "{Manual CLI | Bash Scripted | Terraform/Ansible}"
tags: [runbook, sop, operations]
---

# Runbook: {Judul Prosedur Operasional & Hardening}

## 1. Tujuan & Cakupan Operasional
<!-- OPERATOR VOICE: Jelaskan kapan runbook ini dipakai dan masalah operasional apa yang dicegah/ditangani -->
{Penjelasan singkat fungsi prosedur ini di lingkungan produksi/homelab.}

## 2. Prasyarat & Pemeriksaan Awal (Pre-Flight Checks)
Pastikan kondisi berikut terpenuhi sebelum mengeksekusi perubahan:
* **Target Host & Akses**: {Misal: SSH non-root ke STB ARM64 / Proxmox}
* **Perintah Cek Kondisi Awal**:
  ```bash
  {Perintah untuk memeriksa status sistem sebelum diubah}
  ```

## 3. Langkah Eksekusi Standar (Step-by-Step Execution)
<!-- Cantumkan perintah yang bersifat idempoten beserta alasan teknis di balik flag/parameter yang dipilih -->

### Langkah 3.1: {Nama Langkah Pertama}
* **Alasan Teknis Operator**: {Mengapa perintah/konfigurasi ini dijalankan}
* **Perintah / Konfigurasi**:
  ```bash
  {Perintah CLI atau blok konfigurasi}
  ```

### Langkah 3.2: {Nama Langkah Kedua}
* **Alasan Teknis Operator**: {Mengapa parameter ini dipilih}
* **Perintah / Konfigurasi**:
  ```bash
  {Perintah CLI atau blok konfigurasi}
  ```

## 4. Verifikasi Kesehatan & Validasi (Post-Execution Check)
Jalankan perintah berikut untuk memastikan konfigurasi berjalan normal tanpa merusak layanan lain:
```text
{Perintah verifikasi dan bukti output terminal yang diharapkan}
```

## 5. Prosedur Pemulihan Darurat (Rollback Plan)
<!-- Sangat dinilai tinggi oleh Hiring Manager: Bagaimana jika konfigurasi ini bikin server hang atau putus koneksi? -->
Jika terjadi kegagalan layanan setelah prosedur di atas dijalankan, lakukan langkah pemulihan berikut:
1. {Langkah rollback 1 beserta perintah CLI-nya}
2. {Verifikasi layanan kembali ke kondisi semula}
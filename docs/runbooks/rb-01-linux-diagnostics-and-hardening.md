---
runbook-id: RB-0001
title: "Prosedur Diagnostik Layanan systemd, Penanganan OOM-Kill, dan Penguatan Akses SSH"
date: 2026-11-11
target-system: "Inner Loop WSL2 & Outer Loop Host (STB / Proxmox)"
automation-status: "Manual CLI"
operator: "Eljerr"
reviewed-by: "AI Systems Architect Agent"
tags: [runbook, systemd, cgroups-v2, oom-kill, ssh-hardening, troubleshooting]
---

# Runbook: Diagnostik Layanan systemd, Penanganan OOM-Kill, dan Penguatan Akses SSH

## 1. Tujuan & Cakupan Operasional

Runbook ini adalah Standar Operasional Prosedur (SOP) untuk dua kebutuhan rekayasa sistem:
1. **Investigasi & Triage Layanan `systemd`**: Prosedur cepat mendiagnosis dan menangani layanan latar belakang (*background service*) yang mengalami kegagalan berulang (*auto-restart loop*) akibat terminasi memori kernel Linux (`Result: oom-kill`).
2. **Penguatan Akses Host (SSH Hardening)**: Prosedur standar pembuatan kunci kriptografi eliptis modern (Ed25519), verifikasi izin berkas (*file permissions*), dan pengetatan konfigurasi SSH daemon pada host lingkungan *Outer Loop* (STB Gateway & Proxmox Host) guna memitigasi serangan *brute-force* dan *credential spraying*.

---

## 2. Prasyarat & Pemeriksaan Awal (Pre-Flight Checks)

Sebelum mengeksekusi tindakan diagnostik atau perubahan konfigurasi, pastikan verifikasi berikut:
* **Hak Akses**: Pengguna memiliki akses `sudo` untuk membaca jurnal sistem dan memodifikasi konfigurasi di `/etc/systemd/` dan `/etc/ssh/`.
* **Pemeriksaan Status Awal Sistem**:
  ```bash
  # Cek kapasitas memori fisik dan swap yang tersedia di host
  free -h

  # Cek status cgroups v2 aktif pada hierarki system.slice
  cat /sys/fs/cgroup/system.slice/cgroup.controllers

  # Verifikasi PID 1 dijalankan oleh systemd
  ps -p 1 -o comm=
  ```

---

## 3. Prosedur Investigasi & Triage Layanan systemd (OOM-Kill)

### Langkah 3.1: Deteksi Gejala & Bukti Forensik
Jika suatu layanan dicurigai mengalami *crash*, jalankan perintah investigasi primer:
```bash
# 1. Periksa status unit dan ringkasan metrik puncak (peak metrics)
systemctl status <nama-layanan>.service

# 2. Periksa baris log terakhir yang dicatat sebelum proses mati
journalctl -u <nama-layanan>.service -n 50 --no-pager
```

**Indikator Positif OOM-Kill pada Log**:
* `Active: activating (auto-restart) (Result: oom-kill)`
* `Main PID: <PID> (code=killed, signal=KILL)` — Menunjukkan sinyal `SIGKILL` (Sinyal 9) ditembakkan oleh kernel.
* `Mem peak: <BatasMemory>` — Menunjukkan konsumsi proses menyentuh batas `MemoryMax` cgroups v2.
* Log sistem mencatat: `The kernel OOM killer killed some processes in this unit.`

---

### Langkah 3.2: Alur Keputusan & Respons Operasional Operator

<!-- OPERATOR VOICE (VERBATIM ASSEMBLY) -->
> **Analisis & Pertimbangan Teknis Operator**:
> *"Service yang sudah dialokasikan resource-nya di produksi pasti sudah dites dengan teliti batas minimum dan maksimum penggunaan memory-nya. Jika saya melihat ada service yang dimatikan oleh OOM Killer, dan jika kapasitas resource di host masih tersedia, langkah pertama saya adalah menaikkan batas `MemoryMax`-nya sementara sebagai penanganan cepat (triage). Setelah itu, saya segera memberitahukan temuan log ini kepada developer agar diinvestigasi dan ditangani lebih lanjut (root-cause analysis)."*

**Langkah Mitigasi Sementara (Stopgap)**:
1. Jika kapasitas RAM host masih aman (`free -h`), longgarkan batas di `/etc/systemd/system/<nama-layanan>.service`:
   ```ini
   [Service]
   # Naikkan batas sesuai kalkulasi kapasitas aman host
   MemoryMax=100M
   ```
2. Muat ulang daemon dan restart unit:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart <nama-layanan>.service
   ```
3. Eskalasikan artefak log `journalctl` kepada tim pengembang aplikasi untuk profiling memori (mencegah *memory leak* jangka panjang).

---

### Langkah 3.3: Catatan Teknis: Fenomena ZRAM vs Data Ber-entropi Tinggi

<!-- OPERATOR VOICE (VERBATIM ASSEMBLY) -->
> **Temuan Kritis Mekanisme Kernel (ZRAM vs Random Data)**:
> *"Saat pengujian, alokasi `bytearray` (80 MB) berhasil lolos tanpa dimatikan oleh OOM Killer karena kernel menggunakan algoritma kompresi zstd di ZRAM. Data yang dihasilkan oleh `bytearray` bernilai `0x00` (angka nol semua), yang sangat memungkinkan untuk dikompresi hingga ukurannya menyusut drastis di RAM. Berbeda dengan `os.urandom`, data byte yang dihasilkan tidak beraturan (acak), sehingga mustahil dikompresi oleh zstd. Begitu memori acak ini melampaui `MemoryMax` dan swap dimatikan (`MemorySwapMax=0`), kernel langsung mengeksekusi OOM Killer."*

**Pelajaran Rekayasa**:
Di lingkungan pengujian, jangan pernah menguji batas cgroup memori menggunakan alokasi byte kosong/nol jika host mengaktifkan modul ZRAM/ZSWAP, karena hasil tes akan terdistorsi oleh kompresi transparan kernel Linux.

---

## 4. Prosedur Penguatan Akses Host (SSH Hardening)

### Langkah 4.1: Pembuatan Kunci Kriptografi Asimetris Ed25519 (Inner Loop)
Buat kunci modern berbasis kurva eliptis di mesin lokal:
```bash
ssh-keygen -t ed25519 -C "joy@homelab"
```
* Lokasi Private Key: `~/.ssh/id_ed25519` (Gembok rahasia, dilarang disalin keluar dari mesin lokal).
* Lokasi Public Key: `~/.ssh/id_ed25519.pub` (Gembok publik, dikirim ke server tujuan).

---

### Langkah 4.2: Audit Ketat Hak Akses Berkas (File Permissions)
SSH client akan menolak melakukan otentikasi jika izin berkas kunci terlalu terbuka. Standar izin wajib:
```bash
chmod 700 ~/.ssh
chmod 600 ~/.ssh/id_ed25519
chmod 644 ~/.ssh/id_ed25519.pub
```

Bukti verifikasi (`ls -la ~/.ssh/`):
```text
drwx------ 2 joy joy 4096 Nov 10 07:34 .
-rw------- 1 joy joy  419 Nov 10 07:32 id_ed25519
-rw-r--r-- 1 joy joy  106 Nov 10 07:32 id_ed25519.pub
```

---

### Langkah 4.3: Penguatan Daemon SSH di Host Remote (`/etc/ssh/sshd_config`)
Kirimkan public key ke server tujuan:
```bash
ssh-copy-id -i ~/.ssh/id_ed25519.pub <user>@<ip-server>
```

Terapkan parameter pengerasan berikut pada `/etc/ssh/sshd_config` (atau di `/etc/ssh/sshd_config.d/99-hardening.conf`):
```ini
# Matikan autentikasi kata sandi untuk mencegah brute-force
PasswordAuthentication no

# Larang login langsung sebagai root
PermitRootLogin no

# Aktifkan otentikasi berbasis public key
PubkeyAuthentication yes
```

> [!CAUTION]
> **ATURAN EMAS ANTI-LOCKOUT**:
> Sebelum me-restart layanan `ssh` / `sshd` di server remote, **JANGAN PERNAH menutup jendela sesi terminal yang sedang aktif**. 
> Buka tab terminal baru terpisah dan coba lakukan login: `ssh <user>@<ip-server>`.  
> Jika koneksi gagal, perbaiki konfigurasi pada tab pertama yang masih terbuka agar kamu tidak terkunci di luar server secara permanen!

Muat ulang layanan SSH:
```bash
sudo systemctl reload sshd || sudo systemctl reload ssh
```

---

## 5. Verifikasi Kesehatan & Validasi Pasca-Eksekusi

1. **Validasi Unit systemd**:
   ```bash
   systemctl is-active <nama-layanan>.service
   # Keluaran yang diharapkan: active
   ```
2. **Validasi cgroups v2 Memory Limit**:
   ```bash
   cat /sys/fs/cgroup/system.slice/<nama-layanan>.service/memory.max
   # Keluaran yang diharapkan: angka byte sesuai nilai MemoryMax (misal: 104857600 untuk 100M)
   ```
3. **Validasi SSH Hardening**:
   ```bash
   # Uji coba login dari workstation lokal tanpa prompt password:
   ssh -o PubkeyAuthentication=yes -i ~/.ssh/id_ed25519 <user>@<ip-server> "echo 'SSH Key Authentication Berhasil!'"
   ```

---

## 6. Prosedur Pemulihan Darurat (Rollback Plan)

1. **Rollback Service systemd**:
   Jika penyesuaian `MemoryMax` justru membuat host kehabisan RAM atau memicu ketidakstabilan sistem:
   ```bash
   sudo systemctl stop <nama-layanan>.service
   # Kembalikan file unit ke cadangan awal atau set MemoryMax ke nilai aman
   sudo systemctl daemon-reload
   ```
2. **Rollback SSH Hardening**:
   Jika konfigurasi SSH baru bermasalah (dari sesi aktif yang belum ditutup):
   ```bash
   sudo cp /etc/ssh/sshd_config.bak /etc/ssh/sshd_config
   sudo systemctl reload sshd || sudo systemctl reload ssh
   ```

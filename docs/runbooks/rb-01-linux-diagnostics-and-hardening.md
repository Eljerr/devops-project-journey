---
runbook-id: RB-0001
title: "Prosedur Diagnostik Layanan systemd (cgroups v2, OOM-Kill, start-limit-hit) dan Penguatan Akses SSH (Anti-Lockout)"
date: 2026-11-20
last-tested: 2026-11-20
author: Asylon
target-system: "INNER CachyOS & OUTER Host (Proxmox / STB)"
tested-versions: "CachyOS Linux (Kernel 7.3.3-1-cachyos, systemd 262, OpenSSH_10.6p2, fish 4.9.3)"
automation-status: "Manual CLI"
tags: [runbook, sop, operations, systemd, cgroups-v2, oom-kill, ssh-hardening, fase-1]
---

# Runbook: Diagnostik Layanan `systemd` (`cgroups v2`, OOM-Kill, `start-limit-hit`) dan Penguatan Akses SSH

## 1. Tujuan & Kapan Dipakai
<!-- OPERATOR VOICE -->
* **Tujuan**:
  1. Mendiagnosis layanan `systemd` yang mati berulang kali (`start-limit-hit`) atau terbunuh oleh sinyal `SIGKILL` (`137`), serta memastikan apakah kematiannya benar-benar disebabkan oleh **OOM Killer** lewat `memory.events` dan `journalctl`.
  2. Mengamankan akses SSH menggunakan pasangan kunci **Ed25519** dengan izin file yang ketat (`700`/`600`/`644`) dan mematikan autentikasi password (`PasswordAuthentication no`) tanpa terjebak urutan baca *drop-in* `sshd_config` maupun risiko terkunci di luar server (*lockout*).
* **Pemicu (gejala/alert yang membuat runbook ini dibuka)**:
  * Service berstatus `failed (Result: start-limit-hit)` atau mati dengan kode `137` (`SIGKILL`).
  * Penyiapan awal atau audit keamanan SSH pada node homelab (CachyOS, Proxmox VE, STB Ubuntu ARM64).

---

## 2. Prasyarat & Pre-Flight Check
* **Akses**: Hak `sudo` pada mesin target (atau `--user` untuk layanan tingkat pengguna).
* **Sesi cadangan (WAJIB untuk perubahan SSH)**: Buka minimal **2 tab terminal**. Tab 1 tetap terhubung ke sesi SSH aktif di dalam server sebagai jalur penyelamatan (*lifeline*); Tab 2 digunakan untuk menguji koneksi login baru setelah perubahan. **Jangan pernah menutup Tab 1 sebelum Tab 2 terbukti berhasil login!**
* **Cek kondisi awal**:
  ```fish
  # 1. Pastikan hierarki cgroups v2 aktif dan cek konfigurasi swap/ZRAM
  stat -fc %T /sys/fs/cgroup
  swapon --show

  # 2. Cek batas default StartLimit pada systemd
  systemctl show -p DefaultStartLimitBurst -p DefaultStartLimitIntervalUSec

  # 3. Cek baris aktif pada /etc/ssh/sshd_config (perhatikan posisi Include di baris atas!)
  grep -n -v -E '^\s*#|^\s*$' /etc/ssh/sshd_config
  ```

---

## 3. Backup Sebelum Perubahan
Sebelum mengedit file unit `systemd` atau konfigurasi `/etc/ssh/sshd_config`, **wajib** buat salinan cadangan (`.bak`) agar langkah *Rollback* (Bagian 6) dapat dieksekusi jika terjadi kesalahan:

```fish
# Backup konfigurasi utama sshd_config dengan penanda tanggal
sudo cp -a /etc/ssh/sshd_config /etc/ssh/sshd_config.bak.(date +%F)
```

---

## 4. Langkah Eksekusi

### Langkah 4.1: Diagnostik Service Crash (`start-limit-hit` vs OOM Killer) & Pembatasan `cgroups v2`
* **Alasan teknis**: <!-- OPERATOR VOICE -->
  * `[Unit]` berfungsi mengatur urutan, identitas, dan ketergantungan; `[Service]` mengatur cara menjalankan dan mengelola proses; `[Install]` mengatur pendaftaran autostart ke target boot tertentu.
  * Ketika kita menuliskan konfigurasi kuota (`MemoryMax=`, `MemorySwapMax=`, `CPUQuota=`) pada file `.service`, sebenarnya `systemd` mensinkronisasinya langsung ke file kernel di folder `/sys/fs/cgroup`.
  * Konfigurasi baru hanya dibaca setelah kita mengetik `daemon-reload` (ibarat tombol *save* agar perubahan yang belum selesai diketik tidak langsung merusak sistem secara *real-time*).
  * Kode keluar `137` hanya berarti proses mati karena `SIGKILL` (`128 + 9`). Untuk memastikan pelakunya adalah OOM Killer, kita wajib memeriksa counter `oom_kill` pada `memory.events` atau memeriksa `journalctl`. Selain itu, pada sistem dengan swap ZRAM aktif, `MemorySwapMax=0` wajib dipasang jika kita ingin mencegah kelebihan memori (baik data berulang maupun data acak `os.urandom`) dipindahkan ke swap.
* **Perintah Diagnostik & Penerapan Kuota**:
  ```fish
  # 1. Cek dasbor spesifikasi & status terkini service
  systemctl status <nama-unit>.service

  # 2. Buka buku rekaman log lengkap di journalctl
  journalctl -u <nama-unit>.service -n 30 --no-pager

  # 3. Cek bukti sah OOM Killer pada cgroup unit terkait
  cat /sys/fs/cgroup/system.slice/<nama-unit>.service/memory.events

  # 4. Setelah mengedit file .service (misal menambah MemoryMax / memperbaiki ExecStart):
  sudo systemctl daemon-reload
  sudo systemctl reset-failed <nama-unit>.service
  sudo systemctl restart <nama-unit>.service
  ```
* **Validasi langkah ini**:
  ```fish
  cat /sys/fs/cgroup/system.slice/<nama-unit>.service/memory.max
  cat /sys/fs/cgroup/system.slice/<nama-unit>.service/memory.swap.max
  cat /sys/fs/cgroup/system.slice/<nama-unit>.service/cpu.max
  ```

---

### Langkah 4.2: Pembuatan Kunci SSH Ed25519 & Pengetatan Izin File (`700` / `600` / `644`)
* **Alasan teknis**: <!-- OPERATOR VOICE -->
  * Private key wajib dijaga rahasia di mesin lokal, sedangkan Public key (`.pub`) ditaruh pada server tujuan (`~/.ssh/authorized_keys`).
  * Folder `~/.ssh` wajib berizin `700` (agar hanya pemilik yang bisa masuk, membuat, atau menghapus file di dalamnya), file Private Key wajib `600`, dan Public Key `644`. Jika izin Private Key terlalu longgar (misalnya `644` sehingga bisa di-*read* orang lain), OpenSSH akan menolaknya dengan `WARNING: UNPROTECTED PRIVATE KEY FILE!` karena berisiko dipakai orang lain untuk nyelonong masuk.
* **Perintah**:
  ```fish
  # 1. Buat pasangan kunci Ed25519 (jika belum ada)
  ssh-keygen -t ed25519 -C "joy@homelab"

  # 2. Pastikan izin folder dan file kunci sesuai standar ketat OpenSSH
  chmod 700 ~/.ssh
  chmod 600 ~/.ssh/id_ed25519
  chmod 644 ~/.ssh/id_ed25519.pub
  ```
* **Validasi langkah ini**:
  ```fish
  stat -c '%a %A %n' ~/.ssh ~/.ssh/id_ed25519 ~/.ssh/id_ed25519.pub
  ```

---

### Langkah 4.3: Penguatan `sshd_config` dengan Drop-in `00-hardening.conf` (Aturan *First-Match Wins*)
* **Alasan teknis**: <!-- OPERATOR VOICE -->
  * Pada konfigurasi `sshd`, urutan bacanya adalah dari atas ke bawah sehingga konfigurasi yang ditemukan paling atas selalu diutamakan (*First-Match Wins*).
  * Karena baris `Include /etc/ssh/sshd_config.d/*.conf` berada di baris atas `/etc/ssh/sshd_config`, file bawaan seperti `50-cloud-init.conf` yang berisi `PasswordAuthentication yes` akan mengalahkan baris `PasswordAuthentication no` yang ditulis di bagian bawah `/etc/ssh/sshd_config`.
  * Untuk memenangkannya secara bersih, kita menamai file konfigurasi *drop-in* kita dengan awalan angka kecil yaitu `00-hardening.conf` agar dibaca paling pertama.
  * Sebelum me-*reload* `sshd`, kita wajib menjalankan `sshd -t` (huruf kecil) untuk memverifikasi jika ada kesalahan penulisan/sintaks, dan `sshd -T` (huruf besar) untuk merontgen nilai mana yang benar-benar menang.
* **Perintah**:
  ```fish
  # 1. Tulis aturan penguatan pada drop-in berawalan 00- agar dibaca paling awal
  printf "PasswordAuthentication no\nPermitRootLogin no\nPubkeyAuthentication yes\n" | sudo tee /etc/ssh/sshd_config.d/00-hardening.conf

  # 2. Wajib uji sintaks (-t) dan rontgen konfigurasi efektif (-T) SEBELUM reload!
  sudo sshd -t; echo "Syntax exit: $status"
  sudo sshd -T | grep -iE 'passwordauthentication|permitrootlogin|pubkeyauthentication'

  # 3. Jika dan hanya jika Syntax exit: 0 dan hasil -T sudah sesuai, reload layanan SSH
  # (Di CachyOS/Arch: sshd | Di Ubuntu STB: ssh)
  sudo systemctl reload sshd
  ```

---

## 5. Verifikasi Pasca-Eksekusi
Perintah verifikasi yang dijalankan pada tanggal `last-tested` (`2026-11-20`, lihat bukti lengkap di [lab-1.3-01-systemd-and-ssh-hardening.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.3-01-systemd-and-ssh-hardening.md)):

```fish
# 1. Verifikasi izin kunci SSH
stat -c '%a %A %n' ~/.ssh ~/.ssh/*

# 2. Verifikasi sintaks (-t) dan konfigurasi efektif (-T) setelah pemasangan 00-hardening.conf
/usr/bin/sshd -t -f ~/lab/ssh_lab/sshd_uji.conf -h ~/lab/ssh_lab/host_key; echo "Syntax exit: $status"
/usr/bin/sshd -T -f ~/lab/ssh_lab/sshd_uji.conf -h ~/lab/ssh_lab/host_key | grep -iE 'passwordauthentication|permitrootlogin'
```

Hasil nyata di terminal:
```text
700 drwx------ /home/joy/.ssh
600 -rw------- /home/joy/.ssh/id_ed25519
644 -rw-r--r-- /home/joy/.ssh/id_ed25519.pub

Syntax exit: 0
PermitRootLogin no
PasswordAuthentication no
```

---

## 6. Rollback
1. **Rollback Konfigurasi SSH (dari Tab 1 / sesi SSH cadangan yang masih terbuka)**:
   ```fish
   sudo rm -f /etc/ssh/sshd_config.d/00-hardening.conf
   sudo cp -a /etc/ssh/sshd_config.bak.(date +%F) /etc/ssh/sshd_config
   sudo sshd -t && sudo systemctl reload sshd
   ```
2. **Verifikasi kembali ke kondisi awal**:
   ```fish
   sudo sshd -T | grep -iE 'passwordauthentication|permitrootlogin'
   ```

---

## 7. Eskalasi
* Jika sesi SSH terputus dan *rollback* jarak jauh tidak dapat dilakukan:
  * **OUTER Proxmox VE**: Gunakan konsol fisik laptop Proxmox atau Web GUI Proxmox (`shell` node / konsol LXC).
  * **OUTER STB ARM64**: Hubungkan keyboard USB & monitor HDMI langsung ke STB untuk memperbaiki `/etc/ssh/sshd_config.d/` secara lokal.

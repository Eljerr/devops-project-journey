---
lab-id: LAB-1.3-01
date: 2026-11-20
module-phase: "Fase 1.3"
environment: "INNER CachyOS"
author: Asylon
hint-level-max: "0"
tags: [lab, evidence, fase-1, systemd, cgroups-v2, journalctl, ssh-keygen, ed25519, sshd-config]
---

# Lab: Siklus Hidup Unit `systemd` (`start-limit-hit`), Integrasi `cgroups v2`, Izin Kunci SSH Ed25519, & Urutan Drop-in `sshd_config` (`sshd -T`)

> Berkas ini adalah **bukti mentah** dari sesi terminal asli di CachyOS Linux (shell `fish`, `systemd 262`, `OpenSSH_10.6p2`). Dokumen `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md` mengutip langsung dari berkas ini.

## 1. Pertanyaan yang Ingin Dibuktikan
<!-- OPERATOR VOICE -->
1. Apa yang terjadi jika file `.service` diubah tanpa `systemctl daemon-reload`, dan bagaimana `systemd` menghentikan service yang crash berulang kali (`start-limit-hit`)?
2. Bagaimana `MemoryMax=50M`, `MemorySwapMax=0`, dan `CPUQuota=20%` di dalam file `.service` diterjemahkan ke `/sys/fs/cgroup/`?
3. Mengapa OpenSSH menolak Private Key yang izinnya terlalu longgar (`chmod 644`)?
4. Mengapa baris `PasswordAuthentication no` di bagian bawah `sshd_config` dikalahkan oleh file drop-in `50-cloud.conf` (`Include` di baris atas), dan bagaimana cara memenangkannya dengan `00-hardening.conf` serta memverifikasinya lewat `sshd -t` dan `sshd -T`?

## 2. Kondisi Awal (Versi & Lingkungan)
```text
❯ systemctl cat sshd
# /usr/lib/systemd/system/sshd.service
[Unit]
After=network.target sshdgenkeys.service
Before=ssh-access.target
Description=OpenSSH Daemon
Documentation=man:sshd(8) man:sshd_config(5)
Wants=sshdgenkeys.service ssh-access.target

[Service]
Type=notify-reload
ExecStart=/usr/bin/sshd -D
KillMode=process
Restart=always

[Install]
WantedBy=multi-user.target

❯ systemctl show -p DefaultStartLimitBurst -p DefaultStartLimitIntervalUSec
DefaultStartLimitIntervalUSec=10s
DefaultStartLimitBurst=5
```

## 3. Variabel
* **Eksperimen 1.3a (`systemd` lifecycle)**: Mengubah `ExecStart=/usr/bin/sleep 30` menjadi `/usr/bin/false` dengan `Restart=always`, sebelum vs sesudah `daemon-reload`.
* **Eksperimen 1.3b (`systemd` → `cgroups v2`)**: Menambahkan `MemoryMax=50M`, `MemorySwapMax=0`, dan `CPUQuota=20%` ke dalam `[Service]` pada `latihan.service`.
* **Eksperimen 1.3c (Izin Kunci SSH)**: Mengubah izin `~/lab/kunci_uji` dari `600` ke `644` lalu kembali ke `600`.
* **Eksperimen 1.3d & Chaos 1.3 (`sshd_config` First-Match Wins)**: Menguji `conf.d/50-cloud.conf` (`PasswordAuthentication yes`) melawan `PasswordAuthentication no` di bawah `Include`, lalu menambahkan `conf.d/00-hardening.conf` (`PasswordAuthentication no`).

## 4. Prediksi (ditulis SEBELUM menjalankan)
<!-- OPERATOR VOICE: jangan diedit setelah melihat hasil -->
* **Prediksi 1.3a (edit tanpa `daemon-reload` & `/usr/bin/false`)**: "Tidak terjadi apa-apa; dia akan benar-benar mati dan `systemd` akan menandai `start-limit-hit`."
* **Prediksi 1.3b (nilai file cgroup `latihan.service`)**: "`52428800`; `0`; `20000 100000`."
* **Prediksi 1.3c (`chmod 644` pada Private Key)**: "Ditolak karena izinnya terlalu longgar."
* **Prediksi 1.3d (`50-cloud.conf` di atas vs `PasswordAuthentication no` di bawah)**: "`yes`."

## 5. Perintah & Output Asli

### A. Misi 1.3a — Siklus Hidup `latihan.service`, Peringatan `daemon-reload`, & `start-limit-hit`
```text
 ~/Projects/devops-project-journey  main*
❯ printf "[Unit]\nDescription=Service Latihan DevOps\n\n[Service]\nExecStart=/usr/bin/false\nRestart=always\n" > ~/.config/systemd/user/latihan.service
  systemctl --user status latihan.service
  systemctl --user daemon-reload; systemctl --user restart latihan.service; sleep 2; systemctl --user status latihan.service
Warning: The unit file, source configuration file or drop-ins of latihan.service changed on disk. Run 'systemctl --user daemon-reload' to reload units.
● latihan.service - Service Latihan DevOps
     Loaded: loaded (/home/joy/.config/systemd/user/latihan.service; static)
     Active: active (running) since Thu 2026-11-19 15:12:20 WIB; 21s ago
   Main PID: 203328 (sleep)
     CGroup: /user.slice/user-1000.slice/user@1000.service/app.slice/latihan.service
             └─203328 /usr/bin/sleep 30

× latihan.service - Service Latihan DevOps
     Loaded: loaded (/home/joy/.config/systemd/user/latihan.service; static)
     Active: failed (Result: start-limit-hit) since Thu 2026-11-19 15:13:06 WIB
    Process: 205005 ExecStart=/usr/bin/false (code=exited, status=1/FAILURE)
   Main PID: 205005 (code=exited, status=1/FAILURE)

Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Scheduled restart job, restart counter is at 5.
Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Start request repeated too quickly.
Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Failed with result 'start-limit-hit'.
Nov 19 15:13:06 cachyos systemd[1025]: Failed to start Service Latihan DevOps.
```

### B. Misi 1.3b — Rekaman `journalctl` & Sinkronisasi Unit File ke `/sys/fs/cgroup`
```text
 ~/Projects/devops-project-journey  main*
❯ journalctl --user -u latihan.service -n 20 --no-pager
Nov 19 15:13:05 cachyos systemd[1025]: latihan.service: Scheduled restart job, restart counter is at 1.
Nov 19 15:13:05 cachyos systemd[1025]: Started Service Latihan DevOps.
Nov 19 15:13:05 cachyos systemd[1025]: latihan.service: Main process exited, code=exited, status=1/FAILURE
Nov 19 15:13:05 cachyos systemd[1025]: latihan.service: Failed with result 'exit-code'.
...
Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Scheduled restart job, restart counter is at 5.
Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Start request repeated too quickly.
Nov 19 15:13:06 cachyos systemd[1025]: latihan.service: Failed with result 'start-limit-hit'.

❯ systemctl --user status latihan.service
● latihan.service - Service Latihan DevOps
     Loaded: loaded (/home/joy/.config/systemd/user/latihan.service; static)
     Active: active (running) since Thu 2026-11-19 15:35:26 WIB; 11s ago
   Main PID: 253395 (sleep)
     Memory: 196K (max: 50M, swap max: 0B, available: 49.8M, peak: 2M)
        CPU: 9ms
     CGroup: /user.slice/user-1000.slice/user@1000.service/app.slice/latihan.service
             └─253395 /usr/bin/sleep 120

❯ cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/latihan.service/memory.max
  cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/latihan.service/memory.swap.max
  cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/latihan.service/cpu.max
52428800
0
20000 100000
```

### C. Misi 1.3c — Audit Izin `~/.ssh` & Uji `UNPROTECTED PRIVATE KEY FILE`
```text
 ~/Projects/devops-project-journey  main*
❯ stat -c '%a %A %n' ~/.ssh ~/.ssh/*
700 drwx------ /home/joy/.ssh
600 -rw------- /home/joy/.ssh/id_ed25519
644 -rw-r--r-- /home/joy/.ssh/id_ed25519.pub
600 -rw------- /home/joy/.ssh/known_hosts
600 -rw------- /home/joy/.ssh/known_hosts.old

❯ ssh-keygen -t ed25519 -f ~/lab/kunci_uji -N ""
  chmod 644 ~/lab/kunci_uji; ssh-keygen -y -f ~/lab/kunci_uji
  chmod 600 ~/lab/kunci_uji; ssh-keygen -y -f ~/lab/kunci_uji
...
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@         WARNING: UNPROTECTED PRIVATE KEY FILE!          @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
Permissions 0644 for '/home/joy/lab/kunci_uji' are too open.
It is required that your private key files are NOT accessible by others.
This private key will be ignored.
Load key "/home/joy/lab/kunci_uji": bad permissions
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGMzoJC0KHxjn20JBw0L+1Qwc9xZZSA+BQdX6rEtCN+Q joy@cachyos
```

### D. Misi 1.3d & Tantangan Chaos 1.3 — Aturan *First-Match Wins* `sshd_config` & `sshd -T`
```text
❯ grep -n -v -E '^\s*#|^\s*$' /etc/ssh/sshd_config
2:Include /etc/ssh/sshd_config.d/*.conf
43:AuthorizedKeysFile   .ssh/authorized_keys
113:Subsystem   sftp    /usr/lib/ssh/sftp-server

# Sebelum 00-hardening.conf (50-cloud.conf berisi PasswordAuthentication yes menang atas baris bawah):
❯ /usr/bin/sshd -T -f ~/lab/ssh_lab/sshd_uji.conf -h ~/lab/ssh_lab/host_key | grep -iE 'passwordauthentication|permitrootlogin'
PermitRootLogin no
PasswordAuthentication yes

# Tantangan Chaos 1.3 — Menambahkan 00-hardening.conf untuk mengalahkan 50-cloud.conf:
❯ echo "PasswordAuthentication no" > ~/lab/ssh_lab/conf.d/00-hardening.conf
  /usr/bin/sshd -t -f ~/lab/ssh_lab/sshd_uji.conf -h ~/lab/ssh_lab/host_key; echo "Syntax exit: $status"
  /usr/bin/sshd -T -f ~/lab/ssh_lab/sshd_uji.conf -h ~/lab/ssh_lab/host_key | grep -iE 'passwordauthentication|permitrootlogin'
Syntax exit: 0
PermitRootLogin no
PasswordAuthentication no
```

## 6. Prediksi vs Kenyataan
| Eksperimen | Prediksi | Kenyataan | Cocok? |
| :--- | :--- | :--- | :--- |
| Edit `.service` tanpa `daemon-reload` lalu `/usr/bin/false` | Tidak terjadi apa-apa sampai reload; lalu mati `start-limit-hit` | Muncul `Warning:` sampai di-reload; setelah reload gagal 5× lalu `Result: start-limit-hit` | ✅ |
| Nilai `memory.max`, `memory.swap.max`, `cpu.max` | `52428800; 0; 20000 100000` | `52428800`, `0`, `20000 100000` | ✅ |
| `chmod 644` pada Private Key | Ditolak karena izinnya terlalu longgar | Ditolak dengan `WARNING: UNPROTECTED PRIVATE KEY FILE!` | ✅ |
| `50-cloud.conf` (`yes`) di atas vs `no` di bawah | `yes` menang | `PasswordAuthentication yes` menang | ✅ |
| Chaos 1.3: Tambah `00-hardening.conf` (`no`) | `no` menang atas `50-cloud.conf` | `Syntax exit: 0`, `PasswordAuthentication no` menang | ✅ |

## 7. Kesimpulan & Fakta yang Terbukti
<!-- OPERATOR VOICE -->
* **Anatomi Unit & `daemon-reload`**: "`[Unit]` berfungsi mengatur urutan, identitas, dan ketergantungan; `[Service]` itu cara menjalankan dan mengelola proses; `[Install]` pendaftaran autostart ke target boot tertentu. Konfigurasi sistem baru dibaca setelah kita ketik `daemon-reload`, bisa dibilang ga realtime, ibarat save—kalau misalnya realtime bisa kacau juga."
* **Sinkronisasi `systemd` ke `cgroups v2`**: "Ketika kita menuliskan konfigurasi pada file `.service`, sebenarnya kita benar-benar menuliskannya di folder `/sys/fs/cgroup`, `systemd` mensinkronisasinya."
* **Keamanan Kunci SSH**: "Yang private key wajib dijaga, untuk `.pub` itu ditaruh pada server tujuan. `~/.ssh` `700` biar bisa buat key dan delete, private `600`, pub `644`, karena kalau izinnya terlalu longgar dan bisa di-read, nanti banyak yang nyelonong masuk."
* **Urutan `sshd_config` & `sshd -t` vs `sshd -T`**: "Karena pada config `sshd` urutan bacanya dari atas ke bawah, sehingga konfigurasi yang terdapat paling atas diutamakan; kita bisa menamai file konfigurasi drop-in dengan menggunakan angka kecil misal `00`, `01`, `02`. `-t` itu memverifikasi jika ada kesalahan penulisan/sintaks dan `-T` besar untuk melihat nilai mana yang benar-benar menang."

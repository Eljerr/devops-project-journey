# STATE.md — Memori Status Pembelajaran & Kompas Misi Kecil

> Agen: baca berkas ini PERTAMA. Bagian "1. Posisi Saat Ini" dan "2. Peta Misi Kecil" adalah kompas utamamu.
> Update berkas ini setiap kali satu Misi Kecil selesai (`☑`) atau saat operator mengetik `!selesai`.

---

## 1. Posisi Saat Ini

| Kunci | Nilai |
| :--- | :--- |
| **Fase aktif** | Fase 2 — Jaringan, Keamanan Host, & Reverse Proxy |
| **Sub-modul aktif** | 2.1 Packet Flow, Siklus Hidup TCP, Socket, & DNS |
| **Misi Kecil aktif** | **🎯 Tantangan Chaos 2.1 & Penyusunan Otomatis Dokumen (`01-tcp-lifecycle...`)** |
| **Langkah aktif** | **Tantangan Chaos (Level 3) — Menangkap Paket TCP (`SYN`, `RST`, `FIN`) secara Langsung** |
| **Mode Belajar** | `MULAI DARI NOL (Nol-Asumsi · Bedah Sintaks · Misi Kecil)` |
| **Pembaruan terakhir** | 2026-11-20 |

### Instruksi Spesifik untuk Agen di Sesi Berikutnya
1. **Patuhi `AGENTS.md` (v3)**: Ngobrol santai seperti mentor manusia, **1 baris Kompas Posisi di atas**, dan **1 baris `📚 Referensi` di paling bawah** (jangan taruh tag `[DOK]`/`[PACK]` di tengah kalimat, dan jangan tampilkan tabel rubrik R1–R3).
2. **Nol-Asumsi**: Bedah setiap perintah dan flag baru sebelum dijalankan.
3. **Jalankan 🎯 Tantangan Chaos 2.1 (Level 3)**:
   * Gabungkan ilmu `ip`, `ss`, dan TCP handshake/RST: tangkap paket TCP asli di interface `lo` menggunakan `tcpdump` (atau jika tidak ada, `ss`), membuktikan `SYN` → `RST` saat port `5999` tutup vs `SYN` → `SYN-ACK` → `ACK` → `FIN` saat port buka.
   * Setelah selesai dan disimpulkan operator, **agen otomatis menyusun `docs/labs/lab-2.1-01-packet-flow-tcp-dns.md` dan memperbarui `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md`** (audit A10–A11).

---

## 2. Peta Misi Kecil Fase 1 & Fase 2 (Kompas Anti-Nyasar Agen)

> Agen wajib menjalankan misi kecil di bawah ini **satu per satu secara berurutan**.
> Setiap misi melewati: **Langkah 1 (Konsep & Bedah Alat)** → **Langkah 2 (Coba Level 1 & Level 2 di Terminal + Bedah Output)** → **Langkah 3 (Klik & Centang `☑`)**.

### Sub-modul 1.1: Namespaces, cgroups v2, Sinyal, & OOM Killer (`✅ SELESAI`)
- [x] **Misi 1.1a — Proses (`/proc`) & Linux Namespaces**
- [x] **Misi 1.1b — Pagar Resource cgroups v2 (`/sys/fs/cgroup`)**
- [x] **Misi 1.1c — Sinyal OS (`SIGTERM` vs `SIGKILL`) & Kutukan PID 1**
- [x] **Misi 1.1d — OOM Killer, `memory.events`, & Peran Swap/ZRAM**
- [x] **🎯 Tantangan Chaos 1.1 & Penyusunan Otomatis Dokumen**

### Sub-modul 1.2: File Descriptor, Inode, & Misteri Disk Penuh (`✅ SELESAI`)
- [x] **Misi 1.2a — File Descriptor (`/proc/<PID>/fd`) & Offset (`fdinfo`)**
- [x] **Misi 1.2b — Nama File (Dentry) vs Koper Data (Inode) & `unlink()`**
- [x] **Misi 1.2c — Anomali `df` vs `du` (Deleted Open Files) & `lsof +L1`**
- [x] **Misi 1.2d — Inode Exhaustion (`df -i`) vs Disk Block Penuh (`df -h`)**
- [x] **🎯 Tantangan Chaos 1.2 & Penyusunan Otomatis Dokumen (`pm-01` Bagian 2)**

### Sub-modul 1.3: Manajemen Layanan `systemd` & Penguatan SSH (`✅ SELESAI`)
- [x] **Misi 1.3a — Anatomi Unit File `systemd` & Siklus Hidup Service**
- [x] **Misi 1.3b — Integrasi `systemd` ke cgroups v2 (`MemoryMax`, `CPUQuota`) & `journalctl`**
- [x] **Misi 1.3c — Kunci Kriptografi SSH Ed25519 & Izin File (`chmod 700/600/644`)**
- [x] **Misi 1.3d — Penguatan `sshd_config`, Jebakan Urutan Drop-in, & `sshd -T`**
- [x] **🎯 Tantangan Chaos 1.3 & Penyusunan Otomatis Runbook (`rb-01`)**

### Fase 2 — Sub-modul 2.1: Packet Flow, TCP, Socket, & DNS
- [x] **Misi 2.1a — Identitas Jaringan (`ip -br addr`) & Tabel Rute Kernel (`ip route`)**:
  * *Level 1 (Observasi)*: Memeriksa antarmuka `lo` (`127.0.0.1/8`), `enp8s0` (`192.168.18.90/24`), `docker0` (`172.17.0.1/16`) dan tabel rute `default via 192.168.18.1`.
  * *Level 2 (Eksperimen)*: Membuktikan keputusan routing kernel dengan `ip route get` ke `127.0.0.1` (`dev lo`), `192.168.18.50` (`dev enp8s0` langsung), dan `1.1.1.1` (`via 192.168.18.1`).
- [x] **Misi 2.1b — Siklus Hidup Koneksi TCP (3-Way Handshake, `SYN`, `ACK`, `FIN`, `RST`) & `ss -tlpn`**:
  * *Level 1 (Observasi)*: Mengintip proses pemilik socket `LISTEN` dengan `sudo ss -tlpn` (`pid=...`, `fd=...`, `127.0.0.1` vs `0.0.0.0`).
  * *Level 2 (Eksperimen)*: Menguji `curl -v http://127.0.0.1:5999` saat loket tutup (`Connection refused` dalam `0 ms` akibat paket `RST`) vs setelah dibuka dengan `python3 -m http.server 5999` (`Established connection` dari port sementara `43448` → `HTTP/1.0 200 OK`).
- [x] **Misi 2.1c — Port Sementara (*Ephemeral Ports*), Status `TIME_WAIT` vs `CLOSE_WAIT`**:
  * *Level 1 (Observasi)*: Membaca `/proc/sys/net/ipv4/ip_local_port_range` (`32768 60999`) dan melihat timer 60 detik `TIME_WAIT` (`timer:(timewait,59sec,0)`).
  * *Level 2 (Eksperimen)*: Mensimulasikan bug `CLOSE_WAIT` dengan proses Python yang menahan socket `fd=11` tanpa `close()` setelah `curl` mengirim `FIN` (`CLOSE-WAIT 79 0 127.0.0.1:5999 127.0.0.1:54064`).
- [x] **Misi 2.1d — Resolusi Nama Domain (`/etc/resolv.conf` & `dig`)**:
  * *Level 1 (Observasi)*: Memeriksa `/etc/hosts` (berisi `# BEGIN_STUDY_DISCIPLINE_BLOCK`) dan `/etc/resolv.conf` (`nameserver 127.0.0.53` milik `systemd-resolved`).
  * *Level 2 (Eksperimen)*: Membuktikan `getent ahosts reddit.com` (`127.0.0.1`/`::1`) vs `dig github.com` (`20.205.243.166`, `Query time: 0 msec` dari cache RAM `127.0.0.53`), dan melompati `systemd-resolved` dengan `dig @1.1.1.1 +noall +answer reddit.com`.
- [ ] **🎯 Tantangan Chaos 2.1 & Penyusunan Otomatis Dokumen (`01-tcp-lifecycle...`)**

---

## 3. Fakta Lingkungan (Terverifikasi)

> Agen: gunakan HANYA fakta ini. Jika butuh fakta lain, minta operator menjalankan perintah deteksi lalu tambahkan di sini.

### Inner Loop — Workstation (dideteksi agen 2026-11-12)
| Item | Nilai | Perintah deteksi |
| :--- | :--- | :--- |
| OS | CachyOS Linux (Arch-based), **native, bukan WSL2** | `cat /etc/os-release` |
| Kernel | 7.3.3-1-cachyos | `uname -r` |
| Shell | **fish 4.9.3** (bukan bash — kode keluar pakai `$status`, bukan `$?`) | `echo $SHELL; fish --version` |
| Filesystem `/` | ext4 | `findmnt -no FSTYPE /` |
| Filesystem `/tmp` | **tmpfs** (pakai RAM, jangan dipakai untuk lab disk penuh ext4) | `findmnt -no FSTYPE /tmp` |
| cgroup | v2 (`cgroup2fs`); controller user: `cpuset cpu io memory pids` | `stat -fc %T /sys/fs/cgroup` |
| Swap | zram 15.5G (zstd), `vm.swappiness=150` | `swapon --show; cat /proc/sys/vm/swappiness` |
| systemd | 262 | `systemctl --version` |
| Python | 3.12.16 | `python3 --version` |
| Firewall | `iptables v1.8.13 (nf_tables)`, `nftables 1.1.7`, `ufw 0.36.2` (enabled) | `iptables -V; nft --version; systemctl is-enabled ufw` |
| Docker | Engine 29.9.3 (enabled) | `docker version` |
| Port sementara | 32768–60999 | `cat /proc/sys/net/ipv4/ip_local_port_range` |

### Outer Loop — Proxmox & STB (Belum dideteksi ulang)
Akan dideteksi bersama saat masuk Misi 1.3c/1.3d (SSH) dan Fase 2.

---

## 4. Poin Audit Lama yang Akan Dibuktikan Sambil Jalan

> Agen: jangan dijadikan beban ujian di depan. Buktikan secara alami saat masuk Misi Kecil yang bersangkutan.

| Kode | Topik terkait | Hal yang akan dibuktikan bareng di terminal | Status |
| :--- | :--- | :--- | :--- |
| A3 | Misi 1.1c & 1.1d | Kode keluar `137` berasal dari sinyal `SIGKILL` (9) dari mana pun; bukti spesifik OOM Killer ada di `oom_kill` pada `memory.events`. | ☑ |
| A7 | Chaos 1.1 | Apakah `bytearray` lolos OOM karena kompresi zstd atau karena swap (`MemorySwapMax`) belum dikunci ke 0. | ☑ |
| A5 & A6 | Misi 1.2c | Efek `truncate` terhadap posisi tulis (`pos` di `fdinfo`) dan perbedaan ukuran tampak vs blok disk pada *sparse file*. | ☑ |
| A2 | Misi 1.2d | Perbedaan nyata antara *deleted open files* (`df -h` penuh) dan *inode exhaustion* (`df -i` penuh). | ☑ |
| A8 & A9 | Misi 1.3d | Aturan "nilai pertama menang" pada `sshd_config` (verifikasi pakai `sudo sshd -T`) dan wajibnya langkah backup sebelum ubah konfigurasi. | ☑ |

---

## 5. Kamus Analogi Operator

> Analogi yang sudah akrab bagi operator. Agen boleh memakainya kembali sambil memperjelas batas analoginya.

* **Linux Namespaces**: *"Kacamata kuda / sekat ruangan kernel"* — proses tetap jalan di mesin host yang sama, tapi pandangannya dibatasi hanya ke ruangannya sendiri.
* **cgroups v2**: *"Pagar pembatas / jatah ransum fisik"* kuota CPU & RAM di tingkat kernel.
* **Dentry vs Inode & File Descriptor**: *"Stiker nama di pintu"* vs *"koper data di loker"*. `rm` hanya mencopot stiker nama; koper belum dibuang selama masih ada tangan proses (FD) yang memegangnya.

---

## 6. Status Artefak `docs/`

| Berkas | Status |
| :--- | :--- |
| `docs/decisions/0001-hybrid-hardware-architecture.md` | ✅ TERVERIFIKASI (Telah dirapikan ke templat MADR v2 oleh agen). |
| `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` | ✅ TERVERIFIKASI (`[LAB ✓]` Bagian 1 & Bagian 2 telah disusun ulang otomatis oleh agen). |
| `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md` | ✅ TERVERIFIKASI (`[LAB ✓]` Telah disusun ulang otomatis oleh agen setelah Sub-modul 1.3 tuntas). |
| `docs/concepts/02-linux-namespaces.md` s/d `05-file-descriptors-and-inodes.md` | ✅ TERVERIFIKASI (`[LAB ✓]` 4 catatan konsep Fase 1 lengkap). |
| `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md` | ⚠️ PERLU REVISI — akan diperbarui saat masuk Fase 2. |

---

## 7. Riwayat Sesi

| Tanggal | Misi / Tahap | Catatan |
| :--- | :--- | :--- |
| 2026-11-12 | — | Inisialisasi struktur v2 (anti-halusinasi & knowledge pack). |
| 2026-11-16 | Evaluasi Workflow (`/grill-me`) | Merombak sistem ke **v3**: hapus kekakuan tabel ujian & tag inline, terapkan **Bedah Anatomi Sintaks (Anti-Copas Buta)**, **Alur Misi Kecil 3 Langkah (Mulai dari Nol)**, dan **Penyusunan Dokumen Otomatis oleh Agen**. Posisi di-reset bersih ke **Misi 1.1a Langkah 1**. |

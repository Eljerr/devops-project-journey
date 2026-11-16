# STATE.md — Memori Status Pembelajaran & Kompas Misi Kecil

> Agen: baca berkas ini PERTAMA. Bagian "1. Posisi Saat Ini" dan "2. Peta Misi Kecil" adalah kompas utamamu.
> Update berkas ini setiap kali satu Misi Kecil selesai (`☑`) atau saat operator mengetik `!selesai`.

---

## 1. Posisi Saat Ini

| Kunci | Nilai |
| :--- | :--- |
| **Fase aktif** | Fase 1 — Fondasi Kernel Linux, Proses, Filesystem, & systemd |
| **Sub-modul aktif** | 1.1 Namespaces, cgroups v2, Sinyal OS, & OOM Killer |
| **Misi Kecil aktif** | **Misi 1.1a: Kenalan sama Proses (`/proc`) & Linux Namespaces** |
| **Langkah aktif** | **Langkah 1/3 — Konsep Dasar & Bedah Alat (Mulai dari Nol)** |
| **Mode Belajar** | `MULAI DARI NOL (Nol-Asumsi · Bedah Sintaks · Misi Kecil)` |
| **Pembaruan terakhir** | 2026-11-16 |

### Instruksi Spesifik untuk Agen di Sesi Berikutnya
1. **Patuhi `AGENTS.md` (v3)**: Ngobrol santai seperti mentor manusia, **1 baris Kompas Posisi di atas**, dan **1 baris `📚 Referensi` di paling bawah** (jangan taruh tag `[DOK]`/`[PACK]` di tengah kalimat, dan jangan tampilkan tabel rubrik R1–R3).
2. **Nol-Asumsi**: Anggap operator benar-benar mulai dari nol. Jangan mengasumsikan operator sudah tahu apa itu `/proc`, `PID`, `symlink`, atau `unshare`.
3. **Mulai di Misi 1.1a — Langkah 1/3**:
   * Jelaskan konsep dasar: apa itu proses di Linux, di mana kernel menyimpan "KTP" setiap proses (`/proc`), dan bagaimana konsep **Namespaces** ("kacamata kuda / sekat ruangan") bekerja.
   * **Bedah anatomi perintah Level 1** kata demi kata: cara mengintip daftar namespace proses shell saat ini (`ls -l /proc/self/ns` — jelaskan arti `ls`, flag `-l`, apa itu `/proc/self`, dan apa itu folder `ns`).
   * Ajak operator menjalankan perintah pengamatan dasar tersebut di terminal dan menempelkan hasilnya. **Berhenti & tunggu respons operator.**

---

## 2. Peta Misi Kecil Fase 1 (Kompas Anti-Nyasar Agen)

> Agen wajib menjalankan misi kecil di bawah ini **satu per satu secara berurutan**.
> Setiap misi melewati: **Langkah 1 (Konsep & Bedah Alat)** → **Langkah 2 (Coba Level 1 & Level 2 di Terminal + Bedah Output)** → **Langkah 3 (Klik & Centang `☑`)**.

### Sub-modul 1.1: Namespaces, cgroups v2, Sinyal, & OOM Killer
- [ ] **Misi 1.1a — Proses (`/proc`) & Linux Namespaces**:
  * *Level 1 (Observasi)*: Mengintip `/proc/self/ns` dengan `ls -l` & membaca nomor inode ruangan dengan `readlink`.
  * *Level 2 (Eksperimen)*: Membuat ruangan PID baru dengan `unshare` (bedah flag `--pid`, `--fork`, `--mount-proc`) dan membandingkan `ps` di dalam vs di luar ruangan.
- [ ] **Misi 1.1b — Pagar Resource cgroups v2 (`/sys/fs/cgroup`)**:
  * *Level 1 (Observasi)*: Mengecek tipe filesystem cgroup v2 (`stat -fc %T /sys/fs/cgroup`) dan melihat isi folder cgroup milik user saat ini (`memory.current`, `memory.max`, `cpu.max`).
  * *Level 2 (Eksperimen)*: Menjalankan proses di dalam pagar memori sementara dengan `systemd-run --user --scope` (bedah tiap flag-nya!).
- [ ] **Misi 1.1c — Sinyal OS (`SIGTERM` vs `SIGKILL`) & Kutukan PID 1**:
  * *Level 1 (Observasi)*: Melihat daftar sinyal (`kill -l`) dan cara baca kode keluar di shell **fish** (`$status`, aturan `128 + N`).
  * *Level 2 (Eksperimen)*: Mengirim sinyal sopan `SIGTERM` (15) vs paksa `SIGKILL` (9), lalu menguji mengapa PID 1 di dalam namespace `unshare` mengabaikan `SIGTERM`.
- [ ] **Misi 1.1d — OOM Killer, `memory.events`, & Peran Swap/ZRAM**:
  * *Level 1 (Observasi)*: Melihat konfigurasi swap ZRAM di CachyOS (`swapon --show`) dan membaca counter `memory.events`.
  * *Level 2 (Eksperimen)*: Memicu OOM Killer dengan sengaja (`MemoryMax=50M` + `MemorySwapMax=0`) dan membedakan bukti `137` biasa vs bukti asli `oom_kill`.
- [ ] **🎯 Tantangan Chaos 1.1 & Penyusunan Otomatis Dokumen**:
  * *Chaos (Level 3)*: Eksperimen terukur membandingkan pengaruh `MemorySwapMax` vs jenis data (`bytearray` nol vs `os.urandom` acak) satu variabel per langkah.
  * *Dokumentasi Otomatis*: Agen menyusun Catatan Konsep (`docs/concepts/`) & Bagian 1 Post-Mortem (`pm-01`) dari obrolan chat.

### Sub-modul 1.2: File Descriptor, Inode, & Misteri Disk Penuh
- [ ] **Misi 1.2a — File Descriptor (`/proc/<PID>/fd`) & Offset (`fdinfo`)**
- [ ] **Misi 1.2b — Nama File (Dentry) vs Koper Data (Inode) & `unlink()`**
- [ ] **Misi 1.2c — Anomali `df` vs `du` (Deleted Open Files) & `lsof +L1`**
- [ ] **Misi 1.2d — Inode Exhaustion (`df -i`) vs Disk Block Penuh (`df -h`)**
- [ ] **🎯 Tantangan Chaos 1.2 & Penyusunan Otomatis Dokumen (`pm-01` Bagian 2)**

### Sub-modul 1.3: Manajemen Layanan `systemd` & Penguatan SSH
- [ ] **Misi 1.3a — Anatomi Unit File `systemd` & Siklus Hidup Service**
- [ ] **Misi 1.3b — Integrasi `systemd` ke cgroups v2 (`MemoryMax`, `CPUQuota`) & `journalctl`**
- [ ] **Misi 1.3c — Kunci Kriptografi SSH Ed25519 & Izin File (`chmod 700/600/644`)**
- [ ] **Misi 1.3d — Penguatan `sshd_config`, JebakanUrutan Drop-in, & `sshd -T`**
- [ ] **🎯 Tantangan Chaos 1.3 & Penyusunan Otomatis Runbook (`rb-01`)**

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
| A3 | Misi 1.1c & 1.1d | Kode keluar `137` berasal dari sinyal `SIGKILL` (9) dari mana pun; bukti spesifik OOM Killer ada di `oom_kill` pada `memory.events`. | ☐ |
| A7 | Chaos 1.1 | Apakah `bytearray` lolos OOM karena kompresi zstd atau karena swap (`MemorySwapMax`) belum dikunci ke 0. | ☐ |
| A5 & A6 | Misi 1.2c | Efek `truncate` terhadap posisi tulis (`pos` di `fdinfo`) dan perbedaan ukuran tampak vs blok disk pada *sparse file*. | ☐ |
| A2 | Misi 1.2d | Perbedaan nyata antara *deleted open files* (`df -h` penuh) dan *inode exhaustion* (`df -i` penuh). | ☐ |
| A8 & A9 | Misi 1.3d | Aturan "nilai pertama menang" pada `sshd_config` (verifikasi pakai `sudo sshd -T`) dan wajibnya langkah backup sebelum ubah konfigurasi. | ☐ |

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
| `docs/decisions/0001-hybrid-hardware-architecture.md` | ⚠️ PERLU REVISI — akan dirapikan ke templat MADR oleh agen. |
| `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` | ⚠️ PERLU REVISI — akan disusun ulang otomatis oleh agen setelah Sub-modul 1.1 & 1.2 tuntas. |
| `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md` | ⚠️ PERLU REVISI — akan disusun ulang otomatis oleh agen setelah Sub-modul 1.3 tuntas. |
| `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md` | ⚠️ PERLU REVISI — akan diperbarui saat masuk Fase 2. |

---

## 7. Riwayat Sesi

| Tanggal | Misi / Tahap | Catatan |
| :--- | :--- | :--- |
| 2026-11-12 | — | Inisialisasi struktur v2 (anti-halusinasi & knowledge pack). |
| 2026-11-16 | Evaluasi Workflow (`/grill-me`) | Merombak sistem ke **v3**: hapus kekakuan tabel ujian & tag inline, terapkan **Bedah Anatomi Sintaks (Anti-Copas Buta)**, **Alur Misi Kecil 3 Langkah (Mulai dari Nol)**, dan **Penyusunan Dokumen Otomatis oleh Agen**. Posisi di-reset bersih ke **Misi 1.1a Langkah 1**. |

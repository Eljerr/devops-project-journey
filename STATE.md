# STATE.md — Memori Status Pembelajaran & Konteks Sesi Aktif

## Status Eksekusi Saat Ini

* **Fase Aktif**: Fase 2 — Jaringan Komputer, Keamanan Host, dan Arsitektur Reverse Proxy
* **Sub-Modul Aktif**: **2.2: Linux Netfilter, Firewalling (UFW/iptables), & Gateway Rebuild** (Target: `archive/v1/core-infrastructure/gateway-stb/`)
* **Status Langkah Siklus**: Selesai Sub-Modul 2.1 (Konsep Note Terbit) — Siap Masuk Tahap 2 Teori Mikro Sub-Modul 2.2
* **Pembaruan Terakhir**: 2026-11-12

---

## Topologi & Status Infrastruktur Homelab

1. **Workstation / Sandbox — Inner Loop (PC WSL2 / Host Linux)**:
   * *Spesifikasi*: AMD Ryzen 5 2600 (6C/12T), RAM 16 GB (`x86_64`), CachyOS Linux (Kernel 7.3.3-cachyos, ZRAM zstd enabled).
   * *Peran*: Lingkungan eksperimen cepat (isolasi kernel, cgroups, simulasi proses, build Docker lokal, verifikasi socket/firewall, dan penulisan kode Go/Python).
2. **Compute & Staging Cluster — Outer Loop (Proxmox VE Host)**:
   * *Spesifikasi*: Laptop Rekondisi (CPU 2 Cores, RAM 8 GB, `x86_64`).
   * *Peran*: Menjalankan kontainer LXC *unprivileged* untuk klaster K3s (1 Master, 2 Worker) dan node Monitoring. Status: Menyala, menunggu audit ulang di Fase 5.
3. **Edge Ingress Gateway — Outer Loop (STB Amlogic S905X4)**:
   * *Spesifikasi*: `ARM64`, RAM 2 GB, OS Ubuntu.
   * *Peran*: Menjalankan Nginx Reverse Proxy, Cloudflared Tunnel, Node Exporter, dan cAdvisor. Status: Target audit jaringan & firewall di Fase 2.

---

## Log Kompetensi Terverifikasi (Sudah Dikuasai Operator)

* **[LULUS] 1.1.1 Isolasi Proses & Namespaces**: Memahami secara mendalam konsep isolasi kernel (`pid`, `net`, `mnt`), hierarki `/proc/<PID>/ns`, serta membuktikan bahwa kontainer bukanlah mesin virtual melainkan proses host biasa yang dibatasi visibilitasnya.
* **[LULUS] 1.1.2 cgroups v2 & Manajemen Kuota Sumber Daya**: Memahami struktur hierarki direktori `/sys/fs/cgroup/`, fungsi `memory.max`, `memory.swap.max`, `cpu.max`, serta `cpuset`. Berhasil mempraktikkan pembatasan memori manual dengan memasukkan PID ke `cgroup.procs`.
* **[LULUS] 1.1.3 Sinyal OS & Forensik OOM Killer**: Memahami perbedaan fundamental `SIGTERM` (15) vs `SIGKILL` (9), alasan mengapa proses PID 1 mengabaikan `SIGTERM` tanpa signal handler (*kutukan PID 1*), serta memverifikasi langsung eksekusi OOM Killer dengan bukti kode keluar `137` ($128 + 9$) dan counter `memory.events`.
* **[LULUS] 1.2 File Descriptors, Inode, & Misteri Disk Penuh (deleted open files)**: Memahami secara mendalam pemisahan Dentry vs Inode, filosofi *Everything is a file*, struktur File Descriptor di `/proc/<PID>/fd/`, mekanisme `unlink()`, anomali selisih pelaporan `df` vs `du`, deteksi forensik dengan `lsof +L1`, serta implikasi kritis *file write offset* Open File Table vs `truncate` pada mode `O_APPEND` vs non-append.
* **[LULUS] 1.3 Manajemen Layanan systemd & Penguatan Host (SSH Hardening)**: Memahami arsitektur `systemd` sebagai PID 1, deklarasi unit service, integrasi deklaratif cgroups v2, kompresi ZRAM zstd vs penembakan OOM Killer fisik via `os.urandom`, serta pengamanan kriptografi asimetris Ed25519 (hak akses `600`/`644`).
* **[LULUS] 2.1 Packet Flow, Routing, Socket Lifecycle, & DNS**:
  * Memahami bahwa socket jaringan di Linux dipegang sebagai File Descriptor (`fd`) oleh proses aplikasi (dibuktikan dari output `ss -tulpn` `antigravity` dan `systemd-resolve`).
  * Memahami perlunya TCP 3-Way Handshake sebagai batas minimum matematis pembuktian koneksi dua arah yang andal.
  * Memverifikasi penolakan instan 0 ms kernel Linux via paket `RST` (`Connection refused`) ketika port tujuan kosong.
  * Memahami dan memverifikasi siklus penutupan TCP: kehati-hatian `TIME_WAIT` ($2 \times \text{MSL}$) untuk membunuh paket nyasar (didukung mekanisme TTL IP) vs anomali kebocoran socket `CLOSE_WAIT` akibat kode aplikasi lupa memanggil `close()`.
  * Membedakan konsumsi sumber daya: aplikasi inbound hanya memakan File Descriptor (FD), sedangkan reverse proxy outbound (Nginx) memakan File Descriptor dan Port Sementara (Ephemeral Port).
  * Memahami arsitektur *Outbound Tunnel* (`cloudflared`) yang meniadakan kebutuhan port forwarding modem, serta optimasi *connection pooling* Nginx `keepalive 32;` (analogi garasi armada mobil standby).

---

## Status Artefak Portofolio (`docs/`)

* `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md`:
  * **Status**: **LENGKAP (PUBLISHED)** — Analisis OOM Killer cgroups v2 & Deleted Open Files.
* `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`:
  * **Status**: **LENGKAP (PUBLISHED)** — SOP Triage OOM-Kill & SSH Hardening Ed25519.
* `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md`:
  * **Status**: **LENGKAP (PUBLISHED)** — Catatan atomik Obsidian mengenai siklus hidup TCP, socket states (`TIME_WAIT` vs `CLOSE_WAIT`), dan arsitektur reverse proxy L4 vs L7.

---

## Kamus Model Mental & Analogi Operator (Konteks Lintas Sesi)

* **Linux Namespaces**: *"Kacamata kuda kernel"* — proses tetap berada di host yang sama, tetapi pandangannya dibatasi.
* **cgroups v2**: *"Pagar pembatas / ransum fisik"* kuota CPU/RAM di tingkat kernel.
* **Dentry vs Inode & File Descriptor**: *"Stiker label pintu"* vs *"koper di loker stasiun"*. `rm` hanya mencopot stiker (`unlink`). Kernel menolak membuang koper selama masih ada tangan proses yang memegang FD.
* **Socket Jaringan**: *"Nomor tiket gagang telepon"* resmi dari kernel ke aplikasi untuk mengirim dan menerima paket.
* **3-Way Handshake**: Batas minimal 3 langkah pembuktian jalur dua arah antara dua pihak di jaringan tidak andal.
* **Paket RST**: Satpam kernel yang langsung menolak ketukan di pintu kamar port kosong dalam 0 ms (`Connection refused`).
* **TIME_WAIT**: Bentuk kehati-hatian pihak pertama pamit menunggu 60 detik agar paket nyasar (dibatasi TTL IP) mati sebelum port dipakai lagi.
* **CLOSE_WAIT**: Bug aplikasi lupa memanggil `close()` pada socket FD penerima; kernel menolak menutup paksa tanpa instruksi aplikasi sehingga memicu socket/FD leak.
* **Inbound vs Outbound Resource**: Aplikasi inbound hanya makan FD; Gateway/Reverse proxy outbound makan FD + Port Sementara (Ephemeral Port).
* **Nginx Keepalive**: Garasi 32 mobil standby untuk dipakai bergilir melayani request paralel tanpa 3-way handshake berulang.

---

## Kendala Teknis & Catatan Sesi Terakhir (Active Blockers)

* Sub-Modul 2.1 resmi **TUNTAS 100%**.
* Kode legacy terkait gateway tersimpan di `archive/v1/core-infrastructure/gateway-stb/`.

---

## Instruksi Eksplisit untuk Sesi Berikutnya

1. **Patuhi Aturan Turn-Gate (`AGENTS.md`)**: Satu tahap per giliran.
2. **Masuk ke Sub-Modul 2.2: Linux Netfilter, Firewalling (UFW/iptables), & Gateway Rebuild**:
   * **Tahap 2 (Elaborasi Teori Mikro)**:
     - Jelaskan rantai Netfilter di Linux kernel (`PREROUTING`, `INPUT`, `FORWARD`, `OUTPUT`, `POSTROUTING`).
     - Jelaskan **Jebakan Docker vs UFW**: Mengapa perintah `docker run -p 80:80` mem-bypass aturan UFW karena Docker memanipulasi rantai `iptables` secara langsung.
     - Ajukan 1-2 pertanyaan konseptual sokratik dan **BERHENTI (HARD STOP)** menunggu jawaban operator sebelum masuk ke praktik.

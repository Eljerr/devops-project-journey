# STATE.md — Memori Status Pembelajaran & Konteks Sesi Aktif

## Status Eksekusi Saat Ini

* **Fase Aktif**: Fase 1 — Fundamental Kernel Linux, Sistem Operasi, dan Manajemen Proses
* **Sub-Modul Aktif**: 1.1 Hirarki Sistem Berkas, Izin Akses, dan Isolasi Proses (cgroups & namespaces) -> Masuk ke Sub-modul 1.2
* **Status Langkah Siklus**: Tahap 4 (Dokumentasi Post-Mortem Selesai & Verifikasi Log Kompetensi)
* **Pembaruan Terakhir**: 2026-11-06

## Status Infrastruktur Homelab (Verifikasi Ulang)

* **Workstation / Sandbox (PC WSL2 Ubuntu)**: AMD Ryzen 5 2600 (6C/12T), RAM 16 GB. Berperan sebagai lingkungan *Inner Loop* (pengujian cepat, build Docker lokal, kompilasi, dan sandbox eksperimen).
* **Proxmox Host (8GB RAM)**: Menyala, target lingkungan *Staging / Compute Cluster* (LXC K3s).
* **STB Gateway (2GB RAM ARM64)**: Menyala, target *Edge Ingress & Reverse Proxy*.

## Log Kompetensi Terverifikasi (Sudah Dikuasai Kembali)

* **1.1.1 Isolasi Proses & Namespaces**: Memahami secara mendalam konsep kacamata kuda kernel (`pid`, `net`, `mnt`), hierarki `/proc/<PID>/ns`, serta membuktikan bahwa kontainer bukanlah ruang virtual melainkan proses host yang diisolasi visibilitasnya.
* **1.1.2 cgroups v2 & Manajemen Kuota Sumber Daya**: Memahami struktur hierarki direktori `/sys/fs/cgroup/`, fungsi `memory.max`, `memory.swap.max`, `cpu.max`, serta `cpuset`. Berhasil mempraktikkan pembatasan memori manual dengan `cgroup.procs`.
* **1.1.3 Sinyal OS & Forensik OOM Killer**: Memahami perbedaan fundamental `SIGTERM` (15) vs `SIGKILL` (9), misteri kutukan PID 1 dan kegagalan *graceful shutdown*, serta memverifikasi langsung terjadinya eksekusi OOM Killer dengan bukti kode keluar `137` ($128 + 9$) dan counter `memory.events`.
* **Artefak Portofolio Terbit**:
  - `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` (Dokumentasi investigasi chaos OOM Killer berbasis analisis murni operator).

## Kendala Teknis & Catatan Sesi Terakhir (Active Blockers)

* Tidak ada kendala aktif. Pacing belajar diperlambat sesuai kebutuhan operator untuk memastikan pemahaman mendalam (*first-principles thinking*).

## Instruksi untuk Agen AI pada Sesi Berikutnya / Langkah Aktif

1. **Jaga Ritme & Pacing Pembelajaran (Sesuai Konstitusi `AGENTS.md`)**:
   - Jangan buru-buru atau menumpuk istilah teknis.
   - Pakai analogi visual *first-principles* sampai operator mendapatkan momen *"Aha!"*.
   - Operator yang menulis seluruh analisis dan kesimpulan dokumentasi portofolio.
2. **Mulai Sub-Modul 1.2 & 1.3**:
   - Bahas konsep **File Descriptors** di `/proc/<PID>/fd`, misteri hilangnya ruang harddisk saat file dihapus tapi masih dipegang proses (`deleted open files`), dan perbedaan mendasar antara `df -h` vs `du`.
   - Lakukan eksperimen terminal dengan perkakas `lsof` dan `strace`.
   - Latihan pembuatan unit service `systemd` mandiri dengan pembatasan kuota memori/CPU (`MemoryMax`, `CPUQuota`).
   - Pandu operator menyusun runbook: `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`.

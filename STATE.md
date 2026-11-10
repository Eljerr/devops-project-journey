# STATE.md — Memori Status Pembelajaran & Konteks Sesi Aktif

## Status Eksekusi Saat Ini

* **Fase Aktif**: Fase 1 — Fundamental Kernel Linux, Sistem Operasi, dan Manajemen Proses
* **Sub-Modul Aktif**: 1.3 Manajemen Layanan `systemd` & Penguatan Host (SSH Hardening)
* **Status Langkah Siklus**: Tahap 1 (Selesai Sub-Modul 1.2 — Siap masuk Tahap 2 Elaborasi Teori Mikro Sub-Modul 1.3)
* **Pembaruan Terakhir**: 2026-11-10

---

## Topologi & Status Infrastruktur Homelab

1. **Workstation / Sandbox — Inner Loop (PC WSL2 Ubuntu)**:
   * *Spesifikasi*: AMD Ryzen 5 2600 (6C/12T), RAM 16 GB (`x86_64`).
   * *Peran*: Lingkungan eksperimen cepat (isolasi kernel, cgroups, simulasi proses, build Docker lokal, dan penulisan kode Go/Python).
2. **Compute & Staging Cluster — Outer Loop (Proxmox VE Host)**:
   * *Spesifikasi*: Laptop Rekondisi (CPU 2 Cores, RAM 8 GB, `x86_64`).
   * *Peran*: Menjalankan kontainer LXC *unprivileged* untuk klaster K3s (1 Master, 2 Worker) dan node Monitoring. Status: Menyala, menunggu audit ulang di Fase 1.3 & Fase 5.
3. **Edge Ingress Gateway — Outer Loop (STB Amlogic S905X4)**:
   * *Spesifikasi*: `ARM64`, RAM 2 GB, OS Ubuntu.
   * *Peran*: Menjalankan Nginx Reverse Proxy, Cloudflared Tunnel, Node Exporter, dan cAdvisor. Status: Menyala, menunggu audit firewall & jaringan di Fase 1.3 & Fase 2.

---

## Log Kompetensi Terverifikasi (Sudah Dikuasai Operator)

* **[LULUS] 1.1.1 Isolasi Proses & Namespaces**: Memahami secara mendalam konsep isolasi kernel (`pid`, `net`, `mnt`), hierarki `/proc/<PID>/ns`, serta membuktikan bahwa kontainer bukanlah mesin virtual melainkan proses host biasa yang dibatasi visibilitasnya.
* **[LULUS] 1.1.2 cgroups v2 & Manajemen Kuota Sumber Daya**: Memahami struktur hierarki direktori `/sys/fs/cgroup/`, fungsi `memory.max`, `memory.swap.max`, `cpu.max`, serta `cpuset`. Berhasil mempraktikkan pembatasan memori manual dengan memasukkan PID ke `cgroup.procs`.
* **[LULUS] 1.1.3 Sinyal OS & Forensik OOM Killer**: Memahami perbedaan fundamental `SIGTERM` (15) vs `SIGKILL` (9), alasan mengapa proses PID 1 mengabaikan `SIGTERM` tanpa signal handler (*kutukan PID 1*), serta memverifikasi langsung eksekusi OOM Killer dengan bukti kode keluar `137` ($128 + 9$) dan counter `memory.events`.
* **[LULUS] 1.2 File Descriptors, Inode, & Misteri Disk Penuh (deleted open files)**: Memahami secara mendalam pemisahan Dentry vs Inode, filosofi *Everything is a file*, struktur File Descriptor di `/proc/<PID>/fd/`, mekanisme `unlink()`, anomali selisih pelaporan `df` vs `du`, deteksi forensik dengan `lsof +L1`, serta implikasi kritis *file write offset* Open File Table vs `truncate` pada mode `O_APPEND` vs non-append.

---

## Status Artefak Portofolio (`docs/`)

* `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md`:
  * **Status**: **LENGKAP (PUBLISHED)** — Menggabungkan Bagian 1 (Investigasi OOM Killer cgroups v2) dan Bagian 2 (Anomali Deleted Open Files) berbasis data forensik terminal dan analisis otentik operator.
* `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`:
  * **Status**: Belum dibuat (Target penyelesaian di akhir Sub-Modul 1.3).

---

## Kamus Model Mental & Analogi Operator (Konteks Lintas Sesi)

* **Linux Namespaces**: Dianalogikan sebagai *"kacamata kuda kernel"* — proses tetap berada di lapangan (host) yang sama, tetapi pandangannya dibatasi sehingga mengira dirinya sendirian.
* **cgroups v2**: Dianalogikan sebagai *"pagar pembatas/jatah ransum fisik"* di tingkat kernel agar proses tidak memonopoli sumber daya host.
* **Dentry vs Inode & File Descriptor**: Dianalogikan sebagai *"stiker label"* di pintu folder vs *"koper di loker resmi stasiun"*. Perintah `rm` hanya mencopot stiker (`unlink`). Kernel menolak membuang koper fisik selama masih ada tangan proses yang memegang File Descriptor (FD) ke Inode tersebut.
* **File Offset & Truncation**: Operator secara mandiri membedah bahwa `truncate` ke file descriptor aktif bisa memicu *sparse file* dengan *hole offset* jika proses menulis tanpa flag `O_APPEND`.
* **Gaya Belajar Operator**: Sangat menyukai pembedahan *first-principles* (membuktikan lewat direktori virtual `/proc` dan `/sys` secara langsung di terminal, bukan sekadar menghafal perintah).

---

## Kendala Teknis & Catatan Sesi Terakhir (Active Blockers)

* Tidak ada kendala teknis. Sub-Modul 1.2 selesai dengan pemahaman mendalam melampaui ekspektasi dasar.

---

## Instruksi Eksplisit untuk Agen AI pada Sesi Berikutnya

1. **Patuhi Aturan Turn-Gate (`AGENTS.md`)**: Jangan menggabungkan teori (Tahap 2) dan instruksi praktik (Tahap 3) dalam satu pesan.
2. **Eksekusi Sesi Berikutnya (Sub-Modul 1.3: Manajemen Layanan systemd & SSH Hardening)**:
   * **Tahap 2 (Elaborasi Teori Mikro)**:
     - Jelaskan arsitektur `systemd` sebagai PID 1, anatomi berkas unit (`[Unit]`, `[Service]`, `[Install]`).
     - Jelaskan integrasi `systemd` dengan `cgroups v2` (`MemoryMax`, `CPUQuota`).
     - Jelaskan konsep dasar SSH hardening (asymmetric cryptography Ed25519, penonaktifan password auth & root login) yang akan diterapkan di Proxmox/STB.
     - Ajukan 1-2 pertanyaan konseptual dan **BERHENTI (HARD STOP)** menunggu jawaban operator.
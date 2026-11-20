---
adr-id: ADR-0001
status: "accepted"
date: 2026-10-03
review-date: 2026-11-20
author: Asylon
hardware-constraint: "Hybrid (INNER CachyOS 16GB | OUTER Proxmox 2C/8GB | OUTER STB ARM64 2GB)"
tags: [adr, architecture, fase-1, homelab, proxmox, stb, k3s, lxc]
---

# ADR-0001: Gunakan Arsitektur Hibrida (STB ARM64 sebagai Edge Gateway & Laptop Proxmox sebagai Compute Cluster)

## 1. Konteks dan Masalah
<!-- OPERATOR VOICE: masalah nyata di homelab + batas fisiknya -->
Seluruh proyek homelab ini berjalan di atas perangkat keras pribadi yang dimanfaatkan kembali (*repurposed hardware*) tanpa penyedia cloud berbayar:
1. **Workstation Utama (Inner Loop)**: PC CachyOS Linux native (Ryzen 5 2600, 16 GB RAM, `x86_64`) untuk eksperimen kernel, pengembangan kode, dan pengujian lokal.
2. **Laptop Proxmox VE (Outer Loop — Compute)**: Laptop bekas (2 Core, 8 GB RAM, `x86_64`) yang juga masih dipakai kuliah dengan dual-boot CachyOS.
3. **STB Amlogic S905X4 (Outer Loop — Edge)**: Perangkat STB Android yang sudah di-flash Ubuntu Linux (`ARM64`, 2 GB RAM) dan sangat hemat daya.

Masalah utama yang harus dipecahkan adalah bagaimana membagi beban kerja (*reverse proxy*, *tunneling*, *monitoring*, dan klaster orkestrasi kontainer) agar tidak melampaui batas fisik RAM pada STB (2 GB) maupun Laptop Proxmox (8 GB).

## 2. Pendorong Keputusan (Decision Drivers)
* Total pemakaian RAM di STB Amlogic S905X4 (`ARM64`, 2 GB RAM) tidak boleh melebihi ~1,5 GB dan wajib hemat daya karena menyala terus sebagai pintu masuk jaringan.
* Total beban seluruh LXC/VM di Laptop Proxmox (2 Core, 8 GB RAM) harus dijaga ≤ ~6 GB agar host Proxmox tidak mengalami OOM.
* Isolasi keamanan berlapis (*blast-radius minimization*): memisahkan titik masuk trafik internet publik dari klaster komputasi utama.

## 3. Opsi yang Dipertimbangkan
* **Opsi 1**: Menjalankan seluruh workload (Gateway, Reverse Proxy, Monitoring, dan Klaster K3s) sekaligus di dalam satu mesin Laptop Proxmox.
* **Opsi 2**: Memisahkan peran secara hibrida — menjadikan **STB ARM64** khusus sebagai **Edge Gateway / Reverse Proxy** (`cloudflared` + Nginx), dan mendedikasikan **Laptop Proxmox** khusus untuk **Compute Cluster** (K3s di dalam kontainer LXC).

## 4. Analisis Trade-off
<!-- OPERATOR VOICE. Klaim teknis diberi label [LAB ✓]/[DOK]/[INGATAN] -->

### Opsi 1: Semua Workload di Laptop Proxmox
* **Kelebihan**: Hanya satu mesin server yang perlu dikelola dan semua berjalan di arsitektur CPU yang sama (`x86_64`).
* **Kekurangan / alasan ditolak**: Memakan jatah RAM 8 GB Proxmox yang sudah terbatas, tidak memanfaatkan STB yang sangat hemat daya, serta jika laptop sedang dimatikan atau dipakai dual-boot untuk kuliah, tidak ada gerbang depan yang tetap aktif.

### Opsi 2: STB sebagai Edge Gateway & Laptop Proxmox sebagai Compute Cluster
* **Kelebihan**:
  * STB saya putuskan sebagai Edge Gateway karena keterbatasan RAM (2 GB), selain itu memiliki beberapa keuntungan sebagai Proxy di mana jika terjadi peretasan hanya sampai pada STB dan meminimalisasi area dampak eksploitasi tersebut sehingga memberikan waktu untuk proses perlindungan kepada Cluster server dan memudahkan *filtering*.
  * Laptop Proxmox memiliki RAM yang lebih lega (8 GB) sehingga sangat memungkinkan untuk menjalankan beberapa LXC di mana bisa dibagi untuk node K3s dan database. Saya memutuskan menggunakan **K3s** dibandingkan K8s penuh karena K3s lebih ramah untuk *low-resource environments* dan fitur-fiturnya tidak kalah lengkap. Saya juga memutuskan memakai **LXC** karena jauh lebih hemat resource dibanding VM penuh (berbagi kernel host melalui Namespaces dan `cgroups v2` `[LAB ✓ → lab-1.1-01-namespaces-cgroups-signals-oom]`).
* **Kekurangan / tantangan**: *Maintenance* sedikit lebih menantang karena infrastruktur terpisah di dua perangkat fisik, jika laptop dipakai kuliah maka klaster compute akan mati sementara, dan terdapat perbedaan arsitektur CPU (`x86_64` di PC/Proxmox vs `arm64` di STB).

## 5. Keputusan & Justifikasi
**Opsi terpilih**: **Opsi 2 (Arsitektur Hibrida: STB sebagai Edge Gateway & Laptop Proxmox sebagai Compute Cluster)**, karena <!-- OPERATOR VOICE --> kebetulan saya mempunyai STB yang tidak terpakai dan sangat hemat daya, sementara Laptop Proxmox yang digunakan sanggup memenuhi kebutuhan klaster LXC saat ini. Keterbatasan arsitektur dan RAM juga dipertimbangkan sehingga semua server aplikasi dijalankan di Proxmox sementara STB fokus sebagai gerbang depan.

### Konsekuensi
* **Positif**: Kontainer LXC di Proxmox bisa dibuat lebih leluasa untuk kebutuhan praktik; jika klaster K3s sedang *down*, proxy di STB tetap hidup dan bisa menyampaikan status kepada internet; serta memudahkan pemisahan *monitoring* dan *debugging*.
* **Utang teknis & mitigasi**: Perbedaan arsitektur CPU (`x86_64` vs `arm64`) mengharuskan setiap biner Go atau image Docker yang akan dijalankan di STB di-*build* dengan target `linux/arm64` (Fase 3 & Fase 4), dan keterbatasan RAM wajib dijaga ketat dengan `MemoryMax` / `limits.memory` pada `cgroups v2` (`[LAB ✓ → lab-1.3-01-systemd-and-ssh-hardening]`).

## 6. Konfirmasi Implementasi
* **Batas anggaran sumber daya yang dikunci**:
  * **INNER CachyOS**: 16 GB RAM, kernel `7.3.3-1-cachyos`, `cgroup2fs`, swap `zram0` 15.5 GB (`[LAB ✓]`).
  * **OUTER Proxmox**: 2 Core, 8 GB RAM (total alokasi LXC ≤ ~6 GB).
  * **OUTER STB**: Amlogic S905X4 `arm64`, 2 GB RAM (total beban ≤ ~1,5 GB).
* **Bukti terkait**:
  * [lab-1.1-01-namespaces-cgroups-signals-oom.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md)
  * [lab-1.3-01-systemd-and-ssh-hardening.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.3-01-systemd-and-ssh-hardening.md)

## 7. Sumber
* Dokumentasi Kernel Linux `cgroups v2` & `namespaces` (`man 7 namespaces`, `docs.kernel.org/admin-guide/cgroup-v2.html`)
* Dokumentasi Resmi Proxmox LXC & K3s System Requirements (`docs.k3s.io`)

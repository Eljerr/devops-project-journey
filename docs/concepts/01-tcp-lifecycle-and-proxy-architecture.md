---
topic: "TCP Lifecycle, Socket States (TIME_WAIT vs CLOSE_WAIT), dan Arsitektur Reverse Proxy"
module-phase: "Fase 2.1"
date: 2026-11-12
related-docs: ["[[rb-01-linux-diagnostics-and-hardening]]", "[[0001-hybrid-hardware-architecture]]"]
tags: [concept, networking, tcp, nginx, first-principles, interview-prep]
---

# Konsep: Siklus TCP, Socket States, dan Arsitektur Reverse Proxy

## 1. Analogi Mental & First-Principles (Momen "Aha!" Operator)

* **Gedung vs Kamar vs Telepon**:
  * **IP Address**: Alamat gedung fisik (misal `192.168.1.15`).
  * **Port**: Nomor kamar/unit di dalam gedung (Port `22` kamar SSH, Port `80` kamar Nginx).
  * **Socket / File Descriptor**: Jalur gagang telepon antara kamar pengirim dan kamar penerima. Di Linux, socket dipegang sebagai File Descriptor (`fd`) oleh proses aplikasi.
* **Tukang Pos (Router L4) vs Resepsionis Cerdas (Reverse Proxy L7 / Nginx)**:
  * *Tukang Pos*: Hanya melempar amplop berdasarkan IP & Port tanpa membaca isi surat (cepat, tapi buta terhadap teks HTTP dan tidak bisa memfilter keamanan).
  * *Resepsionis (Nginx)*: Membuka amplop, mendekripsi HTTPS, membaca path URL, memeriksa izin/keamanan, lalu mengangkat gagang telepon baru via port sementara internal untuk meneruskan data ke Proxmox.
* **Keep-Alive Pool (Garasi Mobil Siap Pakai)**:
  * Konfigurasi `keepalive 32;` di Nginx upstream bukan berarti memakai satu port selamanya, melainkan menyediakan "garasi" berisi 32 koneksi *idle* yang siap pakai. Setiap mobil melayani request secara bergilir tanpa perlu menghancurkan koneksi dan melakukan 3-Way Handshake berulang-ulang dari nol.

---

## 2. Mekanisme Internal (Di Balik Layar Kernel & Socket)

### A. 3-Way Handshake & Penolakan Instan Kernel (`RST`)
* **3 Langkah Wajib**: Komunikasi jaringan membutuhkan minimal 3 langkah (`SYN` $\to$ `SYN-ACK` $\to$ `ACK`) agar kedua belah pihak sama-sama membuktikan bahwa jalur kirim dan terima dua arah berfungsi.
* **Paket `RST` (Reset)**: Jika koneksi dikirim ke port yang tidak memiliki proses yang mendengarkan (`LISTEN`), kernel Linux langsung mengirim paket `RST` seketika (`0 ms`), menghasilkan error `Connection refused`. Berbeda dengan Firewall yang melakukan `DROP` (menggantung hingga timeout).

### B. `TIME_WAIT` vs `CLOSE_WAIT`
* **`TIME_WAIT` (Normal - Sisi Inisiator Pamit)**:
  * Terjadi di pihak yang pertama kali mengirim `FIN` (misalnya browser klien atau Nginx).
  * Kernel menahan port selama $2 \times \text{MSL}$ (~60 detik) untuk memastikan paket `ACK` terakhir sampai dan semua paket terlambat (*ghost packets*) di internet mati.
  * *Mengapa paket tidak bisa nyasar lebih dari 60 detik?* Karena ada mekanisme **TTL (Time to Live)** pada header IP yang berkurang di setiap hop router hingga dibuang jika mencapai 0, serta nomor urut acak 32-bit TCP Sequence Number.
* **`CLOSE_WAIT` (Bahaya / Bug Kode - Sisi Penerima Pamit)**:
  * Terjadi ketika lawan bicara sudah mengirim `FIN` dan kernel lokal sudah membalas `ACK`, tetapi **kode aplikasi backend lupa memanggil fungsi `close()` pada File Descriptor socket tersebut**.
  * Kernel menolak menutup paksa socket tanpa instruksi aplikasi, memicu kebocoran File Descriptor (*socket leak*).

### C. Konsumsi Sumber Daya: Inbound vs Outbound
* **Aplikasi Penerima Saja (Inbound / Database / Backend Pod)**: Hanya mengonsumsi **File Descriptor (FD)** pada port listening yang sama.
* **Aplikasi Reverse Proxy / Gateway (Nginx / Outbound)**: Mengonsumsi **keduanya: File Descriptor (FD) dan Port Sementara (Ephemeral Port)** saat meneruskan koneksi ke upstream.

---

## 3. Pembuktian Langsung di Terminal (One-Liner Proof)

1. **Audit Socket Listening beserta PID dan File Descriptor**:
   ```bash
   sudo ss -tlpn
   # Contoh output nyata:
   # users:(("antigravity",pid=3420,fd=111)) -> Membuktikan socket adalah FD!
   # users:(("systemd-resolve",pid=451,fd=25)) -> Resolver DNS lokal di port 53.
   ```

2. **Membuktikan Penolakan 0 ms Kernel (Paket RST)**:
   ```bash
   curl -v http://127.0.0.1:9999
   # Output: connect to 127.0.0.1 port 9999 ... failed: Connection refused after 0 ms
   ```

3. **Melihat Socket Nyata di Status TIME-WAIT**:
   ```bash
   ss -tan '( sport = :8888 or dport = :8888 )'
   # Output: TIME-WAIT 0 0 127.0.0.1:49870 127.0.0.1:8888
   ```

---

## 4. Jebakan Produksi & Pertanyaan Interview (Gotchas)

* **Pertanyaan Interview**:  
  *"Server backend kami mengalami lonjakan ribuan socket dalam status `CLOSE_WAIT` dan mulai melempar error `Too many open files`. Apakah kamu akan mengubah konfigurasi kernel sysctl atau memeriksa kode aplikasi?"*
* **Jawaban Taktis Operator**:  
  *"Masalah `CLOSE_WAIT` adalah 100% bug di tingkat aplikasi. Status ini menandakan kernel sudah menerima penutupan koneksi dari remote, tetapi aplikasi lokal tidak pernah memanggil `close()` pada socket File Descriptor-nya. Solusinya adalah memperbaiki error-handling atau lifecycle connection di kode aplikasi, bukan mengubah parameter kernel."*

* **Pertanyaan Interview**:  
  *"Apa yang terjadi jika Nginx me-reverse proxy jutaan request ke backend tanpa konfigurasi `keepalive`?"*
* **Jawaban Taktis Operator**:  
  *"Nginx akan membuka dan menutup koneksi TCP baru untuk setiap request. Setiap koneksi yang ditutup oleh Nginx akan tertahan di status `TIME_WAIT` selama 60 detik. Karena jumlah port sementara (ephemeral ports) terbatas (~28.000 port), server akan mengalami Ephemeral Port Exhaustion dan melempar error `Cannot assign requested address (EADDRNOTAVAIL)` meskipun RAM dan CPU masih longgar."*

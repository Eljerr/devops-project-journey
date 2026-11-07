# ROADMAP.md — Kurikulum Fondasi DevOps Zero-to-Production

Kurikulum ini dirancang untuk merekonstruksi pemahaman DevOps dari tingkat paling fundamental (kernel Linux) hingga orkestrasi kontainer, CI/CD, GitOps, dan observabilitas di lingkungan homelab hemat sumber daya (*repurposed hardware*).

---

## Ringkasan Roadmap & Matriks Kompetensi

| Fase | Estimasi | Fokus Kompetensi & Target Direktori | Chaos Drill & Diagnostik (80% Praktik) | Artefak Portofolio Wajib (`docs/`) |
| :--- | :--- | :--- | :--- | :--- |
| **Fase 1** | 3 Minggu | Kernel Linux, File Descriptors, cgroups v2, Namespaces, & systemd (`WSL2, Host Proxmox, & STB`) | Eksekusi OOM Killer (Exit 137), File descriptor leak (`deleted open files`), & systemd resource throttling | `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md`<br>`docs/runbooks/rb-01-linux-diagnostics-and-hardening.md` |
| **Fase 2** | 3 Minggu | Jaringan TCP/IP, DNS, UFW/iptables, Nginx & Cloudflared (`core-infrastructure/gateway-stb/`) | Galat 502 Bad Gateway & 504 Gateway Timeout akibat pemblokiran tabel FORWARD / DNS | `docs/decisions/0001-use-stb-arm64-for-cloudflare-tunnel-gateway.md`<br>`docs/runbooks/rb-02-host-firewall-and-gateway-operations.md` |
| **Fase 3** | 4 Minggu | Otomatisasi Bash, Python asyncio, & Go Konkurensi (`scripts/` & `go-journey/`) | Goroutine leak, race condition (`go test -race`), & kegagalan parsing JSON API | `docs/decisions/0002-select-go-and-python-asyncio-for-custom-tooling.md`<br>Pembaruan README modul |
| **Fase 4** | 3 Minggu | Internals Kontainer, OverlayFS, & Multi-Stage Docker (`src-app/` & `gateway-stb/`) | Startup race condition (API vs MySQL), exit code 137, & penanganan PID 1 signal forwarding | `docs/runbooks/rb-03-container-hardening-and-optimization.md` (metrik reduksi ukuran image) |
| **Fase 5** | 3 Minggu | Virtualisasi Proxmox VE, LXC Unprivileged, & Terraform IaC (`core-infrastructure/proxmox-k3s/`) | Terraform state drift GUI Proxmox & kegagalan layanan K3s kode 127 / modul kernel missing | `docs/decisions/0003-run-k3s-inside-proxmox-lxc-with-terraform.md`<br>`docs/post-mortems/pm-02-k3s-lxc-kernel-module-and-cgroupv2-failure.md` |
| **Fase 6** | 4 Minggu | Orkestrasi K3s, Probes, & Manajemen Memori (`projects/01-hello-nginx/` & `02-nodejs-api/`) | ImagePullBackOff, CrashLoopBackOff (Secret salah), & FailedScheduling (Memory request) | `docs/decisions/0004-optimize-k3s-control-plane-and-workload-limits.md`<br>`docs/runbooks/rb-04-kubernetes-workload-troubleshooting.md` |
| **Fase 7** | 3 Minggu | Pipeline CI/CD GitHub Actions, Trivy/CodeQL, & ArgoCD GitOps (`.github/workflows/` & `platform-tools/`) | Kegagalan Trivy vulnerability gate, sintaks YAML invalid, & ArgoCD OutOfSync / Degraded drift | `docs/decisions/0005-implement-pull-based-gitops-with-argocd.md`<br>`docs/post-mortems/pm-03-cicd-pipeline-gate-failure-and-argocd-drift.md` |
| **Fase 8** | 3 Minggu | Observabilitas Prometheus, PromQL, Grafana, & Respons Insiden (`monitoring-secops/`) | End-to-End Live Incident: Simulasi lonjakan trafik & kebocoran memori hingga auto-healing | `docs/decisions/0006-centralized-observability-architecture-for-hybrid-homelab.md`<br>`docs/post-mortems/pm-04-end-to-end-load-spike-and-memory-saturation-rca.md` |

---

## Detail Rencana Pembelajaran per Fase (Panduan Eksekusi AI Agent)

### Fase 1: Fundamental Kernel Linux, Sistem Operasi, dan Manajemen Proses
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: Eksperimen isolasi/proses di **Inner Loop (WSL2)**; penerapan SSH hardening & audit host di **Outer Loop (Proxmox & STB)**.
* **Sub-Modul & Urutan Progres**:
  * **1.1 Isolasi Kernel, cgroups v2, Sinyal OS, & Forensik OOM Killer**:
    * *Materi*: Kernel vs user space, `/proc/<PID>/ns` (`pid`, `net`, `mnt`), hierarki `/sys/fs/cgroup` (`memory.max`, `memory.swap.max`, `cpu.max`), sinyal `SIGTERM` (15) vs `SIGKILL` (9), kutukan PID 1, dan mekanisme kernel OOM Killer (Exit Code `137`).
    * *Chaos Drill*: Memicu OOM Killer secara terukur pada cgroup terbatas dan menganalisis `memory.events` serta log kernel.
  * **1.2 File Descriptors, Inode, & Misteri Disk Penuh (`deleted open files`)**:
    * *Materi*: Konsep *Everything is a File*, tabel *File Descriptor* di `/proc/<PID>/fd`, perbedaan blok data vs *inode*, mekanisme pembacaan `df -h` (tingkat filesystem/superblock) vs `du -sh` (tingkat traversal direktori).
    * *Chaos Drill*: Membuat proses yang menahan file log besar yang sudah dihapus (`rm`), membuktikan anomali ruang disk antara `df` dan `du`, mendiagnosis dengan `lsof +L1` dan `strace`, serta memulihkan ruang disk tanpa me-reboot server (teknik *file truncation* via `/proc/<PID>/fd/`).
  * **1.3 Manajemen Layanan `systemd` & Penguatan Host (SSH Hardening)**:
    * *Materi*: Anatomi unit file `systemd` (`[Unit]`, `[Service]`, `[Install]`), integrasi `systemd` dengan `cgroups v2` (`MemoryMax`, `CPUQuota`, `Restart=on-failure`), analisis `systemd-cgls` dan `journalctl`, serta pengamanan akses SSH (kunci Ed25519, mematikan `PasswordAuthentication` & `PermitRootLogin` di STB/Proxmox).
    * *Chaos Drill*: Membuat layanan *dummy* yang gagal start berulang kali dan terkena limitasi CPU/RAM `systemd`, lalu mendiagnosisnya via `systemctl status` dan `journalctl -xeu`.
* **Pertanyaan Sokratik Wajib (Uji Konsep Tahap 2)**:
  1. Mengapa proses dengan PID 1 di dalam namespace kontainer mengabaikan `SIGTERM` jika tidak memiliki *signal handler* eksplisit?
  2. Mengapa menghapus file (`rm`) yang sedang ditulis oleh proses aktif tidak langsung mengembalikan ruang kosong di `df -h`, dan apa bahayanya jika kita langsung membunuh proses tersebut dengan `SIGKILL` di produksi?
  3. Bagaimana `systemd` menerjemahkan baris `MemoryMax=100M` di unit file menjadi aturan di tingkat kernel Linux?
* **DoD (Definition of Done)**:
  * Catatan Konsep (`docs/concepts/`): Isolasi Namespaces, Cgroups v2, Sinyal OS, dan File Descriptors.
  * Post-Mortem Terbit: `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` (Bagian 1: OOM Killer & Bagian 2: Deleted Open Files / Inode Exhaustion).
  * Runbook Operasional Terbit: `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`.

---

### Fase 2: Jaringan Komputer, Keamanan Host, dan Arsitektur Reverse Proxy
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: **Outer Loop (STB ARM64 & Proxmox Host)** — meremediasi direktori `core-infrastructure/gateway-stb/`.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * `core-infrastructure/gateway-stb/docker-compose.yml` & konfigurasi Nginx.
  * Menghapus catatan penangguhan keamanan UFW (*deferred UFW security*) di `README.md` utama melalui implementasi firewall nyata.
* **Sub-Modul**:
  * **2.1 Packet Flow, Routing, & DNS**: Model TCP/IP, CIDR subnetting, tabel routing (`ip route`), TCP 3-way handshake & state (`TIME_WAIT`, `CLOSE_WAIT`), resolusi DNS (`/etc/resolv.conf`), serta bedah alur paket: Internet -> Cloudflare Edge -> `cloudflared` (STB) -> Nginx Reverse Proxy (STB) -> LAN -> Proxmox K3s `NodePort`.
  * **2.2 Linux Netfilter, Firewalling (UFW/iptables), & Gateway Rebuild**: Rantai netfilter (`INPUT`, `FORWARD`, `PREROUTING`, `POSTROUTING`, SNAT/DNAT), jebakan Docker vs UFW (manipulasi langsung tabel `iptables` oleh Docker yang mem-bypass UFW), serta penerapan *default deny incoming* di STB dan Proxmox tanpa memutus trafik klaster.
  * **2.3 Chaos & Network Troubleshooting**: Simulasi 502 Bad Gateway (upstream mati/port terblokir) dan 504 Gateway Timeout (paket di-DROP di tabel `FORWARD` atau kegagalan DNS) menggunakan `tcpdump`, `ss -tulpn`, `dig`, dan `curl -v`.
* **Pertanyaan Sokratik Wajib**:
  1. Apa perbedaan kondisi kegagalan di tingkat TCP/aplikasi yang menyebabkan Nginx mengembalikan `502 Bad Gateway` dibandingkan `504 Gateway Timeout`?
  2. Mengapa mengaktifkan UFW saja sering kali gagal memblokir port kontainer Docker yang dipublikasikan dengan `-p`, dan bagaimana cara mengatasinya di rantai `iptables`?
* **DoD**:
  * ADR: `docs/decisions/0001-use-stb-arm64-for-cloudflare-tunnel-gateway.md`.
  * Runbook: `docs/runbooks/rb-02-host-firewall-and-gateway-operations.md`.

---

### Fase 3: Otomatisasi Sistem dengan Bash dan Pemrograman Konkuren (Python dan Go)
* **Durasi**: 4 Minggu
* **Lingkungan Eksekusi**: Pengembangan & pengujian di **Inner Loop (WSL2)**; deployment skrip backup & biner ARM64 di **Outer Loop (Proxmox & STB)**.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * `scripts/python/01-monitoring-scripts/` & `go-journey/01-fundamentals/` serta `02-concurrency-network/`.
* **Sub-Modul**:
  * **3.1 Defensive & Idempotent Bash Scripting**: Aturan `set -euo pipefail`, *exit codes* (`$?`), *trap signals* untuk pembersihan file sementara, pemrosesan teks (`awk`, `sed`, `grep`, `jq`), dan penulisan skrip backup konfigurasi `/etc/pve/lxc/` yang idempoten dengan jadwal `systemd timer`.
  * **3.2 Python Asynchronous I/O (`asyncio` & `aiohttp`)**: Mekanisme *single-threaded event loop*, *coroutines* (`async`/`await`), dan bedah ulang skrip pemantau layanan di `scripts/python/01-monitoring-scripts/` dari nol.
  * **3.3 Go Systems & Concurrency CLI**: Bedah mekanisme *M:N scheduler* Go (*goroutines* vs OS threads), sinkronisasi `sync.WaitGroup`, komunikasi `channels`, pembatalan via `context.WithTimeout`, *structured logging* (`log/slog`), dan *cross-compilation* (`GOOS=linux GOARCH=arm64`) untuk dijalankan di STB.
  * **3.4 Chaos & Concurrency Bugs**: Injeksi *goroutine leak* (endpoint menggantung tanpa timeout), *race condition* pada map/slice yang dideteksi dengan `go test -race`, serta respons JSON rusak.
* **Pertanyaan Sokratik Wajib**:
  1. Apa yang terjadi pada konsumsi memori dan jumlah *goroutine* di agen pemantau jika salah satu server target mengalami *blackhole* (tidak merespons paket sama sekali) dan kita tidak menggunakan `context.WithTimeout`?
  2. Kapan seorang engineer harus memilih menulis otomasi dengan Bash, kapan dengan Python `asyncio`, dan kapan dengan Go?
* **DoD**:
  * ADR: `docs/decisions/0002-select-go-and-python-asyncio-for-custom-tooling.md`.
  * Pembaruan dokumentasi teknis & pengujian pada modul `scripts/` dan `go-journey/`.

---

### Fase 4: Internals Kontainerisasi dan Optimasi Image Docker
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: **Inner Loop (WSL2)** untuk eksperimen `OverlayFS` dan multi-stage build; verifikasi beban di **Outer Loop**.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * `src-app/Dockerfile`, `src-app/server.js`, dan `core-infrastructure/gateway-stb/`.
* **Sub-Modul**:
  * **4.1 Container Runtime & OverlayFS Internals**: Arsitektur `containerd` -> `runc`, bedah cara kerja `OverlayFS` (`lowerdir`, `upperdir`, `workdir`, `merged`) secara manual menggunakan perintah `mount -t overlay`.
  * **4.2 Dockerfile Refactoring & Hardening**: Optimasi *layer caching*, *Multi-Stage Builds*, eksekusi sebagai pengguna non-root (`USER node`), penanganan sinyal PID 1 (`tini` vs `node` langsung), dan pengukuran reduksi ukuran image (`docker history`).
  * **4.3 Multi-Container Stack & Chaos**: Orkestrasi lokal (Node.js API, MySQL, Nginx) dengan `healthcheck` (`service_healthy`), batas memori kontainer, serta simulasi *startup race condition* (API menyala sebelum MySQL siap menerima koneksi) dan *OOM exit code 137*.
* **Pertanyaan Sokratik Wajib**:
  1. Bagaimana `OverlayFS` menangani modifikasi pada file yang berasal dari lapisan `lowerdir` (*Copy-on-Write*), dan mengapa menghapus file besar dengan `rm` pada layer `RUN` berikutnya di Dockerfile tidak mengecilkan ukuran akhir image?
  2. Apa risiko keamanan di tingkat kernel host (terutama karena container berbagi kernel yang sama) jika proses di dalam container berjalan sebagai `UID 0 (root)`?
* **DoD**:
  * Runbook: `docs/runbooks/rb-03-container-hardening-and-optimization.md` (dilengkapi tabel metrik ukuran image sebelum vs sesudah optimasi).

---

### Fase 5: Virtualisasi Proxmox, LXC Unprivileged, dan Infrastructure as Code (Terraform)
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: **Outer Loop (Proxmox VE Host 8GB RAM)** — meremediasi `core-infrastructure/proxmox-k3s/`.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * Seluruh file `.tf` di `core-infrastructure/proxmox-k3s/` dan konfigurasi `/etc/pve/lxc/<ID>.conf` di host Proxmox.
* **Sub-Modul**:
  * **5.1 Hypervisor vs LXC & Proxmox Kernel Tuning**: Perbedaan KVM/QEMU vs LXC, mekanisme *UID/GID mapping* pada Unprivileged LXC (UID 0 kontainer -> UID 100000 host), serta alasan teknis mengapa K3s di dalam LXC membutuhkan modul kernel `overlay`, `br_netfilter`, sysctl `net.ipv4.ip_forward=1`, `/dev/kmsg`, dan delegasi `cgroup2` di host Proxmox.
  * **5.2 Terraform IaC Lifecycle & Refactoring**: Siklus hidup Terraform (`init`, `plan`, `apply`, `destroy`), struktur `terraform.tfstate`, keamanan kredensial/variabel, dan evaluasi kelemahan idempotensi pada `remote-exec` / SSH provisioner di `core-infrastructure/proxmox-k3s/`.
  * **5.3 Chaos & State Drift**: Simulasi perubahan manual alokasi RAM LXC lewat GUI Proxmox (*state drift* reconciliation) dan simulasi kegagalan K3s di LXC (*status 127* / gagal memuat modul `overlay` atau `/dev/kmsg`).
* **Pertanyaan Sokratik Wajib**:
  1. Mengapa menjalankan K3s di dalam kontainer LXC *unprivileged* jauh lebih hemat RAM dibandingkan menggunakan VM KVM penuh di laptop 8 GB, tetapi mengapa secara default K3s langsung gagal berjalan di LXC standar tanpa modifikasi konfigurasi host Proxmox?
  2. Apa yang terjadi jika seorang engineer mengubah spesifikasi server langsung melalui GUI Proxmox tanpa memperbarui kode Terraform, lalu keesokan harinya pipeline menjalankan `terraform apply`?
* **DoD**:
  * ADR: `docs/decisions/0003-run-k3s-inside-proxmox-lxc-with-terraform.md`.
  * Post-Mortem: `docs/post-mortems/pm-02-k3s-lxc-kernel-module-and-cgroupv2-failure.md`.

---

### Fase 6: Orkestrasi Kubernetes Ringan (K3s) pada Lingkungan Terbatas
* **Durasi**: 4 Minggu
* **Lingkungan Eksekusi**: **Outer Loop (Klaster K3s di Proxmox LXC)** — meremediasi `projects/01-hello-nginx/` dan `projects/02-nodejs-api/`.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * Manifest YAML di `projects/01-hello-nginx/`, `projects/02-nodejs-api/`, dan `platform-tools/monitoring/node-exporter-daemonset.yaml`.
* **Sub-Modul**:
  * **6.1 Arsitektur K3s & Penganggaran Memori**: Komponen *Control Plane* (API Server, Scheduler, Controller Manager, SQLite) vs *Worker* (Kubelet, Kube-proxy, Flannel CNI), serta perhitungan penghematan memori dari flag `--disable traefik` dan `--disable servicelb`.
  * **6.2 Workload Hardening & Probes**: Bedah manifest `Pod`, `Deployment`, `ReplicaSet`, `Service` (`ClusterIP` vs `NodePort`), `ConfigMap`, `Secret`, dan `DaemonSet`. Penambahan `requests` & `limits` memori/CPU yang terkalibrasi serta konfigurasi `livenessProbe` dan `readinessProbe`.
  * **6.3 Chaos Tiga Serangkai K8s**: Diagnostik sistematis menggunakan `kubectl describe`, `kubectl get events`, dan `kubectl logs --previous` untuk menangani: (1) `ImagePullBackOff`, (2) `CrashLoopBackOff` akibat kesalahan `Secret`/`ConfigMap` `DB_HOST`, dan (3) `FailedScheduling` (`Pending`) akibat *memory request* melebihi kapasitas node.
* **Pertanyaan Sokratik Wajib**:
  1. Apa perbedaan tindakan yang diambil oleh Kubelet dan `kube-proxy` ketika `readinessProbe` sebuah pod gagal dibandingkan ketika `livenessProbe`-nya gagal?
  2. Apa yang terjadi pada node K3s berkapasitas RAM terbatas jika kita men-deploy pod tanpa mendefinisikan `resources.requests` dan `resources.limits` sama sekali?
* **DoD**:
  * ADR: `docs/decisions/0004-optimize-k3s-control-plane-and-workload-limits.md`.
  * Runbook: `docs/runbooks/rb-04-kubernetes-workload-troubleshooting.md`.

---

### Fase 7: Rekayasa Pipeline CI/CD, Quality Gates, dan Alur Kerja GitOps
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: **GitHub Actions + Klaster K3s (ArgoCD)** — meremediasi `.github/workflows/ci-devsecops.yml`, `src-app/.trivyignore`, `Jenkinsfile`, dan `platform-tools/argocd/`.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * `.github/workflows/ci-devsecops.yml`, `src-app/.trivyignore`, `platform-tools/argocd/bootstrap-argocd.yml`, dan pembersihan/remediasi `Jenkinsfile` placeholder di root repo.
* **Sub-Modul**:
  * **7.1 CI Quality Gates & Vulnerability Management**: Pemisahan tanggung jawab CI vs CD, penambahan *automated unit/endpoint test* sebelum tahap Docker build, cara kerja SAST CodeQL & pemindaian image Trivy, serta audit justifikasi pengecualian CVE di `.trivyignore`.
  * **7.2 Pull-Based GitOps & ArgoCD App-of-Apps**: Perbandingan *Push-based CD* vs *Pull-based GitOps*, cara kerja rekonsiliasi ArgoCD, pola *App of Apps* pada `bootstrap-argocd.yml`, serta mekanisme pencegahan *infinite trigger loop* saat *GitOps Handoff* melakukan commit otomatis ke branch `main` (`paths: ['src-app/**']`).
  * **7.3 Chaos Pipeline & GitOps Drift**: Simulasi kegagalan *gate* Trivy akibat CVE baru, injeksi kesalahan manifest YAML yang memicu status `OutOfSync` / `Degraded` di ArgoCD, serta prosedur *rollback* produksi berbasis Git.
* **Pertanyaan Sokratik Wajib**:
  1. Mengapa langkah *GitOps Handoff* di GitHub Actions yang melakukan `git push` balik ke branch `main` tidak memicu pipeline berjalan berulang kali tanpa henti (*infinite loop*)?
  2. Mengapa kita tidak boleh sekadar memasukkan semua CVE yang merah ke dalam `.trivyignore` agar pipeline hijau, dan apa kriteria teknis yang sah untuk menerima risiko (*risk acceptance*) sebuah CVE?
* **DoD**:
  * ADR: `docs/decisions/0005-implement-pull-based-gitops-with-argocd.md`.
  * Post-Mortem: `docs/post-mortems/pm-03-cicd-pipeline-gate-failure-and-argocd-drift.md`.

---

### Fase 8: Observabilitas Sistem, Metrik Prometheus, dan Simulasi Insiden Produksi
* **Durasi**: 3 Minggu
* **Lingkungan Eksekusi**: **Outer Loop (Node Monitoring LXC, Klaster K3s, & STB)** — meremediasi `monitoring-secops/` dan `platform-tools/monitoring/`.
* **Target Audit Kode Eksisting (Unverified Legacy)**:
  * Konfigurasi Docker Compose di `monitoring-secops/` dan DaemonSet di `platform-tools/monitoring/node-exporter-daemonset.yaml`.
* **Sub-Modul**:
  * **8.1 Arsitektur Prometheus & Optimasi Penyimpanan TSDB**: Perbedaan *Monitoring* vs *Observability*, mekanisme *pull-based scraping*, perbedaan permukaan metrik `Node Exporter` (kernel/host) vs `cAdvisor` (kontainer), tipe metrik (`Counter`, `Gauge`, `Histogram`), serta konfigurasi retensi penyimpanan (`--storage.tsdb.retention.time` & `.size`).
  * **8.2 PromQL & Visualisasi Four Golden Signals**: Penulisan ekspresi PromQL (`rate()`, `irate()`, `sum by`) dari nol untuk memantau *Latency*, *Traffic*, *Errors*, dan *Saturation*, serta ekspor JSON dasbor Grafana ke repositori.
  * **8.3 Grand Finale — End-to-End Live Incident Simulation**: Simulasi lonjakan trafik menggunakan tool konkuren (dari Fase 3) ke endpoint publik Cloudflare Tunnel disertai injeksi kebocoran memori pada `src-app/server.js`. Deteksi via Grafana, observasi pod `OOMKilled`, pengiriman *hotfix* lewat pipeline CI/CD GitHub Actions, sinkronisasi ArgoCD, dan verifikasi *zero-downtime rolling update*.
* **Pertanyaan Sokratik Wajib**:
  1. Mengapa kita harus menggunakan fungsi `rate()` atau `increase()` saat mengkueri metrik bertipe `Counter` di PromQL, dan apa jadinya jika kita langsung menampilkan nilai mentah `Counter` di grafik Grafana?
  2. Apa perbedaan sumber data antara metrik memori yang dilaporkan oleh `Node Exporter` dibandingkan yang dilaporkan oleh `cAdvisor`?
* **DoD**:
  * ADR: `docs/decisions/0006-centralized-observability-architecture-for-hybrid-homelab.md`.
  * Post-Mortem Penutup: `docs/post-mortems/pm-04-end-to-end-load-spike-and-memory-saturation-rca.md`.
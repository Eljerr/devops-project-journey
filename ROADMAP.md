# ROADMAP.md — Kurikulum Fondasi DevOps Zero-to-Production

Kurikulum ini dirancang untuk merekonstruksi pemahaman DevOps dari tingkat paling fundamental (kernel Linux) hingga orkestrasi kontainer, CI/CD, GitOps, dan observabilitas di lingkungan homelab hemat sumber daya (*repurposed hardware*).

---

## Ringkasan Roadmap & Matriks Kompetensi

| Fase | Estimasi | Fokus Kompetensi & Target Direktori | Chaos Drill & Diagnostik (80% Praktik) | Artefak Portofolio Wajib (`docs/`) |
| :--- | :--- | :--- | :--- | :--- |
| **Fase 1** | 3 Minggu | Kernel Linux, systemd, cgroups v2, & Namespaces (`Host Proxmox & STB`) | File descriptors leak (`deleted open files`) & pemicuan OOM Killer pada layanan terbatas | `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`<br>`docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` |
| **Fase 2** | 3 Minggu | Jaringan TCP/IP, DNS, UFW/iptables, Nginx & Cloudflared (`core-infrastructure/gateway-stb/`) | Galat 502 Bad Gateway & 504 Gateway Timeout akibat pemblokiran tabel FORWARD / DNS | `docs/decisions/0001-use-stb-arm64-for-cloudflare-tunnel-gateway.md`<br>`docs/runbooks/rb-02-host-firewall-and-gateway-operations.md` |
| **Fase 3** | 4 Minggu | Otomatisasi Bash, Python asyncio, & Go Konkurensi (`scripts/` & `go-journey/`) | Goroutine leak, race condition (`go test -race`), & kegagalan parsing JSON API | `docs/decisions/0002-select-go-and-python-asyncio-for-custom-tooling.md`<br>Pembaruan README modul |
| **Fase 4** | 3 Minggu | Internals Kontainer, OverlayFS, & Multi-Stage Docker (`src-app/` & `gateway-stb/`) | Startup race condition (API vs MySQL), exit code 137, & penanganan PID 1 signal forwarding | `docs/runbooks/rb-03-container-hardening-and-optimization.md` (metrik reduksi ukuran image) |
| **Fase 5** | 3 Minggu | Virtualisasi Proxmox VE, LXC Unprivileged, & Terraform IaC (`core-infrastructure/proxmox-k3s/`) | Terraform state drift GUI Proxmox & kegagalan layanan K3s kode 127 / modul kernel missing | `docs/decisions/0003-run-k3s-inside-proxmox-lxc-with-terraform.md`<br>`docs/post-mortems/pm-02-k3s-lxc-kernel-module-and-cgroupv2-failure.md` |
| **Fase 6** | 4 Minggu | Orkestrasi K3s, Probes, & Manajemen Memori (`projects/01-hello-nginx/` & `02-nodejs-api/`) | ImagePullBackOff, CrashLoopBackOff (Secret salah), & FailedScheduling (Memory request) | `docs/decisions/0004-optimize-k3s-control-plane-and-workload-limits.md`<br>`docs/runbooks/rb-04-kubernetes-workload-troubleshooting.md` |
| **Fase 7** | 3 Minggu | Pipeline CI/CD GitHub Actions, Trivy/CodeQL, & ArgoCD GitOps (`.github/workflows/` & `platform-tools/`) | Kegagalan Trivy vulnerability gate, sintaks YAML invalid, & ArgoCD OutOfSync / Degraded drift | `docs/decisions/0005-implement-pull-based-gitops-with-argocd.md`<br>`docs/post-mortems/pm-03-cicd-pipeline-gate-failure-and-argocd-drift.md` |
| **Fase 8** | 3 Minggu | Observabilitas Prometheus, PromQL, Grafana, & Respons Insiden (`monitoring-secops/`) | End-to-End Live Incident: Simulasi lonjakan trafik & kebocoran memori hingga auto-healing | `docs/decisions/0006-centralized-observability-architecture-for-hybrid-homelab.md`<br>`docs/post-mortems/pm-04-end-to-end-load-spike-and-memory-saturation-rca.md` |

---

## Detail Rencana Pembelajaran per Fase

### Fase 1: Fundamental Kernel Linux, Sistem Operasi, dan Manajemen Proses
* **Durasi**: 3 Minggu
* **Fokus**: Arsitektur kernel vs user space, syscalls, pseudofilesystems (`/proc`, `/sys`, `/dev`), process lifecycle (`fork`, `exec`, zombie, orphan), sinyal OS (`SIGTERM` kill -15 vs `SIGKILL` kill -9), Linux Namespaces (`pid`, `net`, `mnt`, `uts`, `ipc`, `user`), cgroups v2 (`/sys/fs/cgroup`), systemd service & slice management.
* **Sub-Modul**:
  * **1.1 Hirarki Berkas & Isolasi Kernel**: Eksplorasi `/proc/<PID>/ns`, `/sys/fs/cgroup`, manipulasi process state.
  * **1.2 systemd & Host Hardening**: Unit files, pembatasan `MemoryMax=100M`, `CPUQuota=20%`, SSH hardening (Ed25519, disable root password).
  * **1.3 Diagnostik & Chaos**: File descriptor leakage (`lsof +L1`), OOM killer trigger (`dmesg -T | grep -i oom`).
* **DoD (Definition of Done)**:
  * Lulus ujian konseptual (SIGTERM vs SIGKILL, cgroups vs namespaces).
  * Runbook operasional terverifikasi: `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`.
  * Post-mortem insiden selesai: `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md`.

### Fase 2: Jaringan Komputer, Keamanan Host, dan Arsitektur Reverse Proxy
* **Durasi**: 3 Minggu
* **Fokus**: TCP/IP model, CIDR subnetting, routing tables, TCP 3-way handshake & states (`TIME_WAIT`, `CLOSE_WAIT`), DNS resolution (`/etc/resolv.conf`), HTTP status codes (502, 503, 504), Linux netfilter & iptables/nftables (`INPUT`, `FORWARD`, `PREROUTING`, `POSTROUTING`, NAT), Nginx reverse proxy, Cloudflare Tunnel (`cloudflared`).
* **Sub-Modul**:
  * **2.1 Packet Flow & Routing**: Trace paket dari internet -> Cloudflare -> STB -> Proxmox K3s NodePort.
  * **2.2 Firewalling & Gateway Rebuild**: UFW default deny incoming, iptables bridge forwarding, Nginx upstream proxy.
  * **2.3 Chaos & Network Troubleshooting**: Simulasi 502 Bad Gateway / 504 Gateway Timeout dengan `tcpdump`, `ss -tulpn`, `dig`.
* **DoD**:
  * ADR: `docs/decisions/0001-use-stb-arm64-for-cloudflare-tunnel-gateway.md`.
  * Runbook: `docs/runbooks/rb-02-host-firewall-and-gateway-operations.md`.

### Fase 3: Otomatisasi Sistem dengan Bash dan Pemrograman Konkuren (Python dan Go)
* **Durasi**: 4 Minggu
* **Fokus**: Production Bash (`set -euo pipefail`, exit codes `$?`, idempotent logic, awk/sed/grep/jq), Python asyncio & aiohttp non-blocking event loops, Go concurrency (goroutines, channels, `sync.WaitGroup`, `context.WithTimeout`, `log/slog`).
* **Sub-Modul**:
  * **3.1 Idempotent Bash Automation**: Script backup konfigurasi `/etc/pve/lxc/` terjadwal via systemd timer.
  * **3.2 Python Async Tooling**: Tool pemantau status layanan asinkron di `scripts/python/01-monitoring-scripts/`.
  * **3.3 Go Concurrency CLI**: Tool diagnostik konkuren di `go-journey/` dengan cross-compile `GOOS=linux GOARCH=arm64`.
  * **3.4 Chaos & Concurrency Bugs**: Deteksi goroutine leak & race condition (`go test -race`).
* **DoD**:
  * ADR: `docs/decisions/0002-select-go-and-python-asyncio-for-custom-tooling.md`.
  * README terverifikasi untuk tool Python dan CLI Go.

### Fase 4: Internals Kontainerisasi dan Optimasi Image Docker
* **Durasi**: 3 Minggu
* **Fokus**: Container runtime (containerd, runc), OverlayFS storage architecture (`lowerdir`, `upperdir`, `merged`), layer caching, multi-stage builds, rootless container security, init process / PID 1 signal handling.
* **Sub-Modul**:
  * **4.1 Dockerfile Refactoring**: Multi-stage build Node.js Alpine, non-root user execution, Trivy CVE mitigation.
  * **4.2 Multi-Container Local Stack**: Docker Compose (Node.js, MySQL, Nginx) dengan `healthcheck` & dependency order.
  * **4.3 Chaos**: Startup race conditions, exit code 137 (OOM), penanganan sinyal PID 1 dengan `tini`.
* **DoD**:
  * Runbook: `docs/runbooks/rb-03-container-hardening-and-optimization.md` dengan bukti metrik reduksi ukuran image.

### Fase 5: Virtualisasi Proxmox, LXC Unprivileged, dan Infrastructure as Code (Terraform)
* **Durasi**: 3 Minggu
* **Fokus**: KVM/QEMU Hypervisor vs LXC System Container, UID/GID mapping pada unprivileged LXC (UID 0 mapped to 100000+), Terraform lifecycle (init, plan, apply, state locking), Proxmox kernel tuning (`ip_forward`, `br_netfilter`, `overlay`, `/dev/kmsg`, cgroup2).
* **Sub-Modul**:
  * **5.1 Proxmox Host Pre-configuration**: Modul kernel host & sysctl parameters.
  * **5.2 Terraform Refactoring**: Manajemen node K3s (1 Master, 2 Worker) + Monitoring LXC di `core-infrastructure/proxmox-k3s/`.
  * **5.3 Chaos**: Terraform state drift recovery (modifikasi manual GUI) & K3s LXC failure exit code 127.
* **DoD**:
  * ADR: `docs/decisions/0003-run-k3s-inside-proxmox-lxc-with-terraform.md`.
  * Post-Mortem: `docs/post-mortems/pm-02-k3s-lxc-kernel-module-and-cgroupv2-failure.md`.

### Fase 6: Orkestrasi Kubernetes Ringan (K3s) pada Lingkungan Terbatas
* **Durasi**: 4 Minggu
* **Fokus**: Arsitektur biner K3s (Control Plane vs Worker), memory-saving flags (`--disable traefik`, `--disable servicelb`), Core Objects (Pod, Deployment, ReplicaSet, Service ClusterIP vs NodePort, ConfigMap, Secret, DaemonSet), Probe mechanics (`livenessProbe` vs `readinessProbe`).
* **Sub-Modul**:
  * **6.1 Workload Hardening**: Refactoring manifest `projects/01-hello-nginx/` dan `02-nodejs-api/` dengan probes & strict memory limits.
  * **6.2 Networking & Storage In-Cluster**: Flannel CNI, ClusterIP to NodePort routing, service discovery.
  * **6.3 Chaos**: Tiga serangkai error: `ImagePullBackOff`, `CrashLoopBackOff` (DB_HOST Secret), dan `FailedScheduling` (Memory requests exhaustion).
* **DoD**:
  * ADR: `docs/decisions/0004-optimize-k3s-control-plane-and-workload-limits.md`.
  * Runbook: `docs/runbooks/rb-04-kubernetes-workload-troubleshooting.md`.

### Fase 7: Rekayasa Pipeline CI/CD, Quality Gates, dan Alur Kerja GitOps
* **Durasi**: 3 Minggu
* **Fokus**: CI responsibilities (unit test, CodeQL SAST, container build, Trivy scan) vs CD responsibilities (Push vs Pull GitOps), ArgoCD App of Apps pattern, GitOps handoff without infinite trigger loops (`paths: ['src-app/**']`), risk exception tracking di `.trivyignore`.
* **Sub-Modul**:
  * **7.1 GitHub Actions Enhancement**: Automated test runner & security scanner di `.github/workflows/ci-devsecops.yml`.
  * **7.2 GitOps ArgoCD Synchronization**: App of Apps pattern di `platform-tools/argocd/`, DORA Lead Time tracking.
  * **7.3 Chaos**: Trivy CRITICAL gate failure, syntax error YAML in GitOps repo causing OutOfSync/Degraded, Git rollback.
* **DoD**:
  * ADR: `docs/decisions/0005-implement-pull-based-gitops-with-argocd.md`.
  * Post-Mortem: `docs/post-mortems/pm-03-cicd-pipeline-gate-failure-and-argocd-drift.md`.

### Fase 8: Observabilitas Sistem, Metrik Prometheus, dan Simulasi Insiden Produksi
* **Durasi**: 3 Minggu
* **Fokus**: Monitoring vs Observability, Prometheus pull architecture, Node Exporter vs cAdvisor metric surface, PromQL expressions (`rate()`, `irate()`, `sum by`) untuk Four Golden Signals (Latency, Traffic, Errors, Saturation).
* **Sub-Modul**:
  * **8.1 Prometheus & Retention Tuning**: Storage retention flags (`--storage.tsdb.retention.time` dan `.size`) pada LXC monitoring.
  * **8.2 Grafana Dashboards**: Golden Signals dashboard export (JSON & screenshots).
  * **8.3 Grand Finale Chaos Drill**: Live incident simulation (load spike + memory leak -> Grafana alert -> OOMKilled pod -> zero-downtime rolling update via GitOps).
* **DoD**:
  * ADR: `docs/decisions/0006-centralized-observability-architecture-for-hybrid-homelab.md`.
  * Post-Mortem: `docs/post-mortems/pm-04-end-to-end-load-spike-and-memory-saturation-rca.md`.

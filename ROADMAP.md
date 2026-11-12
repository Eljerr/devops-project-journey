# ROADMAP.md — Kurikulum Fondasi DevOps Zero-to-Production (v2)

Kurikulum ini membangun pemahaman DevOps dari tingkat paling dasar (kernel Linux) sampai orkestrasi kontainer, CI/CD, GitOps, dan observabilitas di homelab hemat sumber daya (*repurposed hardware*).

---

## Aturan Umum Semua Fase (untuk Agen)

1. **Sebelum fase dimulai**, harus ada paket pengetahuan `ai-context/knowledge/fase-<N>.md` (buat dari `_template.md` jika belum ada, lihat `ai-context/knowledge/README.md`).
2. Setiap sub-modul mengikuti mesin status T0–T5 di `AGENTS.md`.
3. Setiap lab dicatat di `docs/labs/` memakai `docs/templates/lab-session-template.md` (bukti mentah = sumber kebenaran).
4. Fase ditutup dengan **Blind Drill** (`ai-context/03-asesmen.md` bagian 4) dan checklist Gate di `STATE.md`.
5. **Lokasi kode**: kode lama tetap di `archive/v1/` sebagai pembanding (Unverified). Kode hasil bangun ulang ditulis di root repo pada lokasi di kolom "Target Bangun Ulang". Operator boleh mengubah lokasinya lewat ADR.
6. Pertanyaan Sokratik di bawah adalah **contoh minimal**. Kunci jawabannya ada di knowledge pack, bukan di ingatan agen.

---

## Ringkasan Roadmap

| Fase | Estimasi | Fokus | Lingkungan | Kode Lama (audit) → Target Bangun Ulang | Artefak Wajib (`docs/`) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **1** | 3 mgg | Kernel, namespaces, cgroups v2, sinyal, FD/inode, systemd, SSH | INNER CachyOS; SSH di Proxmox & STB | — (bukti lab di `docs/labs/`) | 4 catatan konsep, `pm-01`, `rb-01` |
| **2** | 3 mgg | TCP/IP, DNS, netfilter/UFW, Docker vs firewall, Nginx, Cloudflared, TLS dasar | INNER CachyOS → OUTER STB | `archive/v1/core-infrastructure/gateway-stb/` → `infra/gateway-stb/` | ADR-0002, `rb-02` |
| **3** | 4 mgg | Bash defensif, Python asyncio, Go konkurensi | INNER → deploy ke OUTER | `archive/v1/scripts/`, `archive/v1/go-journey/` → `tools/` | ADR-0003, README tools |
| **4** | 3 mgg | Internals kontainer, OverlayFS, Dockerfile multi-stage, Compose | INNER | `archive/v1/src-app/` → `apps/src-app/` | `rb-03` (metrik ukuran image) |
| **5** | 3 mgg | Proxmox, LXC unprivileged, Terraform, backup | OUTER Proxmox | `archive/v1/core-infrastructure/proxmox-k3s/` + `terraform/k3s-lxc/` → `terraform/k3s-lxc/` | ADR-0004, `pm-02` |
| **6** | 4 mgg | K3s, probes, requests/limits | OUTER K3s | `archive/v1/projects/` → `k8s/` | ADR-0005, `rb-04` |
| **7** | 3 mgg | Git workflow, GitHub Actions, Trivy/CodeQL, ArgoCD, secrets | GitHub + K3s | `archive/v1/.github/`, `archive/v1/platform-tools/` → `.github/`, `platform/` | ADR-0006, `pm-03` |
| **8** | 3 mgg | Prometheus, PromQL, Grafana, (Loki), respons insiden | OUTER Monitoring LXC + STB | `archive/v1/monitoring-secops/` → `monitoring/` | ADR-0007, `pm-04` |

> Penomoran ADR: `0001` sudah dipakai untuk arsitektur hardware hibrida. ADR fase berikutnya dimulai dari `0002`.

---

## Detail per Fase

### Fase 1: Fundamental Kernel Linux, Proses, Filesystem, systemd
* **Durasi**: 3 minggu · **Status**: diulang dalam mode **RE-VERIFIKASI** (lihat `STATE.md`).
* **Lingkungan**: eksperimen di **INNER CachyOS** (kernel asli, shell fish, swap zram); SSH hardening di **OUTER Proxmox & STB**.
* **Knowledge pack**: `ai-context/knowledge/fase-1.md`
* **Sub-modul**:
  * **1.1 Namespaces, cgroups v2, Sinyal, OOM Killer**: `/proc/<PID>/ns`, `unshare`, hierarki `/sys/fs/cgroup`, `memory.max`/`memory.swap.max`/`memory.events`, `cpu.max`, `systemd-run`, SIGTERM vs SIGKILL, aturan sinyal PID 1, kode keluar 128+N, peran zram.
    * *Chaos*: OOM terukur; **matriks 2×2 swap × jenis data** (membuktikan ulang klaim zram); PID 1 yang "tuli".
  * **1.2 File Descriptor, Inode, Disk Penuh**: `/proc/<PID>/fd` dan `fdinfo`, inode vs nama file, `unlink()`, `df` vs `du`, `lsof +L1`, offset & `O_APPEND`, sparse file, inode exhaustion.
    * *Chaos*: deleted open file di ext4 (`~/lab/`, bukan `/tmp` yang tmpfs); **inode habis sungguhan** di loop device ext4.
  * **1.3 systemd & SSH Hardening**: anatomi unit, `MemoryMax`/`CPUQuota` → file cgroup, `Restart=` & `start-limit-hit`, `journalctl`, Ed25519, izin file, aturan "nilai pertama menang" `sshd_config`, `sshd -t`/`sshd -T`, rollback anti-lockout.
    * *Chaos*: service crash berulang; drop-in sshd yang mengalahkan konfigurasi utama (di STB, dengan sesi cadangan).
* **Pertanyaan Sokratik**:
  1. Mengapa PID 1 di PID namespace tidak mati oleh `SIGTERM` tanpa handler, tetapi tetap mati oleh `SIGKILL` dari host?
  2. Mengapa `rm` pada file yang sedang ditulis tidak mengembalikan ruang di `df`? Apa beda masalah ini dengan inode habis?
  3. Bagaimana `MemoryMax=100M` di unit file berubah menjadi aturan kernel, dan file apa yang membuktikannya?
* **DoD**:
  * Catatan konsep: `docs/concepts/02-linux-namespaces.md`, `03-cgroups-v2.md`, `04-os-signals-and-pid1.md`, `05-file-descriptors-and-inodes.md`.
  * Post-mortem v2: `docs/post-mortems/pm-01-oom-killer-and-deleted-open-files.md` (menggantikan versi lama).
  * Runbook v2: `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md`.

---

### Fase 2: Jaringan, Keamanan Host, Reverse Proxy
* **Durasi**: 3 minggu
* **Lingkungan**: konsep & firewall lokal di **INNER CachyOS** (ufw + Docker tersedia) → penerapan di **OUTER STB** (ARM64) & Proxmox.
* **Knowledge pack**: `ai-context/knowledge/fase-2.md`
* **Audit kode lama**: `archive/v1/core-infrastructure/gateway-stb/` (`docker-compose.yaml`, `nginx/nginx.conf`, `nginx/conf.d/default.conf`).
* **Sub-modul**:
  * **2.1 Packet Flow, TCP, Socket, DNS**: model TCP/IP, CIDR, `ip route`, handshake, RST vs timeout, `TIME_WAIT`/`CLOSE_WAIT`, port sementara, `/etc/resolv.conf` & `dig`, alur Internet → Cloudflare → `cloudflared` → Nginx → LAN → K3s.
  * **2.2 Netfilter, UFW, Docker**: 5 hook netfilter, tabel filter/nat, DNAT/SNAT, conntrack, iptables-nft vs nftables, UFW, Docker DNAT melewati FORWARD (menembus UFW), `DOCKER-USER`, publish ke `127.0.0.1`, *default deny incoming* tanpa memutus klaster.
  * **2.3 Chaos 502 vs 504**: upstream mati (RST) → 502; paket di-`DROP` → timeout → 504; nama upstream salah → `nginx -t` gagal. Alat: `tcpdump`, `ss`, `dig`, `curl -v`, `error.log`.
  * **2.4 TLS Dasar (baru)**: rantai sertifikat, `openssl s_client`, masa berlaku, di mana TLS berakhir pada Cloudflare Tunnel.
* **Pertanyaan Sokratik**:
  1. Kondisi TCP apa yang membuat Nginx mengembalikan 502, dan kondisi apa yang membuat 504?
  2. Mengapa `ufw deny 8080` tidak menutup port kontainer `-p 8080:80`? Di rantai mana masalahnya bisa diatasi?
* **DoD**:
  * ADR: `docs/decisions/0002-host-firewall-strategy-docker-vs-ufw.md`.
  * Runbook: `docs/runbooks/rb-02-host-firewall-and-gateway-operations.md`.
  * Revisi: `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md` (audit A10–A11).

---

### Fase 3: Otomatisasi Bash & Pemrograman Konkuren (Python, Go)
* **Durasi**: 4 minggu
* **Lingkungan**: kembangkan di **INNER**; deploy skrip backup & biner ARM64 ke **OUTER**.
* **Audit kode lama**: `archive/v1/scripts/python/01-monitoring-scripts/`, `archive/v1/go-journey/`.
* **Sub-modul**:
  * **3.1 Bash Defensif & Idempoten**: `set -euo pipefail` (beserta kelemahannya), exit code, `trap`, `awk`/`sed`/`grep`/`jq`, backup `/etc/pve/lxc/` dengan `systemd timer`. Ingat: shell interaktif operator = fish; skrip ditulis untuk bash dengan shebang.
  * **3.2 Python asyncio**: event loop, coroutine, `aiohttp`, bangun ulang skrip pemantau dari nol.
  * **3.3 Go Concurrency CLI**: goroutine vs thread OS, `sync.WaitGroup`, channel, `context.WithTimeout`, `log/slog`, cross-compile `GOOS=linux GOARCH=arm64`.
  * **3.4 Chaos**: goroutine leak, race (`go test -race`), JSON rusak.
* **Pertanyaan Sokratik**:
  1. Apa yang terjadi pada memori & jumlah goroutine jika target *blackhole* dan tidak ada `context.WithTimeout`?
  2. Kapan memilih Bash, Python asyncio, atau Go?
* **DoD**: ADR `docs/decisions/0003-select-go-and-python-asyncio-for-custom-tooling.md` + README & test di `tools/`.

---

### Fase 4: Internals Kontainer & Optimasi Image Docker
* **Durasi**: 3 minggu
* **Lingkungan**: **INNER** (OverlayFS manual, build multi-stage); uji beban ringan di **OUTER**.
* **Audit kode lama**: `archive/v1/src-app/` (`Dockerfile`, `server.js`), `archive/v1/core-infrastructure/gateway-stb/`.
* **Sub-modul**:
  * **4.1 Runtime & OverlayFS**: `containerd` → `runc`, `mount -t overlay` manual (`lowerdir`, `upperdir`, `workdir`, `merged`), hubungkan ke namespace/cgroup Fase 1.
  * **4.2 Dockerfile Hardening**: layer caching, multi-stage, `USER` non-root, PID 1 (`tini`/`--init` vs langsung), `docker history`, build multi-arch untuk STB.
  * **4.3 Multi-Container & Chaos**: Compose (API, MySQL, Nginx), `healthcheck` + `service_healthy`, batas memori, *startup race*, exit 137 (buktikan OOM vs `docker stop` timeout dengan `docker inspect`).
* **Pertanyaan Sokratik**:
  1. Bagaimana OverlayFS menangani perubahan file dari `lowerdir` (copy-on-write), dan mengapa `rm` di layer `RUN` berikutnya tidak mengecilkan image?
  2. Apa risiko kernel host jika proses kontainer berjalan sebagai UID 0?
* **DoD**: Runbook `docs/runbooks/rb-03-container-hardening-and-optimization.md` (tabel ukuran image sebelum vs sesudah).

---

### Fase 5: Proxmox, LXC Unprivileged, Terraform, Backup
* **Durasi**: 3 minggu
* **Lingkungan**: **OUTER Proxmox** (2 core / 8 GB — total LXC ≤ ~6 GB).
* **Audit kode lama**: `archive/v1/core-infrastructure/proxmox-k3s/*.tf`, `terraform/k3s-lxc/*.tf`, `/etc/pve/lxc/<ID>.conf`.
* **Sub-modul**:
  * **5.1 Hypervisor vs LXC & Kebutuhan Kernel K3s**: KVM vs LXC, UID/GID mapping unprivileged, mengapa K3s di LXC butuh `overlay`, `br_netfilter`, `ip_forward`, `/dev/kmsg`, delegasi cgroup2.
  * **5.2 Terraform Lifecycle**: `init`/`plan`/`apply`/`destroy`, state, kredensial (jangan di Git), kelemahan `remote-exec`.
  * **5.3 Chaos & Drift**: ubah RAM LXC via GUI → `terraform plan`; K3s gagal karena modul kernel / `/dev/kmsg`.
  * **5.4 Backup & Restore (baru)**: `vzdump` LXC + uji restore; backup datastore K3s (metode bergantung SQLite vs etcd — cek dokumentasi K3s versi terpasang).
  * *Opsional*: Ansible dasar untuk konfigurasi host yang idempoten.
* **Pertanyaan Sokratik**:
  1. Mengapa K3s di LXC unprivileged lebih hemat RAM daripada VM, tetapi gagal tanpa konfigurasi host?
  2. Apa yang terjadi jika spesifikasi diubah lewat GUI lalu `terraform apply` dijalankan?
* **DoD**: ADR `docs/decisions/0004-run-k3s-inside-proxmox-lxc-with-terraform.md`, Post-mortem `docs/post-mortems/pm-02-k3s-lxc-kernel-module-and-cgroupv2-failure.md`.

---

### Fase 6: Orkestrasi K3s di Lingkungan Terbatas
* **Durasi**: 4 minggu
* **Lingkungan**: **OUTER K3s di Proxmox LXC**.
* **Audit kode lama**: `archive/v1/projects/01-hello-nginx/`, `archive/v1/projects/02-nodejs-api/`, `archive/v1/platform-tools/monitoring/node-exporter-daemonset.yaml`.
* **Sub-modul**:
  * **6.1 Arsitektur K3s & Anggaran Memori**: control plane vs worker, datastore, ukur dampak `--disable traefik`/`servicelb`.
  * **6.2 Workload Hardening & Probes**: Pod, Deployment, Service, ConfigMap, Secret, DaemonSet, `requests`/`limits`, `livenessProbe` vs `readinessProbe`.
  * **6.3 Chaos Tiga Serangkai**: `ImagePullBackOff`, `CrashLoopBackOff` (Secret/ConfigMap salah), `Pending` (request melebihi kapasitas). Alat: `kubectl describe`, `get events`, `logs --previous`.
  * *Opsional*: Kustomize/Helm dasar.
* **Pertanyaan Sokratik**:
  1. Apa beda tindakan kubelet & endpoint Service saat readiness gagal vs liveness gagal?
  2. Apa yang terjadi pada node RAM kecil jika pod tanpa `requests`/`limits`?
* **DoD**: ADR `docs/decisions/0005-optimize-k3s-control-plane-and-workload-limits.md`, Runbook `docs/runbooks/rb-04-kubernetes-workload-troubleshooting.md`.

---

### Fase 7: CI/CD, Quality Gates, GitOps, Secrets
* **Durasi**: 3 minggu
* **Lingkungan**: **GitHub Actions + K3s (ArgoCD)**.
* **Audit kode lama**: `archive/v1/.github/workflows/ci-devsecops.yml`, `archive/v1/platform-tools/argocd/`, `archive/v1/platform-tools/bootstrap-argocd.yml`, `archive/v1/Jenkinsfile`.
* **Sub-modul**:
  * **7.0 Git Workflow (baru)**: branch, PR, rebase vs merge, commit atomik, proteksi branch.
  * **7.1 CI Quality Gates**: CI vs CD, test sebelum build, CodeQL, Trivy, justifikasi pengecualian CVE.
  * **7.2 Pull-Based GitOps**: push vs pull, rekonsiliasi ArgoCD, App-of-Apps, mencegah loop pemicu dari commit otomatis.
  * **7.3 Chaos**: gate Trivy gagal, manifest rusak → `OutOfSync`/`Degraded`, rollback via Git.
  * **7.4 Secrets di GitOps (baru)**: mengapa Secret base64 bukan enkripsi; SOPS+age atau Sealed Secrets (pilih lewat ADR).
* **Pertanyaan Sokratik**:
  1. Mengapa commit otomatis kembali ke `main` tidak memicu pipeline tanpa henti?
  2. Mengapa tidak boleh memasukkan semua CVE ke `.trivyignore`, dan apa kriteria *risk acceptance* yang sah?
* **DoD**: ADR `docs/decisions/0006-implement-pull-based-gitops-with-argocd.md`, Post-mortem `docs/post-mortems/pm-03-cicd-pipeline-gate-failure-and-argocd-drift.md`.

---

### Fase 8: Observabilitas & Simulasi Insiden
* **Durasi**: 3 minggu
* **Lingkungan**: **OUTER** (Monitoring LXC, K3s, STB).
* **Audit kode lama**: `archive/v1/monitoring-secops/`, `archive/v1/platform-tools/monitoring/`.
* **Sub-modul**:
  * **8.1 Prometheus & TSDB**: monitoring vs observability, pull scraping, Node Exporter vs cAdvisor, tipe metrik, retensi.
  * **8.2 PromQL & Four Golden Signals**: `rate()`, `irate()`, `sum by`, dasbor Grafana sebagai JSON di repo.
  * **8.3 Grand Finale**: lonjakan trafik + memory leak → deteksi Grafana → `OOMKilled` → hotfix via CI/CD → ArgoCD → rolling update tanpa downtime.
  * *Opsional*: Loki + Promtail untuk log terpusat (cek anggaran RAM dulu).
* **Pertanyaan Sokratik**:
  1. Mengapa Counter harus dikueri dengan `rate()`/`increase()`?
  2. Apa beda sumber data memori Node Exporter vs cAdvisor?
* **DoD**: ADR `docs/decisions/0007-centralized-observability-architecture-for-hybrid-homelab.md`, Post-mortem `docs/post-mortems/pm-04-end-to-end-load-spike-and-memory-saturation-rca.md`.

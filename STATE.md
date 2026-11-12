# STATE.md — Memori Status Pembelajaran

> Agen: baca berkas ini PERTAMA. Bagian "Posisi Saat Ini" menentukan apa yang kamu lakukan.
> Update berkas ini di akhir sesi atau saat operator mengetik `!selesai`.

---

## 1. Posisi Saat Ini

| Kunci | Nilai |
| :--- | :--- |
| **Fase aktif** | Fase 1 v2 — Re-verifikasi Fondasi Kernel Linux |
| **Sub-modul aktif** | 1.1 Namespaces, cgroups v2, Sinyal, OOM Killer |
| **Tahap aktif** | T2 (mode re-verifikasi — lihat bagian 2) |
| **Mode** | `RE-VERIFIKASI` |
| **Pembaruan terakhir** | 2026-11-12 |

### Instruksi untuk sesi berikutnya (agen)
1. Ikuti `AGENTS.md` bagian 0 (urutan membaca).
2. Lewati T0 di sesi pertama ini (belum ada modul LULUS-v2).
3. Mulai Sub-modul 1.1 di **T2**: ajukan 2 pertanyaan berbasis `ai-context/knowledge/fase-1.md` (1.1), termasuk minimal 1 miskonsepsi dari bagian 5 di bawah. HARD STOP.

---

## 2. Mode Re-Verifikasi (Fase 1 v2)

**Alasan:** Kelulusan v1 dinilai oleh agen yang bisa berhalusinasi. Audit `ai-context/audit-2026-11-12.md` menemukan beberapa kesalahan di dokumen v1. Operator memutuskan mengulang Fase 1 dengan standar baru.

**Aturan mode ini:**
1. Setiap sub-modul **dimulai di T2** (uji konsep langsung, karena operator sudah pernah belajar materinya).
2. Jika LULUS R1–R3 → lanjut T3 dengan **chaos baru** yang belum pernah dilakukan (lihat "Ide Chaos" di pack).
3. Jika BELUM → kembali ke T1 normal untuk konsep yang salah saja.
4. Semua butir audit A1–A9 harus dibuktikan/dibantah lewat lab sebelum Gate Fase 1.
5. PM-01 dan RB-01 ditulis ulang (v2) dengan templat baru dan suara operator.

---

## 3. Fakta Lingkungan (Terverifikasi)

> Agen: gunakan HANYA fakta ini. Jika butuh fakta lain, minta operator menjalankan perintah deteksi lalu tambahkan di sini.

### Inner Loop — Workstation (dideteksi agen 2026-11-12)
| Item | Nilai | Perintah deteksi |
| :--- | :--- | :--- |
| OS | CachyOS Linux (Arch-based), **native, bukan WSL2** | `cat /etc/os-release` |
| Kernel | 7.3.3-1-cachyos | `uname -r` |
| Shell | **fish 4.9.3** (bukan bash) | `echo $SHELL; fish --version` |
| Filesystem `/` | ext4 | `findmnt -no FSTYPE /` |
| Filesystem `/tmp` | **tmpfs** (pakai RAM) | `findmnt -no FSTYPE /tmp` |
| cgroup | v2 (`cgroup2fs`); controller user: `cpuset cpu io memory pids` | `stat -fc %T /sys/fs/cgroup` |
| Swap | zram 15.5G (zstd), `vm.swappiness=150` | `swapon --show; cat /proc/sys/vm/swappiness` |
| systemd | 262 | `systemctl --version` |
| Python | 3.12.16 | `python3 --version` |
| Firewall | `iptables v1.8.13 (nf_tables)`, `nftables 1.1.7`, `ufw 0.36.2` (enabled) | `iptables -V; nft --version; systemctl is-enabled ufw` |
| Docker | Engine 29.9.3 (enabled) | `docker version` |
| Port sementara | 32768–60999 | `cat /proc/sys/net/ipv4/ip_local_port_range` |

### Outer Loop — Proxmox (BELUM dideteksi)
Jalankan via SSH lalu isi: `pveversion; uname -r; free -h; pct list; stat -fc %T /sys/fs/cgroup`

### Outer Loop — STB (BELUM dideteksi)
Jalankan via SSH lalu isi: `cat /etc/os-release; uname -rm; free -h; nginx -v; cloudflared --version; docker version; sudo sshd -T | grep -Ei 'passwordauth|permitroot'; ls /etc/ssh/sshd_config.d/`

---

## 4. Kompetensi

### Status v2 (standar baru, rubrik `ai-context/03-asesmen.md`)
| Sub-modul | Status | Petunjuk tertinggi | Terakhir diulang | Bukti |
| :--- | :--- | :--- | :--- | :--- |
| 1.1 Namespaces, cgroups, sinyal, OOM | ⏳ Belum | — | — | — |
| 1.2 FD, inode, disk penuh | ⏳ Belum | — | — | — |
| 1.3 systemd & SSH | ⏳ Belum | — | — | — |
| 2.1 Packet flow, TCP, DNS | ⏳ Belum (v1: lulus) | — | — | — |
| 2.2 Netfilter, UFW, Docker | ⏳ Belum | — | — | — |
| 2.3 Chaos 502/504 | ⏳ Belum | — | — | — |
| 2.4 TLS dasar | ⏳ Belum | — | — | — |

### Arsip klaim v1 (BELUM diverifikasi ulang — jangan dianggap dikuasai)
* 1.1.1 Isolasi namespaces · 1.1.2 cgroups v2 · 1.1.3 Sinyal & OOM (exit 137, `memory.events`)
* 1.2 FD, inode, deleted open files, `lsof +L1`, offset & `O_APPEND`
* 1.3 systemd unit, cgroups via systemd, ZRAM vs OOM, SSH Ed25519
* 2.1 Socket = FD, 3-way handshake, RST, `TIME_WAIT` vs `CLOSE_WAIT`, ephemeral port, cloudflared outbound, nginx keepalive

---

## 5. Log Miskonsepsi Operator

> Agen: catat kalimat **persis** operator yang keliru + koreksinya. Gunakan untuk pertanyaan T0/T2.

| Tanggal | Modul | Kalimat operator | Koreksi | Sumber | Sudah diperbaiki? |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 2026-11-06 | 1.1 | "137 menandakan exit code kehabisan memori" | 137 = SIGKILL dari mana pun; OOM dibuktikan `oom_kill` | pack 1.1-F13 | ☐ |
| 2026-11-10 | 1.2 | "truncate itu tidak efektif karena bisa mengubah offset di Open File Table" | Truncate tidak mengubah offset; offset yang tetap itulah masalahnya | pack 1.2-F8 | ☐ |
| v1 | 1.3 | "bytearray lolos OOM karena zstd mengompres nol" | Dugaan: penentunya swap diizinkan/tidak. Belum dibuktikan | pack 1.1 Chaos 2 | ☐ |
| v1 | 2.1 | "paket nyasar dibatasi TTL ~60 detik" | TTL = hop, bukan detik | pack 2.1 | ☐ |

---

## 6. Kamus Analogi Operator

> Milik operator. Agen boleh menambahkan **batas analogi**, tidak boleh mengganti analoginya.

| Konsep | Analogi operator | Batas analogi (di mana analogi ini tidak akurat) |
| :--- | :--- | :--- |
| Namespaces | Kacamata kuda kernel | Bukan hanya "pandangan": `net` namespace punya tumpukan jaringan sendiri sungguhan |
| cgroups v2 | Pagar pembatas / ransum fisik | — |
| Dentry vs Inode & FD | Stiker label pintu vs koper di loker stasiun | — |
| Socket | Nomor tiket gagang telepon dari kernel | — |
| Paket RST | Satpam kernel menolak ketukan di kamar kosong | — |
| TIME_WAIT | Menunggu paket nyasar mati sebelum port dipakai lagi | Paket "mati" karena batas hop (TTL), bukan stopwatch |
| CLOSE_WAIT | Aplikasi lupa `close()` | — |
| Nginx keepalive | Garasi 32 mobil standby | Garasi per worker; hanya mobil yang sedang menganggur |

---

## 7. Gate Fase (Checklist)

### Gate Fase 1
- [ ] 1.1, 1.2, 1.3 LULUS v2
- [ ] Butir audit A1–A9 dibuktikan/dibantah di lab
- [ ] `docs/concepts/` : namespaces, cgroups-v2, sinyal-os, file-descriptor (4 catatan, suara operator)
- [ ] `docs/post-mortems/pm-01-...` v2 (templat baru)
- [ ] `docs/runbooks/rb-01-...` v2 (templat baru, ada backup & `sshd -T`)
- [ ] Blind Drill Fase 1 lulus
- [ ] Kalimat bertanda `[INGATAN]` di dokumen Fase 1 sudah dibuktikan / dicek silang

### Gate Fase 2
- [ ] 2.1–2.4 LULUS v2
- [ ] `docs/concepts/01-tcp-...` direvisi (audit A10, A11)
- [ ] ADR `docs/decisions/0002-host-firewall-strategy-docker-vs-ufw.md`
- [ ] Runbook `docs/runbooks/rb-02-host-firewall-and-gateway-operations.md`
- [ ] Blind Drill Fase 2 lulus

---

## 8. Status Artefak `docs/`

| Berkas | Status |
| :--- | :--- |
| `docs/decisions/0001-hybrid-hardware-architecture.md` | ⚠️ PERLU REVISI — belum ikut templat MADR, setengah bahasa Inggris. Isi bagian konteks sudah suara operator (bagus). |
| `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md` | ⚠️ PERLU REVISI — lihat audit A1–A7. Akan diganti v2. |
| `docs/runbooks/rb-01-linux-diagnostics-and-hardening.md` | ⚠️ PERLU REVISI — lihat audit A7–A9. Akan diganti v2. |
| `docs/concepts/01-tcp-lifecycle-and-proxy-architecture.md` | ⚠️ PERLU REVISI — lihat audit A10–A11. |

---

## 9. Riwayat Sesi

| Tanggal | Fase/Sub | Tahap dicapai | Catatan |
| :--- | :--- | :--- | :--- |
| 2026-11-12 | — | — | Sistem belajar dirombak ke v2 (anti-halusinasi, rubrik, knowledge pack). Fase 1 diulang dalam mode re-verifikasi. |

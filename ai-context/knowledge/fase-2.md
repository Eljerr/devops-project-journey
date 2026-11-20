# Paket Pengetahuan — Fase 2: Jaringan, Firewall, Reverse Proxy

> Status: Draf v1 (2026-11-12) · Penyusun: model besar + cek web + cek langsung di mesin operator
> Lingkungan (cek-agen 2026-11-12): `iptables v1.8.13 (nf_tables)`, `nftables v1.1.7`, `ufw 0.36.2` (enabled), Docker Engine `29.9.3` (enabled), `ip_local_port_range = 32768 60999`.
> STB (Ubuntu ARM64) dan Proxmox: **versi belum dideteksi** — minta operator menjalankan perintah deteksi di `STATE.md`.

---

## Sub-modul 2.1: Packet Flow, TCP, Socket, DNS

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 2.1-F1 | ☑ | Socket adalah FD milik proses. `ss -tlpn` menunjukkan proses pemilik socket LISTEN. | [DOK: ss(8), man 7 socket] [LAB ✓] | `sudo ss -tlpn` → `users:(("containerd",pid=743,fd=13))` |
| 2.1-F2 | ☑ | Koneksi ke port tanpa listener → kernel membalas `RST` → klien langsung dapat `Connection refused`. | [DOK: RFC 9293] [LAB ✓] | `curl -v http://127.0.0.1:5999` → `after 0 ms: Connection refused` vs dengan `python3 -m http.server 5999` → `200 OK` |
| 2.1-F3 | ☐ | Firewall `DROP` tidak membalas apa pun → klien menunggu sampai timeout. `REJECT` membalas (ICMP atau RST) → gagal cepat. | [DOK: iptables-extensions(8) REJECT] | Ide Chaos 1 |
| 2.1-F4 | ☑ | `TIME_WAIT` terjadi di pihak yang **pertama** mengirim `FIN`. | [DOK: RFC 9293 state diagram] [LAB ✓] | `ss -tan state time-wait` setelah `curl` (`127.0.0.1:5999`) |
| 2.1-F5 | ☑ | Di Linux, lama `TIME_WAIT` adalah 60 detik dan **tidak** bisa diubah lewat sysctl. `tcp_fin_timeout` mengatur hal lain (FIN_WAIT_2). | [DOK: docs.kernel.org networking/ip-sysctl `tcp_fin_timeout`] [LAB ✓] | `ss -tano state time-wait` → `timer:(timewait,59sec,0)` |
| 2.1-F6 | ☑ | `CLOSE_WAIT` = pihak lawan sudah menutup, tetapi aplikasi lokal belum memanggil `close()`. Jika menumpuk → bug aplikasi (FD bocor). | [DOK: RFC 9293] [LAB ✓] | Python menahan `conn` dengan `time.sleep(15)` setelah `curl --max-time 1` → `CLOSE-WAIT 79 0 127.0.0.1:5999` |
| 2.1-F7 | ☑ | Koneksi keluar memakai port sementara dari `ip_local_port_range`. Batasnya berlaku **per tujuan** (IP+port tujuan yang sama), bukan global. | [DOK: ip-sysctl `ip_local_port_range`] [LAB ✓] | `cat /proc/sys/net/ipv4/ip_local_port_range` → `32768 60999` |
| 2.1-F8 | ☐ | `keepalive N;` di blok `upstream` nginx = maksimal N koneksi **idle** yang disimpan **per worker process**. Butuh `proxy_http_version 1.1;` dan `proxy_set_header Connection "";`. | [DOK: nginx.org ngx_http_upstream_module#keepalive] | di STB: `ss -tn dst <ip-upstream>` saat ada lalu lintas |
| 2.1-F9 | ☐ | `cloudflared` membuat koneksi **keluar** ke edge Cloudflare (port 7844), jadi tidak butuh port forwarding di modem. | [DOK: developers.cloudflare.com Tunnel firewall] | di STB: `sudo ss -tunp \| grep cloudflared` |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan |
| :--- | :--- |
| "TTL membatasi paket hidup 60 detik" | TTL adalah **jumlah hop**, bukan detik. MSL adalah asumsi di RFC, bukan nilai yang diukur. Linux memilih 60 detik untuk TIME_WAIT. |
| "keepalive 32 = 32 koneksi total" | Per worker, dan hanya koneksi idle (F8). |
| "Ephemeral port habis setelah ~28 ribu koneksi keluar" | Per tujuan yang sama (F7). |

### Sumber Primer
* https://www.rfc-editor.org/rfc/rfc9293
* https://docs.kernel.org/networking/ip-sysctl.html
* https://nginx.org/en/docs/http/ngx_http_upstream_module.html#keepalive
* https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/configure-tunnels/tunnel-with-firewall/

---

## Sub-modul 2.2: Netfilter, UFW, Docker

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 2.2-F1 | ☐ | Netfilter punya 5 hook IPv4: `PREROUTING`, `INPUT`, `FORWARD`, `OUTPUT`, `POSTROUTING`. | [DOK: netfilter.org, iptables(8)] | `sudo iptables -L -n -v` (filter) dan `sudo iptables -t nat -L -n -v` |
| 2.2-F2 | ☐ | Jalur paket: **masuk untuk host** = PREROUTING → keputusan routing → INPUT. **Diteruskan** = PREROUTING → routing → FORWARD → POSTROUTING. **Dibuat lokal** = OUTPUT → routing → POSTROUTING. | [DOK: iptables(8)] | Gunakan counter `-v` sebelum/sesudah `curl` |
| 2.2-F3 | ☐ | DNAT terjadi di `nat PREROUTING` (dan `OUTPUT` untuk lalu lintas lokal). SNAT/MASQUERADE di `nat POSTROUTING`. | [DOK: iptables-extensions(8)] | `sudo iptables -t nat -S` |
| 2.2-F4 | ☐ | NAT diputuskan pada paket **pertama** sebuah koneksi; paket berikutnya mengikuti entri conntrack. | [DOK: netfilter conntrack docs] | `sudo conntrack -L` (paket `conntrack-tools`) |
| 2.2-F5 | ☐ | Di sistem modern, perintah `iptables` sering hanya penerjemah ke nftables. Semua aturan terlihat lewat `nft list ruleset`. | [DOK: iptables-nft(8)] (cek-agen: `(nf_tables)`) | `iptables -V`; `sudo nft list ruleset \| head -50` |
| 2.2-F6 | ☐ | UFW menulis aturannya terutama di rantai INPUT/OUTPUT (rantai `ufw-*`). | [DOK: ufw(8), ufw-framework(8)] | `sudo iptables -S \| grep ufw \| head` |
| 2.2-F7 | ☐ | `docker run -p 8080:80` membuat DNAT di nat PREROUTING ke IP kontainer. Karena tujuan berubah ke IP kontainer, paket lewat **FORWARD**, bukan INPUT → aturan `ufw deny 8080` tidak berlaku. | [DOK: docs.docker.com/engine/network/packet-filtering-firewalls/] | Ide Chaos 2 |
| 2.2-F8 | ☐ | Docker menyediakan rantai `DOCKER-USER` (dipanggil dari FORWARD sebelum aturan Docker) untuk aturan admin. Di sana port tujuan **sudah** di-DNAT; untuk mencocokkan port host asli pakai `-m conntrack --ctorigdstport`. | [DOK: Docker packet filtering docs] | `sudo iptables -L DOCKER-USER -n -v` |
| 2.2-F9 | ☐ | Alternatif paling sederhana: publish hanya ke loopback (`-p 127.0.0.1:8080:80`) jika hanya nginx lokal yang butuh akses. | [DOK: Docker docs] | `ss -tlpn \| grep 8080` → `127.0.0.1:8080` |
| 2.2-F10 | ☐ | Uji firewall harus dari **mesin lain**. Uji dari host sendiri sering melewati jalur berbeda (OUTPUT/loopback). | [DOK: Docker docs] + logika F2 | `curl` dari HP/PC lain di LAN |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan |
| :--- | :--- |
| "UFW aktif = semua port tertutup" | Port yang di-publish Docker lewat FORWARD (F7). |
| "Matikan iptables Docker (`\"iptables\": false`) saja" | Merusak NAT/akses internet kontainer. Docker tidak menganjurkannya. |
| "Uji `curl localhost:8080` gagal = firewall jalan" | Jalur lokal berbeda dari jalur dari luar (F10). |

### Jebakan Lingkungan Operator
* CachyOS punya **ufw dan Docker sama-sama enabled** (cek-agen). Lab 2.2 bisa dilakukan di Inner Loop tanpa menyentuh STB dulu.
* Perubahan firewall di STB/Proxmox via SSH bisa mengunci operator. Wajib: sesi SSH kedua terbuka + perintah rollback berjadwal (misal `sudo systemd-run --on-active=5min ufw disable`) sebelum mengaktifkan aturan.
* Versi Docker berpengaruh pada detail aturan (Docker 28+ banyak mengubah aturan firewall; Docker 29 punya backend nftables eksperimental) — [INGATAN]. Cek `docker info` dan dokumen resmi versi tersebut sebelum menjelaskan detail.

### Ide Chaos
1. **DROP vs REJECT** (satu variabel): blokir port yang sama dengan `DROP` lalu dengan `REJECT --reject-with tcp-reset`; ukur `time curl`.
2. **Docker menembus UFW**: `ufw default deny incoming`, jalankan kontainer `-p 8080:80`, akses dari mesin lain → tetap terbuka. Lacak counter FORWARD/DOCKER-USER.

### Sumber Primer
* https://docs.docker.com/engine/network/packet-filtering-firewalls/
* https://man7.org/linux/man-pages/man8/iptables-extensions.8.html
* https://wiki.nftables.org/
* https://manpages.ubuntu.com/manpages/noble/man8/ufw.8.html

---

## Sub-modul 2.3: Chaos 502 vs 504

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 2.3-F1 | ☐ | Upstream menolak koneksi (RST, `connect() failed (111: Connection refused)`) → nginx mengembalikan **502**. | [DOK: nginx proxy module] + error.log | matikan backend → `curl -v` + `tail error.log` |
| 2.3-F2 | ☐ | Upstream tidak menjawab dalam batas waktu (`proxy_connect_timeout`/`proxy_read_timeout`, default 60s) → **504** (`upstream timed out (110: ...)`). | [DOK: nginx.org ngx_http_proxy_module] | `DROP` paket ke upstream → `curl -v` |
| 2.3-F3 | ☐ | Jadi: firewall `REJECT` ke upstream cenderung → 502, firewall `DROP` → 504. | turunan F1 + F2 + 2.1-F3 | Ide Chaos 1 |
| 2.3-F4 | ☐ | Jika nama host di `proxy_pass`/`upstream` gagal di-resolve **saat nginx memuat konfigurasi**, nginx gagal start/reload (`host not found in upstream`). Itu bukan 504. | [DOK: nginx docs] + [INGATAN untuk perilaku versi baru dengan `resolver`/`resolve`] | ganti nama upstream ke host palsu → `sudo nginx -t` |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan |
| :--- | :--- |
| "DNS gagal → 504" (ditulis di ROADMAP lama) | Dengan nama statis, nginx gagal saat `nginx -t`/reload. Resolusi saat runtime (`resolver` + variabel) bergantung versi dan konfigurasi — buktikan, jangan ditebak. |

### Ide Chaos
1. Upstream mati → 502. 2. `DROP` ke upstream → 504 (setelah ~60s; turunkan `proxy_connect_timeout` agar cepat). 3. Nama upstream salah → `nginx -t` gagal.
Alat: `curl -v`, `ss -tan`, `sudo tcpdump -ni any port <p>`, `tail -f /var/log/nginx/error.log`, `dig`.

### Sumber Primer
* https://nginx.org/en/docs/http/ngx_http_proxy_module.html
* https://nginx.org/en/docs/http/ngx_http_upstream_module.html

---

## Sub-modul 2.4 (Baru): TLS & Sertifikat Dasar

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 2.4-F1 | ☐ | Sertifikat mengikat public key ke nama domain dan ditandatangani CA. Klien memverifikasi rantai sampai ke root CA yang dipercaya. | [DOK: RFC 5280] | `openssl s_client -connect example.com:443 -showcerts </dev/null` |
| 2.4-F2 | ☐ | Dengan Cloudflare Tunnel, TLS pengunjung berakhir di edge Cloudflare; koneksi tunnel ke `cloudflared` terenkripsi terpisah. | [DOK: Cloudflare Tunnel docs] | bandingkan sertifikat yang terlihat browser vs di STB |
| 2.4-F3 | ☐ | Sertifikat punya masa berlaku (`notBefore`/`notAfter`); kedaluwarsa adalah penyebab insiden yang umum. | [DOK: RFC 5280] | `openssl x509 -noout -dates -in cert.pem` |

### Ide Chaos
* Buat CA + sertifikat lokal dengan `openssl`, pasang di nginx Inner Loop → uji sertifikat kedaluwarsa/nama salah → baca error `curl -v`.

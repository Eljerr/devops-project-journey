# Paket Pengetahuan — Fase 1: Kernel Linux, Proses, Filesystem, systemd

> Status: Draf v1 (2026-11-12) · Penyusun: model besar + cek web + cek langsung di mesin operator
> Lingkungan: CachyOS native, kernel 7.3.x, ext4 di `/`, swap = zram (zstd), shell **fish 4.x**, systemd 262.
> Fakta bertanda **(cek-agen)** sudah dijalankan agen di mesin operator pada 2026-11-12, tetapi tetap `☐` sampai operator membuktikannya sendiri.

---

## Sub-modul 1.1: Namespaces, cgroups v2, Sinyal, OOM Killer

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN (perintah → harapan) |
| :--- | :--- | :--- | :--- | :--- |
| 1.1-F1 | ☐ | Linux punya 8 jenis namespace: `mnt`, `pid`, `net`, `uts`, `ipc`, `user`, `cgroup`, `time`. | [DOK: man 7 namespaces] | `ls -l /proc/self/ns` → 8 jenis (+ varian `*_for_children`) |
| 1.1-F2 | ☐ | Dua proses berada di namespace yang sama jika nomor inode symlink `/proc/<PID>/ns/<jenis>` sama. | [DOK: man 7 namespaces] | `readlink /proc/self/ns/pid` di dua terminal → angka sama |
| 1.1-F3 | ☐ | Proses pertama di PID namespace baru melihat dirinya sebagai PID 1, tetapi di host punya PID biasa. | [DOK: man 7 pid_namespaces] | `sudo unshare --pid --fork --mount-proc bash` lalu `ps -ef` → bash = PID 1; dari terminal lain `ps -ef \| grep unshare` → PID besar |
| 1.1-F4 | ☐ | Sistem memakai cgroup v2 (hierarki tunggal) jika `/sys/fs/cgroup` bertipe `cgroup2fs`. | [DOK: docs.kernel.org cgroup-v2] (cek-agen: `cgroup2fs`) | `stat -fc %T /sys/fs/cgroup` → `cgroup2fs` |
| 1.1-F5 | ☐ | `memory.max` = batas keras memori. Melewatinya → kernel mencoba reclaim; jika gagal → OOM killer **di dalam cgroup itu**. | [DOK: cgroup-v2, bagian Memory] | lihat Ide Chaos 1 |
| 1.1-F6 | ☐ | Swap dibatasi terpisah lewat `memory.swap.max`. Memori yang di-swap **tidak** dihitung di `memory.current`. | [DOK: cgroup-v2] | saat lab: bandingkan `memory.current` vs `memory.swap.current` |
| 1.1-F7 | ☐ | `memory.events` berisi penghitung `low`, `high`, `max`, `oom`, `oom_kill`, `oom_group_kill`. Bukti OOM yang sah = `oom_kill` naik. | [DOK: cgroup-v2] | `cat .../memory.events` setelah chaos |
| 1.1-F8 | ☐ | `cpu.max` berformat `"<kuota> <periode>"` dalam mikrodetik. Default `max 100000`. `20000 100000` = 20% satu CPU. | [DOK: cgroup-v2, bagian CPU] | `cat /sys/fs/cgroup/user.slice/cpu.max` |
| 1.1-F9 | ☐ | `SIGKILL` (9) dan `SIGSTOP` tidak bisa ditangkap, diblokir, atau diabaikan. `SIGTERM` (15) bisa ditangkap. | [DOK: man 7 signal] | `bash -c 'trap "echo dapat TERM" TERM; sleep 100 & wait'` lalu kirim `kill -TERM` vs `kill -KILL` |
| 1.1-F10 | ☐ | PID 1 sebuah PID namespace hanya menerima sinyal yang **punya handler**. Pengecualian: `SIGKILL`/`SIGSTOP` dari namespace leluhur tetap dikirim paksa. | [DOK: man 7 pid_namespaces] | Di dalam `unshare` (F3): `kill -TERM 1` → tidak mati; dari host `sudo kill -KILL <pid-host>` → mati |
| 1.1-F11 | ☐ | Jika PID 1 sebuah namespace mati, kernel mengirim `SIGKILL` ke semua proses lain di namespace itu. | [DOK: man 7 pid_namespaces] | Jalankan `sleep 1000 &` di dalam unshare, matikan bash PID 1 → sleep ikut hilang |
| 1.1-F12 | ☐ | Kode keluar `128+N` berarti proses mati oleh sinyal N. 137 = 128+9 (SIGKILL). Di bash dibaca lewat `$?`, di fish lewat `$status`. | [DOK: bash(1) EXIT STATUS] (cek-agen: fish & bash sama-sama 137) | fish: `bash -c 'kill -9 $$'; echo $status` → `137` |
| 1.1-F13 | ☐ | **137 tidak otomatis berarti OOM.** `kill -9` manual juga menghasilkan 137. OOM dibuktikan oleh `oom_kill` di `memory.events` atau log kernel. | [DOK: gabungan F7 + F12] | `journalctl -k \| grep -i "out of memory"` |
| 1.1-F14 | ☐ | zram adalah perangkat swap di RAM yang terkompresi. Bagi cgroup, halaman yang dipindah ke zram dihitung sebagai **swap**, bukan memori. | [DOK: docs.kernel.org admin-guide/blockdev/zram] + [INGATAN untuk detail charging] | `swapon --show` → `/dev/zram0` (cek-agen: ada, 15.5G) |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan | Sumber |
| :--- | :--- | :--- |
| "Kontainer = VM kecil" | Kontainer = proses host biasa + namespace + cgroup. Kernel sama. | man 7 namespaces |
| "PID 1 mengabaikan SIGTERM karena dia root" | Karena aturan kernel khusus init namespace: sinyal tanpa handler tidak dikirim. Tidak ada hubungannya dengan root. | man 7 pid_namespaces |
| "Exit 137 = kehabisan memori" | 137 = mati oleh SIGKILL, dari mana pun. | F13 |
| "zram/kompresi mencegah OOM" | Yang menentukan adalah **apakah swap diizinkan** (`memory.swap.max`). Kompresi hanya memengaruhi berapa banyak RAM fisik yang terpakai zram. | F6, F14 — **wajib dibuktikan dengan Ide Chaos 2** |
| "`bytearray(100MB)` Python tidak memakai RAM karena isinya nol" | [INGATAN] `bytearray(n)` mengisi nol secara eksplisit sehingga halaman benar-benar disentuh. Berbeda dengan `bytes(n)` yang bisa memakai `calloc`. | Buktikan: ukur `memory.current`/`memory.peak` saat keduanya dialokasikan |

### Jebakan Lingkungan Operator
* **fish ≠ bash**: `$?` tidak ada di fish (error). `$$` di fish adalah `$fish_pid`. Loop: `for i in (seq 10); ...; end`. Untuk skrip bash, jalankan `bash -c '...'`.
* **zram + swappiness 150** (cek-agen): CachyOS agresif memindahkan halaman ke zram. Lab OOM **wajib** mengatur `memory.swap.max=0` / `MemorySwapMax=0`, jika tidak proses akan pindah ke swap dan tidak terbunuh.
* **systemd mengelola cgroup**: membuat folder manual di `/sys/fs/cgroup/` boleh untuk belajar, tapi cara yang dianjurkan adalah `systemd-run`. Memory controller sudah didelegasikan ke sesi user (cek-agen: `cpuset cpu io memory pids`), jadi bisa tanpa sudo:
  `systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c '...'`
  [DOK: systemd.io/CGROUP_DELEGATION, systemd-run(1)]

### Ide Chaos (satu variabel per eksperimen)
1. **OOM terukur**: scope `MemoryMax=50M`, `MemorySwapMax=0`, alokasi 100 MB → mati, `$status` 137, `oom_kill` naik, log kernel ada.
2. **Matriks 2×2 untuk membuktikan peran swap vs kompresi** (menguji ulang klaim lama di RB-01):
   | | `MemorySwapMax=0` | `MemorySwapMax=infinity` |
   | :--- | :--- | :--- |
   | Data nol (`bytearray`) | ? | ? |
   | Data acak (`os.urandom`) | ? | ? |
   Operator memprediksi 4 sel dulu, lalu menjalankan. Ubah **satu** variabel per percobaan.
3. **PID 1 tuli**: `unshare` → kirim `SIGTERM` ke PID 1 dari dalam → tidak mati → bandingkan dengan `SIGKILL` dari host.

### Sumber Primer
* https://man7.org/linux/man-pages/man7/namespaces.7.html
* https://man7.org/linux/man-pages/man7/pid_namespaces.7.html
* https://man7.org/linux/man-pages/man7/signal.7.html
* https://docs.kernel.org/admin-guide/cgroup-v2.html
* https://docs.kernel.org/admin-guide/blockdev/zram.html
* https://systemd.io/CGROUP_DELEGATION/

---

## Sub-modul 1.2: File Descriptor, Inode, Disk Penuh

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 1.2-F1 | ☐ | Setiap proses punya tabel FD di `/proc/<PID>/fd/`. 0, 1, 2 = stdin, stdout, stderr. | [DOK: man 5 proc] | `ls -l /proc/self/fd` |
| 1.2-F2 | ☐ | `/proc/<PID>/fdinfo/<FD>` menampilkan `pos` (offset) dan `flags` (oktal). `O_APPEND` = bit `02000`. | [DOK: man 5 proc] (cek-agen: append → `0102001`, tanpa append → `0100001`) | `bash -c 'exec 3>>/tmp/a; cat /proc/$$/fdinfo/3'` |
| 1.2-F3 | ☐ | Nama file (entri direktori) terpisah dari inode (metadata + blok data). Satu inode bisa punya banyak nama (hard link). | [DOK: man 7 inode] | `touch a; ln a b; stat -c '%i %h' a b` → inode sama, link 2 |
| 1.2-F4 | ☐ | `rm` memanggil `unlink()`. Data baru dibebaskan jika jumlah link = 0 **dan** tidak ada FD yang masih membuka file. | [DOK: man 2 unlink] | Ide Chaos 1 |
| 1.2-F5 | ☐ | `df` bertanya ke filesystem lewat `statfs()`. `du` menjumlahkan blok dari file yang **bisa ditemukan lewat direktori**. File terhapus-tapi-terbuka hanya terlihat oleh `df`. | [DOK: df(1), du(1), statfs(2)] | Ide Chaos 1 |
| 1.2-F6 | ☐ | `lsof +L1` menampilkan file terbuka dengan link count < 1 (sudah dihapus). | [DOK: lsof(8)] | `lsof +L1` saat chaos |
| 1.2-F7 | ☐ | Dengan `O_APPEND`, setiap `write()` dimulai di akhir file. | [DOK: man 2 open] | lihat F8 |
| 1.2-F8 | ☐ | Jika file di-*truncate* ke 0 sementara penulis **tanpa** `O_APPEND` masih di offset lama, tulisan berikutnya membuat **sparse file**: ukuran tampak (`ls -l`) besar, tetapi blok di disk (`du`) kecil. Ruang disk **tetap** bebas; yang menipu adalah ukuran tampaknya. | [DOK: man 2 lseek "holes"] | `ls -ls file` dan `du --apparent-size file` vs `du file` |
| 1.2-F9 | ☐ | ext4 punya jumlah inode **tetap** sejak `mkfs`. Inode bisa habis walau ruang masih ada → `No space left on device`. | [DOK: mke2fs(8) opsi -N/-i] | Ide Chaos 2, cek `df -i` |
| 1.2-F10 | ☐ | "File terhapus tapi masih terbuka" dan "inode habis" adalah **dua masalah berbeda**. | [DOK: F4 + F9] | Bandingkan `df -h` vs `df -i` di kedua chaos |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan |
| :--- | :--- |
| "`rm` menghapus data" | `rm` hanya melepas nama. Data dibebaskan setelah link 0 dan FD terakhir ditutup (F4). |
| "Truncate mengubah offset di Open File Table" | **Salah.** Truncate mengubah ukuran **inode**. Offset milik penulis **tidak berubah** — itulah sumber masalahnya (F8). |
| "Setelah truncate tanpa O_APPEND, disk penuh lagi" | Bagian awal menjadi *hole* yang tidak memakan blok. Disk tidak penuh lagi; hanya ukuran tampaknya besar (F8). |
| "Deleted open files = inode exhaustion" | Berbeda (F10). |

### Jebakan Lingkungan Operator
* `/` = ext4, **`/tmp` = `tmpfs`** (cek-agen 2026-11-12, `findmnt -no FSTYPE /tmp`). Lab disk penuh di `tmpfs` memakai RAM/swap, bukan disk! Jalankan lab deleted-open-file di direktori ext4, misal `~/lab/`. (Catatan: PM-01 lama menunjukkan perangkat `259,9` untuk `/tmp` — kemungkinan konfigurasi `/tmp` berubah sejak itu; cek ulang.)
* fish: `> file` tanpa perintah tidak sama dengan bash. Pakai `truncate -s 0 <path>`.

### Ide Chaos
1. **Deleted open file**: proses bash memegang FD ke file 500 MB di ext4 (`~/lab/`, bukan `/tmp`) → `rm` → bandingkan `df`/`du` → `lsof +L1` → bebaskan via `truncate -s 0 /proc/<PID>/fd/<FD>` → cek `fdinfo` (pos) dan `ls -ls`.
2. **Inode exhaustion sungguhan** (aman, pakai loop device):
   `truncate -s 20M /tmp/img; mkfs.ext4 -N 64 /tmp/img; sudo mount -o loop /tmp/img /mnt/lab` → buat 100 file kosong → `df -h` vs `df -i`.

### Sumber Primer
* https://man7.org/linux/man-pages/man2/unlink.2.html
* https://man7.org/linux/man-pages/man2/open.2.html
* https://man7.org/linux/man-pages/man2/lseek.2.html
* https://man7.org/linux/man-pages/man5/proc.5.html
* https://man7.org/linux/man-pages/man7/inode.7.html

---

## Sub-modul 1.3: systemd & SSH Hardening

### Fakta Inti
| # | Status | Klaim | Label/Sumber | BUKTIKAN |
| :--- | :--- | :--- | :--- | :--- |
| 1.3-F1 | ☐ | Setiap unit service berjalan di cgroup sendiri, misal `/sys/fs/cgroup/system.slice/<unit>.service/`. | [DOK: systemd.resource-control(5)] | `systemctl show -p ControlGroup <unit>` |
| 1.3-F2 | ☐ | `MemoryMax=` → `memory.max`; `MemorySwapMax=` → `memory.swap.max`; `CPUQuota=20%` → `cpu.max` = `20000 100000`. | [DOK: systemd.resource-control(5)] | set di unit → `cat` file cgroup-nya |
| 1.3-F3 | ☐ | Default batas start: 5 kali dalam 10 detik (`StartLimitBurst`, `StartLimitIntervalSec`). Melewati → unit gagal dengan hasil `start-limit-hit`. | [DOK: systemd.unit(5)] (cek-agen: 5 / 10s) | `systemctl show -p DefaultStartLimitBurst -p DefaultStartLimitIntervalUSec` |
| 1.3-F4 | ☐ | Setelah mengubah file unit, wajib `systemctl daemon-reload`. | [DOK: systemctl(1)] | ubah unit tanpa reload → `systemctl status` memberi peringatan |
| 1.3-F5 | ☐ | Di `sshd_config`, untuk tiap kata kunci, **nilai pertama yang ditemukan dipakai**. `Include sshd_config.d/*.conf` biasanya ada di atas, jadi drop-in menang atas baris di bawahnya. | [DOK: sshd_config(5)] | `sudo sshd -T \| grep -i passwordauthentication` |
| 1.3-F6 | ☐ | `sshd -t` memeriksa sintaks; `sshd -T` menampilkan konfigurasi **efektif**. Selalu jalankan sebelum reload. | [DOK: sshd(8)] | `sudo sshd -t && echo OK` |
| 1.3-F7 | ☐ | Klien SSH menolak private key yang izinnya terlalu terbuka (`UNPROTECTED PRIVATE KEY FILE`). Server (StrictModes) memeriksa izin `~/.ssh` dan `authorized_keys`. | [DOK: ssh(1), sshd_config(5) StrictModes] | `chmod 644 ~/.ssh/id_ed25519` (salinan uji!) → coba login → kembalikan 600 |
| 1.3-F8 | ☐ | Nama service berbeda: Ubuntu `ssh`, Arch/CachyOS `sshd`. Ubuntu baru memakai socket activation (`ssh.socket`). | [DOK: Ubuntu Server docs] + [INGATAN untuk versi pastinya] | di STB: `systemctl status ssh.socket ssh.service` |

### Miskonsepsi Umum
| Miskonsepsi | Kenyataan |
| :--- | :--- |
| "Tambahkan `PasswordAuthentication no` di bawah sshd_config = aman" | Bisa kalah oleh drop-in (misal `50-cloud-init.conf`) karena aturan nilai pertama (F5). Verifikasi dengan `sshd -T`. |
| "Rollback cukup `cp sshd_config.bak`" | Hanya jika backup **dibuat sebelumnya**. Runbook wajib punya langkah backup. |

### Ide Chaos
1. Service dummy yang crash terus → amati `start-limit-hit` di `systemctl status` dan `journalctl -u <unit>`.
2. Service dengan `MemoryMax` kecil → buktikan pemetaan F2 dengan `cat` file cgroup.
3. Di STB: buat drop-in `00-test.conf` berisi `PasswordAuthentication yes` → buktikan baris `no` di bawah kalah (`sshd -T`). **Jangan tutup sesi SSH aktif.**

### Sumber Primer
* https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html
* https://www.freedesktop.org/software/systemd/man/latest/systemd.unit.html
* https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
* https://man.openbsd.org/sshd_config
* https://man.openbsd.org/sshd

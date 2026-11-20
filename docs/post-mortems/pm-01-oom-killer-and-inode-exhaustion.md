---
incident-id: PM-0001
date: 2026-11-19
module-phase: "Fase 1.1 & 1.2"
environment: "INNER CachyOS"
author: Asylon
severity: "High"
status: "Resolved"
verification: "Semua klaim LAB/DOK"
evidence: ["docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md", "docs/labs/lab-1.2-01-fd-inode-disk-full.md"]
tags: [post-mortem, troubleshooting, fase-1, cgroups-v2, oom-killer, zram, deleted-open-files, inode-exhaustion, sparse-file]
---

# Post-Mortem: Investigasi Terminasi Proses OOM Killer (cgroups v2 & ZRAM) dan Anomali Disk Penuh (*Deleted Open Files*, *Inode Exhaustion*, & *Sparse File*)

---

## Bagian 1: Investigasi Terminasi Proses Akibat Batasan `cgroups v2`, Lazy Allocation, dan Swap ZRAM (`[LAB ✓]`)

### 1. Ringkasan & Dampak
<!-- OPERATOR VOICE: ceritakan dengan bahasamu sendiri -->
Eksperimen *chaos* terukur dilakukan di workstation **CachyOS Linux native** (kernel `7.3.3-1-cachyos`, `cgroup2fs`, swap `/dev/zram0` `15.5G`, shell `fish`) menggunakan `systemd-run --user --scope` untuk menguji kapan sebuah proses yang meminta memori 100 MB di dalam pagar `MemoryMax=50M` benar-benar dieksekusi oleh **OOM Killer**, kapan ia lolos karena *lazy allocation* (`calloc()`), dan kapan ia selamat karena dipindahkan ke swap ZRAM.

### 2. Timeline
| Waktu | Kejadian |
| :--- | :--- |
| 09:29 | Pemeriksaan awal `/sys/fs/cgroup/user.slice/memory.events` mencatat `max 35`, `oom 1`, `oom_kill 1`. |
| 09:37 | Hipotesis 1 diuji: alokasi `bytearray(100 MB)` dengan `MemoryMax=50M` dan `MemorySwapMax=0` justru keluar dengan status `0` (selamat, `oom_kill` tetap `1`). |
| 09:41 | Akar masalah lolosnya `bytearray` ditemukan (*lazy allocation* halaman nol oleh `calloc()`). Beban diubah menjadi penulisan fisik `"a" * 100 MB` → proses langsung dimatikan (`$status` = `137`, `max` melonjak ke `91`, `oom_kill` naik menjadi `2`). |
| 08:31 (+1d) | Matriks 2×2 diuji pada `"a" * 100 MB` vs `os.urandom(100 MB)` dengan dan tanpa `MemorySwapMax=0` untuk membuktikan peran izin swap terhadap ZRAM. |

### 3. Gejala Awal & Bukti Log
<!-- RAW ARTIFACT: tempel output asli tanpa diedit -->
```text
 ~/Projects/devops-project-journey  main*
❯ swapon --show
NAME       TYPE       SIZE USED PRIO
/dev/zram0 partition 15.5G 1.2M  100

 ~/Projects/devops-project-journey  main*
❯ systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c 'x = bytearray(100 * 1024 * 1024)'; echo $status
Running as unit: run-p148196-i148273.scope; invocation ID: c90c1e5dd39e4eabb133eda069e62b17
0

 ~/Projects/devops-project-journey  main*
❯ systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c 'x = "a" * (100 * 1024 * 1024)'; echo $status
Running as unit: run-p156657-i152894.scope; invocation ID: 77f4b1b6cb2e43d4ae5b50e5031deba0
137
fish: Job 1, 'systemd-run --user --scope -p M…' terminated by signal SIGKILL (Forced quit)

 ~/Projects/devops-project-journey  main*
❯ cat /sys/fs/cgroup/user.slice/memory.events
low 0
high 0
max 91
oom 2
oom_kill 2
oom_group_kill 0
sock_throttled 0
```

Bukti Matriks 2×2 (Peran `MemorySwapMax` vs Kompresi ZRAM):
```text
❯ systemd-run --user --scope -p MemoryMax=50M python3 -c 'x = "a" * (100 * 1024 * 1024)'; echo $status
Running as unit: run-p31562-i39429.scope; invocation ID: 4b61a0bd7de94336a8170cd7087d7592
0

 ~/Projects/devops-project-journey  main*
❯ systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c 'import os; x = os.urandom(100 * 1024 * 1024)'; echo $status
Running as unit: run-p31643-i22916.scope; invocation ID: 6f22d888072a4fd59e7f21a7bea1fa29
137
fish: Job 1, 'systemd-run --user --scope -p M…' terminated by signal SIGKILL (Forced quit)

 ~/Projects/devops-project-journey  main*
❯ systemd-run --user --scope -p MemoryMax=50M python3 -c 'import os; x = os.urandom(100 * 1024 * 1024)'; echo $status
Running as unit: run-p32942-i10749.scope; invocation ID: d795861f119e4b98a3ed8a8d48d08a52
0
```

### 4. Kronologi Investigasi
<!-- OPERATOR VOICE: termasuk tebakan yang salah -->
1. **Hipotesis awal (`bytearray` 100 MB)**: Mengalokasikan `x = bytearray(100 * 1024 * 1024)` di dalam scope `MemoryMax=50M` dan `MemorySwapMax=0` diprediksi akan langsung mati (`137`). → Hasil: keluar `0` dan *"lah kok OOM-nya ga nambah ya"* (`oom_kill` tetap `1`).
2. **Temuan menengah (*Lazy Allocation* `calloc()`)**: Pada kasus angka nol, Python hanya menjanjikan tempat untuk memory namun tidak benar-benar menuliskan datanya pada memory dengan menggunakan fungsi `calloc()`. Sebelum memberikan blok fisik RAM pada program baru, kernel Linux akan membersihkan terlebih dahulu blok yang tadinya dipakai dengan angka 0 biar program baru tidak bisa mengintip sisa data pada program lama. Jadi pada dasarnya menuliskan angka nol itu tidak benar-benar menuliskan data ke blok fisik, karena dari awal nilainya sudah nol.
3. **Titik kegagalan & Pembuktian Matriks 2×2**: Saat data benar-benar ditulis ke blok fisik (`"a" * 100MB` sebagai data mudah dikompresi, dan `os.urandom(100MB)` sebagai data acak anti-kompresi), keduanya **mati (`137`)** ketika `MemorySwapMax=0` (Sel 1A & 2A), dan keduanya **selamat (`0`)** ketika `MemorySwapMax` tidak dikunci (Sel 1B & 2B), karena kelebihan memory-nya masih bisa dipindahkan ke swap.

### 5. Akar Masalah (5 Whys)
<!-- OPERATOR VOICE. Setiap "Why" yang menyebut mekanisme kernel diberi label [LAB ✓]/[DOK]/[INGATAN] -->
* **Gejala**: Proses Python `run-p156657-i152894.scope` dan `run-p31643-i22916.scope` berhenti mendadak dengan `$status` = `137` (`terminated by signal SIGKILL`).
  * **Why 1**: Mengapa keluar kode `137`? Karena proses menerima sinyal paksa `SIGKILL` (sinyal nomor 9 → `128 + 9 = 137`). [DOK: `man 7 signal`] [LAB ✓]
  * **Why 2**: Siapa yang mengirim `SIGKILL` tersebut (bukan `kill -9` manual)? Kita buktikan dengan memeriksa `memory.events` dari cgroup atau melihat `journalctl`, di mana nilai `oom_kill` naik dari `1` menjadi `2`. [DOK: `cgroup-v2.html`] [LAB ✓]
  * **Why 3**: Mengapa OOM Killer terpicu? Karena proses menulis 100 MB data non-nol langsung ke halaman fisik RAM melampaui pagar `memory.max` (`52428800` Byte / 50 MB), dan usaha *reclaim* kernel (tercatat `max` melonjak dari `47` ke `91`) gagal karena pintu swap dikunci dengan `MemorySwapMax=0`. [DOK: `cgroup-v2.html`] [LAB ✓]
  * **Mekanisme internal**: Pada cgroups v2, memori RAM (`memory.max`) dan swap (`memory.swap.max`) dihitung di dua kantong terpisah. Halaman memori yang dipindahkan ke `/dev/zram0` ditagihkan ke kantong swap cgroup (`memory.swap.current`), sehingga tanpa `MemorySwapMax=0`, bahkan data acak `os.urandom(100 MB)` yang tidak dapat dikompresi ZRAM tetap lolos dari OOM Killer selama kapasitas swap masih tersedia. [DOK: `cgroup-v2.html`, `blockdev/zram.html`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]

---

## Bagian 2: Investigasi Tiga "Ilusi Disk" — *Deleted Open Files*, *Inode Exhaustion*, & *Sparse File* (`[LAB ✓]`)

### 1. Ringkasan & Dampak
<!-- OPERATOR VOICE: ceritakan dengan bahasamu sendiri -->
Tiga skenario anomali ruang penyimpanan direkonstruksi langsung pada filesystem **`ext4`** (`~/lab` di `/dev/nvme0n1p6` dan *loop device* `/dev/loop0`) untuk membedakan secara presisi tiga fenomena "ilusi disk":
1. **Deleted Open Files**: File log `500 MB` (`~/lab/hantu.log`) sudah dihapus dengan `rm` (`du -sh` turun dari `501M` ke `4.0K`), tetapi kapasitas disk di `df -BM` tetap tertahan di `52156M` karena Inode-nya (`4849708`) masih dipegang oleh proses pada File Descriptor 3 (`3w`).
2. **Inode Exhaustion**: Pembuatan file kosong (`0 Byte`) gagal dengan `No space left on device` meskipun ruang blok `df -h` baru terpakai `1%` (`48K` dari `19M`), karena stok Inode pada `df -i` sudah terpakai `100%` (`72/72`).
3. **Sparse File akibat `truncate -s 0` tanpa `O_APPEND`**: Memangkas file `100 MB` yang sedang dibuka dengan satu panah `>` (`flags: 0100001`) mengosongkan blok fisik disk ke `0`, tetapi kursor `pos` di `fdinfo/3` tetap nyangkut di `104857600`, sehingga penulisan `"halo\n"` (5 Byte) berikutnya menghasilkan ukuran tampak `101M` (`du -h --apparent-size`) padahal isi fisik di disk (`du -h`) hanya `4.0K`.

### 2. Timeline
| Waktu | Kejadian |
| :--- | :--- |
| 07:10 | Kondisi awal `/dev/nvme0n1p6` diverifikasi: `~/lab` berukuran `4.0K`. |
| 07:15 | File `~/lab/hantu.log` (`500 MB`) ditahan dengan `exec 3>` lalu dihapus dengan `rm` → `du -sh` melaporkan `4.0K`, tetapi `df -BM` tetap `52156M`. `lsof +L1 \| grep lab` menangkap proses `sleep` (PID `74329` & `78069`) yang menahan Inode `4849708` (`NLINK = 0`). |
| 07:25 | Proses `74329` dan `78069` dihentikan dengan `kill` → `df -BM` langsung turun 500 MB menjadi `51656M`. |
| 07:36 | Partisi loop `ext4` 20 MB dengan `72` Inode diuji dengan 100 file kosong → gagal mulai `file_61.txt` (`No space left on device`, `df -h` `1%` vs `df -i` `100%`). |
| 09:03 | Uji `truncate -s 0` pada proses penulis 100 MB tanpa `O_APPEND` membuktikan `pos` tetap di `104857600` → `104857605`, fisik `4.0K` vs tampak `101M`. |

### 3. Gejala Awal & Bukti Log
<!-- RAW ARTIFACT: tempel output asli tanpa diedit -->

#### A. Bukti *Deleted Open Files* (`df -BM` vs `du -sh` & `lsof +L1`)
```text
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 52156M    61586M  46% /
501M    /home/joy/lab
# Sesudah rm ~/lab/hantu.log:
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 52156M    61586M  46% /
4.0K    /home/joy/lab
sleep     74329  joy   3w   REG  259,9 524288000     0 4849708 /home/joy/lab/hantu.log (deleted)
sleep     78069  joy   3w   REG  259,9 524288000     0 4849708 /home/joy/lab/hantu.log (deleted)

❯ kill 74329 78069; df -BM ~/lab
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 51656M    62086M  46% /
```

#### B. Bukti *Inode Exhaustion* (`df -h` vs `df -i`)
```text
touch: cannot touch '/home/joy/lab/mnt_inode/file_61.txt': No space left on device
...
Filesystem      Size  Used Avail Use% Mounted on
/dev/loop0       19M   48K   18M   1% /home/joy/lab/mnt_inode
Filesystem     Inodes IUsed IFree IUse% Mounted on
/dev/loop0         72    72     0  100% /home/joy/lab/mnt_inode
```

#### C. Bukti `truncate -s 0` & *Sparse File* (`du -h` vs `du -h --apparent-size`)
```text
--- 1. Sebelum truncate ---
pos:    104857600
100M    /home/joy/lab/sparse.log
100M    /home/joy/lab/sparse.log
--- 2. Sesudah truncate -s 0 (sebelum tulis lagi) ---
pos:    104857600
0       /home/joy/lab/sparse.log
0       /home/joy/lab/sparse.log
--- 3. Sesudah tulis halo (5 Byte) tanpa O_APPEND ---
pos:    104857605
4.0K    /home/joy/lab/sparse.log
101M    /home/joy/lab/sparse.log
```

### 4. Kronologi Investigasi & Analisis Operator
<!-- OPERATOR VOICE: termasuk tebakan yang salah -->
1. **Hipotesis awal `rm` pada file terbuka**: Awalnya diprediksi setelah `rm ~/lab/hantu.log`, `du` masih `500M` dan `df` turun `500M`. Ternyata hasil di terminal berkebalikan: `du` hanya memindai bagian luarnya doang yang masih ada Dentry-nya (sehingga langsung turun ke `4.0K`), sedangkan `df` memeriksa langsung ke pembukuan Inode-nya yang masih dipegang oleh proses (`52156M`).
2. **Pelacakan dengan `lsof +L1`**: `lsof` memindai file yang masih dipegang oleh proses sistem dan flag `+L1` menyaring file yang `Links`-nya kurang dari 1 (`0`). Setelah PID `74329` dan `78069` dimatikan dengan `kill`, ruang `500 MB` langsung kembali (`51656M`).
3. **Perbandingan dengan *Inode Exhaustion***: Keduanya sama-sama bikin disk ilusi: *Deleted Open Files* itu Inode-nya masih dipegang oleh proses, sedangkan *Inode Exhaustion* adalah kondisi di mana disk kehabisan nomor Inode (`df -i` 100% walaupun `df -h` baru 1%).
4. **Efek `truncate -s 0` tanpa `O_APPEND`**: Kalau proses membuka file pakai 1 panah (`>`), setelah `truncate -s 0` posisi kursor `pos`-nya tidak balik ke 0 (tetap di `104857600`). Begitu ditulis `"halo"` (5 Byte), tergantung bagaimana sistem melihatnya: kalau dilihat dari alokasi tampaknya (`du -h --apparent-size`) masih terlihat penuh (`101M`), tapi kalau sistem melihat benar-benar isi fisiknya (`du -h`) hanya `4.0K` saja (*sparse file*).

### 5. Action Items
| # | Tindakan | Jenis (cegah/deteksi/mitigasi) | Status |
| :--- | :--- | :--- | :--- |
| 1 | Gunakan `df -h` bersama `df -i` saat menerima error `No space left on device` untuk segera membedakan penuhnya blok data vs habisnya stok Inode. | Deteksi | ☑ |
| 2 | Gunakan `lsof +L1` ketika nilai `df` jauh lebih tinggi daripada `du` akibat file log yang dihapus dengan `rm` tanpa menutup proses. | Deteksi | ☑ |
| 3 | Gunakan `truncate -s 0` (atau `logrotate` / sinyal re-open log) ketimbang `rm` pada file log yang sedang aktif, dan pahami bahwa tanpa `O_APPEND`, ukuran tampak (`ls -lh`) bisa menjadi *sparse file* meskipun blok disk fisik (`du -h` / `df -h`) sudah bebas. | Mitigasi | ☑ |

### 6. Refleksi Operator
<!-- OPERATOR VOICE: apa yang kamu pelajari, apa yang akan kamu lakukan berbeda -->
* Dentry hanyalah label yang menunjuk ke nomor Inode, sedangkan Inode menyimpan seluruh informasi dan blok data file. Menghapus file dengan `rm` hanya menghapus Dentry (`Links` berkurang), belum membebaskan blok data selama masih ada File Descriptor yang memegangnya.
* `du` memeriksa dari luar lewat Dentry yang terlihat di folder, sedangkan `df` memeriksa pemakaian blok/Inode langsung ke filesystem.
* Ketika memotong file log aktif dengan `truncate -s 0` pada proses yang tidak memakai `O_APPEND`, kursor `pos` di `/proc/<PID>/fdinfo/<FD>` tidak kembali ke `0`. Penulisan berikutnya membentuk *sparse file* yang ukuran tampaknya besar (`101M`), tetapi pemakaian blok fisik di disk sebenarnya sudah terbebaskan (`4.0K`).

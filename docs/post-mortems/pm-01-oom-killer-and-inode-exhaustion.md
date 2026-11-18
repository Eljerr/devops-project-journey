---
incident-id: PM-0001
date: 2026-11-18
module-phase: "Fase 1.1 & 1.2"
environment: "INNER CachyOS"
author: Asylon
severity: "High"
status: "In Progress (Bagian 1 Terverifikasi LAB ✓ · Bagian 2 Menunggu Sub-modul 1.2)"
verification: "Bagian 1: Semua klaim LAB/DOK | Bagian 2: Menunggu Re-verifikasi Sub-modul 1.2"
evidence: ["docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md"]
tags: [post-mortem, troubleshooting, fase-1, cgroups-v2, oom-killer, zram, deleted-open-files]
---

# Post-Mortem: Investigasi Terminasi Proses OOM Killer (cgroups v2 & ZRAM) dan Anomali Disk Penuh

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
| 08:31 (+1h) | Matriks 2×2 diuji pada `"a" * 100 MB` vs `os.urandom(100 MB)` dengan dan tanpa `MemorySwapMax=0` untuk membuktikan peran izin swap terhadap ZRAM. |

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

### 6. Resolusi & Bukti Pemulihan
* **Perbaikan**: Memastikan setiap pengujian batas memori keras atau *deployment* layanan yang tidak menginginkan pelarian memori ke ZRAM menyertakan konfigurasi eksplisit `MemorySwapMax=0` (atau batas swap yang terukur), serta memverifikasi insiden OOM melalui `memory.events` (`oom_kill`), bukan sekadar menebak dari kode keluar `137`.

### 7. Action Items
| # | Tindakan | Jenis (cegah/deteksi/mitigasi) | Status |
| :--- | :--- | :--- | :--- |
| 1 | Selalu periksa `memory.events` (`oom_kill`) atau `journalctl` saat mendapati exit code `137` agar tidak tertukar dengan `SIGKILL` manual/timeout. | Deteksi | ☑ |
| 2 | Pasang `MemorySwapMax` secara eksplisit bersama `MemoryMax` pada unit `systemd` (Sub-modul 1.3) di host yang memiliki ZRAM aktif. | Cegah | ☐ |

### 8. Refleksi Operator
<!-- OPERATOR VOICE: apa yang kamu pelajari, apa yang akan kamu lakukan berbeda -->
* `SIGTERM` itu sinyal yang dikirim untuk mematikan sebuah proses secara *safe* dan `SIGKILL` mematikan sebuah proses secara paksa (`137`). Untuk memastikan `137` berasal dari OOM Killer, kita wajib memeriksa `memory.events` dari sebuah proses atau melihat `journalctl` dari aplikasi/service tersebut.
* Pada kasus alokasi angka nol (`bytearray`), Python hanya menjanjikan tempat dengan `calloc()` tanpa benar-benar menulis ke blok fisik RAM karena kernel Linux sudah membersihkan blok bekas pakai dengan angka 0 demi keamanan.
* Pada sistem dengan swap ZRAM, proses yang melewati `MemoryMax=50M` (baik data berulang maupun `os.urandom` acak) selamat ketika `MemorySwapMax` tidak dikunci semata-mata karena memory-nya masih bisa dipindahkan ke swap.

---

## Bagian 2: Investigasi Anomali Ruang Disk Akibat Deleted Open Files & Inode Exhaustion

> [!WARNING]
> **MENUNGGU SUB-MODUL 1.2** — Bagian 2 di bawah ini masih berasal dari draf v1 dan akan segera disusun ulang secara otomatis menggunakan bukti terminal asli CachyOS (`~/lab/` di filesystem `ext4`) begitu Misi Kecil 1.2a s/d 1.2d selesai.

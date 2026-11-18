---
topic: "Control Groups (cgroups v2), OOM Killer, & ZRAM Swap"
module-phase: "Fase 1.1"
date: 2026-11-18
author: Asylon
verification: "Semua klaim LAB/DOK"
related-docs: ["[[lab-1.1-01-namespaces-cgroups-signals-oom]]", "[[pm-01-oom-killer-and-inode-exhaustion]]"]
tags: [concept, first-principles, interview-prep, cgroups-v2, oom-killer, zram]
---

# Konsep: Control Groups (`cgroups v2`), OOM Killer, & ZRAM Swap

## 1. Analogi Mental (Momen "Aha!")
<!-- OPERATOR VOICE: analogi versimu sendiri, bisa diingat dalam 10 detik -->
Kalau Namespace itu ruangan yang dibuat untuk mengisolasi proses dari luar, maka `cgroups` merupakan aturan kuota fisik (seperti batas penggunaan CPU dan batas penggunaan memory) yang dibuat untuk proses di dalam ruangan tersebut.

**Batas analogi** — di mana analogi ini tidak lagi akurat:
<!-- OPERATOR VOICE -->
Memasang batas `MemoryMax=50M` saja tidak otomatis membuat proses yang meminta 100 MB langsung dibunuh oleh OOM Killer:
1. Jika proses hanya memesan memori berisi angka nol (`bytearray`), kernel menggunakan *lazy allocation* (`calloc`) tanpa langsung menulis ke blok fisik RAM.
2. Jika pintu swap (`MemorySwapMax`) tidak ikut dikunci ke `0`, kelebihan memori masih bisa dipindahkan oleh kernel ke kantong swap (`/dev/zram0`) sehingga proses tetap selamat.

## 2. Mekanisme Internal
<!-- Setiap klaim diberi label: [LAB ✓ → lab-id] / [DOK: sumber] / [INGATAN] -->
* Sistem memakai **cgroups v2** (*unified hierarchy*) jika `/sys/fs/cgroup` bertipe `cgroup2fs`. Seluruh pengendali sumber daya (`cpu`, `memory`, `io`, `pids`) disusun dalam satu pohon direktori terpadu. [DOK: `docs.kernel.org/admin-guide/cgroup-v2.html`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* File `cpu.max` berformat `"<kuota> <periode>"` dalam satuan mikrodetik. Nilai `20000 100000` berarti proses diizinkan memakai CPU maksimal `20.000` mikrodetik (20%) setiap periode `100.000` mikrodetik (`0,1` detik). [DOK: `cgroup-v2.html` bagian CPU] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* File `memory.max` menyimpan batas keras RAM fisik dalam satuan Byte (contoh: `MemoryMax=50M` ditulis kernel sebagai `52428800`). Pemakaian swap diatur secara terpisah di `memory.swap.max`. [DOK: `cgroup-v2.html` bagian Memory] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Saat pemakaian RAM melampaui `memory.max`, kernel berusaha melakukan *reclaim* / memindahkan halaman ke swap (penghitung `max` di `memory.events` bertambah). Jika swap habis atau dikunci (`MemorySwapMax=0`), kernel memanggil **OOM Killer** di dalam cgroup tersebut, mengirim `SIGKILL` (9), dan menaikkan penghitung `oom` serta `oom_kill` di `memory.events`. [DOK: `cgroup-v2.html`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Bagi cgroup, halaman yang dipindahkan ke perangkat ZRAM (`/dev/zram0`) dihitung sebagai pemakaian **swap** (`memory.swap.current`), bukan memori RAM (`memory.current`). Oleh karena itu, faktor penentu apakah alokasi 100 MB selamat di dalam kandang `MemoryMax=50M` adalah apakah swap diizinkan (`MemorySwapMax`), bukan apakah datanya bisa dikompresi oleh ZRAM. [DOK: `cgroup-v2.html`, `blockdev/zram.html`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]

## 3. Pembuktian di Terminal
```fish
# 1. Buat scope dengan batas CPU 20% dan RAM 50 MB
systemd-run --user --scope -p MemoryMax=50M -p CPUQuota=20% bash
cat /proc/self/cgroup
# Output: 0::/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope
cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope/cpu.max
# Output: 20000 100000
cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope/memory.max
# Output: 52428800

# 2. Memicu OOM Killer (MemoryMax=50M + MemorySwapMax=0 + tulis 100 MB fisik)
systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c 'x = "a" * (100 * 1024 * 1024)'; echo $status
# Output: 137 (oom_kill di /sys/fs/cgroup/user.slice/memory.events naik dari 1 menjadi 2)
```
Bukti lengkap: [lab-1.1-01-namespaces-cgroups-signals-oom.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md)

## 4. Miskonsepsi Awal Saya
<!-- OPERATOR VOICE: apa yang dulu kamu kira benar, dan apa yang membuktikan itu salah -->
* **Dulu saya kira**: Menjalankan `x = bytearray(100 * 1024 * 1024)` di dalam kandang `MemoryMax=50M` dan `MemorySwapMax=0` akan langsung memicu OOM Killer. → **Ternyata**: Proses tersebut keluar dengan kode `0` (selamat) dan `oom_kill` tidak bertambah! Pada kasus angka nol, Python hanya menjanjikan tempat untuk memory namun tidak benar-benar menuliskan datanya pada memory dengan menggunakan fungsi `calloc()`. Sebelum memberikan blok fisik RAM pada program baru, kernel Linux membersihkan terlebih dahulu blok yang tadinya dipakai dengan angka 0 biar program baru tidak bisa mengintip sisa data pada program lama. Jadi pada dasarnya menuliskan angka nol itu tidak benar-benar menuliskan data ke blok fisik, karena dari awal nilainya sudah nol.
* **Dulu saya kira**: Di sistem dengan ZRAM (`zstd`), proses yang memakan 100 MB di dalam cgroup 50 MB bisa selamat semata-mata karena datanya berhasil dikompresi oleh ZRAM menjadi di bawah 50 MB. → **Ternyata**: Saat diuji dengan matriks 2×2 antara data mudah dikompresi (`"a" * 100MB`) dan data acak anti-kompresi (`os.urandom(100MB)`), keduanya sama-sama mati (`137`) ketika `MemorySwapMax=0`, dan keduanya sama-sama selamat (`0`) ketika `MemorySwapMax` tidak dikunci, karena kelebihan memory-nya masih bisa dipindahkan ke swap.

## 5. Jebakan Produksi & Pertanyaan Interview
* **Pertanyaan**: Jika sebuah service atau container mati dengan kode keluar `137`, bagaimana kamu memastikan apakah itu benar-benar akibat OOM Killer, dan mengapa di mesin dengan ZRAM proses bisa lolos dari batas `MemoryMax` jika konfigurasi tidak lengkap?
* **Jawaban saya**: <!-- OPERATOR VOICE --> Kode `137` hanya menandakan proses mati terkena sinyal `SIGKILL` (9) dari mana saja (`128 + 9`). Untuk membuktikan bahwa pelakunya adalah OOM Killer, kita harus memeriksa counter `oom_kill` pada file `memory.events` milik cgroup proses tersebut atau memeriksa log kernel lewat `journalctl`. Selain itu, jika kita hanya membatasi `MemoryMax` tanpa mengunci `MemorySwapMax=0`, proses yang melewati batas RAM masih bisa memindahkan kelebihannya ke swap (ZRAM) sehingga tidak terbunuh oleh OOM Killer.

## 6. Sumber Primer
* Kernel Documentation cgroup v2 — https://docs.kernel.org/admin-guide/cgroup-v2.html
* Kernel Documentation ZRAM — https://docs.kernel.org/admin-guide/blockdev/zram.html
* `man 5 systemd.resource-control` — https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html
* `man 3 calloc` — https://man7.org/linux/man-pages/man3/malloc.3.html

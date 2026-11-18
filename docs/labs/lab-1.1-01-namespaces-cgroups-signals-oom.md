---
lab-id: LAB-1.1-01
date: 2026-11-18
module-phase: "Fase 1.1"
environment: "INNER CachyOS"
author: Asylon
hint-level-max: "0"
tags: [lab, evidence, fase-1, namespaces, cgroups-v2, signals, oom-killer, zram]
---

# Lab: Isolasi Namespaces, Pembatasan cgroups v2, Sinyal PID 1, dan Matriks OOM Killer vs Swap ZRAM

> Berkas ini adalah **bukti mentah** dari sesi terminal asli di CachyOS Linux (shell `fish`). Dokumen lain (`docs/concepts/` dan `docs/post-mortems/pm-01`) mengutip langsung dari berkas ini.

## 1. Pertanyaan yang Ingin Dibuktikan
<!-- OPERATOR VOICE -->
1. Bagaimana `unshare` mengisolasi ruang penomoran proses (PID namespace) dari host?
2. Bagaimana `systemd-run` menerapkan batasan `MemoryMax` dan `CPUQuota` ke file kernel `/sys/fs/cgroup`?
3. Mengapa proses `PID 1` di dalam namespace baru kebal terhadap `SIGTERM` dan `SIGKILL` dari dalam ruangannya sendiri?
4. Mengapa alokasi memori berisi angka nol (`bytearray`) lolos dari OOM Killer, dan apakah proses bermuatan 100 MB di dalam pagar `MemoryMax=50M` selamat karena kompresi ZRAM atau karena pintu Swap (`MemorySwapMax`) dibuka?

## 2. Kondisi Awal (Versi & Lingkungan)
```text
 ~/Projects/devops-project-journey  main*
❯ stat -fc %T /sys/fs/cgroup
cgroup2fs

 ~/Projects/devops-project-journey  main*
❯ cat /sys/fs/cgroup/user.slice/cpu.max
max 100000

 ~/Projects/devops-project-journey  main*
❯ swapon --show
NAME       TYPE       SIZE USED PRIO
/dev/zram0 partition 15.5G 1.2M  100
```

## 3. Variabel
* **Eksperimen 1.1a (Namespaces)**: Menjalankan `bash` di luar vs di dalam `sudo unshare --pid --fork --mount-proc bash`.
* **Eksperimen 1.1b (cgroups v2)**: Mengubah properti scope dengan `-p MemoryMax=50M -p CPUQuota=20%`.
* **Eksperimen 1.1c (Sinyal & PID 1)**: Mengirim `kill -9 $$` di proses biasa vs `kill -TERM 1` dan `kill -KILL 1` pada `PID 1` di dalam PID namespace.
* **Eksperimen 1.1d & Chaos (Matriks OOM Killer)**: Mengubah **satu variabel per langkah** antara jenis data (`bytearray` nol vs `"a" * 100MB` mudah dikompresi vs `os.urandom(100MB)` acak anti-kompresi) dan batas swap (`MemorySwapMax=0` vs default terbuka).

## 4. Prediksi (ditulis SEBELUM menjalankan)
<!-- OPERATOR VOICE: jangan diedit setelah melihat hasil -->
* **Prediksi PID pertama di `unshare` (1.1a)**: "Nomor 1 lah, kita kan daftarin bash sebagai proses pertama."
* **Prediksi isi `cpu.max` dan `memory.max` (1.1b)**: "`20000 100000`; `50M`."
* **Prediksi kode keluar sinyal (1.1c)**: "Kalau SIGKILL angkanya 137, kalau SIGTERM dia 143." Saat ditembak ke PID 1 di dalam `unshare`: "Dia tidak terjadi apa-apa."
* **Prediksi Matriks 2×2 Chaos 1.1**: "1B dan 2B selamat, yang keluar 137 itu 2A, dan 2B selamat karena memory-nya masih bisa dipindahkan ke swap."

## 5. Perintah & Output Asli

### A. Misi 1.1a — `/proc/self/ns`, `readlink`, dan `unshare`
```text
 ~/Projects/devops-project-journey  main
❯ ls /proc/self/ns
lrwxrwxrwx - joy 16 Nov 09:31  cgroup -> cgroup:[4026531835]
lrwxrwxrwx - joy 16 Nov 09:31  ipc -> ipc:[4026531839]
lrwxrwxrwx - joy 16 Nov 09:31  mnt -> mnt:[4026531832]
lrwxrwxrwx - joy 16 Nov 09:31  net -> net:[4026531833]
lrwxrwxrwx - joy 16 Nov 09:31  pid -> pid:[4026531836]
lrwxrwxrwx - joy 16 Nov 09:31  pid_for_children -> pid:[4026531836]
lrwxrwxrwx - joy 16 Nov 09:31  time -> time:[4026531834]
lrwxrwxrwx - joy 16 Nov 09:31  time_for_children -> time:[4026531834]
lrwxrwxrwx - joy 16 Nov 09:31  user -> user:[4026531837]
lrwxrwxrwx - joy 16 Nov 09:31  uts -> uts:[4026531838]

❯ type ls
ls is a function with definition
# Defined via `source`
function ls --wraps='eza -al --color=always --group-directories-first --icons=always' --description 'alias ls=eza -al --color=always --group-directories-first --icons=always'
    eza -al --color=always --group-directories-first --icons=always $argv
end

❯ readlink /proc/self/ns/pid
'pid:[4026531836]'

❯ sudo unshare --pid --fork --mount-proc bash
[sudo] password for joy: 
[root@cachyos devops-project-journey]# readlink /proc/self/ns/pid
'pid:[4026532921]'
[root@cachyos devops-project-journey]# ps -ef
UID          PID    PPID  C STIME TTY          TIME CMD
root           1       0  0 09:50 pts/3    00:00:00 bash
root           3       1  0 09:53 pts/3    00:00:00 ps -ef
```

Perbandingan `ps -ef` di host CachyOS (di luar `unshare`):
```text
joy       426144  425340  0 10:01 pts/1    00:00:00 ps -ef
```

### B. Misi 1.1b — `systemd-run --user --scope` & File cgroups v2
```text
[joy@cachyos devops-project-journey]$ cat /proc/self/cgroup
0::/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope
[joy@cachyos devops-project-journey]$ cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope/cpu.max
20000 100000
[joy@cachyos devops-project-journey]$ cat /sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/app.slice/run-p14101-i33821.scope/memory.max
52428800
```

### C. Misi 1.1c — Sinyal OS, Kode Keluar `$status`, & Kekebalan `PID 1`
```text
 ~/Projects/devops-project-journey  main*
❯ kill -l
HUP INT QUIT ILL TRAP ABRT IOT BUS FPE KILL USR1 SEGV USR2 PIPE ALRM
TERM STKFLT CHLD CLD CONT STOP TSTP TTIN TTOU URG XCPU XFSZ VTALRM PROF
WINCH IO POLL PWR SYS RT<N> RTMIN+<N> RTMAX-<N>

 ~/Projects/devops-project-journey  main*
❯ bash -c 'kill -9 $$'; echo $status
137
fish: Job 1, 'bash -c 'kill -9 $$'' terminated by signal SIGKILL (Forced quit)

[root@cachyos devops-project-journey]# echo $$
kill -TERM 1
kill -KILL 1
1
```

### D. Misi 1.1d & Chaos 1.1 — Lazy Allocation `calloc()` vs OOM Killer vs Swap ZRAM
```text
 ~/Projects/devops-project-journey  main*
❯ cat /sys/fs/cgroup/user.slice/memory.events
low 0
high 0
max 35
oom 1
oom_kill 1
oom_group_kill 0
sock_throttled 0

 ~/Projects/devops-project-journey  main*
❯ systemd-run --user --scope -p MemoryMax=50M -p MemorySwapMax=0 python3 -c 'x = bytearray(100 * 1024 * 1024)'; echo $status
Running as unit: run-p148196-i148273.scope; invocation ID: c90c1e5dd39e4eabb133eda069e62b17
0

 ~/Projects/devops-project-journey  main*
❯ cat /sys/fs/cgroup/user.slice/memory.events
low 0
high 0
max 47
oom 1
oom_kill 1
oom_group_kill 0
sock_throttled 0

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

Hasil Matriks 2×2 (Sel 1B, Sel 2A, Sel 2B):
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

## 6. Prediksi vs Kenyataan
| Eksperimen | Prediksi | Kenyataan | Cocok? |
| :--- | :--- | :--- | :--- |
| PID `bash` di dalam `unshare` | `PID 1` | `PID 1` (inode `pid:[4026532921]`) | ✅ |
| Nilai `cpu.max` & `memory.max` | `20000 100000` & `50M` | `20000 100000` & `52428800` (dalam Byte) | ✅ |
| `bytearray(100MB)` + `MemorySwapMax=0` | Mati `137` | Keluar `0` (lolos karena `calloc()` *lazy allocation*) | ❌ (Temuan baru!) |
| `"a" * 100MB` + `MemorySwapMax=0` (Sel 1A) | Mati `137` | Mati `137`, `oom_kill` naik `1` → `2` | ✅ |
| `"a" * 100MB` tanpa batas swap (Sel 1B) | Selamat (`0`) | Selamat (`0`) | ✅ |
| `os.urandom(100MB)` + `MemorySwapMax=0` (Sel 2A) | Mati `137` | Mati `137` | ✅ |
| `os.urandom(100MB)` tanpa batas swap (Sel 2B) | Selamat (`0`) karena dipindah ke swap | Selamat (`0`) | ✅ |

## 7. Kesimpulan & Fakta yang Terbukti
<!-- OPERATOR VOICE -->
* **Namespaces vs cgroups**: "Namespace itu ruangan yang dibuat untuk mengisolasi proses dari luar, sedangkan cgroups merupakan aturan yang dibuat untuk proses tersebut."
* **Lazy Allocation `calloc()` pada data nol**: "Pada kasus angka nol, Python hanya menjanjikan tempat untuk memory namun tidak benar-benar menuliskan datanya pada memory dengan menggunakan fungsi `calloc()`. Dan setelah gw cari tahu, sebelum memberikan blok fisik RAM pada program baru, kernel Linux akan membersihkan terlebih dahulu blok yang tadinya dipakai dengan angka 0 biar program baru tidak bisa mengintip sisa data pada program lama. Jadi pada dasarnya menuliskan angka nol itu tidak benar-benar menuliskan data ke blok fisik, karena dari awal nilainya sudah nol."
* **Peran Swap vs Kompresi ZRAM**: "1B dan 2B selamat, yang keluar 137 itu 2A, dan 2B selamat karena memory-nya masih bisa dipindahkan ke swap."

## 8. Pertanyaan Baru yang Muncul
* Bagaimana kernel memisahkan antara nama file di direktori dengan blok data di disk ketika file dihapus saat masih dibuka oleh sebuah proses (Sub-modul 1.2)?

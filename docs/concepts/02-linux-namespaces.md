---
topic: "Linux Namespaces & Isolasi Proses (/proc)"
module-phase: "Fase 1.1"
date: 2026-11-18
author: Asylon
verification: "Semua klaim LAB/DOK"
related-docs: ["[[lab-1.1-01-namespaces-cgroups-signals-oom]]", "[[pm-01-oom-killer-and-inode-exhaustion]]"]
tags: [concept, first-principles, interview-prep, namespaces, procfs]
---

# Konsep: Linux Namespaces & Isolasi Proses (`/proc`)

## 1. Analogi Mental (Momen "Aha!")
<!-- OPERATOR VOICE: analogi versimu sendiri, bisa diingat dalam 10 detik -->
Namespace itu ruangan yang dibuat untuk mengisolasi pandangan sebuah proses dari proses lain di luar. Saat kita membuat ruangan baru dengan `unshare`, ia membuat namespace baru yang terpisah dari host yang sudah mencetak ribuan proses anakan sejak komputer dinyalakan, sehingga penghuni pertama di ruangan baru tersebut kembali mendapat nomor antrean `PID 1`.

**Batas analogi** — di mana analogi ini tidak lagi akurat:
<!-- OPERATOR VOICE -->
Meskipun proses di dalam namespace baru merasa punya ruangan sendiri dan tidak bisa melihat proses di luar, proses tersebut sebenarnya masih berjalan di atas kernel host yang sama dan masih bisa memakan habis RAM/CPU host jika belum dipasangi aturan pembatas `cgroups`. Selain itu, jika kita tidak memasang ulang `/proc` dengan `--mount-proc`, perintah seperti `ps` di dalam ruangan baru tetap akan mengintip daftar proses milik host.

## 2. Mekanisme Internal
<!-- Setiap klaim diberi label: [LAB ✓ → lab-id] / [DOK: sumber] / [INGATAN] -->
* Linux modern memiliki 8 jenis namespace yang terekspos sebagai *symbolic link* di `/proc/<PID>/ns/`: `cgroup`, `ipc`, `mnt`, `net`, `pid`, `time`, `user`, dan `uts` (serta varian `pid_for_children` dan `time_for_children`). [DOK: `man 7 namespaces`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Dua proses berada di namespace yang sama jika nomor inode pada target symlink `/proc/<PID>/ns/<jenis>` bernilai sama (contoh di host CachyOS: `pid:[4026531836]`). Semua proses baru secara alami mewarisi namespace induknya kecuali dipisahkan dengan *system call* `unshare(2)` atau `clone(2)`. [DOK: `man 7 namespaces`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Proses pertama yang lahir di dalam PID namespace baru melihat dirinya sendiri sebagai `PID 1` dengan `PPID 0` (karena induknya berada di luar namespace tersebut), sementara dari sudut pandang host proses tersebut tetap memiliki nomor PID reguler yang besar. [DOK: `man 7 pid_namespaces`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Utilitas `ps` membaca data proses dari *virtual filesystem* `/proc`, sehingga pembuatan PID namespace baru wajib diikuti pemasangan instans `procfs` baru (`--mount-proc`) agar `ps` menampilkan proses di namespace tersebut. [DOK: `man 1 unshare`, `man 7 pid_namespaces`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]

## 3. Pembuktian di Terminal
```fish
# 1. Cek nomor inode ruangan PID di host CachyOS
readlink /proc/self/ns/pid
# Output: 'pid:[4026531836]'

# 2. Buat ruangan PID baru dan pasang /proc baru
sudo unshare --pid --fork --mount-proc bash

# 3. Verifikasi nomor ruangan baru & daftar proses di dalamnya
readlink /proc/self/ns/pid
# Output: 'pid:[4026532921]'
ps -ef
# Output: bash sebagai PID 1 (PPID 0) dan ps -ef sebagai PID 3 (PPID 1)
```
Bukti lengkap: [lab-1.1-01-namespaces-cgroups-signals-oom.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md)

## 4. Miskonsepsi Awal Saya
<!-- OPERATOR VOICE: apa yang dulu kamu kira benar, dan apa yang membuktikan itu salah -->
* **Dulu saya kira**: Dua tab terminal memiliki nomor namespace `pid:[4026531836]` yang sama karena keduanya sama-sama menjalankan program `fish` (*"sama lah prosesnya sama-sama fish bukannya?"*). → **Ternyata**: Nomor `4026531836` bukan ID program `fish`, melainkan nomor ruangan utama host yang otomatis diwarisi oleh semua proses anak (baik `fish`, `bash`, maupun `python`) sampai dipisahkan secara eksplisit dengan `unshare` (bukti: setelah `sudo unshare --pid --fork --mount-proc bash`, nomor ruangan berubah menjadi `pid:[4026532921]`).
* **Dulu saya kira**: `ps -ef` bekerja seperti buku log riwayat sistem, dan proses yang sama jika dijalankan ulang akan memakai kembali nomor PID lamanya. → **Ternyata**: `ps -ef` hanya memotret proses yang masih hidup detik itu juga, dan setiap kali perintah dijalankan ulang, kernel selalu melahirkan proses baru dengan nomor antrean PID baru yang terus bertambah (bukti: di dalam `unshare`, `readlink` memakai `PID 2` lalu selesai dan hilang, sehingga `ps -ef` berikutnya mendapat `PID 3`; di host, `ps -ef` mendapat `PID 426144`).

## 5. Jebakan Produksi & Pertanyaan Interview
* **Pertanyaan**: Mengapa ketika kita membuat PID namespace baru dengan `unshare`, kita wajib menambahkan flag `--fork` dan `--mount-proc`?
* **Jawaban saya**: <!-- OPERATOR VOICE --> `--fork` dibutuhkan karena proses yang sudah telanjur lahir tidak bisa langsung pindah ruangan PID, jadi harus melahirkan proses anakan baru (`bash` sebagai proses pertama / nomor 1). Sedangkan `--mount-proc` gunanya untuk membuat salinan sistem `/proc` yang baru khusus untuk namespace tersebut, supaya waktu kita jalankan `ps -ef` yang terbaca bukan lagi `/proc` milik host yang sudah mencetak ratusan ribu proses sejak komputer dinyalakan.

## 6. Sumber Primer
* `man 7 namespaces` — https://man7.org/linux/man-pages/man7/namespaces.7.html
* `man 7 pid_namespaces` — https://man7.org/linux/man-pages/man7/pid_namespaces.7.html
* `man 1 unshare` — https://man7.org/linux/man-pages/man1/unshare.1.html
* `man 5 proc` — https://man7.org/linux/man-pages/man5/proc.5.html

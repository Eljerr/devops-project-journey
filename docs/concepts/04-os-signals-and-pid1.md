---
topic: "Sinyal OS (SIGTERM vs SIGKILL), Kode Keluar 128+N, & Aturan PID 1"
module-phase: "Fase 1.1"
date: 2026-11-18
author: Asylon
verification: "Semua klaim LAB/DOK"
related-docs: ["[[lab-1.1-01-namespaces-cgroups-signals-oom]]", "[[pm-01-oom-killer-and-inode-exhaustion]]"]
tags: [concept, first-principles, interview-prep, signals, pid1, docker-stop]
---

# Konsep: Sinyal OS (`SIGTERM` vs `SIGKILL`), Kode Keluar `128+N`, & Aturan `PID 1`

## 1. Analogi Mental (Momen "Aha!")
<!-- OPERATOR VOICE: analogi versimu sendiri, bisa diingat dalam 10 detik -->
`SIGTERM` (15) itu sinyal yang dikirim untuk mematikan sebuah proses secara *safe* (sopan, memberi waktu beres-beres), sedangkan `SIGKILL` (9) mematikan sebuah proses secara paksa langsung dari kernel. Khusus untuk proses yang menjadi `PID 1` di dalam sebuah namespace, ia hanya menanggapi `SIGTERM` yang sudah dipasang *handler*-nya; selain itu semua sinyal diabaikan kecuali ditembak `SIGKILL` paksa dari sistem induk (host).

**Batas analogi** — di mana analogi ini tidak lagi akurat:
<!-- OPERATOR VOICE -->
Kekebalan `PID 1` terhadap sinyal bukan karena proses tersebut dijalankan oleh user `root`, melainkan aturan perlindungan kernel khusus untuk proses *init* (`PID 1`) di setiap ruangan namespace. Bahkan jika kita menjadi `root` di dalam ruangan `unshare` dan menembak `kill -KILL 1` ke diri sendiri, `PID 1` tetap tidak akan mati.

## 2. Mekanisme Internal
<!-- Setiap klaim diberi label: [LAB ✓ → lab-id] / [DOK: sumber] / [INGATAN] -->
* `SIGTERM` (sinyal 15) adalah sinyal penghentian standar yang dapat ditangkap (*caught*), ditangani (*handled*), atau diabaikan (*ignored*) oleh proses aplikasi. Sebaliknya, `SIGKILL` (sinyal 9) dan `SIGSTOP` langsung dieksekusi oleh kernel dan tidak dapat ditangkap, diblokir, maupun diabaikan oleh proses biasa. [DOK: `man 7 signal`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Proses yang berakhir karena sinyal nomor `N` menghasilkan kode keluar dengan rumus `128 + N`: sinyal `SIGKILL` (9) menghasilkan `128 + 9 = 137`, sedangkan `SIGTERM` (15) menghasilkan `128 + 15 = 143`. Di shell `fish`, kode keluar ini disimpan pada variabel `$status` (berbeda dengan `bash` yang memakai `$?`). [DOK: `man 1 bash`, dokumentasi `fish`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Kode keluar `137` **tidak otomatis berarti kehabisan memori (OOM)**: perintah manual `kill -9 <PID>` juga menghasilkan kode keluar `137`. [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]
* Proses pertama di dalam PID namespace (`PID 1`) memiliki perlakuan khusus di kernel: kernel tidak mengirimkan sinyal yang tidak memiliki *signal handler* terpasang pada `PID 1` (termasuk `SIGTERM` dan `SIGKILL` yang dikirim dari dalam namespace itu sendiri). Hanya `SIGKILL` atau `SIGSTOP` yang dikirim dari **namespace leluhur (host)** yang dapat mematikan paksa `PID 1`. [DOK: `man 7 pid_namespaces`] [LAB ✓ → `lab-1.1-01-namespaces-cgroups-signals-oom`]

## 3. Pembuktian di Terminal
```fish
# 1. Membuktikan rumus 128 + 9 = 137 tanpa OOM di shell fish
bash -c 'kill -9 $$'; echo $status
# Output: 137
# fish: Job 1, 'bash -c 'kill -9 $$'' terminated by signal SIGKILL (Forced quit)

# 2. Membuktikan kekebalan PID 1 di dalam PID namespace terhadap SIGTERM dan SIGKILL lokal
sudo unshare --pid --fork --mount-proc bash
echo $$
# Output: 1
kill -TERM 1
kill -KILL 1
echo $$
# Output: 1 (proses PID 1 tetap hidup dan tidak terpengaruh)
```
Bukti lengkap: [lab-1.1-01-namespaces-cgroups-signals-oom.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.1-01-namespaces-cgroups-signals-oom.md)

## 4. Miskonsepsi Awal Saya
<!-- OPERATOR VOICE: apa yang dulu kamu kira benar, dan apa yang membuktikan itu salah -->
* **Dulu saya kira**: Kode keluar `137` sudah pasti berarti proses mati karena kehabisan memori (OOM Killer), dan `PID 1` kebal sinyal karena berjalan sebagai user `root`. → **Ternyata**: Menjalankan `bash -c 'kill -9 $$'; echo $status` saat RAM masih sangat lega tetap menghasilkan `137` karena `137` murni berarti mati oleh sinyal `SIGKILL` (`128 + 9`). Selain itu, `PID 1` di dalam namespace mengabaikan `SIGTERM` karena aturan kernel untuk *init process* yang hanya menanggapi sinyal yang dipasang *handler*-nya, dan hanya bisa dimatikan paksa dengan `SIGKILL` dari ruangan induk/host.

## 5. Jebakan Produksi & Pertanyaan Interview
* **Pertanyaan**: Mengapa sebuah container Docker sering kali membutuhkan waktu tepat 10 detik (seolah macet) saat dihentikan dengan `docker stop`, lalu akhirnya keluar dengan kode `137`?
* **Jawaban saya**: <!-- OPERATOR VOICE --> Saat `docker stop` dijalankan, Docker mengirim sinyal `SIGTERM` (15) ke proses `PID 1` di dalam container agar berhenti secara *safe*. Namun karena `PID 1` pada namespace hanya menanggapi `SIGTERM` yang dipasang *handler*-nya, jika aplikasi tidak memasang *handler* untuk `SIGTERM`, sinyal itu akan diabaikan oleh kernel. Setelah menunggu batas waktu 10 detik (*timeout*), Docker dari sistem induk (host) akhirnya menembak paksa proses tersebut dengan `SIGKILL` (9), sehingga keluar dengan kode `137` (`128 + 9`).

## 6. Sumber Primer
* `man 7 signal` — https://man7.org/linux/man-pages/man7/signal.7.html
* `man 7 pid_namespaces` — https://man7.org/linux/man-pages/man7/pid_namespaces.7.html
* `man 1 kill` — https://man7.org/linux/man-pages/man1/kill.1.html

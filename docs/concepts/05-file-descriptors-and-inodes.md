---
topic: "File Descriptors (fdinfo), Inodes vs Dentries, Deleted Open Files, & Sparse Files"
module-phase: "Fase 1.2"
date: 2026-11-19
author: Asylon
verification: "Semua klaim LAB/DOK"
related-docs: ["[[lab-1.2-01-fd-inode-disk-full]]", "[[pm-01-oom-killer-and-inode-exhaustion]]"]
tags: [concept, first-principles, interview-prep, file-descriptors, inodes, lsof, sparse-file]
---

# Konsep: File Descriptors (`fdinfo`), Inodes vs Dentries, *Deleted Open Files*, & *Sparse Files*

## 1. Analogi Mental (Momen "Aha!")
<!-- OPERATOR VOICE: analogi versimu sendiri, bisa diingat dalam 10 detik -->
Dentry itu label (stiker nama di folder) yang menunjuk ke nomor Inode tempat file itu benar-benar berada, sedangkan Inode berisi semua informasi dan blok data mengenai filenya. Sementara itu, `pos` pada `/proc/<PID>/fdinfo/<FD>` berfungsi sebagai pembatas buku (titik terakhir kursor setiap kali proses menulis ke file tersebut).

**Batas analogi** — di mana analogi ini tidak lagi akurat:
<!-- OPERATOR VOICE -->
Menghapus file dengan `rm` hanya mencopot label Dentry (`Links` berkurang 1), tetapi tidak langsung membebaskan blok data di SSD jika masih ada proses yang memegang File Descriptor ke Inode tersebut. Sebaliknya, memangkas file dengan `truncate -s 0` mengosongkan blok data di Inode menjadi `0`, tetapi jika proses membuka file dengan satu panah `>` (tanpa `O_APPEND`), pembatas buku (`pos`) di `fdinfo` tidak ikut balik ke `0`, sehingga penulisan berikutnya menciptakan *sparse file* (ukuran tampak `101M`, tetapi isi fisik di disk hanya `4.0K`).

## 2. Mekanisme Internal
<!-- Setiap klaim diberi label: [LAB ✓ → lab-id] / [DOK: sumber] / [INGATAN] -->
* Setiap proses memiliki tabel File Descriptor di `/proc/<PID>/fd/` (dengan `0` = `stdin`, `1` = `stdout`, `2` = `stderr`) dan metadata posisi baca/tulis (`pos`) serta flag pembukaan file (`flags` dalam oktal) di `/proc/<PID>/fdinfo/<FD>`. Flag `O_APPEND` direpresentasikan oleh bit oktal `02000` (misal `0102001` vs `0100001`). [DOK: `man 5 proc`, `man 2 open`] [LAB ✓ → `lab-1.2-01-fd-inode-disk-full`]
* Direktori di Linux hanyalah pemetaan nama (*Dentry*) ke nomor **Inode**. Beberapa nama file dapat menunjuk ke Inode yang sama melalui *hard link* (`ln`). Perintah `rm` memanggil *system call* `unlink(2)` yang hanya mengurangi `Links` (`nlink`). Blok data Inode baru dibebaskan oleh kernel ketika **`Links == 0` DAN tidak ada lagi proses yang memegang File Descriptor terbuka** ke Inode tersebut. [DOK: `man 7 inode`, `man 2 unlink`] [LAB ✓ → `lab-1.2-01-fd-inode-disk-full`]
* `du` menghitung ruang dengan menelusuri pohon direktori (*directory traversal* atas Dentry yang masih ada), sedangkan `df` membaca langsung alokasi blok pada *superblock* filesystem (`statfs(2)`). Oleh karena itu, file yang sudah di-`rm` tetapi masih dipegang oleh proses (*deleted open file*) tidak terlihat oleh `du`, namun tetap dihitung penuh oleh `df`, dan dapat dilacak menggunakan `lsof +L1` (mencari file terbuka dengan `NLINK < 1`). [DOK: `man 1 df`, `man 1 du`, `man 8 lsof`] [LAB ✓ → `lab-1.2-01-fd-inode-disk-full`]
* Ketika file yang sedang dibuka **tanpa** `O_APPEND` dipangkas ke 0 byte (`truncate -s 0`), ukuran pada Inode menjadi `0`, tetapi offset `pos` pada *Open File Description* milik proses tidak berubah. Penulisan berikutnya pada offset lama menciptakan **sparse file** (*lubang hampa / hole* dari byte 0 hingga offset `pos` tersebut): ukuran tampak (`du --apparent-size` / `ls -lh`) terlihat besar (`101M`), tetapi blok disk fisik yang dipakai (`du -h`) hanya sebesar data baru yang ditulis (`4.0K`). [DOK: `man 2 lseek`] [LAB ✓ → `lab-1.2-01-fd-inode-disk-full`]
* Jumlah total Inode pada filesystem `ext4` ditetapkan secara tetap saat format (`mkfs.ext4`). Jika seluruh Inode habis terpakai oleh banyak file kecil (`df -i` menunjukkan `IUse% 100%`), pembuatan file baru akan gagal dengan `No space left on device` meskipun kapasitas blok disk (`df -h`) masih kosong (`Use% 1%`). [DOK: `man 8 mke2fs`] [LAB ✓ → `lab-1.2-01-fd-inode-disk-full`]

## 3. Pembuktian di Terminal
```fish
# 1. Membuktikan Deleted Open File (du vs df & lsof +L1)
bash -c 'exec 3> ~/lab/hantu.log; fallocate -l 500M ~/lab/hantu.log; sleep 600' &
rm ~/lab/hantu.log
du -sh ~/lab       # Output: 4.0K (Dentry sudah hilang)
df -BM ~/lab       # Output: Used tetap 52156M (Inode masih ditahan FD 3)
lsof +L1 | grep lab
# Output: sleep 74329 joy 3w REG 259,9 524288000 0 4849708 /home/joy/lab/hantu.log (deleted)

# 2. Membuktikan Sparse File akibat truncate tanpa O_APPEND
# Setelah menulis 100M (pos: 104857600) -> truncate -s 0 -> echo "halo" >&3:
# pos: 104857605 | du -h: 4.0K (fisik) | du -h --apparent-size: 101M (tampak)
```
Bukti lengkap: [lab-1.2-01-fd-inode-disk-full.md](file:///home/joy/Projects/devops-project-journey/docs/labs/lab-1.2-01-fd-inode-disk-full.md)

## 4. Miskonsepsi Awal Saya
<!-- OPERATOR VOICE: apa yang dulu kamu kira benar, dan apa yang membuktikan itu salah -->
* **Dulu saya kira**: Setelah file 500 MB yang sedang dipegang proses dihapus dengan `rm ~/lab/hantu.log`, `du -sh` masih akan menampilkan `500M` dan `df -BM` akan langsung turun `500M`. → **Ternyata**: Kebalikannya! `du` hanya memindai bagian luarnya doang yang masih ada Dentry-nya (sehingga langsung turun ke `4.0K`), sedangkan `df` memeriksa langsung ke pembukuan Inode-nya yang masih dipegang oleh proses (sehingga tetap `52156M` dan baru turun ke `51656M` setelah proses pemegang FD-nya dimatikan dengan `kill`).
* **Dulu saya kira**: `echo "halo"` hanya menggeser `pos` sebanyak 4 byte. → **Ternyata**: Nilainya `5` karena `echo` selalu menyisipkan 1 karakter Enter (`\n`) pada akhir kalimat.

## 5. Jebakan Produksi & Pertanyaan Interview
* **Pertanyaan**: Apa perbedaan antara insiden *Deleted Open Files* dan *Inode Exhaustion* ketika server mengalami masalah disk penuh, dan bagaimana cara mendiagnosis masing-masing?
* **Jawaban saya**: <!-- OPERATOR VOICE --> Keduanya sama-sama bikin disk ilusi: *Deleted Open Files* itu kondisi di mana file sudah di-`rm` (sehingga tidak terlihat di `du`), tetapi Inode-nya masih dipegang oleh proses aktif sehingga kapasitas blok di `df -h` masih penuh; ini kita lacak dengan `lsof +L1` untuk menyaring file terbuka yang `Links`-nya kurang dari 1. Sedangkan *Inode Exhaustion* adalah kondisi di mana kapasitas ruang di `df -h` masih longgar, tetapi disk kehabisan nomor Inode (terlihat `IUse% 100%` dan `IFree 0` saat dicek dengan `df -i`).

## 6. Sumber Primer
* `man 5 proc` (`/proc/[pid]/fd` & `/proc/[pid]/fdinfo`) — https://man7.org/linux/man-pages/man5/proc.5.html
* `man 7 inode` — https://man7.org/linux/man-pages/man7/inode.7.html
* `man 2 unlink` — https://man7.org/linux/man-pages/man2/unlink.2.html
* `man 2 lseek` ("Seeking file data and holes") — https://man7.org/linux/man-pages/man2/lseek.2.html
* `man 8 lsof` — https://man7.org/linux/man-pages/man8/lsof.8.html

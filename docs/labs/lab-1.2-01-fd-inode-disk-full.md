---
lab-id: LAB-1.2-01
date: 2026-11-19
module-phase: "Fase 1.2"
environment: "INNER CachyOS"
author: Asylon
hint-level-max: "0"
tags: [lab, evidence, fase-1, file-descriptor, inode, dentry, deleted-open-files, inode-exhaustion, sparse-file]
---

# Lab: File Descriptor (`fdinfo`), Hard Link Inode, Anomali `df` vs `du` (`lsof +L1`), Inode Exhaustion, & *Sparse File*

> Berkas ini adalah **bukti mentah** dari sesi terminal asli di CachyOS Linux (filesystem `/` dan `~/lab` bertipe `ext4`, shell `fish`). Dokumen lain (`docs/concepts/05-file-descriptors-and-inodes.md` dan Bagian 2 `docs/post-mortems/pm-01-oom-killer-and-inode-exhaustion.md`) mengutip langsung dari berkas ini.

## 1. Pertanyaan yang Ingin Dibuktikan
<!-- OPERATOR VOICE -->
1. Bagaimana kernel mencatat posisi kursor baca/tulis (`pos`) dan mode `O_APPEND` (`flags: 0102001`) di `/proc/<PID>/fdinfo/<FD>`?
2. Apa perbedaan antara label nama file di direktori (Dentry) dengan Inode ketika dibuat *hard link* (`ln`) lalu dihapus dengan `rm` (`unlink`)?
3. Mengapa file berukuran 500 MB yang sudah di-`rm` tetap menahan ruang di `df -BM` padahal `du -sh` sudah melaporkan kosong (`4.0K`), dan bagaimana `lsof +L1` menangkap prosesnya?
4. Mengapa disk bisa menolak pembuatan file baru (`No space left on device`) ketika `df -h` baru terpakai `1%` (*Inode Exhaustion* pada `df -i`)?
5. Apa efek `truncate -s 0` terhadap kursor `pos` dan pemakaian blok fisik SSD (`du -h` vs `du -h --apparent-size`) pada proses yang membuka file tanpa `O_APPEND`?

## 2. Kondisi Awal (Versi & Lingkungan)
```text
 ~/Projects/devops-project-journey  main*
❯ rm -f ~/lab/koper_b.txt; df -h ~/lab; du -sh ~/lab
Filesystem      Size  Used Avail Use% Mounted on
/dev/nvme0n1p6  118G   51G   61G  46% /
4.0K    /home/joy/lab

 ~/Projects/devops-project-journey  main*
❯ df -h ~/lab; df -i ~/lab
Filesystem      Size  Used Avail Use% Mounted on
/dev/nvme0n1p6  118G   51G   61G  46% /
Filesystem      Inodes  IUsed   IFree IUse% Mounted on
/dev/nvme0n1p6 7831552 816310 7015242   11% /
```

## 3. Variabel
* **Eksperimen 1.2a (`fdinfo`)**: Membuka FD 3 menggunakan `>` (tanpa `O_APPEND`) vs `>>` (dengan `O_APPEND`).
* **Eksperimen 1.2b (Dentry vs Inode)**: Menambah *hard link* dengan `ln`, lalu menghapus stiker nama pertama dengan `rm`.
* **Eksperimen 1.2c (*Deleted Open Files*)**: Menghapus file 500 MB (`rm ~/lab/hantu.log`) saat proses `sleep` masih menahan FD 3, lalu mematikan proses tersebut dengan `kill`.
* **Eksperimen 1.2d (*Inode Exhaustion*)**: Membuat loop device `ext4` 20 MB dengan batas `-N 64` (`72` Inode) dan menyuntikkan 100 file berukuran `0 Byte`.
* **Eksperimen Chaos 1.2 (*Sparse File*)**: Memangkas file 100 MB dengan `truncate -s 0` saat proses masih memegang FD 3 tanpa `O_APPEND`, lalu menulis `"halo"` (5 Byte).

## 4. Prediksi (ditulis SEBELUM menjalankan)
<!-- OPERATOR VOICE: jangan diedit setelah melihat hasil -->
* **Prediksi 1.2a (`pos` & `flags`)**: "Percobaan A = 4 dan Percobaan B = pos-nya 5 dan flags-nya beda."
* **Prediksi 1.2b (`ln` & `rm`)**: "Nomor inode-nya akan sama; masih ada dan angka link-nya satu."
* **Prediksi 1.2c (`du` vs `df` setelah `rm`)**: "500M; turun 500M."
* **Prediksi 1.2d (`mkfs.ext4 -N 64` + 100 file kosong)**: "Tidak, no space left; `IUse%` dan `IUsed` penuh dan `IFree` kosong."
* **Prediksi Chaos 1.2 (`truncate -s 0` tanpa `O_APPEND`)**: "Kalau pake 1 panah gw rasa ga balik deh (`pos`-nya); nah ini gw masih kurang tau, tergantung bagaimana sistem melihatnya, kalau cuma dialokasikan mungkin masih penuh, tapi kalau sistem melihat benar-benar isi fisiknya mungkin 4kb saja."

## 5. Perintah & Output Asli

### A. Misi 1.2a — `/proc/self/fd` & `/proc/$$/fdinfo/3`
```text
 ~/Projects/devops-project-journey  main*
❯ ls /proc/self/fd
lrwx------ - joy 18 Nov 08:47  0 -> /dev/pts/0
lrwx------ - joy 18 Nov 08:47  1 -> /dev/pts/0
lrwx------ - joy 18 Nov 08:47  2 -> /dev/pts/0
...
lr-x------ - joy 18 Nov 08:47  61 -> '/dev/shm/.org.chromium.Chromium.1uim5W (deleted)'

 ~/Projects/devops-project-journey  main*
❯ bash -c 'exec 3> /tmp/uji_fd.txt; echo "halo" >&3; cat /proc/$$/fdinfo/3'
pos:    5
flags:  0100001
mnt_id: 59
ino:    1352

 ~/Projects/devops-project-journey  main
❯ bash -c 'exec 3>> /tmp/uji_fd.txt; echo "halo" >&3; cat /proc/$$/fdinfo/3'
pos:    10
flags:  0102001
mnt_id: 59
ino:    1352
```

### B. Misi 1.2b — Dentry vs Inode (`stat`, `ln`, `rm`)
```text
❯ mkdir -p ~/lab; touch ~/lab/koper_a.txt; stat ~/lab/koper_a.txt
  File: /home/joy/lab/koper_a.txt
  Size: 0               Blocks: 0          IO Block: 4096   regular empty file
Device: 259,9   Inode: 4849708     Links: 1
Access: (0644/-rw-r--r--)  Uid: ( 1000/     joy)   Gid: ( 1000/     joy)

 ~/Projects/devops-project-journey  main*
❯ echo "rahasia negara" > ~/lab/koper_a.txt; ln ~/lab/koper_a.txt ~/lab/koper_b.txt; stat -c '%n -> Inode: %i | Links: %h' ~/lab/koper_a.txt ~/lab/koper_b.txt
/home/joy/lab/koper_a.txt -> Inode: 4849708 | Links: 2
/home/joy/lab/koper_b.txt -> Inode: 4849708 | Links: 2

 ~/Projects/devops-project-journey  main*
❯ rm ~/lab/koper_a.txt; stat -c '%n -> Inode: %i | Links: %h' ~/lab/koper_b.txt; cat ~/lab/koper_b.txt
/home/joy/lab/koper_b.txt -> Inode: 4849708 | Links: 1
rahasia negara
```

### C. Misi 1.2c — *Deleted Open Files* (`df -BM` vs `du -sh` & `lsof +L1`)
```text
 ~/Projects/devops-project-journey  main*
❯ df -BM ~/lab
  bash -c 'exec 3> ~/lab/hantu.log; fallocate -l 500M ~/lab/hantu.log; sleep 600' &
  df -BM ~/lab; du -sh ~/lab
  rm ~/lab/hantu.log; df -BM ~/lab; du -sh ~/lab
  lsof +L1 | grep lab
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 52156M    61586M  46% /
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 52156M    61586M  46% /
501M    /home/joy/lab
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 52156M    61586M  46% /
4.0K    /home/joy/lab
sleep     74329  joy   3w   REG  259,9 524288000     0 4849708 /home/joy/lab/hantu.log (deleted)
sleep     78069  joy   3w   REG  259,9 524288000     0 4849708 /home/joy/lab/hantu.log (deleted)

❯ kill 74329 78069; df -BM ~/lab
Filesystem     1M-blocks   Used Available Use% Mounted on
/dev/nvme0n1p6   119876M 51656M    62086M  46% /
fish: Job 2, 'bash -c 'exec 3> ~/lab/hantu.lo…' terminated by signal SIGTERM (Polite quit request)
fish: Job 1, 'bash -c 'exec 3> ~/lab/hantu.lo…' terminated by signal SIGTERM (Polite quit request)
```

### D. Misi 1.2d — *Inode Exhaustion* (`mkfs.ext4 -N 64` & `df -i`)
```text
 ~/Projects/devops-project-journey  main*
❯ truncate -s 20M /tmp/disk_mini.img; mkfs.ext4 -N 64 /tmp/disk_mini.img
  mkdir -p ~/lab/mnt_inode; sudo mount -o loop /tmp/disk_mini.img ~/lab/mnt_inode; sudo chown joy:joy ~/lab/mnt_inode
  for i in (seq 100); touch ~/lab/mnt_inode/file_$i.txt; end
  df -h ~/lab/mnt_inode; df -i ~/lab/mnt_inode
mke2fs 1.47.4 (6-Mar-2025)
Creating filesystem with 20480 1k blocks and 72 inodes
...
touch: cannot touch '/home/joy/lab/mnt_inode/file_61.txt': No space left on device
...
touch: cannot touch '/home/joy/lab/mnt_inode/file_100.txt': No space left on device
Filesystem      Size  Used Avail Use% Mounted on
/dev/loop0       19M   48K   18M   1% /home/joy/lab/mnt_inode
Filesystem     Inodes IUsed IFree IUse% Mounted on
/dev/loop0         72    72     0  100% /home/joy/lab/mnt_inode
```

### E. Tantangan Chaos 1.2 — `truncate -s 0` vs Offset `pos` & *Sparse File*
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

## 6. Prediksi vs Kenyataan
| Eksperimen | Prediksi | Kenyataan | Cocok? |
| :--- | :--- | :--- | :--- |
| `pos` & `flags` (`echo "halo"`) | `pos` = 4 & 5; `flags` beda | `pos` = 5 & 10 (karena karakter `\n`); `flags` beda (`0100001` vs `0102001`) | Sebagian (belajar `\n`) |
| `ln` dan `rm` pada `koper_a.txt` | Inode sama, `koper_b.txt` tetap ada (`Links: 1`) | Inode sama (`4849708`), `Links` turun `2` → `1`, data utuh | ✅ |
| `du` vs `df` setelah `rm` file terbuka | `du` `500M`, `df` turun `500M` | Kebalikannya: `du` turun ke `4.0K`, `df` tetap `52156M` sampai proses di-`kill` | ❌ (Momen Aha!) |
| 100 file kosong di partisi 72 Inode | Gagal `No space left`, `IUse%` 100% | Gagal mulai `file_61.txt`, `df -h` `1%`, `df -i` `100%` | ✅ |
| `truncate -s 0` tanpa `O_APPEND` + tulis 5 Byte | `pos` tidak balik; tampak penuh (~101M) tapi isi fisik `4.0K` | `pos` tetap `104857600` → `104857605`; fisik `du -h` = `4.0K`, tampak `--apparent-size` = `101M` | ✅ |

## 7. Kesimpulan & Fakta yang Terbukti
<!-- OPERATOR VOICE -->
* **`pos` dan `O_APPEND`**: "`pos` pada `/proc/<PID>/fdinfo/<FD>` berfungsi sebagai pembatas buku titik terakhir kursor setiap kali menulis, dan `echo halo` itu nilainya 5 karena sintaks tersebut selalu menyisipkan 1 karakter yaitu enter (`\n`) pada akhir kalimat. Ketika file dibuka pertama kali, mode `>` memposisikan kursornya pada titik awal, sedangkan dengan `O_APPEND` (`>>`) posisi kursor selalu berada pada akhir."
* **Dentry vs Inode**: "Dentry itu label yang menunjuk ke inode tempat file itu benar-benar berada, inode itu berisi tentang semua informasi mengenai filenya."
* **Anomali `du` vs `df` & `lsof +L1`**: "`du` hanya memindai bagian luarnya doang yang masih ada Dentry-nya, sedangkan `df` memeriksa langsung ke nomor inode-nya yang mana masih dipegang; `lsof` memindai file yang masih dipegang oleh proses sistem dan flag `+L1` menyaring file yang links-nya kurang dari 1."
* ***Deleted Open Files* vs *Inode Exhaustion***: "Keduanya sama-sama bikin disk ilusi: Deleted Open Files itu inode-nya masih dipegang oleh proses, sedangkan Inode Exhaustion kondisi di mana disk kehabisan nomor inode."
* **Efek `truncate` & *Sparse File***: "Kalau pake 1 panah (`>`) `pos`-nya ga balik; kalau cuma dilihat alokasi tampaknya (`--apparent-size`) masih penuh (`101M`), tapi kalau sistem melihat benar-benar isi fisiknya (`du -h`) cuma `4.0K` saja."

---
incident-id: PM-0001
date: 2026-11-06
module-phase: "Fase 1.1"
severity: "High"
status: "Resolved"
operator: "Eljerr"
reviewed-by: "AI Systems Architect Agent"
---

# Post-Mortem: Investigasi Terminasi Proses OOM Killer dan Anomali Disk Penuh (Deleted Open Files)

---

## Bagian 1: Investigasi Terminasi Proses Akibat Batasan cgroups v2 dan OOM Killer

## 1. Ringkasan Eksekutif & Dampak Sistem

Eksperimen terkontrol (*chaos drill*) dilakukan di lingkungan sandbox (WSL2 Ubuntu) untuk memverifikasi mekanisme pembatasan sumber daya kernel Linux menggunakan cgroups v2. Sebuah proses alokasi memori sengaja dipaksa meminta sumber daya melebihi kuota keras (*hard limit*) yang telah dikonfigurasi.

Dampak yang teramati: Proses langsung dimatikan secara mendadak oleh kernel Linux segera setelah alokasi memori melampaui batas yang diizinkan, menghasilkan pemutusan eksekusi seketika (*instant termination*).

## 2. Gejala Awal & Bukti Log (Detection)

Perintah eksekusi:
```bash
sudo sh -c 'echo $$ > /sys/fs/cgroup/lab-oom/cgroup.procs && python3 -c "print(\"Sedang memakan RAM...\"); x = bytearray(100 * 1024 * 1024); print(\"Berhasil lolos!\")"'
```

Keluaran terminal:
```text
Sedang memakan RAM...
Killed
```

Pemeriksaan status kode keluar (*exit code*):
```bash
$ echo $?
137
```

Catatan forensik cgroups v2 (`memory.events`):
```text
oom 2
oom_kill 2
```

## 3. Kronologi Investigasi & Diagnostik (Troubleshooting Flow)

1. **Pembuatan Grup Kontrol (cgroup)**:  
   Membuat sub-folder `/sys/fs/cgroup/lab-oom` sebagai wadah isolasi sumber daya.
2. **Konfigurasi Kuota Batas Memori**:  
   - Menetapkan batas keras memori: `echo "50M" > /sys/fs/cgroup/lab-oom/memory.max`
   - Menonaktifkan alokasi swap file: `echo "0" > /sys/fs/cgroup/lab-oom/memory.swap.max`
3. **Penyuntikan Beban (Load Injection)**:  
   Menjalankan skrip Python sederhana yang meminta alokasi memori sebesar 100 MB di dalam namespace cgroup tersebut.
4. **Verifikasi Forensik Terminasi**:  
   Mengecek kode keluar `echo $?` yang menghasilkan nilai `137` dan memvalidasi counter `oom_kill` di `memory.events`.

## 4. Analisis Akar Masalah (Root Cause Analysis - Analisis Operator)

> *Penjelasan langsung oleh Operator:*  
> "Saya membuat sebuah folder di dalam cgroup dengan mengonfigurasi `memory.max`-nya 50M dan swap-nya 0, lalu saya membuat skrip sederhana untuk melihat bagaimana OOM Killer itu benar-benar bekerja. Saya menjalankan skrip Python untuk mengalokasikan memori sebanyak 100 MB. Setelah skrip dijalankan, output yang keluar adalah `Killed`, di mana ini menandakan OOM Killer bekerja untuk mematikan seketika program tersebut. Dan ketika menjalankan perintah `echo $?`, yang keluar adalah angka `137` yang menandakan exit code kehabisan memori (+9 SIGKILL)."

### Catatan Mekanisme Kernel:
- Batas maksimal cgroup: **50 MB**.
- Alokasi yang diminta proses Python: **100 MB**.
- Karena alokasi swap bernilai 0, kernel tidak dapat memindahkan memori ke disk.
- Kernel mengaktifkan fungsi `out_of_memory()`, menandai proses pemakan memori terbesar di cgroup tersebut, dan mengirimkan sinyal pembunuhan paksa **`SIGKILL` (Sinyal 9)**.
- Kode keluar dihitung berdasarkan formula standar POSIX:  
  $$\text{Exit Code} = 128 + 9 = 137$$

## 5. Resolusi, Mitigasi & Prinsip Rekayasa (Action Items)

> *Kesimpulan Operator:*  
> "Lebih hati-hati dalam menggunakan dan memperhitungkan resource infrastruktur yang kita miliki."

### Tindakan Operasional di Lingkungan Produksi Homelab:
1. **Perhitungan Budgeting Memori yang Ketat**:  
   Mengingat Proxmox VE hanya memiliki RAM total 8 GB dan STB hanya 2 GB, setiap kontainer LXC dan pod Kubernetes wajib memiliki batas `limits.memory` yang realistis dan diuji beban (*load test*) terlebih dahulu sebelum dirilis.
2. **Monitoring Pra-OOM**:  
   Menerapkan observabilitas metrik Prometheus pada metrik `container_memory_working_set_bytes` untuk mendeteksi tren kebocoran memori sebelum menyentuh batas `memory.max`.
3. **Pembersihan Lingkungan Lab**:  
   Menghapus direktori isolasi cgroups setelah sesi selesai:
   ```bash
   sudo rmdir /sys/fs/cgroup/lab-oom
   ```

---

## Bagian 2: Investigasi Anomali Ruang Disk Akibat Deleted Open Files

### 1. Ringkasan Eksekutif & Skenario Insiden

Eksperimen *chaos drill* dilakukan di lingkungan **Inner Loop (WSL2 Ubuntu)** untuk merekonstruksi fenomena hilangnya ruang disk fisik tanpa adanya file yang terdeteksi di direktori (*deleted open files anomaly*).

Sebuah proses penulisan log latar belakang dialokasikan menahan file berukuran 500 MB melalui **File Descriptor 3**. File tersebut sengaja dihapus menggunakan perintah `rm` saat proses masih berjalan aktif.

Dampak yang teramati: File menghilang dari sistem penelusuran direktori biasa (`ls` dan `du`), namun alokasi blok data pada sistem berkas tidak berkurang, menyebabkan selisih pelaporan kapasitas ruang disk antara `df` dan `du`.

### 2. Gejala Awal & Bukti Forensik (Detection & Evidence)

Injeksi proses latar belakang pemegang file 500 MB:
```bash
bash -c 'exec 3> /tmp/insiden_produksi.log; fallocate -l 500M /tmp/insiden_produksi.log; while true; do echo "$(date) - menulis log..." >&3; sleep 2; done' &
```

Penghapusan file via antarmuka direktori:
```bash
rm /tmp/insiden_produksi.log
```

Verifikasi di tingkat direktori:
```bash
$ ls -l /tmp/insiden_produksi.log
ls: cannot access '/tmp/insiden_produksi.log': No such file or directory
```

Deteksi berkas tertahan menggunakan `lsof +L1`:
```text
COMMAND   PID USER   FD   TYPE DEVICE SIZE/OFF NLINK     NODE NAME
bash   695177  joy    3w   REG  259,9 524288000     0 11683373 /tmp/insiden_produksi.log (deleted)
```

Inspeksi langsung pada tabel *File Descriptor* proses (`/proc/<PID>/fd/3`):
```text
❯ stat /proc/695177/fd/3
  File: /proc/695177/fd/3 -> '/tmp/insiden_produksi.log (deleted)'
  Size: 64              Blocks: 0          IO Block: 1024   symbolic link
Device: 0,26    Inode: 11683373    Links: 1
Access: (0300/l-wx------)  Uid: ( 1000/     joy)   Gid: ( 1000/     joy)
Access: 2026-11-10 15:24:49.804137648 +0700
Modify: 2026-11-10 15:21:36.009744069 +0700
Change: 2026-11-10 15:21:36.009744069 +0700
```

Terminasi proses pengunci berkas:
```text
❯ kill 695177
fish: Job 1, 'bash -c 'exec 3> /tmp/insiden_p…' terminated by signal SIGTERM (Polite quit request)
```

### 3. Analisis Akar Masalah (Root Cause Analysis - Analisis Operator)

> *Penjelasan langsung oleh Operator:*  
> 1. "Kernel menambahkan keterangan `deleted` karena Dentry sudah di-`rm`, yang mana `rm` itu menghapus link (*unlink*)."  
> 2. "Tool `lsof +L1` memang dirancang khusus untuk mendeteksi file yang link count-nya kurang dari 1, sedangkan `ls` dan `du` itu mendeteksi link-nya (*directory traversal*)."  
> 3. "Jujur setelah saya cari tahu, perintah `truncate` itu tidak efektif karena bisa mengubah *offset* di Open File Table. Jika diasumsikan aplikasi cukup pintar dalam menulis log sehingga bisa dengan pendekatan `truncate`, ya tidak ada masalah berarti."

### 4. Catatan Mekanisme Kernel & Arsitektur Sistem

1. **Pemisahan Dentry vs Inode**:
   - Perintah `rm` memicu *system call* `unlink()`, yang hanya memutus asosiasi antara string nama file di tabel direktori (dentry) dengan nomor Inode (`11683373`), sehingga `nlink` bernilai 0.
   - Kernel Linux mempertahankan blok data selama `Open File Description` pada tabel kernel masih memiliki referensi aktif (`ref_count > 0`). Kernel menambahkan tag `(deleted)` pada *symlink* di `/proc/<PID>/fd/` untuk menandai bahwa berkas fisik sudah tidak memiliki pintu masuk dentry.
2. **Perbedaan Fundamental `du` vs `df` vs `lsof +L1`**:
   - **`du`**: Melakukan penelusuran hierarki direktori (*dentry tree*). Ketika dentry terhapus, `du` gagal menemukan berkas dan melaporkan ukuran yang berkurang.
   - **`df`**: Membaca alokasi blok fisik pada *superblock* sistem berkas. Blok data masih berstatus *in-use*, sehingga `df` tetap melaporkan partisi penuh.
   - **`lsof +L1`**: Memindai langsung struktur tabel berkas terbuka kernel dan memfilter Inode yang memiliki nilai `nlink < 1`.
3. **Pembedahan File Offset vs `truncate` (Temuan Lanjutan Operator)**:
   - Ketika file yang sedang dibuka di-*truncate* menjadi 0 byte (`> /proc/<PID>/fd/3`), ukuran berkas pada Inode berubah menjadi 0.
   - Namun, jika aplikasi membuka berkas tanpa flag **`O_APPEND`**, penunjuk posisi tulis (*file write offset*) pada struktur *Open File Description* di kernel tetap berada di posisi offset terakhir (misal 500 MB). Penulisan berikutnya akan menghasilkan *sparse file* dengan lubang kosong (*null-byte hole*), sehingga ruang disk dapat kembali terisi secara anomali.
   - Jika aplikasi menggunakan flag `O_APPEND`, setiap pemanggilan *system call* `write()` secara otomatis memaksa kernel memindahkan offset ke posisi akhir berkas terkini (`SEEK_END`), sehingga pemangkasan (*truncation*) dapat berjalan aman dan efektif.

### 5. Resolusi & Mitigasi Praktik Rekayasa Produksi (Action Items)

1. **Standardisasi Log Rotation**:  
   Jangan pernah menghapus file log aktif menggunakan `rm`. Gunakan utilitas standar `logrotate` dengan strategi yang tepat:
   - Mode `copytruncate`: Menyalin isi log lalu memangkas ukuran berkas ke 0 secara aman untuk aplikasi yang mendukungnya.
   - Mode Signal Re-open: Memindahkan berkas log lalu mengirimkan sinyal reload (`SIGHUP` atau `SIGUSR1`) kepada aplikasi (seperti Nginx/Logstash) agar aplikasi menutup FD lama dan membuka FD baru secara bersih.
2. **Prosedur Diagnostik Insiden Disk Penuh**:  
   Jika terjadi *alert* disk 100% penuh namun `du` tidak menemukan berkas besar di direktori:
   ```bash
   # Segera jalankan lsof untuk mencari open deleted files
   lsof +L1
   ```
3. **Pembersihan Tanpa Downtime**:  
   Sebelum memutuskan tindakan *restart* layanan kritis, periksa apakah aplikasi menggunakan penulisan append. Jika aman, lakukan pembebasan instan melalui manipulasi *file descriptor* `/proc/<PID>/fd/<FD>` atau lakukan *graceful restart* via `systemctl reload` / sinyal `SIGUSR1`.

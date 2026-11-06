---
incident-id: PM-0001
date: 2026-11-06
module-phase: "Fase 1.1"
severity: "High"
status: "Resolved"
operator: "Eljerr"
reviewed-by: "AI Systems Architect Agent"
---

# Post-Mortem: Investigasi Terminasi Proses Akibat Batasan cgroups v2 dan OOM Killer

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

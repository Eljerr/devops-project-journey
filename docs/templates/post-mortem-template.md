---
incident-id: PM-XXXX
date: YYYY-MM-DD
module-phase: "Fase X.Y"
environment: "{Inner Loop WSL2 | Outer Loop Proxmox 8GB | Outer Loop STB ARM64}"
severity: "{Critical | High | Medium}"
status: "Resolved"
tags: [post-mortem, troubleshooting, fase-x]
---

# Post-Mortem: {Judul Kegagalan Sistem / Studi Kasus}

## 1. Ringkasan Eksekutif & Dampak Sistem
<!-- OPERATOR VOICE: Ceritakan dengan bahasamu sendiri layanan apa yang mati, dampaknya, dan konteks eksperimen/insidennya -->
{Deskripsi singkat mengenai layanan yang terganggu, gejala yang terlihat oleh pengguna/sistem, dan durasi gangguan pada lingkungan homelab.}

## 2. Gejala Awal & Bukti Log (Detection)
<!-- RAW ARTIFACT: Tempelkan output terminal/log asli tanpa diedit isinya -->
```text
{Tempelkan potongan log error asli dari journalctl, kubectl logs, dmesg, lsof, atau GitHub Actions yang menjadi petunjuk awal.}
```

## 3. Kronologi Investigasi & Diagnostik (Troubleshooting Flow)
<!-- OPERATOR VOICE: Tuliskan alur pikir saat menginvestigasi masalah, termasuk tebakan awal yang salah jika ada -->
1. **Hipotesis Awal**: {Apa yang pertama kali dicurigai dan perintah CLI apa yang dijalankan untuk mengujinya}.
2. **Temuan Menengah**: {Hasil keluaran perintah yang menggugurkan atau memperkuat hipotesis}.
3. **Identifikasi Titik Kegagalan**: {Perintah spesifik yang mengungkap sumber masalah sebenarnya}.

## 4. Analisis Akar Masalah (Root Cause Analysis — First Principles & 5 Whys)
<!-- OPERATOR VOICE: Bedah sampai ke mekanisme internal kernel/jaringan/kontainer -->
* **Gejala Permukaan**: {Masalah yang terlihat}
  * **Why 1**: Mengapa itu terjadi? -> {Jawaban}
  * **Why 2**: Mengapa hal di atas terjadi? -> {Jawaban}
  * **Why 3**: Mengapa sistem merespons demikian? -> {Jawaban}
  * **Why 4 / Mekanisme Internal Kernel/OS**: -> {Penjelasan tingkat lanjut, misal cara kerja inode, cgroup, atau netfilter}

## 5. Resolusi, Bukti Pemulihan, & Mitigasi Permanen
* **Perbaikan Langsung (Fix)**: {Perubahan konfigurasi atau perintah CLI yang memulihkan layanan}.
* **Bukti Verifikasi (Post-Fix Output)**:
  ```text
  {Tempelkan output terminal/metrik yang membuktikan masalah sudah tuntas}
  ```
* **Pencegahan & Refleksi Operator (Lessons Learned)**: {Aturan alerting, health check, atau validasi yang ditambahkan + catatan kritis pribadimu agar insiden tidak terulang}.
---
incident-id: PM-XXXX
date: YYYY-MM-DD
module-phase: "Fase X.Y"
severity: "{Critical | High | Medium}"
status: "Resolved"
---

# Post-Mortem: {Judul Kegagalan Sistem / Studi Kasus}

## 1. Ringkasan Eksekutif & Dampak Sistem

{Deskripsi singkat mengenai layanan yang terganggu, gejala yang terlihat oleh pengguna/sistem, dan durasi gangguan pada lingkungan homelab.}

## 2. Gejala Awal & Bukti Log (Detection)

{Tempelkan potongan log error asli dari journalctl, kubectl logs, dmesg, atau GitHub Actions yang menjadi petunjuk awal.}

## 3. Kronologi Investigasi & Diagnostik (Troubleshooting Flow)

1. Hipotesis Awal: {Apa yang pertama kali dicurigai dan perintah CLI apa yang dijalankan untuk mengujinya}.
2. Temuan Menengah: {Hasil keluaran perintah yang menggugurkan atau memperkuat hipotesis}.
3. Identifikasi Titik Kegagalan: {Perintah spesifik yang mengungkap sumber masalah sebenarnya}.

## 4. Analisis Akar Masalah (Root Cause Analysis - 5 Whys)

{Penjelasan mekanisme internal tingkat OS/Jaringan/Kontainer mengapa kegagalan tersebut terjadi.}

## 5. Resolusi & Mitigasi Permanen (Action Items)

* Perbaikan Langsung (Fix): {Perubahan konfigurasi atau kode yang memulihkan layanan}.
* Pencegahan (Prevention): {Aturan alerting, health check, atau validasi CI yang ditambahkan agar insiden tidak terulang}.

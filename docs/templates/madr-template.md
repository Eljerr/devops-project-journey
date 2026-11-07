---
adr-id: ADR-XXXX
status: "{proposed | accepted | superseded by ADR-YYYY}"
date: YYYY-MM-DD
decision-makers: [Eliazer Raditya Gunawan]
hardware-constraint: "{Proxmox VE 2C/8GB RAM | STB Amlogic S905X4 2GB RAM | Hybrid}"
tags: [adr, architecture, fase-x]
---

# Gunakan {Opsi Terpilih} untuk {Masalah yang Diselesaikan}

## 1. Konteks dan Pernyataan Masalah
<!-- OPERATOR VOICE: Jelaskan masalah nyata di homelab dan batasan fisiknya dengan bahasamu sendiri -->
{Jelaskan kondisi infrastruktur saat ini, batasan perangkat keras (misal: batas RAM 8GB pada Proxmox Laptop atau 2GB pada STB ARM64), serta masalah teknis yang harus diselesaikan.}

## 2. Pendorong Keputusan (Decision Drivers)
* {Batasan 1, misal: Konsumsi memori idle wajib di bawah 512 MB}
* {Batasan 2, misal: Kompatibilitas arsitektur CPU x86_64 / ARM64}
* {Batasan 3, misal: Keamanan jaringan tanpa IP Publik statis & nol biaya sewa cloud}

## 3. Opsi yang Dipertimbangkan
* **Opsi 1**: {Nama Opsi 1}
* **Opsi 2**: {Nama Opsi 2}
* **Opsi 3**: {Nama Opsi 3}

## 4. Analisis Perbandingan & Trade-offs Tiap Opsi
<!-- OPERATOR VOICE: Tuliskan pertimbangan kritis pribadimu untuk tiap opsi -->

### {Opsi 1}
* **Kelebihan**: {Alasan teknis mengapa opsi ini bagus}.
* **Kekurangan / Alasan Ditolak (jika ditolak)**: {Mengapa opsi ini tidak cocok untuk batasan homelab kita}.

### {Opsi 2}
* **Kelebihan**: {Alasan teknis}.
* **Kekurangan / Alasan Ditolak (jika ditolak)**: {Alasan teknis}.

## 5. Hasil Keputusan & Justifikasi Operator
**Opsi Terpilih**: **"{Opsi X}"**, karena {justifikasi teknis murni dari sudut pandang operator yang mengacu langsung pada Decision Drivers}.

### Konsekuensi & Kompromi Teknis (Trade-offs)
* **Dampak Positif**: {Keuntungan nyata yang didapat}.
* **Utang Teknis / Mitigasi**: {Kerumitan tambahan yang harus ditangani akibat memilih opsi ini dan cara memitigasinya}.

## 6. Konfirmasi Implementasi & Kuantifikasi Metrik
* **Metrik Terukur (Before vs After)**: {Contoh: Penggunaan RAM turun dari X MB ke Y MB / Ukuran image dari X MB ke Y MB}.
* **Perintah Verifikasi CLI**:
  ```bash
  {Perintah CLI dan hasil output singkat untuk memverifikasi implementasi di homelab}
  ```
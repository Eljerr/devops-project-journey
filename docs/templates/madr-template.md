---
adr-id: ADR-XXXX
status: "{proposed | accepted | superseded by ADR-YYYY}"
date: YYYY-MM-DD
review-date: YYYY-MM-DD
author: Asylon
hardware-constraint: "{INNER CachyOS 16GB | OUTER Proxmox 2C/8GB | OUTER STB ARM64 2GB | Hybrid}"
tags: [adr, architecture, fase-x]
---

# ADR-XXXX: Gunakan {Opsi Terpilih} untuk {Masalah}

## 1. Konteks dan Masalah
<!-- OPERATOR VOICE: masalah nyata di homelab + batas fisiknya -->
{…}

## 2. Pendorong Keputusan (Decision Drivers)
* {misal: RAM idle < 512 MB di STB}
* {misal: harus jalan di ARM64}
* {misal: tanpa IP publik & tanpa biaya cloud}

## 3. Opsi yang Dipertimbangkan
* **Opsi 1**: {…}
* **Opsi 2**: {…}
* **Opsi 3**: {…}

## 4. Analisis Trade-off
<!-- OPERATOR VOICE. Klaim teknis diberi label [LAB ✓]/[DOK]/[INGATAN] -->

### {Opsi 1}
* **Kelebihan**: {…}
* **Kekurangan / alasan ditolak**: {…}

### {Opsi 2}
* **Kelebihan**: {…}
* **Kekurangan / alasan ditolak**: {…}

## 5. Keputusan & Justifikasi
**Opsi terpilih**: **{Opsi X}**, karena <!-- OPERATOR VOICE --> {…}

### Konsekuensi
* **Positif**: {…}
* **Utang teknis & mitigasi**: {…}

## 6. Konfirmasi Implementasi
* **Metrik sebelum vs sesudah**: {…}
* **Perintah verifikasi**:
  ```bash
  {…}
  ```
* **Bukti**: `docs/labs/{lab-id}.md`

## 7. Sumber
* {dokumentasi resmi yang dipakai}

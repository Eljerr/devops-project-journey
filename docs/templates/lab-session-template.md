---
lab-id: LAB-<fase>.<sub>-<nn>
date: YYYY-MM-DD
module-phase: "Fase X.Y"
environment: "{INNER CachyOS | OUTER Proxmox | OUTER STB}"
author: Asylon
hint-level-max: "{0 | 1 | 2 | 3 | 4}"
tags: [lab, evidence, fase-x]
---

# Lab: {Judul Eksperimen}

> Berkas ini adalah **bukti mentah**. Isinya tidak dipoles. Dokumen lain (concept/PM/runbook) mengutip dari sini.

## 1. Pertanyaan yang Ingin Dibuktikan
<!-- OPERATOR VOICE -->
{Satu kalimat: apa yang mau dibuktikan? Contoh: "Apakah data nol lolos OOM karena kompresi, atau karena swap diizinkan?"}

## 2. Kondisi Awal (Versi & Lingkungan)
```text
{Tempel output: uname -r, versi alat terkait, dan kondisi awal sistem}
```

## 3. Variabel
* **Yang diubah (hanya satu)**: {…}
* **Yang dijaga tetap**: {…}

## 4. Prediksi (ditulis SEBELUM menjalankan)
<!-- OPERATOR VOICE: jangan diedit setelah melihat hasil -->
{Prediksi + alasan.}

## 5. Perintah & Output Asli
```text
{Tempel perintah dan output utuh, apa adanya}
```

## 6. Prediksi vs Kenyataan
| Prediksi | Kenyataan | Cocok? |
| :--- | :--- | :--- |
| {…} | {…} | ✅/❌ |

## 7. Kesimpulan & Fakta yang Terbukti
<!-- OPERATOR VOICE -->
* {Fakta} → naikkan status di `ai-context/knowledge/fase-N.md` menjadi ☑ dengan tautan ke berkas ini.

## 8. Pertanyaan Baru yang Muncul
* {…}

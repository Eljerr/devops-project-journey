---
incident-id: PM-XXXX
date: YYYY-MM-DD
module-phase: "Fase X.Y"
environment: "{INNER CachyOS | OUTER Proxmox | OUTER STB}"
author: Asylon
severity: "{Critical | High | Medium}"
status: "{Investigating | Resolved}"
verification: "{Semua klaim LAB/DOK | Masih ada klaim INGATAN}"
evidence: ["docs/labs/{lab-id}.md"]
tags: [post-mortem, troubleshooting, fase-x]
---

# Post-Mortem: {Judul Kegagalan}

## 1. Ringkasan & Dampak
<!-- OPERATOR VOICE: ceritakan dengan bahasamu sendiri -->
{Layanan apa yang terganggu, gejalanya, durasinya.}

## 2. Timeline
| Waktu | Kejadian |
| :--- | :--- |
| HH:MM | Kerusakan disuntikkan / gejala pertama |
| HH:MM | Hipotesis pertama diuji |
| HH:MM | Akar masalah ditemukan |
| HH:MM | Pulih & terverifikasi |

## 3. Gejala Awal & Bukti Log
<!-- RAW ARTIFACT: tempel output asli tanpa diedit -->
```text
{Log asli}
```

## 4. Kronologi Investigasi
<!-- OPERATOR VOICE: termasuk tebakan yang salah -->
1. **Hipotesis awal**: {…} → perintah uji: `{…}` → hasil: {…}
2. **Temuan menengah**: {…}
3. **Titik kegagalan**: {…}

## 5. Akar Masalah (5 Whys)
<!-- OPERATOR VOICE. Setiap "Why" yang menyebut mekanisme kernel diberi label [LAB ✓]/[DOK]/[INGATAN] -->
* **Gejala**: {…}
  * **Why 1**: {…}
  * **Why 2**: {…}
  * **Why 3**: {…}
  * **Mekanisme internal**: {…} [label]

## 6. Resolusi & Bukti Pemulihan
* **Perbaikan**: {…}
* **Bukti**:
  ```text
  {Output setelah perbaikan}
  ```

## 7. Action Items
| # | Tindakan | Jenis (cegah/deteksi/mitigasi) | Status |
| :--- | :--- | :--- | :--- |
| 1 | {…} | {…} | ☐ |

## 8. Refleksi Operator
<!-- OPERATOR VOICE: apa yang kamu pelajari, apa yang akan kamu lakukan berbeda -->
{…}

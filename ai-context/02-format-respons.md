# 02 — Format Respons Wajib per Tahap

> Model AI kecil lebih konsisten jika mengikuti cetakan. Isi bagian `<...>`. Jangan menambah bagian baru tanpa alasan.

---

## Header Status (WAJIB di baris paling atas SETIAP balasan)

```
📍 Fase <N> · Sub-modul <N.M> <nama singkat> · Tahap <T0–T5> · Target: <INNER: CachyOS | OUTER: Proxmox | OUTER: STB | —>
📚 Sumber dibaca: STATE.md ✓ · knowledge/fase-<N>.md ✓ · <lainnya>
```

Jika tahap tidak jelas, tulis `Tahap: ?` lalu tanya operator. Jangan menebak.

---

## T0 — Pemanasan (awal sesi saja)

```
<Header Status>

### 🔁 Pemanasan
Pertanyaan recall dari <modul lama>: <1 pertanyaan>
(Jawab dengan kata-katamu sendiri, tanpa membuka catatan.)
```
HARD STOP.

---

## T1 — Teori Mikro

```
<Header Status>

### 🧠 Konsep: <nama konsep>  (maks 1–2 konsep, ±400 kata)

**Analogi:** <analogi visual singkat>

**Mekanisme:**
KLAIM: <...> [label]
BUKTIKAN: <perintah> → yang diharapkan: <...>
(ulangi untuk 2–4 klaim inti)

**Batas analogi:** <di mana analogi tadi tidak lagi akurat>

### ❓ Uji Pemahaman
1. <pertanyaan "mengapa">
2. <pertanyaan prediksi: "Jika kamu menjalankan X, apa yang akan terjadi dan mengapa?">
```
HARD STOP. Jangan memberi jawaban.

---

## T2 — Uji Konsep (menilai jawaban)

```
<Header Status>

### 📝 Penilaian
| Kriteria rubrik | Status | Catatan |
| :--- | :--- | :--- |
| R1 Kata-kata sendiri | ✅/❌ | <...> |
| R2 Mekanisme benar | ✅/❌ | <bagian yang salah + koreksi [label]> |
| R3 Prediksi | ✅/❌ | <...> |

**Keputusan:** LULUS → lanjut T3 di balasan berikutnya
            | BELUM → <penjelasan ulang dengan analogi berbeda> + <1 pertanyaan baru>
**Miskonsepsi dicatat:** <jika ada, kalimat persis operator + koreksinya>
```
HARD STOP.

---

## T3 — Lab & Chaos

```
<Header Status>

### 🧪 Lab: <judul>
**Target:** [INNER/OUTER: ...] · **Estimasi RAM:** <...> · **Risiko:** <apa yang bisa rusak> · **Rollback:** <cara kembali>

**Tujuan yang harus dibuktikan:** <1–2 kalimat>

**Langkah 0 — Cek lingkungan:** <perintah cek versi/kondisi awal>
**Langkah 1 — Prediksi (tulis sebelum menjalankan):** <pertanyaan prediksi>
**Langkah 2 — Bangun:** <tujuan per langkah + petunjuk perintah/man page, BUKAN konfigurasi jadi>
**Langkah 3 — 💥 Chaos:** <kerusakan yang disuntikkan atau dicari>
**Langkah 4 — Diagnosis:** <alat yang boleh dipakai, tanpa memberi jawaban>

**Kirimkan:** prediksimu + output terminal asli (tempel utuh, jangan diringkas).
```
HARD STOP.

---

## T4 — Verifikasi Bukti

```
<Header Status>

### 🔍 Prediksi vs Kenyataan
| Prediksi operator | Output asli | Cocok? | Penjelasan [label] |

### ✅ Fakta yang Kini Terbukti [LAB ✓]
- <fakta> (bukti: <baris output>)

### 🎤 Wawancara Reflektif
1. Apa hipotesismu sebelum menjalankan?
2. Bagaimana analisismu terhadap kegagalan ini?
3. Apa pertimbangan teknis pribadimu jika ini terjadi di STB/Proxmox?
```
HARD STOP.

---

## T5 — Dokumentasi

```
<Header Status>

### 📄 Kerangka Dokumen
Templat yang dibaca: docs/templates/<nama>.md ✓
Lokasi berkas: docs/<folder>/<nama>.md
<kerangka dengan bagian OPERATOR VOICE dikosongkan atau berisi kutipan verbatim operator>

### 🗂️ Pembaruan STATE.md
<diff singkat: tahap, kompetensi, miskonsepsi, gate>
```

---

## Batas Panjang

* T1: ±400 kata. T3: ±500 kata. Lainnya: sesingkat mungkin.
* Jika materi terlalu panjang, pecah menjadi beberapa putaran T1. Lebih baik pelan daripada membanjiri.

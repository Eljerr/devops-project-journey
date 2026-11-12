# 03 — Asesmen, Tangga Petunjuk, dan Perintah Operator

---

## 1. Rubrik Kelulusan

### Lulus Uji Konsep (T2) — wajib R1 + R2 + R3

| Kode | Kriteria | Contoh LULUS | Contoh BELUM |
| :--- | :--- | :--- | :--- |
| **R1** | Menjelaskan dengan **kata-kata sendiri** (bukan menyalin istilah agen) | "rm cuma nyopot nama file dari direktori, datanya tetap ada selama masih ada proses yang megang" | "unlink menghapus dentry sehingga nlink 0" (benar, tapi hanya mengulang jargon) |
| **R2** | Mekanisme **benar** sesuai `[DOK]`/`[PACK]` | Menyebut syarat data dibebaskan: link 0 **dan** tidak ada FD terbuka | "Datanya hilang setelah reboot aja" |
| **R3** | Bisa **memprediksi** perilaku pada skenario baru | "Kalau prosesnya di-kill, df langsung turun" | Tidak bisa menjawab skenario baru |

### Lulus Sub-modul (setelah T4) — wajib semua

- [ ] R1–R3 terpenuhi.
- [ ] **R4 Bukti lab**: ada output terminal asli yang mendukung, tersimpan di chat atau `docs/labs/`.
- [ ] **R5 Diagnosis chaos**: operator menemukan penyebab kerusakan sendiri dengan maks. petunjuk Level 2.
- [ ] **R6 Transfer**: operator bisa menjelaskan bagaimana konsep ini muncul di STB/Proxmox/Kubernetes.

Jika operator butuh petunjuk Level 3–4 untuk chaos, sub-modul **belum LULUS**. Ulangi dengan skenario chaos lain.

### Lulus Fase (Gate)

- [ ] Semua sub-modul LULUS.
- [ ] Semua artefak DoD di `ROADMAP.md` terbit dan bagian OPERATOR VOICE diisi operator.
- [ ] **Blind Drill** lulus (lihat bagian 4).
- [ ] Cek silang dokumen fase oleh model AI lain atau operator sudah membuktikan semua klaim berlabel `[INGATAN]`.

---

## 2. Tangga Petunjuk (Hint Ladder)

Saat operator mentok, naikkan **satu level per permintaan**. Jangan langsung lompat ke jawaban.

| Level | Isi petunjuk | Contoh |
| :--- | :--- | :--- |
| 1 — Arah | Bagian sistem mana yang perlu dilihat | "Coba lihat sisi kernel, bukan sisi aplikasi." |
| 2 — Alat | Nama perintah/berkas yang relevan | "Ada perintah yang bisa mencari file terbuka yang sudah dihapus." |
| 3 — Sintaks | Potongan perintah dengan bagian kosong | "`lsof +L___`" |
| 4 — Jawaban | Jawaban lengkap + penjelasan | Dicatat di STATE.md sebagai "dibantu penuh". |

---

## 3. Pemanasan & Pengulangan (Spaced Repetition)

* Di awal setiap sesi (T0), ambil 1 pertanyaan dari modul LULUS. Prioritaskan:
  1. Modul dengan miskonsepsi tercatat di `STATE.md`.
  2. Modul yang paling lama tidak diulang (lihat kolom "Terakhir diulang").
* Jika operator gagal menjawab: catat di STATE.md, jelaskan singkat, lalu lanjut. Jika gagal 2 kali di modul yang sama → jadwalkan **sesi ulang** modul tersebut sebelum lanjut.

---

## 4. Blind Drill (Ujian Akhir Fase)

Tujuan: membuktikan operator bisa bekerja **tanpa AI**.

1. Agen menulis skrip/instruksi kerusakan dalam satu blok, operator menjalankannya **tanpa membaca isinya** (atau minta orang/agen lain menyiapkan).
2. Operator mendiagnosis **tanpa bertanya ke AI** dengan batas waktu (misal 30 menit).
3. Operator mencatat semua perintah dan hipotesis di `docs/labs/` memakai `lab-session-template.md`.
4. Setelah selesai, agen menilai: penyebab ditemukan? perbaikan benar? jalur diagnosis logis?

---

## 5. Perintah Operator

Agen wajib mengenali perintah berikut di awal pesan operator:

| Perintah | Aksi agen |
| :--- | :--- |
| `!status` | Tampilkan posisi dari STATE.md (fase, sub-modul, tahap, gate). |
| `!buktikan` | Beri sumber `[DOK]` + perintah terminal untuk klaim terakhir. |
| `!audit` | Tabel semua klaim di balasan terakhir: klaim · label · yakin (T/S/R) · cara membuktikan. |
| `!hint` | Naik 1 level di tangga petunjuk. |
| `!ulang` | Jelaskan konsep yang sama dengan analogi yang **berbeda**. |
| `!lebih-dalam` | Turun satu lapis mekanisme (misal dari perintah → syscall → struktur kernel). |
| `!review` | Kuis 3 pertanyaan dari modul yang sudah LULUS. |
| `!skip-teori` | Hanya jika operator sudah LULUS konsep ini sebelumnya: langsung T2 (uji konsep). |
| `!selesai` | Tutup sesi: perbarui STATE.md (tahap, kompetensi, miskonsepsi, tanggal). |

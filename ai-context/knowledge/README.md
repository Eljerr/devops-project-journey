# Paket Pengetahuan (Knowledge Packs)

Paket ini adalah **kunci jawaban** untuk agen AI: fakta inti, miskonsepsi umum, dan eksperimen pembuktian per fase.
Tujuannya agar model AI kecil tidak mengandalkan ingatannya sendiri.

## Aturan Pemakaian

1. Agen membaca `fase-<N>.md` untuk fase aktif **sebelum** mengajar.
2. Agen **tidak** menyalin pack ke operator mentah-mentah. Pack dipakai untuk memeriksa jawaban dan merancang pertanyaan.
3. Setiap fakta punya kotak status:
   * `☐` = belum dibuktikan operator di lab (masih `[DOK]`/`[INGATAN]`).
   * `☑` = sudah dibuktikan operator, bukti ditautkan (`docs/labs/...`). Status naik menjadi `[LAB ✓]`.
4. Jika agen atau operator menemukan fakta yang salah: **jangan hapus diam-diam**. Coret (`~~teks~~`), tulis koreksi + sumber + tanggal.
5. Fakta berlabel `[INGATAN]` adalah prioritas pembuktian. Jangan dipakai sebagai dasar kelulusan.

## Membuat Pack untuk Fase Baru

Sebelum fase baru dimulai, agen membuat `fase-<N>.md` dari `_template.md`:
1. Gunakan tool dokumentasi (Context7 / web search / `man`) untuk setiap fakta. Tulis sumbernya.
2. Fakta tanpa sumber → label `[INGATAN]`.
3. Operator membaca sekilas, lalu menjalankan semua perintah `BUKTIKAN` sepanjang fase.
4. Disarankan: minta model AI lain mengaudit pack ("cari kesalahan faktual, beri sumber resmi").

## Status Pack

| Berkas | Status | Catatan |
| :--- | :--- | :--- |
| `fase-1.md` | Draf v1 (2026-11-12) | Disusun model besar + cek web. Semua fakta masih `☐` sampai dibuktikan ulang di Fase 1 v2. |
| `fase-2.md` | Draf v1 (2026-11-12) | Sama. Fakta jaringan tingkat lanjut (Docker/nftables) bergantung versi — cek versi dulu. |
| `fase-3+` | Belum ada | Dibuat oleh agen di awal fase terkait. |

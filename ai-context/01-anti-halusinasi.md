# 01 — Protokol Anti-Halusinasi

> Tujuan: operator bisa **membuktikan sendiri** setiap penjelasan, tanpa harus percaya kepada AI.
> Prinsip utama: **Bukti mengalahkan otoritas.** Termasuk otoritas AI dan otoritas berkas ini.

---

## 1. Hierarki Kebenaran

Jika dua sumber bertentangan, sumber yang lebih tinggi menang:

| Peringkat | Sumber | Label |
| :--- | :--- | :--- |
| 1 | Output terminal **asli** dari mesin operator (sudah ditempel di chat atau di `docs/labs/`) | `[LAB ✓]` |
| 2 | Dokumen resmi: man page (`man 7 cgroups`), docs.kernel.org, dokumentasi resmi proyek (nginx.org, docs.docker.com, freedesktop.org/systemd, RFC) | `[DOK]` + nama/URL |
| 3 | Paket pengetahuan `ai-context/knowledge/fase-N.md` | `[PACK]` |
| 4 | Ingatan model AI | `[INGATAN]` — wajib diberi "cara membuktikan" |

Aturan turunan:
* Klaim `[INGATAN]` **tidak boleh** dipakai sebagai dasar kelulusan atau dokumentasi sebelum naik menjadi `[DOK]` atau `[LAB ✓]`.
* Jika `[PACK]` bertentangan dengan `[DOK]`/`[LAB ✓]`: jangan diam-diam diabaikan. Laporkan ke operator dan usulkan koreksi pack.

---

## 2. Format Klaim Teknis

Setiap klaim mekanisme ditulis seperti ini:

```
KLAIM: <satu kalimat>  [label]
BUKTIKAN: <perintah yang bisa dijalankan operator>  → yang diharapkan: <apa yang harus terlihat>
```

Contoh (benar):
```
KLAIM: Setiap proses punya daftar file descriptor yang bisa dilihat di /proc/<PID>/fd.  [DOK: man 5 proc]
BUKTIKAN: ls -l /proc/self/fd  → terlihat 0, 1, 2 (stdin, stdout, stderr) sebagai symlink.
```

Jika tidak bisa menulis baris `BUKTIKAN`, klaim itu harus diberi label `[INGATAN]` dan kalimat "belum bisa saya buktikan sekarang".

---

## 3. Sepuluh Aturan Anti-Halusinasi

1. **Jangan mengarang output.** Setiap contoh output yang bukan dari operator diberi label `ILUSTRASI — bukan output asli`.
2. **Jangan mengarang flag/path.** Sebelum menyebut flag, file sysfs, atau opsi config yang tidak 100% pasti, minta operator mengecek dengan `--help`, `man`, atau `ls`.
3. **Cek versi dulu.** Sintaks yang bergantung versi (systemd, Docker, nginx, k3s, Terraform) baru boleh diberikan setelah versi diketahui (`STATE.md` atau perintah `--version`).
4. **Cek lingkungan dulu.** Jangan menganggap Ubuntu, bash, ext4, atau WSL2. Baca "Fakta Lingkungan" di `STATE.md`.
5. **Satu variabel per eksperimen.** Jika dua hal berubah sekaligus (misal jenis data *dan* swap), kesimpulan tidak sah. Rancang ulang eksperimennya.
6. **Prediksi sebelum eksekusi.** Operator menulis prediksi → menjalankan → membandingkan. Selisih prediksi vs hasil = bahan belajar paling berharga.
7. **Bedakan "terjadi" dan "mengapa terjadi".** Output membuktikan *apa* yang terjadi. Penjelasan *mengapa* tetap butuh `[DOK]` atau eksperimen pembanding.
8. **Katakan "saya tidak yakin".** Ini jawaban yang sah dan dihargai. Menebak dengan nada yakin adalah pelanggaran.
9. **Jangan menjilat (anti-sycophancy).** Jika jawaban operator salah, katakan "Ini belum tepat" + alasannya. Jika operator membantah tanpa bukti, jangan langsung setuju; ajak membuktikan.
10. **Konsisten dengan riwayat.** Jika penjelasanmu bertentangan dengan dokumen lama di repo, sebutkan kontradiksinya secara eksplisit. Jangan diam-diam menimpa.

---

## 4. Prosedur Saat Ragu (Decision Tree)

```
Apakah saya bisa menyebut sumber resmi untuk klaim ini?
├─ Ya → tulis [DOK] + sumber + baris BUKTIKAN
└─ Tidak → Apakah ada tool dokumentasi (Context7 / web search)?
     ├─ Ya → cari dulu, lalu jawab dengan [DOK]
     └─ Tidak → tulis [INGATAN] + "saya tidak yakin" + perintah agar operator mengecek sendiri
```

---

## 5. Pemeriksaan Diri Sebelum Mengirim Balasan

Agen wajib mengecek daftar ini dalam hati sebelum mengirim:

- [ ] Apakah Header Status ada dan tahapnya benar?
- [ ] Apakah balasan ini hanya berisi **satu** tahap?
- [ ] Apakah setiap klaim teknis punya label sumber?
- [ ] Apakah ada output terminal yang saya karang tanpa label ILUSTRASI?
- [ ] Apakah ada flag/path yang saya tidak yakin ada?
- [ ] Apakah sintaks cocok dengan shell **fish** dan OS yang tercatat?
- [ ] Apakah saya berhenti di titik HARD STOP yang benar?

---

## 6. Perintah Audit dari Operator

Operator dapat mengetik perintah ini kapan saja (daftar lengkap di `03-asesmen.md`):

* `!buktikan` → agen memberikan sumber `[DOK]` dan perintah terminal untuk klaim terakhir.
* `!audit` → agen membuat tabel semua klaim di balasan terakhir: klaim, label, tingkat yakin (Tinggi/Sedang/Rendah), cara membuktikan. Klaim berlabel Rendah wajib dibuktikan sebelum lanjut.

---

## 7. Cek Silang Berkala (Disarankan)

Karena model AI kecil bisa salah dengan nada yakin:
* Di akhir setiap sub-modul, operator boleh menyalin isi dokumen konsep ke **model AI lain** dengan perintah: *"Cari kesalahan faktual di teks ini. Untuk tiap kesalahan, beri sumber resmi."*
* Perbedaan pendapat antar-model diselesaikan dengan **eksperimen di terminal**, bukan dengan memilih model mana yang lebih dipercaya.

---
status: "{proposed | accepted | superseded by ADR-XXXX}"
date: YYYY-MM-DD
decision-makers: [Eljerr]
consulted: [AI Systems Architect Agent]
---

# Gunakan {Opsi Terpilih} untuk {Masalah yang Diselesaikan}

## Konteks dan Pernyataan Masalah

{Jelaskan kondisi infrastruktur saat ini, batasan perangkat keras (misal: batas RAM 8GB pada Proxmox Laptop atau 2GB pada STB ARM64), serta masalah teknis yang harus diselesaikan tanpa menyisipkan kesimpulan awal.}

## Pendorong Keputusan (Decision Drivers)

* {Batasan 1, misal: Konsumsi memori idle harus di bawah 512 MB}
* {Batasan 2, misal: Kompatibilitas arsitektur CPU x86_64 / ARM64}
* {Batasan 3, misal: Keamanan jaringan tanpa IP Publik statis}

## Opsi yang Dipertimbangkan

* {Opsi 1}
* {Opsi 2}
* {Opsi 3}

## Hasil Keputusan

Opsi terpilih: "{Opsi X}", karena {justifikasi teknis yang mengacu langsung pada Decision Drivers}.

### Konsekuensi

* Positif: {Dampak baik terukur, misal: Penghematan RAM sebesar X MB, waktu boot Y detik}.
* Negatif / Trade-offs: {Kompromi teknis yang harus diterima dan dimitigasi, misal: Konfigurasi manual cgroup2 pada host Proxmox diperlukan}.

### Konfirmasi Implementasi

{Perintah CLI atau metrik Prometheus yang digunakan untuk memverifikasi bahwa keputusan ini berjalan sesuai spesifikasi di homelab.}

## Analisis Kelebihan dan Kekurangan Tiap Opsi

### {Opsi 1}

* Baik, karena {alasan teknis}.
* Buruk, karena {alasan mengapa opsi ini membebani sistem atau kurang tepat}.

### {Opsi 2}

* Baik, karena {alasan teknis}.
* Buruk, karena {alasan teknis}.

# ADR-0001: Peran STB sebagai Edge Gateway dan laptop Proxmox sebagai Compute Cluster

* **Status:** Accepted
* **Tanggal:** 2026-10-03
* **Penulis:** Asylon

## 1. Konteks & Masalah (Context & Problem Statement)
STB saya putuskan sebagai Edge Gateway karena keterbatasan ram selain itu memiliki beberapa Keuntungan sebagai Proxy dimana jika terjadi peretasan hanya sampai pada STB dan minimalisasi area dampak eksploitasi tersebut sehingga memberikan waktu untuk proses perlindungan kepada Cluster server dan memudahkan Filtering. 

Laptop Proxmox memiliki ram yang lebih lega sehingga sangat memungkinkan untuk menjalankan beberapa LXC dimana bisa dibagi 1-Server (2-Agent) dan database. Saya memutuskan menggunakan k3s dibandingkan dengan k8s karena k3s lebih ramah untuk low-resource environments dan fitur-fiturnya tidak kalah lengkap. Saya juga memutuskan memakai LXC karena hemat resource dibanding VM.

> [!NOTE]
> Laptop yang saya install Proxmox hanya sekadar kebutuhan praktek tidak digunakan untuk mendeploy aplikasi dan juga Laptop tersebut masih dipakai untuk kebutuhan kuliah dengan dualboot CachyOS

## 🖥️ Infrastructure Overview

This entire project runs on **personal/repurposed hardware** — no cloud provider needed.

| Role | Device | Specs |
|---|---|---|
| **Proxmox Server** | Used Laptop | 2 Cores, 8 GB RAM — hosts K3s cluster inside LXC containers |
| **Gateway** | STB Amlogic S905X4 | 2 GB RAM, Ubuntu OS — reverse proxy & Cloudflare Tunnel entry point |

---

## 2. Pilihan yang Dipertimbangkan (Considered Options)
* **Opsi 1:** Menjalankan semua workload di Proxmox Laptop
* **Opsi 2:** Menjadikan STB sebagai Gateway/Reverse Proxy serta Monitoring dan Laptop khusus cluster

## 3. Keputusan yang Diambil (Decision Outcome)
Pilihan yang dipilih: **Opsi 2**.
Opsi ini merupakan pilihan terbaik karena kebetulan saya mempunyai STB Android yang tidak terpakai dan sangat hemat daya. Laptop yang digunakan juga masih sanggup untuk memenuhi kebutuhan saya saat ini. Keterbatasan arsitektur juga dipertimbangkan sehingga saya memutuskan agar semua server dijalankan di Proxmox sehingga STB hanya sebagai gerbang dan monitoring

## 4. Konsekuensi & Kompromi (Consequences & Trade-offs)
* ➕ **Positif:** LXC bisa dibuat lebih banyak (kebutuh praktek), jika cluster down proxy tetap bisa menyampaikan informasi kepada internet, memudahkan monitoring, dan memudahkan debbuging.
* ➖ **Negatif/Tantangan:** Maintenance sedikit ribet karena infra terpisah, jika laptop dipakai kuliah server akan mati, dan perbedaan arsitektur CPU.

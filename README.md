# DevOps Project Journey

> Jurnal belajar + portofolio *Proof of Work* DevOps dari nol di homelab sendiri (tanpa cloud berbayar).
> Penulis: Asylon

## Tentang Proyek Ini
<!-- OPERATOR VOICE: kenapa kamu membuat proyek ini, dengan bahasamu sendiri -->
{…}

## Homelab
| Peran | Perangkat | Spesifikasi |
| :--- | :--- | :--- |
| Workstation (Inner Loop) | PC CachyOS | Ryzen 5 2600, 16 GB RAM |
| Compute (Outer Loop) | Laptop Proxmox VE | 2 core, 8 GB RAM |
| Edge Gateway (Outer Loop) | STB Amlogic S905X4 | ARM64, 2 GB RAM, Ubuntu |

## Cara Saya Belajar
* Kurikulum: [ROADMAP.md](ROADMAP.md) · Posisi saat ini: [STATE.md](STATE.md)
* Setiap konsep dibuktikan lewat eksperimen di terminal ([docs/labs/](docs/labs/)), bukan sekadar dibaca.
* AI dipakai sebagai tutor dengan aturan anti-halusinasi ([AGENTS.md](AGENTS.md), [ai-context/](ai-context/)).

## Portofolio
| Jenis | Dokumen |
| :--- | :--- |
| Keputusan arsitektur (ADR) | [docs/decisions/](docs/decisions/) |
| Post-mortem | [docs/post-mortems/](docs/post-mortems/) |
| Runbook | [docs/runbooks/](docs/runbooks/) |
| Catatan konsep | [docs/concepts/](docs/concepts/) |

## Progres
<!-- OPERATOR VOICE: perbarui di akhir setiap fase -->
- [ ] Fase 1 — Kernel Linux (diulang dengan standar v2)
- [ ] Fase 2 — Jaringan & firewall
- [ ] Fase 3 — Otomatisasi
- [ ] Fase 4 — Kontainer
- [ ] Fase 5 — Proxmox & Terraform
- [ ] Fase 6 — K3s
- [ ] Fase 7 — CI/CD & GitOps
- [ ] Fase 8 — Observabilitas

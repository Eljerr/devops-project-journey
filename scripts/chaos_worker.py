import os
import sys
import time

print(f"[START] Worker menyala dengan PID: {os.getpid()}", flush=True)
print("[INFO] Memberi jeda 2 detik...", flush=True)
time.sleep(2)

bytes_target = 80 * 1024 * 1024
print(
    f"[CHAOS] Mengalokasikan RAM sebesar {bytes_target / (1024 * 1024)} MB...",
    flush=True,
)

makan_ram = os.urandom(bytes_target)

print("[SUKSES] Alokasi berhasil tanpa terbunuh!", flush=True)
time.sleep(5)

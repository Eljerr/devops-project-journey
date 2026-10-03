# 🐍 Asynchronous Service Health Monitor

A lightweight, asynchronous Python monitoring script designed to concurrently check the health status and availability of multiple target servers and web services.

This script demonstrates how to utilize asynchronous IO (`asyncio` and `aiohttp`) in Python for fast network polling and alerting without blocking, a critical technique in DevOps and Reliability Engineering for custom automated health checks.

---

## ⚡ Features

- **Asynchronous Concurrent Polling**: Uses `aiohttp` and `asyncio` to send non-blocking HTTP requests to multiple server endpoints simultaneously.
- **Status Classification & Logging**:
  - **`UP` (HTTP 200)**: Service is healthy and operating normally.
  - **`DOWN` (Non-200 HTTP codes)**: Logs a warning indicating further intervention is needed (`Perlu tindak lebih lanjut...`).
  - **`ERROR` (Network/Connection/DNS Failures)**: Logs an error when a server is unreachable or DNS resolution fails (`Server tidak merespon`).
- **Structured JSON Reporting**: Exports the compiled check results and error details into a formatted JSON report file via command-line argument.

---

## 📁 Folder Structure

```text
01-monitoring-scripts/
├── 📄 monitor.py          # Core asynchronous monitoring script
├── 📄 requirements.txt    # Python dependency definitions (aiohttp, etc.)
├── 📄 healt_report.json   # Sample generated JSON output report
├── 📄 test.json           # Sample test report
├── 📄 .gitignore          # Git ignore configuration for venv & pycache
└── 📁 venv/               # Local Python virtual environment
```

---

## 🚀 Usage Guide

### 1. Set Up the Virtual Environment

Ensure your virtual environment is activated and dependencies are installed:

```bash
# Navigate to the monitoring scripts directory
cd scripts/python/01-monitoring-scripts

# Activate the virtual environment
source venv/bin/activate

# Install dependencies (if not already set up)
pip install -r requirements.txt
```

### 2. Run the Monitoring Script

Execute `monitor.py` and specify the output JSON path using the `--output` argument:

```bash
python monitor.py --output health_report.json
```

### 3. Example Output (`health_report.json`)

```json
[
    {
        "name": "Nginx Web",
        "status": "UP"
    },
    {
        "name": "NodeJS API",
        "status": "DOWN"
    },
    {
        "name": "Payment Gateway",
        "status": "ERROR",
        "detail": "Cannot connect to host..."
    }
]
```

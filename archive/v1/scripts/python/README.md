# 🐍 Python Automation Scripts

This directory houses Python-based automation and monitoring tools built for infrastructure observability and reliability tasks.

---

## 📁 Modules

| Directory | Name | Highlights |
|---|---|---|
| **[01-monitoring-scripts](./01-monitoring-scripts/)** | Asynchronous Service Health Monitor | Concurrent non-blocking HTTP polling with `aiohttp` and `asyncio`, JSON diagnostic exports, error resilience. |

---

## 🛠️ General Setup

We recommend running all Python scripts within an isolated virtual environment:

```bash
# Create a virtual environment
python3 -m venv venv

# Activate the virtual environment
source venv/bin/activate

# Install required packages per script module
pip install -r <module>/requirements.txt
```

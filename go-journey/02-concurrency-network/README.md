# ⚡ Go Concurrency & Networking

This module explores Go's flagship feature — high-performance concurrency — and network programming for site reliability engineering.

---

## 📂 Directory Layout

```text
02-concurrency-network/
├── 📁 exercise/
│   └── 📁 01-server-audit/           # Concurrent server audit with WaitGroup & Context timeout
│       └── 📄 main.go
└── 📁 project-service-monitor/       # Full CLI: Non-blocking HTTP health checker with JSON export
    ├── 📄 main.go
    ├── 📄 hasil.json                 # Sample generated audit report
    └── 📄 README.md
```

---

## 🔑 Core Concepts Practiced

1. **Goroutines (`go func()`)**: Lightweight execution threads managed by Go runtime.
2. **`sync.WaitGroup`**: Synchronizing multiple concurrent worker routines.
3. **Channels (`chan jobResult`)**: Thread-safe communication and result collection across routines.
4. **Context & Timeouts (`context.WithTimeout`)**: Preventing dangling connections and implementing SLA timeouts.
5. **Structured Logging (`log/slog`)**: JSON-based logging suitable for production log collectors (Loki, Fluentbit).

---

## 🚀 Quick Start

### Run Concurrency Audit Exercise
```bash
cd exercise/01-server-audit
go run main.go
```

### Run Service Health Monitor CLI
```bash
cd project-service-monitor
go run main.go -output hasil.json
```

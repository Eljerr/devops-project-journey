# 📡 Project: Concurrent HTTP Service Monitor (Go)

A fast, concurrent command-line utility built in Go to inspect HTTP microservice health status, log diagnostic data with structured JSON logging (`slog`), and write results into a JSON report.

---

## ⚙️ How It Works

1. Defines targets with endpoint URLs.
2. Spawns a dedicated **goroutine** for every target.
3. Performs HTTP GET queries with a strict 5-second `context.WithTimeout`.
4. Streams results through a buffered channel into an ordered slice.
5. Emits structured log entries (`INFO`, `WARN`, `ERROR`) using Go's built-in `log/slog`.
6. Formats and writes the results to a specified JSON file.

---

## 🚀 Usage

Run the program by specifying the `-output` flag:

```bash
# From this directory:
go run main.go -output hasil.json
```

### Command Flags

| Flag | Type | Description |
|---|---|---|
| `-output` | string | **(Required)** Path to write the output JSON report file |

---

## 📄 Output Schema Example (`hasil.json`)

```json
[
  {
    "name": "Nginx Server",
    "status": "UP",
    "detail": "200 OK"
  },
  {
    "name": "Database",
    "status": "DOWN",
    "detail": "HTTP Status: 500"
  },
  {
    "name": "Payment-gateway",
    "status": "DOWN",
    "detail": "HTTP Status: 404"
  }
]
```

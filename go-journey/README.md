# 🐹 Go Journey for DevOps

This directory chronicles my journey learning the **Go programming language (Golang)**, specifically focused on engineering DevOps tools, automation utilities, and reliability monitoring scripts.

Go is widely recognized as the backbone language for modern cloud-native infrastructure (Docker, Kubernetes, Terraform, Prometheus, etc.). This section documents hands-on exercises and practical micro-projects.

---

## 📁 Learning Modules

| Module | Focus & Topics | Projects / Exercises |
|---|---|---|
| **[01-fundamentals](./01-fundamentals/)** | Basic syntax, structs, methods, maps, error handling, slices | Fleet Status Report CLI |
| **[02-concurrency-network](./02-concurrency-network/)** | Goroutines, channels, `sync.WaitGroup`, `context.WithTimeout`, HTTP client | Concurrency audit & Concurrent Service Health Monitor |

---

## 🛠️ Prerequisites

- **Go**: Version 1.20 or newer
- Check your local installation:
  ```bash
  go version
  ```

---

## 🎯 Key Objectives

1. **Systems Automation**: Writing fast, compiled, single-binary CLI tools for server fleet audits.
2. **Concurrent Observability**: Polling microservices and infrastructure endpoints without blocking.
3. **Resiliency**: Proper context cancellation, timeouts, and structured error handling.

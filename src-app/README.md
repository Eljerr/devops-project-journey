# 📦 Application Source Code (`project-02-api`)

This directory contains the source code for the containerized Node.js (Express) backend API utilized across **Project 02** (`projects/02-nodejs-api/`), the **GitHub Actions DevSecOps pipeline**, and the **Jenkinsfile template**.

---

## 🚀 API Endpoints

| Endpoint | Method | Middleware / Features | Purpose |
|---|---|---|---|
| `/` | `GET` | None | Welcome banner, simulated DB status, and displays current `DB_HOST`. |
| `/healtz` | `GET` | None | **Liveness Probe**: Returns HTTP 200 `OK` immediately to indicate process vitality. |
| `/ready` | `GET` | `rateLimit` (max 60 req/min) | **Readiness Probe**: Verifies live TCP & authentication handshake with MySQL database. Returns HTTP 503 if disconnected. |

---

## 🔐 Environment Variables

The application consumes configuration passed from Kubernetes ConfigMaps & Secrets (or local `.env`):

| Variable | Default Value | Description |
|---|---|---|
| `DB_HOST` | `mysql-service` | Hostname or IP of the MySQL database |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | Database username |
| `DB_PASSWORD` | `secret` | Database password |
| `DB_NAME` | `mysql` | Target database name |

---

## 🧪 Local Testing & Development

```bash
# Install dependencies
npm install

# Run Jest unit tests (Quality Gate in CI)
npm test

# Run application locally
npm start
```

---

## 🐳 Docker Build & DevSecOps

The `Dockerfile` is built on top of `node:20-alpine`:
- Executes `apk update && apk upgrade` to patch system-level CVEs before scanning.
- Exposes port `80`.
- Scanned continuously by **Trivy** and analyzed via **CodeQL** in `.github/workflows/ci-devsecops.yml`.

To build and test the container locally:

```bash
docker build -t project-02-api:local .
docker run -p 8080:80 project-02-api:local
```

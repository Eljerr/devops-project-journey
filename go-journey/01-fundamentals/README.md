# 📘 Go Fundamentals

This module covers core Go syntax, types, and programming concepts applied to DevOps server management scenarios.

---

## 📂 Directory Layout

```text
01-fundamentals/
├── 📁 exercise/                 # Step-by-step progressive syntax exercises
│   ├── 📁 00-setup/             # Basic setup & "Hello, World" entry point
│   ├── 📁 01-server-info/       # Variable declarations & basic formatting
│   ├── 📁 02-server-struct/     # Defining structs & constructor functions
│   └── 📁 03-server-registry/   # Maps, slices, and lookups
└── 📁 project-fleet-report/     # Mini-project: Fleet status report generator
    └── 📄 main.go
```

---

## 🏃 Running the Exercises

Navigate into any exercise directory and run:

```bash
cd exercise/01-server-info
go run main.go
```

---

## 🌟 Featured Project: Fleet Status Report

The `project-fleet-report/` mini-project models an infrastructure fleet:
- Defines server entities with attributes (`Name`, `Role`, `IP`, `Port`, `Online`, `Uptime`).
- Classifies servers by uptime stability (`Excellent`, `Good`, `Warning`).
- Filters server registries by roles (e.g., finding all `database` nodes).
- Produces a formatted audit report to standard output.

Run the project:
```bash
cd project-fleet-report
go run main.go
```

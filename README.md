# verification

**A Go-based verification microservice / project for building secure backend workflows.**

## Overview

`verification` is a Golang project focused on implementing identity or transaction verification workflows. It offers modular components for HTTP endpoints, persistent storage, and verification logic, making it easy to build secure services such as 2FA, email/SMS verification, or token-based access validation.

Built in Go with a clean architecture: separate packages for `database`, `http`, and core logic under `src`, the project is designed for clarity and extensibility.

---

## Features

- **HTTP API** — lightweight RESTful endpoints to handle verification requests and validate tokens.
- **Database Integration** — easily connectable storage backend for verification state.
- **Modular Architecture** — separation of concerns across `database`, `http`, and core logic allows for custom extensions.
- **Go-first** - written in idiomatic Go, compatible with Go modules for seamless use.

---

## Getting Started

### Prerequisites

- Go version 1.20 or higher
- A running datastore (PostgreSQL)

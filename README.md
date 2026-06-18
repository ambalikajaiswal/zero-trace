# ZeroTrace

A full-stack zero-trust secret-sharing application with automatic secret destruction, TTL-based expiration, and least-privilege access control.

## Tech Stack

- **Backend:** Go 1.22, `net/http`, goroutines, thread-safe in-memory storage
- **Frontend:** TypeScript, React 18, Vite, React Router
- **Deployment:** Docker, Docker Compose, Nginx

## Features

- **TTL-Based Expiration** — Secrets auto-destruct after a configurable time (5 min to 7 days)
- **View-Once Secrets** — Automatic destruction after first retrieval
- **Least-Privilege Access** — Each secret requires a unique cryptographic token
- **Zero Trust Architecture** — No stored logs, no session tracking, token in URL hash (never sent to server in referer)
- **Concurrent Goroutines** — Background cleanup routines with thread-safe `sync.RWMutex` storage
- **Rate Limiting** — Per-IP request throttling to prevent abuse

## Architecture

```
┌─────────────────┐         ┌──────────────────────────┐
│  React Frontend │  REST   │       Go Backend          │
│  (TypeScript)   │────────▶│  net/http + goroutines    │
│  Vite + Router  │         │  Thread-safe store        │
└─────────────────┘         │  TTL cleanup goroutine    │
                            └──────────────────────────┘
```

## API Endpoints

| Method | Endpoint            | Description           |
|--------|--------------------|-----------------------|
| POST   | `/api/secrets`     | Create a new secret   |
| GET    | `/api/secrets/:id` | Retrieve & destroy    |
| GET    | `/api/health`      | Health check          |

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker (optional)

### Run Backend

```bash
cd backend
go mod download
go run .
```

Server starts on `http://localhost:8080`

### Run Frontend

```bash
cd frontend
npm install
npm run dev
```

App starts on `http://localhost:3000` (proxies API to backend)

### Run with Docker Compose

```bash
docker-compose up --build
```

Frontend on `:3000`, Backend on `:8080`

## Security Model

1. **Secret creation** returns an `id` + cryptographic `token`
2. **Token is placed in URL hash** (`#token`) — never sent to server in HTTP headers/referer
3. **Retrieval requires both** the ID (path) and token (header)
4. **View-once secrets** are immediately deleted from memory after retrieval
5. **Background goroutine** continuously purges expired secrets
6. **Rate limiting** prevents brute-force token guessing

## License

MIT

# Pulse Polls

A live polling platform built with React, Go/Gin, MongoDB, Redis, and WebSockets.

## Architecture

React is the realtime UI. Go/Gin owns authentication, validation, authorization, persistence, vote rules, and realtime transport. MongoDB is the persistent source of truth. Redis is the live aggregation, Pub/Sub, and distributed rate-limit infrastructure. WebSockets carry Redis events to connected browsers.

The backend is split into `internal/app/config.go`, `models.go`, `validation.go`, `repository.go`, `auth_handlers.go`, `poll_handlers.go`, `realtime.go`, and `server.go`. The handlers still use a repository value directly in a few small paths; the next cleanup should move those calls behind service methods if the application grows further.

## Request flow

1. A client submits `POST /api/polls/:id/vote` with only `optionId`.
2. Gin limits the request body and rate limit, then validates the poll, option, state, expiry, and voter cookie.
3. MongoDB inserts the vote. Its unique `(pollId, voterId)` index is the durable duplicate-vote guard.
4. Redis runs `HINCRBY poll:{pollId}:results {optionId} 1`.
5. Redis publishes JSON on `poll:{pollId}:updates`.
6. WebSocket subscribers receive the event and React replaces its result state without polling or refresh.

MongoDB remains authoritative. Redis hashes are rebuildable acceleration state. If a hash is missing, counts are reconstructed from MongoDB and written back to Redis. A Redis failure after MongoDB insertion is reported as degraded; the next results request can rebuild the counter.

## Data model and indexes

- `users`: email, bcrypt password hash, timestamps. Unique `email` index.
- `polls`: public random ID, owner ID, question, options, status, timestamps, optional expiry. Owner/created index.
- `votes`: poll ID, option ID, anonymous voter ID, timestamp. Unique `(pollId, voterId)` and `pollId` indexes.

Owner IDs are omitted from public JSON responses. MongoDB internal ObjectIDs are not exposed because public IDs are random URL-safe strings.

## Redis

- Results hash: `poll:{id}:results`, mapping option ID to count.
- Pub/Sub channel: `poll:{id}:updates`.
- Rate-limit keys: `rate:auth:{client-ip}` and `rate:vote:{client-ip}`.

Redis is not just a cache: it performs live counter increments, event fan-out, recovery-backed result serving, and distributed request throttling. Pub/Sub is intentionally ephemeral; MongoDB plus hash rebuilds provide correctness after missed events or Redis restart.

## Authentication and security

Signup/login use bcrypt and an HTTP-only JWT cookie. JWTs expire after 24 hours and parsing only accepts HS256. Protected routes enforce authentication, and poll updates/deletes enforce ownership. Authentication limits are 20 requests per IP per minute; voting is limited to 60 requests per IP per minute. Exceeding a limit returns HTTP 429.

Bodies are capped at 64 KiB, CORS origins are explicit, WebSocket origins are checked against `ALLOWED_ORIGINS`, WebSocket connections use deadlines and ping/pong, and internal errors are not returned to clients. The anonymous `voter_id` cookie is a basic duplicate-vote guard, not a fraud-proof identity system: clearing or changing cookies can bypass it.

## Local development

Copy `.env.example` to `.env`, set a random `JWT_SECRET` of at least 32 characters, run MongoDB and Redis, then:

```text
cd backend
go test ./...
go run ./cmd/server

cd frontend
npm ci
npm test
npm run build
npm run dev
```

The frontend defaults to same-origin `/api` and WebSockets. For separate local dev servers, set `VITE_API_URL=http://localhost:8080/api` and `VITE_WS_URL=ws://localhost:8080/api`.

## Docker

Set `JWT_SECRET` in `.env`, then run:

```text
docker compose up --build
```

Open `http://localhost:5173`. Nginx serves the SPA and proxies `/api/` plus WebSocket upgrades to the backend. Compose health checks gate backend startup on MongoDB and Redis readiness. MongoDB and Redis are kept on the Compose network rather than published to the host.

## API

`POST /api/auth/signup`, `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/auth/me`

`POST /api/polls`, `GET /api/polls`, `GET /api/polls/:id`, `PATCH /api/polls/:id`, `DELETE /api/polls/:id`

`POST /api/polls/:id/vote`, `GET /api/polls/:id/results`, `GET /api/polls/:id/ws`

## Deployment

Build the frontend as a static image or deploy it to a static host. Deploy the backend to a WebSocket-capable host. Use MongoDB Atlas and Redis Cloud/Upstash. Set `MONGO_URI`, `MONGO_DATABASE`, `REDIS_URL`, `JWT_SECRET`, `ALLOWED_ORIGINS`, and `COOKIE_SECURE=true`. For separate frontend/backend origins, configure the frontend build URLs and ensure cookies/CORS are compatible with the chosen domains.

## Verification status

Frontend `npm test`, `npm run build`, editor diagnostics, and `npm audit --omit=dev` have been run successfully. Backend `go test ./...`, `go vet ./...`, `go build ./...`, and Docker Compose smoke tests were not verifiable in the current environment because Go is not installed and Docker Desktop's Linux engine is not running. The repository must not be called production-ready until those checks and a real two-client WebSocket test pass.

## Known limitations and next phase

The rate limiter is fixed-window and IP-based, so it is intentionally simple. Cookie-based anonymous voting is not resistant to determined abuse. Redis Pub/Sub does not replay missed events; clients recover by fetching current results. The next engineering phase should add service interfaces and integration tests using disposable MongoDB/Redis containers, then verify deployment on a WebSocket-capable host.

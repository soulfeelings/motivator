# Motivator Monorepo

## Backend Rules (Go / Fiber)

### Logging & Tracing
- Every new HTTP handler or endpoint MUST include request tracing via the `requestid` middleware.
- When adding new routes, always use the `locals:requestid` trace ID in any log output within the handler.
- Use structured log lines that include: timestamp, status, latency, method, path, IP, trace ID, and error.
- Example log format: `${time} | ${status} | ${latency} | ${method} ${path} | ${ip} | trace=${locals:requestid} | ${error}`
- When logging inside handlers (e.g. for debugging or business logic), retrieve the trace ID with `requestid.FromContext(c)` and include it as `trace=<id>` in all log messages.

### Swagger
- All API handlers must have swaggo annotations (`@Summary`, `@Description`, `@Tags`, `@Produce`, `@Success`, `@Failure`, `@Router`).
- Run `make backend-swagger` or `swag init -g cmd/server/main.go -o docs` after adding/changing endpoints.

## Supabase

### Connection
- Project ref: `evfkxiphjhriwaozppsf`
- Region: `eu-west-3` (Paris)
- Database uses **Session Pooler** (not direct connection) because the direct connection is IPv6 only and not compatible with IPv4 networks.
- Pooler URL: `aws-1-eu-west-3.pooler.supabase.com:5432`
- Auth: Supabase Auth handles signup/login — backend only validates JWTs.

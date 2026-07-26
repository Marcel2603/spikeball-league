# Health Endpoints

The service exposes two dedicated health endpoints for use with container orchestrators and load balancers.

## Liveness — `GET /health/live`

Signals that the **process is running**. It has no external dependencies and always returns `200 OK`.

```shell
GET /health/live
```

```json
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

## Readiness — `GET /health/ready`

Signals that the service is **ready to serve traffic**.

```shell
GET /health/ready
```

**All fine reachable:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

## Log Behaviour

Both `/health/live` and `/health/ready` are **excluded from application logs** to prevent polling noise
in production environments.

## Kubernetes Example

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 3000
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 3000
  initialDelaySeconds: 3
  periodSeconds: 5
```

## Docker Compose Example

```yaml
services:
  spikeball-league:
    image: ghcr.io/marcel2603/spikeball-league/spikeball-league:latest
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/health/live"]
      interval: 30s
      timeout: 5s
      retries: 3
```

> **Tip:** Use `health/live` for the Docker Compose `healthcheck` (lightweight).
> Use `health/ready` as the Kubernetes readiness probe so traffic is not routed until all external dependencies are reachable.

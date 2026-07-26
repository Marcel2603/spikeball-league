# What to configure

All options are nested under their respective top-level YAML keys. Below is a full annotated reference.

## Server

Controls the HTTP listener.

```yaml
server:
  host: localhost   # Hostname used for CORS allowed origins
  port: 3000        # Port the service listens on
```

| Key            | Type   | Default     | Description                        |
|----------------|--------|-------------|------------------------------------|
| `server.host`  | string | `localhost` | Allowed CORS origin hostname       |
| `server.port`  | string | `3000`      | HTTP port to listen on             |

---

## Log

Controls the log level.

```yaml
log:
  level: info
```

| Key         | Type   | Default | Values                            |
|-------------|--------|---------|-----------------------------------|
| `log.level` | string | `info`  | `debug`, `info`, `warn`, `error`  |

All logs are emitted as **structured JSON** to stdout with source file and request ID (`req_id`) fields attached.

---

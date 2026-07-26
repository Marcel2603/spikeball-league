# What to configure

All options are nested under their respective top-level YAML keys. Below is a full annotated reference.

## Server

Controls the HTTP listener.

```yaml
server:
  host: localhost   # Hostname used for CORS allowed origins
  port: 3000        # Port the service listens on
  domain: http://localhost:3000 # Domain to use for generated public links
```

| Key            | Type   | Default                 | Description                        |
|----------------|--------|-------------------------|------------------------------------|
| `server.host`  | string | `localhost`             | Allowed CORS origin hostname       |
| `server.port`  | string | `8080`                  | HTTP port to listen on             |
| `server.domain`| string | `http://localhost:8080` | Public domain prefix for generated URLs |

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

## Database

Controls the local SQLite database connection.

```yaml
database:
  path: spikeball.db
```

| Key              | Type   | Default        | Description                       |
|------------------|--------|----------------|-----------------------------------|
| `database.path`  | string | `spikeball.db` | File path to the SQLite database. |

---

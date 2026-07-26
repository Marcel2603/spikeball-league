# Metrics

The service exposes a Prometheus-compatible `/metrics` endpoint to help you monitor the health,
performance, and usage of the application.

This endpoint combines standard HTTP/application metrics with custom metrics specifically designed to track operation performance.

## Endpoint Details

* **Path:** `/metrics`
* **Method:** `GET`
* **Format:** Prometheus Text-based format (`text/plain`)

---

## Custom Metrics

To help monitor the performance and reliability of your server, the service exposes the following custom metric:

* None

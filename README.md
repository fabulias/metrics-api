# Metrics API Documentation

This document outlines the API endpoints and handlers for managing metrics related to devices.

## Handlers and Endpoints

### Register Device

Registers a new device.

- **Handler:** `MetricsHandler.RegisterDevice`
- **Endpoint:** `POST /devices/register`

#### Request
```json
{
    "name": "device_name"
}
```

#### Response
```json
{
    "id": "device_id",
    "name": "device_name"
}
```

### Send Metrics

Sends data of a device related to metrics.

- **Handler:** `MetricsHandler.SendMetrics`
- **Endpoint:** `POST /devices/:id/metrics`

#### Request
```json
{
    "metric": "metric_data"
}
```

#### Response
```json
{
    "message": "Metrics created successfully"
}
```

### Get Latest Metrics

Requests the latest metrics for a device.

- **Handler:** `MetricsHandler.GetLatestMetrics`
- **Endpoint:** `GET /devices/:id/metrics`

#### Response
```json
{
    "metrics": "latest_metrics_data"
}
```

### Get Metrics History

Requests the history of metrics for a device.

- **Handler:** `MetricsHandler.GetMetricsHistory`
- **Endpoint:** `GET /devices/:id/metrics/history`

#### Response
```json
{
    "history": "metrics_history_data"
}
```

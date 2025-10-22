# ZTE OLT Management API Documentation

![Version](https://img.shields.io/badge/Version-1.0.0-blue.svg)
![Framework](https://img.shields.io/badge/Framework-Fiber%20v2-green.svg)

Complete REST API documentation for ZTE OLT Management System.

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Endpoints](#endpoints)
  - [System](#system)
  - [Templates](#templates)
  - [ONU Management](#onu-management)
  - [Batch Operations](#batch-operations)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

## Base URL

```
Development: http://localhost:8080
Production:  https://your-domain.com
```

## Authentication

Currently, the API does not require authentication. However, this may change in future versions.

## Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "1234567890"
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Indicates if the request was successful |
| `data` | object/array | The response data (varies by endpoint) |
| `error` | string | Error message if `success` is false |
| `timestamp` | string | ISO 8601 timestamp of the response |
| `request_id` | string | Unique identifier for the request |

## Endpoints

### System

#### GET `/`
Get API information and available endpoints.

**Response:**
```json
{
  "success": true,
  "data": {
    "name": "ZTE OLT Management API",
    "version": "1.0.0",
    "status": "running",
    "framework": "Fiber v2",
    "endpoints": {
      "health": "/api/v1/health",
      "templates": "/api/v1/templates",
      "add_onu": "/api/v1/onu/add",
      "delete_onu": "/api/v1/onu/delete",
      "check_attenuation": "/api/v1/onu/check-attenuation",
      "check_unconfigured": "/api/v1/onu/check-unconfigured",
      "batch_commands": "/api/v1/batch/commands"
    }
  }
}
```

#### GET `/api/v1/health`
Health check endpoint to verify service status.

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "uptime": "2h 15m 30s",
    "services": {
      "olt": "ok",
      "templates": "ok",
      "fiber": "ok"
    },
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Templates

#### GET `/api/v1/templates`
List all available command templates.

**Response:**
```json
{
  "success": true,
  "data": {
    "templates": [
      "add-onu",
      "delete-onu",
      "check-attenuation",
      "check-unconfigured"
    ],
    "count": 4
  }
}
```

### ONU Management

#### POST `/api/v1/onu/add`
Add/register a new ONU to the OLT.

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "slot": 2,
  "olt_port": 4,
  "onu": 17,
  "serial_number": "HWTC8A24189E",
  "code": "220219123239",
  "render_only": false
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | OLT IP address |
| `port` | integer | Yes | SSH/Telnet port (usually 23) |
| `user` | string | Yes | OLT username |
| `password` | string | Yes | OLT password |
| `slot` | integer | Yes | OLT slot number |
| `olt_port` | integer | Yes | OLT port number |
| `onu` | integer | Yes | ONU ID |
| `serial_number` | string | Yes | ONU serial number |
| `code` | string | Yes | Authorization/registration code |
| `render_only` | boolean | No | If true, only render commands without execution |

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "add-onu",
    "commands": [
      "configure terminal",
      "interface gpon-olt_2/4",
      "onu 17 type HWTC sn HWTC8A24189E",
      "exit",
      "interface gpon-onu_2/4:17",
      "tcont 1 profile 1G",
      "gemport 1 tcont 1",
      "gemport 2 tcont 1",
      "service-port 1 vport 1 user-vlan 10 vlan 10",
      "exit",
      "exit"
    ],
    "output": "ONU added successfully",
    "success": true,
    "time": "2.5s",
    "render_only": false
  }
}
```

#### POST `/api/v1/onu/delete`
Delete/remove an ONU from the OLT.

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "slot": 2,
  "olt_port": 4,
  "onu": 17,
  "render_only": false
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | OLT IP address |
| `port` | integer | Yes | SSH/Telnet port |
| `user` | string | Yes | OLT username |
| `password` | string | Yes | OLT password |
| `slot` | integer | Yes | OLT slot number |
| `olt_port` | integer | Yes | OLT port number |
| `onu` | integer | Yes | ONU ID |
| `render_only` | boolean | No | If true, only render commands without execution |

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "delete-onu",
    "commands": [
      "configure terminal",
      "interface gpon-onu_2/4:17",
      "shutdown",
      "exit",
      "interface gpon-olt_2/4",
      "no onu 17",
      "exit"
    ],
    "output": "ONU deleted successfully",
    "success": true,
    "time": "1.8s",
    "render_only": false
  }
}
```

#### POST `/api/v1/onu/check-attenuation`
Check optical power attenuation for a specific ONU.

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "slot": 2,
  "olt_port": 4,
  "onu": 17,
  "render_only": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "check-attenuation",
    "success": true,
    "time": "3.2s",
    "render_only": false,
    "data": {
      "host": "136.1.1.100",
      "slot": 2,
      "port": 4,
      "onu": 17,
      "direction": "up",
      "olt_rx_power_dbm": -18.5,
      "olt_tx_power_dbm": 2.1,
      "onu_rx_power_dbm": -20.8,
      "onu_tx_power_dbm": 1.5,
      "attenuation_db": 15.6,
      "status": "good",
      "raw_output": "OLT Rx Power: -18.5 dBm..."
    }
  }
}
```

**Attenuation Status Values:**
- `good`: < 20 dB attenuation
- `warning`: 20-25 dB attenuation
- `critical`: > 25 dB attenuation

#### POST `/api/v1/onu/check-unconfigured`
Find all unconfigured ONUs in the OLT.

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "render_only": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "check-unconfigured",
    "success": true,
    "time": "5.1s",
    "render_only": false,
    "data": {
      "host": "136.1.1.100",
      "total_count": 3,
      "onus": [
        {
          "olt_index": "gpon-olt_2/1",
          "model": "HWTC",
          "serial_number": "HWTC8A2418F0",
          "slot": 2,
          "port": 1
        },
        {
          "olt_index": "gpon-olt_2/3",
          "model": "HWTC",
          "serial_number": "HWTC8A2418F1",
          "slot": 2,
          "port": 3
        }
      ],
      "grouped_by_slot": {
        "2": [
          {
            "olt_index": "gpon-olt_2/1",
            "model": "HWTC",
            "serial_number": "HWTC8A2418F0",
            "slot": 2,
            "port": 1
          }
        ]
      },
      "status": "found",
      "raw_output": "Unconfigured ONUs found..."
    }
  }
}
```

### Batch Operations

#### POST `/api/v1/batch/commands`
Execute custom commands on the OLT.

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "commands": [
    "show version",
    "show gpon onu state",
    "display optical-module-info gpon-olt_2/4"
  ]
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | OLT IP address |
| `port` | integer | Yes | SSH/Telnet port |
| `user` | string | Yes | OLT username |
| `password` | string | Yes | OLT password |
| `commands` | array | Yes | List of commands to execute |

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "batch",
    "commands": [
      "show version",
      "show gpon onu state"
    ],
    "output": "ZTE OLT Version...\nONU State Information...",
    "success": true,
    "time": "4.7s"
  }
}
```

## Error Handling

### HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| `200` | Success |
| `400` | Bad Request - Invalid input data |
| `500` | Internal Server Error - OLT connection or execution error |

### Error Response Format

```json
{
  "success": false,
  "data": null,
  "error": "Connection timeout to OLT device",
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "1234567890"
}
```

### Common Error Messages

- `"Invalid request body"` - JSON parsing error
- `"Template rendering failed"` - Command template error
- `"Connection timeout"` - OLT unreachable
- `"Authentication failed"` - Invalid credentials
- `"Command execution failed"` - OLT command error

## Rate Limiting

Currently, there are no rate limits implemented. However, consider implementing rate limiting for production use to prevent OLT device overload.

## Examples

### Complete ONU Provisioning Workflow

```bash
# 1. Check for unconfigured ONUs
curl -X POST http://localhost:8080/api/v1/onu/check-unconfigured \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 23,
    "user": "aba",
    "password": "zte"
  }'

# 2. Add new ONU
curl -X POST http://localhost:8080/api/v1/onu/add \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 23,
    "user": "aba",
    "password": "zte",
    "slot": 2,
    "olt_port": 4,
    "onu": 17,
    "serial_number": "HWTC8A24189E",
    "code": "220219123239"
  }'

# 3. Check attenuation after installation
curl -X POST http://localhost:8080/api/v1/onu/check-attenuation \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 23,
    "user": "aba",
    "password": "zte",
    "slot": 2,
    "olt_port": 4,
    "onu": 17
  }'
```

### Command Preview (Render Only)

```bash
# Preview commands without execution
curl -X POST http://localhost:8080/api/v1/onu/add \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "slot": 2,
    "olt_port": 4,
    "onu": 17,
    "serial_number": "HWTC8A24189E",
    "code": "220219123239",
    "render_only": true
  }'
```

### Using with HTTPie

```bash
# Add ONU with HTTPie
http POST localhost:8080/api/v1/onu/add \
  host=136.1.1.100 \
  port:=23 \
  user=aba \
  password=zte \
  slot:=2 \
  olt_port:=4 \
  onu:=17 \
  serial_number=HWTC8A24189E \
  code=220219123239
```

### Using with PowerShell

```powershell
# Check attenuation with PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/onu/check-attenuation" `
  -Method Post `
  -ContentType "application/json" `
  -Body @{
    host = "136.1.1.100"
    port = 23
    user = "aba"
    password = "zte"
    slot = 2
    olt_port = 4
    onu = 17
  } | ConvertTo-Json -Depth 10
```

## SDK Examples

### Go (using net/http)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    // Add ONU request
    reqBody := map[string]interface{}{
        "host": "136.1.1.100",
        "port": 23,
        "user": "aba",
        "password": "zte",
        "slot": 2,
        "olt_port": 4,
        "onu": 17,
        "serial_number": "HWTC8A24189E",
        "code": "220219123239",
    }

    jsonData, _ := json.Marshal(reqBody)

    resp, err := http.Post(
        "http://localhost:8080/api/v1/onu/add",
        "application/json",
        bytes.NewBuffer(jsonData),
    )

    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    fmt.Printf("Response: %+v\n", result)
}
```

### JavaScript (using fetch)

```javascript
async function addONU() {
    const response = await fetch('http://localhost:8080/api/v1/onu/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            host: '136.1.1.100',
            port: 23,
            user: 'aba',
            password: 'zte',
            slot: 2,
            olt_port: 4,
            onu: 17,
            serial_number: 'HWTC8A24189E',
            code: '220219123239'
        })
    });

    const result = await response.json();
    console.log('Response:', result);
}

addONU();
```

## Best Practices

1. **Use `render_only: true`** for testing and debugging
2. **Implement proper error handling** in your applications
3. **Use appropriate timeouts** for OLT operations
4. **Monitor OLT device health** to avoid overloading
5. **Log all operations** for audit and troubleshooting
6. **Validate input data** before sending requests
7. **Handle network timeouts** gracefully with retries

## Support

For API issues and questions:
- Create GitHub issue
- Check application logs
- Review network connectivity to OLT devices
- Verify OLT credentials and accessibility
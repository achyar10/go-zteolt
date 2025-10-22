# ZTE OLT Management API

REST API untuk management ZTE OLT devices melalui web interface.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently no authentication required (for development).

## Common Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "timestamp": "2025-01-22T10:30:00Z",
  "request_id": "1642848600123456789"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Invalid request parameters",
  "timestamp": "2025-01-22T10:30:00Z",
  "request_id": "1642848600123456789"
}
```

## Endpoints

### 1. Health Check
Check API health status.

**GET** `/api/v1/health`

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "uptime": "0h 5m 30s",
    "services": {
      "olt": "ok",
      "templates": "ok"
    }
  }
}
```

### 2. List Templates
Get available command templates.

**GET** `/api/v1/templates`

**Response:**
```json
{
  "success": true,
  "data": {
    "templates": ["add-onu", "delete-onu", "check-attenuation"],
    "count": 3
  }
}
```

### 3. Add ONU
Add new ONU to ZTE OLT.

**POST** `/api/v1/onu/add`

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

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "mode": "add-onu",
    "commands": ["con t", "interface gpon-olt_1/2/4", ...],
    "output": "=== 136.1.1.100:23 ===\n>>> con t\n...",
    "success": true,
    "time": "8.5s",
    "render_only": false
  }
}
```

### 4. Delete ONU
Delete ONU from ZTE OLT.

**POST** `/api/v1/onu/delete`

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

### 5. Check Attenuation
Check optical power attenuation.

**POST** `/api/v1/onu/check-attenuation`

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
    "output": "show pon power attenuation gpon-onu_1/2/4:17\n...",
    "success": true,
    "time": "3.2s"
  }
}
```

### 6. Batch Commands
Execute custom commands.

**POST** `/api/v1/batch/commands`

**Request Body:**
```json
{
  "host": "136.1.1.100",
  "port": 23,
  "user": "aba",
  "password": "zte",
  "commands": [
    "show version",
    "show interface gpon-olt_1/2/4",
    "display ont info 1/2/4 all"
  ]
}
```

## Usage Examples

### Using cURL
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Add ONU
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

# Check attenuation
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

### Using JavaScript
```javascript
// Add ONU
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
console.log(result);
```

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200         | Success |
| 400         | Bad Request |
| 404         | Not Found |
| 405         | Method Not Allowed |
| 500         | Internal Server Error |

## Development

Start development server:
```bash
make dev
```

Build for production:
```bash
make prod-build
```

Run tests:
```bash
make test
```
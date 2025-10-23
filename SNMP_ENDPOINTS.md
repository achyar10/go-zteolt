# SNMP Monitoring Endpoints Documentation

Project ini sekarang mendukung monitoring OLT menggunakan SNMP protokol selain Telnet yang sudah ada.

## 🚀 Fitur SNMP Monitoring

### Endpoint Baru yang Ditambahkan:

1. **`POST /api/v1/board/{board_id}/pon/{pon_id}/snmp`**
   - Mendapatkan semua ONU di board dan PON tertentu
   - Method: POST (untuk parameter authentication)

2. **`POST /api/v1/board/{board_id}/pon/{pon_id}/onu/{onu_id}/snmp`**
   - Mendapatkan detail ONU spesifik
   - Method: POST (untuk parameter authentication)

3. **`POST /api/v1/board/{board_id}/pon/{pon_id}/empty-slots/snmp`**
   - Mendapatkan slot ONU yang kosong
   - Method: POST (untuk parameter authentication)

## 📋 Contoh Penggunaan

### 1. Get All ONUs in Board 1 PON 7

```bash
curl -X POST http://localhost:8080/api/v1/board/1/pon/7/snmp \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 161,
    "community": "public",
    "timeout": 30
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "board_id": 1,
    "pon_id": 7,
    "total_onus": 12,
    "onus": [
      {
        "board": 1,
        "pon": 7,
        "id": 1,
        "name": "220219123239",
        "serial_number": "HWTC8A24189E",
        "rx_power": "-12.5",
        "tx_power": "2.1",
        "status": "online",
        "ip_address": "192.168.1.100",
        "last_online": "2024-01-23 10:30:15",
        "uptime": "2d 5h 15m"
      }
    ],
    "execution_time": "2.3s",
    "timestamp": "2024-01-23T15:45:30Z"
  },
  "error": null,
  "timestamp": "2024-01-23T15:45:30Z",
  "request_id": "1705995330123456789"
}
```

### 2. Get Specific ONU Details

```bash
curl -X POST http://localhost:8080/api/v1/board/1/pon/7/onu/1/snmp \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 161,
    "community": "public",
    "timeout": 30
  }'
```

### 3. Get Empty ONU Slots

```bash
curl -X POST http://localhost:8080/api/v1/board/1/pon/7/empty-slots/snmp \
  -H "Content-Type: application/json" \
  -d '{
    "host": "136.1.1.100",
    "port": 161,
    "community": "public",
    "timeout": 30
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "host": "136.1.1.100",
    "board_id": 1,
    "pon_id": 7,
    "total_empty": 116,
    "empty_slots": [
      {
        "board": 1,
        "pon": 7,
        "onu_id": 13
      },
      {
        "board": 1,
        "pon": 7,
        "onu_id": 14
      }
    ],
    "execution_time": "1.8s",
    "timestamp": "2024-01-23T15:50:00Z"
  }
}
```

## 🔧 Parameter Authentication

Tidak perlu setup environment variables. Semua parameter SNMP dikirim via request body:

- **host**: IP address OLT device
- **port**: SNMP port (default: 161)
- **community**: SNMP community string (default: "public")
- **timeout**: Timeout dalam detik (optional, default: 30)
- **board_id**: Board ID (1 atau 2)
- **pon_id**: PON ID (1-16)
- **onu_id**: ONU ID (1-128, untuk detail endpoint)

## 🎯 Keunggulan SNMP vs Telnet

### SNMP (Monitoring):
- ✅ **Lebih cepat** - UDP-based protocol
- ✅ **Less overhead** - Tidak perlu full connection
- ✅ **Real-time data** - Power levels, status, serial numbers
- ✅ **Standardized** - Consistent data format
- ✅ **Scalable** - Banyak concurrent requests

### Telnet (Configuration):
- ✅ **Complete control** - Access ke semua CLI commands
- ✅ **Configuration** - Add/delete ONU, templates
- ✅ **Backward compatibility** - Support older devices

## 🏗️ Hybrid Approach

Project ini sekarang support hybrid approach:
- **SNMP**: Untuk monitoring real-time data
- **Telnet**: Untuk konfigurasi dan management operations

## 🐛 Debug Endpoints (Untuk troubleshooting)

### 1. Basic SNMP Connection Test

```bash
curl -X POST http://localhost:8080/api/v1/debug/snmp \
  -H "Content-Type: application/json" \
  -d '{
    "host": "103.249.18.134",
    "port": 161,
    "community": "public",
    "timeout": 30
  }'
```

Response akan menunjukkan:
- Koneksi SNMP berhasil atau tidak
- Hasil test untuk berbagai OID dasar
- Walk results untuk mencari ONU-related OIDs

### 2. ONU OID Discovery

```bash
curl -X POST http://localhost:8080/api/v1/debug/snmp/board/1/pon/6/discover \
  -H "Content-Type: application/json" \
  -d '{
    "host": "103.249.18.134",
    "port": 161,
    "community": "public",
    "timeout": 30
  }'
```

Response akan menunjukkan OID patterns mana yang menghasilkan data ONU.

## 📝 Request Parameters

### SNMPMonitoringRequest
```json
{
  "host": "string",      // required - OLT IP address
  "port": 161,          // required - SNMP port
  "community": "string", // required - SNMP community
  "timeout": 30         // optional - Timeout in seconds
}
```

**Note:** `board_id` dan `pon_id` diambil dari URL parameter, bukan dari request body.

### Response Format
Semua response menggunakan format standar yang konsisten dengan existing API.

## 🚨 Error Handling

- **400 Bad Request**: Invalid parameters
- **404 Not Found**: ONU tidak ditemukan
- **500 Internal Server Error**: SNMP connection/timeout error

## 🗑️ Cleanup SNMP Reference Files

Setelah implementasi selesai, folder `snmp-olt-zte-main/` bisa dihapus karena sudah tidak diperlukan lagi.
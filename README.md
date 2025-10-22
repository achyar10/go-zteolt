# ZTE OLT Management API

![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)

RESTful API untuk management ZTE OLT devices dengan Fiber framework dan Docker support.

## 🚀 Features

- **🌐 Web API**: RESTful API dengan Fiber v2 (high performance)
- **📡 ONU Management**: Add, delete, dan monitoring ONU devices
- **📊 Attenuation Check**: Cek redaman daya optik dengan parsing otomatis
- **🔧 Batch Commands**: Execute custom commands pada OLT
- **📝 Template System**: Flexible command templates dengan Go templates
- **⚡ Real-time Execution**: Live feedback dengan timeout handling
- **🏗️ Clean Architecture**: Proper separation of concerns
- **🐳 Docker Ready**: Multi-stage Docker build dengan security best practices
- **📈 Health Monitoring**: Built-in health checks dan logging
- **🔐 Security**: Non-root user, CORS support, request validation

## 📁 Project Structure

```
go-zteolt/
├── cmd/server/          # Web server entry point
├── internal/            # Private application code
│   ├── api/            # HTTP handlers & routes
│   ├── config/         # Configuration & templates
│   ├── olt/            # OLT business logic
│   └── utils/          # Utility functions
├── templates/          # Command templates
├── docs/              # API documentation
├── bin/               # Built binaries
└── main.go            # Legacy CLI (preserved)
```

## 🛠️ Quick Start

### Prerequisites
- **Go 1.22+** - Latest Go version
- **Docker** - Optional, for containerization
- **Make** - Optional, for build commands

### Installation

1. **Clone repository**
```bash
git clone https://github.com/achyar10/go-zteolt.git
cd go-zteolt
```

2. **Install dependencies**
```bash
make deps
# atau
go mod download && go mod tidy
```

3. **Start development server**
```bash
# Dengan hot reload (recommended untuk development)
make dev

# Tanpa hot reload
make dev-simple

# Atau langsung dengan go run
go run cmd/server/main.go -dev
```

Server akan start pada `http://localhost:8080`

### Quick Test

```bash
# Health check
curl http://localhost:8080/api/v1/health

# API Info
curl http://localhost:8080/

# List available templates
curl http://localhost:8080/api/v1/templates
```

## 🌐 API Endpoints

### Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API information & endpoints list |
| GET | `/api/v1/health` | Health check & service status |
| GET | `/api/v1/templates` | List available command templates |
| POST | `/api/v1/onu/add` | Add/register new ONU |
| POST | `/api/v1/onu/delete` | Delete/remove ONU |
| POST | `/api/v1/onu/check-attenuation` | Check optical power attenuation |
| POST | `/api/v1/onu/check-unconfigured` | Find unconfigured ONUs |
| POST | `/api/v1/batch/commands` | Execute custom commands |

### Example Usage

#### Add ONU
```bash
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
```

#### Check Attenuation
```bash
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

#### Render Commands Only (No Execution)
```bash
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

### Response Format

Semua response menggunakan format standar:
```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "1234567890"
}
```

📖 **Lihat dokumentasi lengkap API di [docs/api.md](docs/api.md)**

## 🏗️ Build & Deployment

### Development
```bash
make dev              # Start dev server
make dev-build        # Build development binary
make test             # Run tests
make lint             # Run linter
```

### Production
```bash
make prod-build       # Build for production (Linux)
make build            # Build all binaries
make install          # Install to system
```

### Docker Deployment

#### Build & Run
```bash
# Build Docker image
docker build -t go-zteolt:latest .

# Run container
docker run -p 8080:8080 go-zteolt:latest

# Atau gun docker-compose (recommended)
docker-compose up -d
```

#### Docker Compose
```bash
# Start dengan docker-compose
docker-compose up --build

# Background mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Production Docker Build
```bash
# Build dengan version info
docker build \
  --build-arg VERSION=$(git describe --tags) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t go-zteolt:latest .
```

## 📚 Documentation

- **📖 API Documentation**: [docs/api.md](docs/api.md) - Complete API reference
- **🔧 Development Guide**: [docs/development.md](docs/development.md) - Contributing guidelines
- **📋 Legacy CLI**: Original CLI tool preserved as `main.go`

## ⚙️ Configuration

Default configuration can be modified in `internal/config/config.go`:

```go
Server:
  Host: "0.0.0.0"
  Port: 8080
  ReadTimeout: 30s
  WriteTimeout: 30s

OLT:
  DefaultTimeout: 8s
  WriteTimeout: 24s
  MaxRetries: 2
```

## 🔧 Development

### Adding New Templates

1. Create template file in `templates/`
2. Add to `templates.go` loader
3. Create corresponding API endpoint
4. Update documentation

### Code Structure

- **Handlers**: HTTP request/response handling
- **Services**: Business logic
- **Models**: Data structures
- **Templates**: Command templates

## 🐛 Troubleshooting

### Common Issues

1. **Port already in use**
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use different port
make dev port=8081
```

2. **Build errors**
```bash
# Clean and rebuild
make clean && make build
```

3. **Template not found**
```bash
# Check templates directory
ls -la templates/
```

## 📄 License

[Your License]

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make your changes
4. Add tests
5. Submit pull request

## 📞 Support

For issues and questions:
- Create GitHub issue
- Check documentation
- Review existing issues
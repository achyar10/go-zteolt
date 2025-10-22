# ZTE OLT Management Tool

Web-based REST API untuk management ZTE OLT devices.

## 🚀 Features

- **Web API**: RESTful API dengan Go native (no heavy framework)
- **ONU Management**: Add, delete, dan monitoring ONU
- **Attenuation Check**: Cek redaman daya optik
- **Batch Commands**: Execute custom commands
- **Template System**: Flexible command templates
- **Real-time Execution**: Live feedback dengan timeout handling
- **Clean Architecture**: Proper separation of concerns

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
- Go 1.21 or higher
- Make (optional, for build commands)

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd go-zteolt
```

2. **Install dependencies**
```bash
make deps
```

3. **Start development server**
```bash
make dev
```

Server will start on `http://localhost:8080`

### Quick Test

```bash
# Health check
curl http://localhost:8080/api/v1/health

# List available templates
curl http://localhost:8080/api/v1/templates
```

## 🌐 API Endpoints

### Core Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/templates` | List templates |
| POST | `/api/v1/onu/add` | Add ONU |
| POST | `/api/v1/onu/delete` | Delete ONU |
| POST | `/api/v1/onu/check-attenuation` | Check attenuation |
| POST | `/api/v1/batch/commands` | Execute custom commands |

### Example: Add ONU

```bash
curl -X POST http://localhost:8080/api/v1/onu/add \
  -H "Content-Type: application/json" \
  -d '{
    "host": "103.249.18.134",
    "port": 2727,
    "user": "aba",
    "password": "@aba1010#",
    "slot": 2,
    "olt_port": 4,
    "onu": 17,
    "serial_number": "HWTC8A24189E",
    "code": "220219123239"
  }'
```

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

### Docker
```bash
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

## 📚 Documentation

- **API Documentation**: [docs/api.md](docs/api.md)
- **Legacy CLI**: Original CLI tool preserved as `main.go`

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
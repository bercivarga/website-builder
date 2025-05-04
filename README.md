# Website Builder Full Stack Project

A modern full-stack web application built with React (frontend), Go (backend), and PostgreSQL (database), containerized with Docker for seamless development and deployment.

## 🏗️ Architecture

```
website-builder-fullstack-project/
├── frontend/               # React application with Vite
├── backend/               # Go API server
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── migrations/        # Database migrations using Goose
├── docker-compose.yml     # Multi-container Docker setup
└── Makefile              # Development commands
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Make

### Getting Started

```bash
# Start all services
make up

# Run database migrations
make migrate

# Stop all services
make down
```

Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

## 🛠️ Development

### Available Commands

```bash
# Start all services
make up

# Stop all services
make down

# Show logs for frontend and backend
make logs

# Show all logs
make logs-all

# Run database migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create

# Check migration status
make migrate-status

# Run tests
make test

# Clean up everything
make clean
```

### Project Structure

#### Frontend
- Built with React and Vite
- TypeScript/JavaScript support
- Hot module reloading for development
- Environment variables: `REACT_APP_API_URL`

#### Backend
- RESTful API built with Go
- Hot reloading with Air
- Database migrations with Goose
- PostgreSQL integration
- Environment-based configuration

#### Database
- PostgreSQL 16 Alpine
- Persistent data storage with Docker volumes
- Automatic migration management

## 🗄️ Database

### Migrations

The project uses Goose for database migrations. Migrations are located in `backend/migrations/`.

```bash
# Create a new migration
make migrate-create

# Apply migrations
make migrate-up

# Rollback
make migrate-down
```

### Schema Example

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🐳 Docker Configuration

### Services

1. **Frontend**
   - Port: 3000
   - Build context: `./frontend`
   - Volumes: source code mount, node_modules cache

2. **Backend**
   - Port: 8080
   - Build context: `./backend`
   - Volumes: source code mount for hot reloading

3. **Database**
   - Port: 5432
   - Image: postgres:16-alpine
   - Volume: persistent data storage

## 🔧 Environment Variables

### Frontend
- `REACT_APP_API_URL`: Backend API URL (default: http://localhost:8080)

### Backend
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `GOOSE_DRIVER`: Database driver for migrations
- `GOOSE_DBSTRING`: Database connection string

## 📝 Development Workflow

1. Start the development environment:
   ```bash
   make up
   ```

2. Watch logs for both frontend and backend:
   ```bash
   make logs
   ```

3. Create new database migrations as needed:
   ```bash
   make migrate-create
   ```

4. Run tests before committing:
   ```bash
   make test
   ```

## 🚢 Production Deployment

For production deployment, create separate Dockerfiles optimized for production builds:

1. Multi-stage builds for smaller images
2. Production environment variables
3. Health checks
4. Security best practices

## 📚 API Documentation

API endpoints will be documented here as they are developed.

Example:
```
GET /api/health - Health check endpoint
```

## 🤝 Contributing

1. Create a new branch for your feature
2. Make your changes
3. Run tests: `make test`
4. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues or questions:
1. Check the logs: `make logs`
2. Review common issues in the troubleshooting section
3. Create an issue in the repository

## 🔍 Troubleshooting

### Common Issues

1. **Frontend not accessible**
   - Check if container is running: `docker-compose ps`
   - Verify Vite is binding to `0.0.0.0`
   - Check logs: `make logs frontend`

2. **Backend database connection issues**
   - Ensure database is running: `docker-compose ps`
   - Verify database credentials in environment variables
   - Run migration status: `make migrate-status`

3. **Hot reloading not working**
   - Add `CHOKIDAR_USEPOLLING=true` to frontend environment
   - Ensure volumes are properly mounted
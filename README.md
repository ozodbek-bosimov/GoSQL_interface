# GoSQL Interface

A web-based admin panel built with Go and PostgreSQL for learning database connectivity and CRUD operations.

## Requirements

- Go 1.20+
- PostgreSQL 14+

## Installation

1. Navigate to project directory:
```bash
cd /Users/ozodbek/projects/GoSQL_interface
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Create database:
```bash
createdb -U ozodbek gosql_db
```

4. Configure environment:
```bash
cp .env.example .env
# Edit .env file with your database credentials
```

5. Run the application:
```bash
go run cmd/server/main.go
```

6. Open in browser:
```
http://localhost:8080/login
```

## Default Credentials

- **Username**: admin
- **Password**: admin123

## Features

- ✅ User management (CRUD)
- ✅ Network interfaces management (CRUD)
- ✅ ACL rules management (CRUD)
- ✅ Search functionality
- ✅ Password security (bcrypt)
- ✅ Session authentication
- ✅ Database schema viewer

## Tech Stack

- **Backend**: Go 1.26 + Gin framework
- **Database**: PostgreSQL 14
- **Frontend**: HTML templates + Vanilla JS + CSS
- **Auth**: Sessions + bcrypt

## Project Structure

```
GoSQL_interface/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── config/                     # Configuration
│   ├── database/                   # DB connection & migrations
│   ├── models/                     # Data models (User, Interface, ACL)
│   ├── handlers/                   # HTTP handlers
│   ├── middleware/                 # Auth middleware
│   └── utils/                      # Utilities (password hashing)
├── migrations/                     # SQL migration files
├── web/
│   ├── templates/                  # HTML templates
│   └── static/                     # CSS & JavaScript
└── .env                           # Configuration file
```

## License

MIT

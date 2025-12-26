
# Active Directory Go Client

This project provides a simple Go client for interacting with Active Directory services. It is designed to help manage users and groups via LDAP.

## Features
- Connect to Active Directory using LDAP
- Manage users and groups
- Configuration via file
- Docker Compose support for easy setup

## Prerequisites
- Go 1.18 or later
- Docker & Docker Compose (for containerized usage)

## Getting Started

### 1. Clone the repository
```sh
git clone git@github.com:GergesHany/Active-Directory.git
cd Active-Directory
```

### 2. Build and Run
#### Run with Go
```sh
go run main.go
```

#### Or use Docker Compose
```sh
docker-compose up --build
```

### 3. Configuration
Edit the configuration files in `pkg/ad/config.go` as needed. Store your domain password in `secrets/domain_password.txt` (do not commit this file).

### 4. Usage
Implement your logic in `main.go` or extend the client in `pkg/ad/`.

## Project Structure
- `main.go` - Entry point
- `pkg/ad/` - Active Directory client code
- `secrets/` - Store sensitive files (e.g., domain password)
- `docker-compose.yml` - Docker Compose setup

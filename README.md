# CLI Auth System with 2FA

A Go CLI application that supports user registration, authentication, 2FA (TOTP), and session management.

## Setup & Running

### Using Docker

Start the PostgreSQL database and run the interactive CLI application:

```bash
docker compose run app
```

To stop and remove containers:

```bash
docker compose down -v
```

### Running Locally

Prerequisites: A running PostgreSQL instance.

1. Set the configuration environment variables:
   ```bash
   export DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable
   ```

2. Start the application (automatically runs migrations on startup):
   ```bash
   go run cmd/app/main.go
   ```

## CLI Commands

### Guest State
- `register` - Create a new account
- `login` - Log in with username/password (and 2FA if enabled)
- `help` - Show available commands
- `exit` - Quit application

### Authenticated State
- `whoami` - Show session details and expiration time
- `enable-2fa` - Generate and display a TOTP QR code to scan
- `disable-2fa` - Disable TOTP
- `logout` - End current session
- `help` - Show available commands

## Configuration

The app uses the following environment variables:
- `DATABASE_URL` (default: `postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable`)
- `SESSION_DURATION` (default: `5m`)
- `LOCKOUT_DURATION` (default: `1m`)

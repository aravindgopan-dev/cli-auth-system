# CLI Auth System with 2FA
<img width="1894" height="750" alt="image" src="https://github.com/user-attachments/assets/b941c5c7-75de-46f5-a1f6-66a883744613" />

<img width="1909" height="748" alt="image" src="https://github.com/user-attachments/assets/f12e679d-8e35-434f-9a79-27e8c758f085" />

<img width="1918" height="547" alt="image" src="https://github.com/user-attachments/assets/eb240201-ef7f-4135-919e-2137e0e4ae0c" />
<img width="1915" height="411" alt="image" src="https://github.com/user-attachments/assets/ef2d7f89-b948-40b4-88dd-90e7e4d929a7" />


A Go CLI application that supports user registration, authentication, 2FA (TOTP), and session management. Built using dependency injection (DI) .

## Setup & Running

### Using Docker

Start the PostgreSQL database and run the interactive CLI application:

```bash
make run
# or: sudo docker compose run app
```

To stop and remove containers:

```bash
make clean
# or: sudo docker compose down -v
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
- `DATABASE_URL` (default: `postgres://postgres:postgres@postgres:5432/auth_db?sslmode=disable`)
- `SESSION_DURATION` (default: `5m`)
- `LOCKOUT_DURATION` (default: `1m`)

## Key Libraries Used

- **[go-prompt](https://github.com/c-bata/go-prompt)**: Generates the interactive shell with dynamic command autocomplete.
- **[qrterminal](https://github.com/mdp/qrterminal)**: Renders compact QR codes directly into the terminal window.
- **[otp/totp](https://github.com/pquerna/otp)**: Generates and validates Google Authenticator compatible TOTP keys.
- **[pgx](https://github.com/jackc/pgx)**: High-performance PostgreSQL driver and toolkit for Go.
- **[bcrypt](https://golang.org/x/crypto/bcrypt)**: Secure password hashing algorithm.
- **[uuid](https://github.com/google/uuid)**: Generates cryptographically secure session IDs.
- **[golang-migrate](https://github.com/golang-migrate/migrate)**: Standard tool container used for running database schema migrations.

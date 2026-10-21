# E-WALLET API

[![Go Version](https://img.shields.io/badge/Go-1.27.1-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.12.0-008ECF?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgx_v5-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-go--redis_v9-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-Swaggo-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![Docker](https://img.shields.io/badge/Docker-Multi--stage-2496ED?style=for-the-badge&logo=docker&logoColor=white)

A RESTful digital wallet API built with Go and the Gin web framework. Developed for the **The Workers E-Wallet** collaboration project, this backend manages authentication, account profiles, wallets, transfers, transaction dashboard summaries, and top-up checkout requests via the DOKU payment gateway. It uses PostgreSQL for persistent data and Redis for token revocation and session-related operations. The API is designed to serve the E-Wallet frontend application.

## Technologies

- **Go 1.27.1** — Backend programming language declared in `go.mod`.
- **Gin 1.12.0** — HTTP router, middleware, request binding, and REST API handlers.
- **PostgreSQL with pgx/v5** — Relational storage and database connection pooling.
- **Redis (go-redis/v9)** — Token blacklist used during authentication and logout.
- **JWT (golang-jwt/jwt/v5)** — Bearer tokens used to protect private endpoints.
- **Swaggo & Swagger UI** — OpenAPI 2.0 (Swagger) documentation generated from Go controller annotations.
- **golang-migrate** — SQL migration workflow exposed through the project's Makefile; install the CLI separately.
- **DOKU API** — Sandbox checkout integration and payment status webhook handling.
- **SMTP / go-mail** — Email delivery for the forgot-password workflow (configured for Gmail SMTP in the source).
- **Docker** — Multi-stage build using Go for compilation and Debian for runtime.

## Features

- **Authentication** — Register, login, logout, request password-reset email, and reset password.
- **Profile management** — Retrieve and update user profiles, change password, upload or delete an avatar.
- **Transaction PIN** — Set and change a six-digit transaction PIN.
- **Dashboard** — Retrieve wallet balance summary (balance, income, expense) and recent transaction history.
- **Transfer** — Search recipients, list recent/favorite contacts, and transfer balance between wallets with PIN verification.
- **Top-up and payments** — List payment methods, initiate top-up via DOKU Checkout, and receive payment status notifications.
- **Swagger API documentation** — Interactive endpoint descriptions, request parameters, HTTP responses, and Bearer authentication.
- **Database migrations and seeds** — Versioned SQL migrations under `db/migration/` and sample-data scripts under `db/seed/`.

## Getting Started

### Prerequisites

Install and prepare the following before running the backend:

- **Git** to clone the repository.
- **Go 1.27.1** (or a Go toolchain capable of honoring the version declared in `go.mod`).
- **PostgreSQL** with a database and credentials for this application.
- **Redis** accessible by the backend.
- **SMTP credentials** if testing the password-reset email flow.
- **DOKU sandbox credentials** if testing online top-up checkout.
- **golang-migrate CLI** and **Make** if using the provided migration commands (a PostgreSQL `psql` CLI is needed for seeding).
- **Docker** (optional, for building/running the API container).

The application starts a PostgreSQL connection, initializes Redis, and creates an SMTP mailer during startup. Configure those services before testing the endpoints.

### Clone the repository

The application's source is maintained in the backend repository. The collaborative development branch is `dev`:

```bash
git clone -b dev https://github.com/rcontrisha/api-workers-ewallet.git
cd api-workers-ewallet
```

If you contribute via fork, create a feature branch from the team's `dev` branch instead of committing directly to the shared branch.

### Configure environment variables

Copy the provided `.env.example` and fill in your local credentials:

```powershell
Copy-Item .env.example .env
```

Example for **PowerShell** (use `cp .env.example .env` on a Unix-like shell). Do not commit `.env`.

The backend reads these values (replace placeholders with your own values):

```dotenv
HOST=0.0.0.0
PORT=8080
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
REDIS_USER=your_redis_user
REDIS_PASSWORD=your_redis_password
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_KEY=your_strong_jwt_signing_secret
JWT_ISSUER=the-workers-ewallet
SMTP_USER=your_smtp_user
SMTP_PASS=your_smtp_password
SMTP_SENDER_EMAIL=sender@example.com
DOKU_CLIENT_ID=your_doku_client_id
DOKU_SECRET_KEY=your_doku_secret_key
DOKU_CALLBACK_URL=https://your-public-callback.example/payment/doku/notify
```

**Important:** The uploaded `.env.example` omits `HOST` and `PORT`. It also uses `SMTP_SENDER`, whereas the code reads `SMTP_SENDER_EMAIL`. Use the keys in the example above until the team's environment template is corrected. `HOST` and `PORT` are read by `main.go`; `PORT=8080` matches the Swagger host annotation. The DOKU client currently uses the DOKU **sandbox** URL in code.

### Run database migrations

The SQL migrations are stored in `db/migration/`. The project's Makefile defines `migrate-up`, `migrate-down`, and `migrate-create` targets.

```bash
make migrate-up
```

Alternatively, if you do not have Make installed, use the `migrate` CLI directly (with your own PostgreSQL connection URL):

```bash
migrate -path db/migration -database "postgres://USER:PASSWORD@localhost:5432/DB_NAME?sslmode=disable" up
```

**Use migrations before seeding or starting the application against a new database.** `make seed` invokes PowerShell and `psql`, so it depends on those commands being available in the environment where Make runs.

### Run in development mode

Install the Go dependencies and start the backend from the repository root:

```bash
go mod download
go run .
```

When running locally with `PORT=8080`, API routes are served at `http://localhost:8080` and Swagger UI is available at:

**http://localhost:8080/swagger/index.html**

The API uses paths without an `/api` prefix, for example `POST /auth/login`, `GET /user/profile`, and `POST /transfer`. Most private endpoints require `Authorization: Bearer <token>`; login and registration do not.

### Generate Swagger documentation

The project already includes Swaggo libraries, the Swagger UI route, and generated files under `docs/`. To regenerate them after changing controller annotations, install the Swaggo CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.6
```

Run the CLI from the backend root:

```bash
swag init -g main.go --parseDependency --parseInternal
```

On Windows PowerShell, if `swag` is not in your PATH, use:

```powershell
& "$(go env GOPATH)\bin\swag.exe" init -g main.go --parseDependency --parseInternal
```

This refreshes `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`. Restart the Go server to reload the generated API specification. Use the Swagger **Authorize** control with the full value `Bearer <token>` for authenticated routes. Do not use real credentials for Swagger examples or commit secrets.

### Run with Docker

Build the image using the provided multi-stage Dockerfile:

```bash
docker build -t workers-ewallet-api .
```

Run the image using an environment file with service hostnames reachable **from the container**:

```bash
docker run --rm --env-file .env -p 8080:8080 workers-ewallet-api
```

The Docker image does **not** start PostgreSQL or Redis. If those services run outside the container, set `DB_HOST` and `REDIS_HOST` to reachable addresses rather than `localhost`. The Dockerfile creates a `public` folder for static assets but does not define a Compose stack. File uploads require persistent storage configuration for long-term use.

### Available commands

| Command                                                  | Description                                                  |
| -------------------------------------------------------- | ------------------------------------------------------------ |
| `go mod download`                                        | Download Go module dependencies                              |
| `go run .`                                               | Run the API locally                                          |
| `go build -o ewallet-api .`                              | Build the Go executable                                      |
| `go test ./...`                                          | Run Go tests (if provided)                                   |
| `swag init -g main.go --parseDependency --parseInternal` | Regenerate Swagger files                                     |
| `make migrate-up`                                        | Apply SQL migrations                                         |
| `make migrate-down`                                      | Roll back SQL migrations (check the affected versions first) |
| `make seed`                                              | Run seed scripts (PowerShell and `psql` required)            |
| `docker build -t workers-ewallet-api .`                  | Build the API container image                                |

## API Routes

The backend exposes the following API endpoints using the
Gin framework.

**Base URL (local):** `http://localhost:8080`

Protected routes require the following header:

`Authorization: Bearer <token>`

### Authentication

| Method | Endpoint                | Description                      | Auth |
| ------ | ----------------------- | -------------------------------- | ---- |
| POST   | `/auth/register`        | Register a new user              | No   |
| POST   | `/auth/login`           | Authenticate user and return JWT | No   |
| POST   | `/auth/logout`          | Log out and invalidate token     | Yes  |
| POST   | `/auth/forgot-password` | Request password reset           | No   |
| POST   | `/auth/reset-password`  | Reset account password           | No   |

### User & Profile

| Method | Endpoint                | Description                    | Auth |
| ------ | ----------------------- | ------------------------------ | ---- |
| GET    | `/user/profile`         | Get authenticated user profile | Yes  |
| PATCH  | `/user/update`          | Update user profile            | Yes  |
| DELETE | `/user/avatar`          | Delete user avatar             | Yes  |
| PATCH  | `/user/change-password` | Change account password        | Yes  |
| POST   | `/user/pin`             | Set transaction PIN            | Yes  |
| PATCH  | `/user/change-pin`      | Change transaction PIN         | Yes  |

### Dashboard

| Method | Endpoint                  | Description                        | Auth |
| ------ | ------------------------- | ---------------------------------- | ---- |
| GET    | `/user/dashboard-summary` | Get wallet and transaction summary | Yes  |
| GET    | `/user/dashboard-history` | Get recent dashboard transactions  | Yes  |

### Transfer

| Method | Endpoint           | Description                            | Auth |
| ------ | ------------------ | -------------------------------------- | ---- |
| GET    | `/transfer/search` | Search users and retrieve contacts     | Yes  |
| POST   | `/transfer`        | Transfer balance using transaction PIN | Yes  |

### Favorites

| Method | Endpoint    | Description                     | Auth |
| ------ | ----------- | ------------------------------- | ---- |
| POST   | `/favorite` | Toggle a user's favorite status | Yes  |

The `/favorite` endpoint adds or removes a favorite
depending on its current state.

### Top-Up & Payments

| Method | Endpoint                    | Description                   | Auth   |
| ------ | --------------------------- | ----------------------------- | ------ |
| GET    | `/payment/method`           | Get available payment methods | No     |
| POST   | `/payment/topup`            | Create a top-up transaction   | Yes    |
| GET    | `/payment/topup/{trx_code}` | Get top-up status and summary | Yes    |
| POST   | `/payment/doku/notify`      | Receive DOKU payment webhook  | No JWT |

### Transaction History

| Method | Endpoint        | Description                          | Auth |
| ------ | --------------- | ------------------------------------ | ---- |
| GET    | `/history`      | Get transaction history              | Yes  |
| DELETE | `/history/{id}` | Hide a transaction from user history | Yes  |

### API Documentation

Swagger UI is available at:

`http://localhost:8080/swagger/index.html`

**Note:** Routes above reflect the Gin router configuration.
A reverse proxy may add an `/api` prefix in deployment.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Related Repositories

- **Frontend Repository:** [fe-workers-ewallet](https://github.com/rcontrisha/fe-workers-ewallet)

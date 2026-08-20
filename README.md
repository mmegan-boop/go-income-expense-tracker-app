# Go Income & Expense Tracker API

A RESTful API for tracking personal income and expenses. Users can register, manage categories, log financial records, generate Excel reports, and view monthly summaries.

## Tech Stack

- **Language:** Go 1.26
- **Framework:** Echo v5
- **Database:** PostgreSQL (via GORM)
- **Authentication:** JWT (HS256) with role-based access control (admin/user)
- **Hot Reload:** Air
- **Excel Export:** Excelize

## Project Structure

```
cmd/
  server/            — Application entrypoint
  generate/          — CLI tool to generate an admin account
docs/
  *.json             — API documentation (Postman collection)
internal/
  api/               — Echo server initialization and dependency wiring
  router/            — Route registrations grouped by domain
  controller/        — HTTP request handlers
  service/           — Business logic layer
  repository/        — Database query layer
  model/             — GORM data models
  dto/               — Request and response structs
  middleware/        — JWT authentication and role-based middleware
  config/            — Configuration loaders
  constant/          — Environment variable key constants
  db/                — PostgreSQL connection and migrations
  validator/         — Custom Echo validator
  utils/             — Helpers (password hashing, env loading)
tests/
  controller/        — Controller unit tests
  service/           — Service unit tests
  repository/        — Repository unit tests
```

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/)

## Environment Variables

Create a `.env` file in the project root with the following variables:

| Variable | Description | Example |
|---|---|---|
| `APP_MODE` | Application mode (`DEV` or `PROD`) | `DEV` |
| `PORT` | Server port | `3000` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USERNAME` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `secret` |
| `DB_NAME` | PostgreSQL database name | `tracker` |
| `JWT_SECRET_KEY` | Secret key for signing JWT tokens | *(generate one)* |
| `JWT_EXPIRE_DURATION` | JWT token expiration in minutes | `1440` |
| `ADMIN_USERNAME` | Default admin username (for `cmd/generate`) | `admin` |
| `ADMIN_EMAIL` | Default admin email | `admin@example.com` |
| `ADMIN_PASSWORD` | Default admin password (min 6 characters) | `Admin123` |

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/mmegan-boop/go-income-expense-tracker-app.git
   cd go-income-expense-tracker-app
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env   # or create .env manually
   ```

   Edit `.env` with your database credentials and a secure JWT secret.

3. **Generate an admin account**

   ```bash
   go run ./cmd/generate
   ```

4. **Run the server**

   ```bash
   go run ./cmd/server
   ```

   Or with hot reload via Air:

   ```bash
   air
   ```

   The server starts on the port defined in `PORT` (default `3000`).

## API Endpoints

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

### Authentication

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/auth/register` | No | Register a new user |
| `POST` | `/api/auth/login` | No | Login and receive a JWT token |
| `POST` | `/api/auth/logout` | Yes | Logout (client-side token discard) |

### Users

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/users/me` | Yes | Get current user profile |
| `PUT` | `/api/users/me` | Yes | Update current user profile |
| `GET` | `/api/users` | Yes (admin) | Get all users |
| `DELETE` | `/api/users/:id` | Yes (admin) | Delete a user by ID |

### Categories

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/categories` | Yes (admin) | Create a new category |
| `GET` | `/api/categories` | Yes | Get all categories |
| `GET` | `/api/categories/:id` | Yes | Get a category by ID |
| `PUT` | `/api/categories/:id` | Yes (admin) | Update a category |
| `DELETE` | `/api/categories/:id` | Yes (admin) | Delete a category |

### Records

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/records` | Yes | Create a new income/expense record |
| `GET` | `/api/records` | Yes | Get all records for the current user |
| `GET` | `/api/records/:id` | Yes | Get a record by ID |
| `PUT` | `/api/records/:id` | Yes | Update a record |
| `DELETE` | `/api/records/:id` | Yes | Delete a record |
| `GET` | `/api/records/report` | Yes | Export records as an Excel file (`?start_date=DD-MM-YYYY&end_date=DD-MM-YYYY`) |
| `GET` | `/api/records/summary` | Yes | Get monthly income/expense summary (`?month=MM-YYYY`) |

## Testing

Run all tests:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./tests/controller
go test ./tests/service
go test ./tests/repository
```

## Postman Collection

A Postman collection with all API endpoints, request examples, and test scripts is available in the `docs/` folder. You can import it directly into Postman:

1. Open Postman
2. Click **Import** > **File**
3. Select `docs/go-income-expense-tracker-app.postman_collection.json`

## License

This project is licensed under the MIT License.

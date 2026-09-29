# OTP Based User Login

A full-stack web application with user registration, OTP-based login, and checkout functionality.

---

## Features

- **Registration** — email + name → generates a random 6-digit OTP displayed to the user
- **Email Recognition** — debounced background check while filling checkout form
- **OTP Login** — modal prompts registered users; correct code shows welcome message
- **Guest Checkout** — skip login and submit as guest
- **Checkout Submission** — stores email, phone, shipping address in PostgreSQL
- **Responsive UI** — works on desktop, tablet, and mobile

---

## Architecture

```
frontend/          ← Static HTML/CSS/JS (Vercel)
backend/           ← Go REST API (Fly.io)
database/          ← PostgreSQL schema (Vercel Postgres / Neon)
```

The frontend communicates with the backend exclusively through REST APIs.  
The backend communicates with PostgreSQL using parameterized queries via `lib/pq`.

---

## Technology Stack

| Layer    | Technology                        |
|----------|-----------------------------------|
| Frontend | HTML5, CSS3, Vanilla JavaScript   |
| Backend  | Go 1.21, chi router, lib/pq       |
| Database | PostgreSQL (Vercel Postgres/Neon) |
| Hosting  | Frontend: Vercel, Backend: Fly.io |

---

## Project Structure

```
OTP_Based/
├── frontend/
│   ├── index.html        # Checkout page
│   ├── register.html     # Registration page
│   ├── styles.css        # Responsive stylesheet
│   ├── config.js         # API base URL config
│   ├── api.js            # API service layer
│   ├── checkout.js       # Checkout page logic
│   └── register.js       # Registration page logic
├── backend/
│   ├── cmd/
│   │   ├── main.go       # Entry point
│   │   └── main_test.go  # Validation tests
│   ├── internal/
│   │   ├── database/db.go
│   │   ├── models/models.go
│   │   ├── repositories/
│   │   │   ├── user_repository.go
│   │   │   └── checkout_repository.go
│   │   ├── services/
│   │   │   ├── auth_service.go
│   │   │   └── checkout_service.go
│   │   └── handlers/
│   │       ├── auth_handler.go
│   │       ├── checkout_handler.go
│   │       └── response.go
│   ├── go.mod
│   ├── Dockerfile
│   └── fly.toml
├── database/
│   ├── schema.sql        # Table definitions
│   └── seed.sql          # Sample data
├── .env.example
├── .gitignore
├── vercel.json
├── README.md
└── prompts.md
```

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [PostgreSQL 14+](https://www.postgresql.org/download/) (local) **or** a Vercel Postgres / Neon connection string
- A modern web browser

---

## Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/OTP_Based.git
cd OTP_Based
```

### 2. Set up environment variables

```bash
cp .env.example backend/.env
# Edit backend/.env and fill in your database credentials
```

### 3. Set up the database

```bash
# Connect to your PostgreSQL instance and run:
psql -U postgres -d otp_login -f database/schema.sql
# Optional seed data:
psql -U postgres -d otp_login -f database/seed.sql
```

Or using a connection string:

```bash
psql "$DATABASE_URL" -f database/schema.sql
```

### 4. Run the backend

```bash
cd backend
go mod download
go run ./cmd/main.go
# Server starts on http://localhost:8080
```

### 5. Run the frontend

Open `frontend/index.html` directly in a browser, **or** use a local server:

```bash
# Using Python
cd frontend
python -m http.server 5500
# Open http://localhost:5500
```

Or use the VS Code Live Server extension.

---

## Environment Variables

| Variable         | Description                                      | Default         |
|------------------|--------------------------------------------------|-----------------|
| `DATABASE_URL`   | Full Postgres connection string (preferred)      | —               |
| `DB_HOST`        | Postgres host (if DATABASE_URL not set)          | `localhost`     |
| `DB_PORT`        | Postgres port                                    | `5432`          |
| `DB_USER`        | Postgres user                                    | `postgres`      |
| `DB_PASSWORD`    | Postgres password                                | —               |
| `DB_NAME`        | Database name                                    | `otp_login`     |
| `DB_SSLMODE`     | SSL mode (`disable` for local, `require` for prod)| `disable`      |
| `PORT`           | HTTP server port                                 | `8080`          |
| `ALLOWED_ORIGIN` | CORS allowed origin (your frontend URL)          | `localhost:*`   |

---

## API Documentation

### POST /api/auth/register

Register a new user.

**Request:**
```json
{ "email": "john@example.com", "firstName": "John", "lastName": "Doe" }
```

**Response 201:**
```json
{ "message": "Registration successful", "code": "482931", "firstName": "John", "lastName": "Doe", "email": "john@example.com" }
```

**Errors:** `400` (validation), `409` (duplicate email), `500`

---

### GET /api/users/recognize?email=john@example.com

Check if an email is registered.

**Response 200:**
```json
{ "registered": true }
```

---

### POST /api/auth/verify

Verify OTP code.

**Request:**
```json
{ "email": "john@example.com", "code": "482931" }
```

**Response 200:**
```json
{ "message": "Login successful", "firstName": "John", "lastName": "Doe", "email": "john@example.com" }
```

**Errors:** `400` (validation), `401` (wrong code), `500`

---

### POST /api/checkout

Submit checkout form.

**Request:**
```json
{ "email": "john@example.com", "phone": "9876543210", "shippingAddress": "Bengaluru, Karnataka" }
```

**Response 201:**
```json
{ "success": true, "message": "Checkout submitted successfully" }
```

**Errors:** `400` (validation), `500`

---

### GET /health

Health check — returns `{"status":"ok"}`.

---

## Testing

### Backend unit tests (no DB required)

```bash
cd backend
go test ./cmd/... -v
```

### Manual end-to-end test

1. Open `http://localhost:5500/register.html`
2. Register with email, first name, last name → note the 6-digit code
3. Open `http://localhost:5500/index.html`
4. Enter the registered email → OTP modal appears
5. Enter the correct code → welcome message appears
6. Fill phone + address → submit → success message
7. Verify in DB: `SELECT * FROM checkout_submissions;`

Also test:
- Wrong OTP → error stays in modal
- Skip → guest checkout works
- Duplicate registration → 409 error shown

---

## Deployment

### Database — Vercel Postgres or Neon

1. Create a project at [vercel.com](https://vercel.com) or [neon.tech](https://neon.tech)
2. Copy the connection string
3. Run `database/schema.sql` against it

### Backend — Fly.io

```bash
cd backend
fly auth login
fly launch          # first time — follow prompts, use existing fly.toml
fly secrets set DATABASE_URL="postgres://..." ALLOWED_ORIGIN="https://your-frontend.vercel.app"
fly deploy
```

### Frontend — Vercel

1. Edit `frontend/config.js` — set `window.API_BASE` to your Fly.io backend URL
2. Push to GitHub
3. Import repo in Vercel dashboard → deploy

Or via CLI:
```bash
npm i -g vercel
vercel --prod
```

---

## Public URLs

> Fill in after deployment:

- **Frontend:** `https://<your-project>.vercel.app`
- **Backend:** `https://otp-login-backend.fly.dev`

---

## Known Limitations

- OTP codes do not expire (no TTL). In production, add an `otp_expires_at` column.
- No rate limiting on OTP attempts. Add middleware for production.
- No HTTPS enforcement locally (handled by Fly.io in production).
- No email/SMS delivery — OTP is displayed in the UI as per assignment requirements.
- Go is required to be installed locally to run the backend.

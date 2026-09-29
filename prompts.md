# LLM Prompts Used During Development

This file documents every prompt/instruction given to an LLM (Amazon Q Developer) while building this application, in chronological order.

---

## Prompt 1 — Full Project Build Instruction

```
You are a senior full-stack software engineer. I need you to build the complete application described below.

IMPORTANT:
- Build the application fully, not just a prototype or UI mockup.
- Follow clean architecture and production-quality coding practices.
- Do not skip any required functionality.
- Before making major architectural decisions, inspect the existing project files.
- If the project is empty, initialize the project from scratch.
- Keep the frontend, API/backend, and database as distinct layers.
- Make the application easy to run locally and easy to deploy publicly.
- Use environment variables for all secrets and database credentials.
- Never hardcode passwords, API keys, JWT secrets, or database credentials.
- Add proper validation and error handling.
- Generate a README.md with complete setup and deployment instructions.
- Generate a prompts.md file containing every prompt/instruction used with an LLM while building this application.

PROJECT: OTP Based User Login

Goal: Build a web application with two main flows:
1. Registration Flow
2. User Recognition & Login Flow

The application also contains a checkout form whose submitted data must be stored in a database.

RECOMMENDED TECHNOLOGY STACK:
- Frontend: HTML/CSS, Responsive design
- Backend/API: Go, REST API, Proper layered architecture, JSON request/response format
- Database: Vercel Postgres
- Hosting: Frontend: Vercel, Backend: suitable public hosting service

ARCHITECTURE REQUIREMENT:
The application MUST have three distinct layers:
1. Frontend
2. API/backend
3. Database

[Full project specification including Flow 1 — User Registration, Flow 2 — User Recognition & Login,
OTP Validation, Logged-in User Experience, Skip Login / Guest Checkout, Checkout Submission,
Database Design, API Requirements, Security, Frontend UX, Error Cases, Testing, Project Structure,
README.md requirements, prompts.md requirements, GitHub Requirements, Deployment, and Acceptance Criteria]
```

This single comprehensive prompt drove the entire implementation. The LLM (Amazon Q Developer) then:

1. Inspected the empty project directory
2. Created the full folder structure
3. Implemented the PostgreSQL schema (`database/schema.sql`, `database/seed.sql`)
4. Implemented the Go backend with layered architecture:
   - `backend/internal/models/models.go`
   - `backend/internal/database/db.go`
   - `backend/internal/repositories/user_repository.go`
   - `backend/internal/repositories/checkout_repository.go`
   - `backend/internal/services/auth_service.go`
   - `backend/internal/services/checkout_service.go`
   - `backend/internal/handlers/auth_handler.go`
   - `backend/internal/handlers/checkout_handler.go`
   - `backend/internal/handlers/response.go`
   - `backend/cmd/main.go`
   - `backend/cmd/main_test.go`
5. Implemented the frontend:
   - `frontend/index.html` (checkout page with OTP modal)
   - `frontend/register.html` (registration page)
   - `frontend/styles.css` (responsive stylesheet)
   - `frontend/config.js` (API base URL configuration)
   - `frontend/api.js` (API service layer)
   - `frontend/checkout.js` (checkout logic with debounced recognition)
   - `frontend/register.js` (registration logic)
6. Created deployment configuration:
   - `backend/Dockerfile`
   - `backend/fly.toml`
   - `vercel.json`
7. Created project configuration:
   - `.env.example`
   - `.gitignore`
   - `README.md`
   - `prompts.md` (this file)

---

## Implementation Notes

All code was generated in a single session by Amazon Q Developer responding to the above prompt.
No additional prompts were required as the initial prompt was comprehensive enough to drive
the complete implementation from architecture through deployment configuration.

# oph26-backend

## Architecture Overview

This backend is built using [Go](https://go.dev/) and [Fiber](https://gofiber.io/) web framework, with [GORM](https://gorm.io/) as the ORM for database interactions. Please refer to the [Golang Clean Architecture](https://github.com/khannedy/golang-clean-architecture) for more details on the architectural patterns used.

## 🛠️ Tech Stack

- **Language**: Go
- **Web Framework**: Fiber
- **ORM**: GORM
- **Database**: PostgreSQL

## 🧑‍💻 Local development

1. **Clone the repository**:
   ```bash
   git clone
   ```

````
2. **Navigate to the project directory**:
   ```bash
   cd oph26-backend
````

3. **Set up environment variables**:
   Copy the `.env.template` file to `.env` and fill in the required values:

   ```bash
   cp .env.example .env
   ```

4. **Run the application**:
   ```bash
   go run main.go
   ```

### Mocking Google Token Validation in Development

In development mode, you can bypass Google token validation for easier testing. To do this, set the `APP_ENV` variable in your `.env` file to `development`:

```env
APP_ENV=development
```

When `APP_ENV` is set to `development`, the authentication middleware will allow you to use a JWT token with any claims, and it will log the claims for debugging purposes. This is intended for development and testing only, and should never be used in production.

You can generate a JWT token to pass to `/api/auth/token` endpoint for testing using this command:

```terminal
$ go run cmd/script/generate_mock_google_token.go -email test@example.com
Email: attendee@example.com
Mock Google OAuth Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImF0dGVuZGVlQGV4YW1wbGUuY29tIn0.kcpt1G9BDldx615oC-9YYk9gH1K9ojWZ4WlkMAXqzh0
```

You'll get a token that you can send to `/api/auth/token` to authenticate as `test@example.com`. The token will be accepted in development mode, and the claims will be logged for debugging.

Here is an example request using JavaScript's `fetch` API:

```javascript
const response = await fetch("http://localhost:8080/api/auth/token", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    idToken: "eyJhbGciOi...AXqzh0", // Replace with your generated token
  }),
});

const data = await response.json();
console.log(data);
// {
//  "accessToken": "eyJhbGciOi...",
//  "refreshToken": "eyJhbGciOi..."
// }
// These tokens can then be used to authenticate subsequent requests to protected endpoints in development mode using the Authorization header:
fetch("http://localhost:8080/api/protected-endpoint", {
  headers: {
    Authorization: "Bearer eyJhbGciOi...", // Replace with the access token you received
  },
});
```

Access tokens is valid for 15 minutes, and refresh tokens is valid for 7 days. You can use the refresh token to get a new access token when it expires by sending a request to `/api/auth/refresh` endpoint.

---

## 🔒 Branch Rules for Main

The `main` branch is protected with the following rules:

1. **Require a pull request before merging**
   - Direct pushes to `main` are blocked.
   - All changes must be submitted via pull request (PR).

2. **Require branches to be up to date before merging**
   - The PR branch must be rebased or merged with the latest `main` before merging.

3. **No bypass allowed**
   - These rules apply to everyone, including administrators.

These rules ensure `main` always contains production-ready, tested, and reviewed code.

---

## 🏷 Branch Naming Scheme

Follow this format:

```

<name>/<type>/<short-description>

```

### Example:

```

arka/feat/resume-button

```

### Types:

- **feat/** – New features (e.g., `jay/feat/add-appointment-page`)
- **fix/** – Bug fixes (e.g., `alex/fix/responsive-layout`)
- **chore/** – Maintenance tasks, dependency updates, config changes (e.g., `sam/chore/update-tailwind-config`)
- **refactor/** – Code refactoring without changing functionality (e.g., `lee/refactor/dashboard-layout`)
- **remove/** – Removing unused code, dependencies, or features (e.g., `jay/remove/prisma`)

### More Examples:

- `mike/feat/user-profile`
- `anna/refactor/api-service`
- `jay/remove/old-api-endpoints`

## Database ER Diagram

```mermaid
erDiagram
    direction TB
    Attendees    ||--o|  Questionnaire : surveys
    Attendees    ||--||  Users: attendee_accounts
    Staffs       |o--o|  Users: staff_accounts
    Attendees    |o--o{  CollectedPieces : collected_pieces
    MyPieces     ||--||  Attendees: user_pieces
    Leaderboards ||--||  Attendees: top_scorers
    Scores       }|--||  Attendees: score
    Attendees    }|--||  Staffs: checkin_staff

    Users {
        string id PK "not null, UUID"
        string email UK "not null"
        string role "not null, enum: attendee, staff"
        string attendee_id FK "nullable, UUID"
        string staff_id FK "nullable, UUID"
        datetime created_at "not null"
        datetime updated_at "not null"
    }

    Attendees {
        string id PK "not null, UUID"
        string user_id UK, FK "not null, UUID"
        string firstname "not null"
        string surname "not null"
        string attendee_type "not null, enum: Matthayom, Prathom, Parent, EducationStaff, Other"
        int age "not null"
        string province "not null"
        string study_level "nullable, for Student type (Matthayom/Prathom)"
        string school_name "nullable, for Student type (Matthayom/Prathom)"
        string[] news_sources_selected "nullable, array of selected values"
        string news_sources_other "nullable, conditional other"
        string initial_first_interested_faculty "not null"
        string[] interested_faculty "not null, array of interested faculty"
        string[] objective_selected "nullable, array of selected values"
        string objective_other "nullable, conditional other"
        string(7) ticket_code UK "not null,<br/>C000000, prefix consult docs"
        datetime created_at "not null"
        datetime updated_at "not null"
        string my_piece_id FK "nullable, uuid"
        string certificate_name "nullable"
        %% checkedin
        datetime checked_in_at "nullable"
        string checkin_staff_id FK "nullable"
        %% workshops
        string[] favorite_workshops "nullable, array of id of attendee's favorite workshops"
    }

    Staffs {
        string id PK "not null, UUID"
        string user_id UK, FK "nullable, UUID"
        string cuid UK "not null, student id or uni staff uid"
        string firstname "not null"
        string surname "not null"
        string nickname "not null"
        string phone "not null"
        string year "not null, ชั้นปี, or อื่น ๆ for shared account"
        string email UK "not null, predefined from collected form"
        string faculty "not null"
        datetime created_at "not null, default current_timestamp"
        datetime updated_at "not null, default current_timestamp"
    }

    MyPieces {
        string id PK "not null, UUID"
        string user_id FK "not null"
        string piece_code UK "not null, shorter unique code to send to user"
        datetime expire_date "not null"
    }

    CollectedPieces {
        string user_id UK, FK, PK "not null"
        string piece_id UK, FK, PK "not null, uuid of piece"
        datetime collected_at "not null"
    }

    Scores {
        string user_id UK, FK, PK "not null"
        int[20] count "not null, default 0, array of pieces count per faculty"
    }


    %% Each row in Leaderboards records, for a given user_id, in which faculties the user is in the top 1% at the given updated_at time.
    Leaderboards {
        string user_id UK,FK,PK "not null, user id"
        bool[20] is_top "not null, default false, array of top1% status per faculty"
        datetime updated_at "not null"
    }

    Questionnaire {
        string user_id UK, FK "not null"
        json_text answers "not null"
        datetime created_at "not null"
    }
```

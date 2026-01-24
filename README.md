# oph26-backend
ฺBackend repository for Chula Openhouse 2026 website and services.

## Architecture Overview
This backend is built using [Go](https://go.dev/) and [Fiber](https://gofiber.io/) web framework, with [GORM](https://gorm.io/) as the ORM for database interactions. Please refer to the [Golang Clean Architecture](https://github.com/khannedy/golang-clean-architecture) for more details on the architectural patterns used.


## Database ER Diagram
```mermaid
erDiagram
    direction TB
    Attendees    ||--o|  Questionnaire : surveys
    Attendees    ||--||  Users: attendee_accounts
    Staffs       |o--o|  Users: staff_accounts
    Attendees    |o--o{  CollectedPieces : collected_pieces
    MyPieces     ||--||  Attendees: user_pices
    Leaderboards ||--||  Attendees: top_scorers
    Scores       }|--||  Attendees: score
    Attendees    }|--||  Staffs: checkin_staff

    Users {
        string id PK "not null, UUID"
        string email UK "not null"
        string role "not null, enum: attendee, staff"
        string attendee_id FK "nullable, UUID"
        string staff_id FK "nullable, UUID"
        string piece_id FK "nullable, UUID"
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
        string study_level "nullable, for Student type (mathayom/prathom)"
        string school_name "nullable, for Student type (mathayom/prathom)"
        string[] news_sources_selected "nullable, array of selected value"
        string news_sources_other "nullable, conditional other"
        string initial_first_interested_faculty "not null"
        string[] interested_faculty "not null, array of interested faculty"
        string[] objective_selected "nullable, array of selected value"
        string objective_other "nullable, conditional other"
        string(7) ticket_code UK "not null,<br/>C000000, prefix consult docs"
        datetime created_at "not null"
        datetime updated_at "not null"
        string my_piece_id FK "nullable, uuid"
        string certificate_name "nullable"
        %% checkedin
        datetime checkined_at "nullable"
        string checkin_staff_id FK "nullable"
        %% workshops
        string[] favorite_workshops "nullable, array of id of attendee's favorite workshops"
        %% certificate
        string certificate_name "nullable"
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
        %%ควรเป็นอันที่เอาไปส่งให้คนอื่น
        string user_id FK "not null"
        string piece_code UK "not null"
        datetime expire_date "not null"
    }

    CollectedPieces {
        string user_id UK, FK, PK "not null"
        string piece_id UK, FK, PK "not null"
        datetime collected_at "not null"
    }

    Scores {
        string user_id UK, FK, PK "not null"
        int[20] count "not null, default 0, array of pieces count per faculty"
    }

    %% Leaderboards is this user_id top1% of any of this faculty
    Leaderboards {
        string user_id UK,FK,PK "not null, user id"
        bool[20] is_top "not null, default false, array of top1% status per faculty"
        datetime updated_at "not null"
    }

    Questionnaire {
        string user_id UK, FK "not null"
        json_text answers "not null"
        datedtime created_at "not null"
    }
```

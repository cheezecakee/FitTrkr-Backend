## Project Structure

```bash
Fitrkr/
├── cmd/
│   └── server/
│       └── main.go                     # Entry point of the application
├── internal/
│   ├── api/                            # Presentation layer (handlers)
│   │   ├── handler/                    # HTTP route handlers
│   │   │   ├── middleware.go           # Middleware (e.g., auth, logging)
│   │   │   ├── user_handler.go         # User Crud handlers 
│   │   │   ├── auth_handler.go         # Auth handlers
│   │   │   ├── plan_handler.go         # Plan Crud handlers 
│   │   │   ├── plan_exercise_handler.go# Plan Exercise Crud handlers 
│   │   │   └── handler.go              # Holds generic handler functions and context key consts (error handling will be moved to errors/ at a later time) 
│   │   ├── api.go                      # Initializes handler structs for use in routes
│   ├── router/                         # Routes and the router setup
│   │   ├── router.go      
│   │   ├── v1/      
│   │   └──  └── routes.go        
│   ├── app/                            # Presentation layer (handlers)
│   │   └── app.go                      # Connects all db repos to services
│   ├── models/                         # Data structs for tables
│   │   ├── user.go                     # User struct
│   │   ├── plan.go                     # Plan struct
│   │   ├── session.go                  # Session struct (and so on)
│   │   ├── log.go                      # Log struct
│   │   └── exercise.go                 # Exercise struct
│   ├── services/                       # Business logic
│   │   ├── auth_service.go             # Authentication logic (e.g., JWT)
│   │   ├── plan_service.go             # Plan-related logic
│   │   ├── plan_exercise_service.go    # Plan-related logic
│   │   ├── session_service.go          # Session-related logic
│   │   ├── session_exercise_service.go # Session-related logic
│   │   ├── exercise_set_service.go     # Session-related logic
│   │   ├── user_service.go             # User-related logic
│   │   └── log_service.go              # Logging logic
│   ├── repository/                     # Database access layer
│   │   ├── db.go                       # DB connection setup
│   │   ├── user_repo.go                # User-related queries
│   │   ├── plan_repo.go                # Plan-related queries
│   │   ├── plan_exercise_repo.go       # Plan-related queries
│   │   ├── session_exercise_repo.go    # Session-related logic
│   │   ├── exercise_set_repo.go        # Session-related queries
│   │   ├── session_repo.go             # Session-related queries (etc.)
│   │   └── log_repo.go                 # Logging logic
│   └── utils/                          # Utilities
│       ├── config/                     # Configuration loading (e.g., env vars)
│       │   └── config.go        
│       ├── auth/                       # JWT helper functions
│       │   └── jwt.go        
│       ├── helper/                     # Helper functions
│       │   └── helper.go        
│       └── transaction/                # Transaction wrapper for repos (Base repo)
│           └── transaction.go        
├── pkg/                                # Optional: reusable external packages (if any) (will probably move them the stuff here to internal/utils/ later on)
│       ├── errors.go/                  # Custom error handler (Currently not in use)
│       │   └── errors.go               # Error setup    
│       └── logger/                     # Custom logger handler (Currently not in use)
│           └── logger.go               # Logging setup     
│                           
├── migrations/                         # SQL migration files for schema
│   │   └── schema/        
│   │        ├── 001__users.sql
│   │        ├── 002__refresh_tokens.sql (currently not in use)
│   │        ├── 003__exercises.sql
│   │        ├── 004__plans.sql
│   │        ├── 005__sessions.sql
│   │        ├── 006__logs.sql
│   │        └── 007__plans.sql (etc.)
├── go.mod                              # Go module file
└── go.sum                              # Go dependencies
```

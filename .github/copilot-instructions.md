# Instagram Clone - AI Coding Instructions

## Architecture Overview
This is a Go-based Instagram clone with a clean architecture using the Repository pattern. The app consists of:
- **Backend**: Go REST API with Gorilla Mux router, GORM for database access, JWT authentication
- **Frontend**: Vanilla HTML/CSS/JavaScript SPA that consumes the REST API
- **Database**: MySQL 8.0 running in Docker
- **Deployment**: Docker Compose for local development

## Key Patterns & Conventions

### Repository Pattern Structure
- **Interfaces**: Define contracts in `repo/` (e.g., `UserRepo`, `PostRepo`)
- **Implementations**: GORM implementations like `repo/gorm_user_repo.go`
- **Dependency Injection**: Repositories injected into handlers via constructor in `main.go`

### Handler Architecture
- **Interface**: `handler/interface.go` defines `InstaHandler` contract
- **Implementation**: Split across multiple files (`login.go`, `post.go`, `user_handler.go`, etc.)
- **Constructor**: `NewInstaHandler()` in `account_handler.go` takes all repos as dependencies

### Authentication Flow
- **JWT Tokens**: Generated in `jwt/jwt.go` with hardcoded secret `"secret"`
- **Middleware**: `middleware/jwt_auth.go` validates Bearer tokens and adds username to context
- **Protected Routes**: Use subrouter with `JWTAuthMiddleware` applied
- **Frontend Storage**: Tokens stored in localStorage, auto-included in API calls

## Critical Development Workflows

### Docker Development (Primary)
```bash
# Use Docker Compose v2 syntax (space, not hyphen)
make up        # Start containers
make ra        # Quick app rebuild (most common during development)  
make logs      # View container logs
make down      # Stop everything
```
**Note**: Makefile recently updated to use `docker compose` instead of `docker-compose`

### Database Management
- **Auto-migration**: GORM AutoMigrate runs on startup via `database/gorm.go`
- **Model Registration**: Models registered in `models/models.go` init function
- **Connection**: Singleton pattern with retry logic for Docker startup timing

### API Development
- **Public Routes**: `/health`, `/createAccount`, `/login` (no auth required)
- **Protected Routes**: All others require `Authorization: Bearer <token>` header
- **CORS**: Currently uses `CustomCORS` middleware (file:// protocol compatibility)

## Project-Specific Conventions

### Error Handling
- **Standard Pattern**: Return specific HTTP status codes with JSON error messages
- **Authentication**: Always return 401 for invalid/missing tokens
- **Validation**: Return 400 for missing required fields with descriptive messages

### Database Models
- **GORM Tags**: Use explicit column names (e.g., `gorm:"column:username"`)
- **Primary Keys**: Use `Id int` with `gorm:"primaryKey;column:id"`
- **Relationships**: Posts have UserID, Comments have PostID/CommentID for nesting

### Frontend Integration
- **API Base**: `http://localhost:8080` (hardcoded in `frontend/script.js`)
- **Authentication**: Frontend automatically includes JWT in headers
- **Error Display**: Toast notifications for user feedback

## Key Files for Understanding
- `main.go`: Dependency injection and routing setup
- `handler/interface.go`: Complete API contract
- `database/gorm.go`: Singleton DB connection with retry logic
- `middleware/jwt_auth.go`: Authentication flow
- `frontend/script.js`: Complete frontend API integration patterns

## Quick Debugging
- **Container Issues**: Check `make logs` for startup errors
- **Database**: Connection retries happen automatically (5 attempts)
- **CORS Issues**: Switch between `CORS` and `CustomCORS` middleware in `main.go`
- **Frontend**: Browser developer tools show API calls and localStorage tokens

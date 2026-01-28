# Auth Service

A lightweight authentication microservice written in Go, providing JWT-based authentication with support for traditional email/password login and OAuth2 providers (Google and GitHub).

## Features

- 🔐 **JWT Authentication** - Secure access tokens with configurable expiration
- 🔄 **Refresh Tokens** - Secure token rotation with hashed storage
- 📧 **Email/Password Registration & Login**
- ✅ **Email Verification** - One-time tokens with expiration
- 🔑 **Password Reset** - One-time tokens with expiration
- 🌐 **OAuth2 Integration**
  - Google (OpenID Connect)
  - GitHub
- 🛡️ **Secure Password Hashing** - Using bcrypt
- 🗄️ **PostgreSQL Database** - With connection pooling via pgx
- 🌍 **CORS Support** - Configurable allowed origins
- 🐳 **Docker Ready** - Docker Compose setup included

## Tech Stack

- **Go 1.25+**
- **PostgreSQL 16**
- **[Chi Router](https://github.com/go-chi/chi)** - Lightweight HTTP router
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver with connection pooling
- **[golang-jwt](https://github.com/golang-jwt/jwt)** - JWT implementation
- **[go-oidc](https://github.com/coreos/go-oidc)** - OpenID Connect client

## Project Structure

```
auth-service/
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── docker-compose.yml          # Docker Compose for PostgreSQL
├── sql/
│   └── schema.sql              # Database schema
└── internal/
    ├── auth/
    │   ├── jwt.go              # JWT signing and verification
    │   ├── password.go         # Password hashing (bcrypt)
    │   └── refresh.go          # Refresh token generation
    ├── config/
    │   └── config.go           # Environment configuration
    ├── db/
    │   └── db.go               # Database connection pool
    ├── http/
    │   ├── router.go           # HTTP routes definition
    │   ├── handler_auth.go     # Auth handlers (register, login, etc.)
    │   ├── handlers_google.go  # Google OAuth handlers
    │   ├── handlers_github.go  # GitHub OAuth handlers
    │   ├── middleware.go       # Auth middleware
    │   ├── cors.go             # CORS middleware
    │   ├── oauth_service.go    # OAuth user upsert logic
    │   └── oauth_types.go      # OAuth store interface
    └── store/
        ├── store.go            # Store interfaces
        └── postgres/
            ├── users.go        # User repository
            ├── tokens.go       # Refresh token repository
            └── oauth.go        # OAuth identity repository
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 16 (or use Docker)
- Google OAuth credentials (for Google login)
- GitHub OAuth credentials (for GitHub login)

### 1. Clone the Repository

```bash
git clone https://github.com/MihajasoaAlain/auth-service-go.git
cd auth-service-go
```

### 2. Start PostgreSQL

Using Docker Compose:

```bash
docker-compose up -d
```

This starts a PostgreSQL instance with:
- **User**: `auth`
- **Password**: `auth`
- **Database**: `authdb`
- **Port**: `5432`

### 3. Initialize the Database

Connect to PostgreSQL and run the schema:

```bash
docker exec -i $(docker-compose ps -q db) psql -U auth -d authdb < sql/schema.sql
```

Or manually:

```bash
psql -h localhost -U auth -d authdb -f sql/schema.sql
```

### 4. Configure Environment Variables

Create a `.env` file in the project root:

```env
# Database
DATABASE_URL=postgres://auth:auth@localhost:5432/authdb?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ISSUER=auth-service

# CORS (comma-separated origins, or * for all)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback

# Web Application Origin (for OAuth redirects)
WEB_ORIGIN=http://localhost:3000
```

### 5. Run the Service

```bash
go run main.go
```

The service will start on `http://localhost:8080`.

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/auth/register` | Register a new user | No |
| `POST` | `/auth/login` | Login with email/password | No |
| `POST` | `/auth/refresh` | Refresh access token | No |
| `POST` | `/auth/logout` | Revoke refresh token | No |
| `POST` | `/auth/verify/start` | Start email verification | No |
| `POST` | `/auth/verify/confirm` | Confirm email verification | No |
| `POST` | `/auth/password/forgot` | Start password reset | No |
| `POST` | `/auth/password/reset` | Reset password with token | No |
| `GET` | `/me` | Get current user info | Yes |

### OAuth2

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/auth/google/start` | Initiate Google OAuth flow |
| `GET` | `/auth/google/callback` | Google OAuth callback |
| `GET` | `/auth/github/start` | Initiate GitHub OAuth flow |
| `GET` | `/auth/github/callback` | GitHub OAuth callback |

## API Usage Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword123"}'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com"
}
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword123"}'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

### Access Protected Route

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."}'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "refresh_token": "bmV3IHJlZnJlc2ggdG9rZW4..."
}
```

### Logout

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."}'
```

### OAuth2 Login Flow

1. **Redirect user to start OAuth:**
   ```
   GET http://localhost:8080/auth/google/start
   GET http://localhost:8080/auth/github/start
   ```

2. **User authenticates with provider**

3. **Callback redirects to your web app:**
   ```
   http://localhost:3000/oauth/success?access=<token>&refresh=<token>
   ```

## Database Schema

The service uses the following tables:

### `users`
Stores user information.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `email` | TEXT | Unique email address |
| `password_hash` | TEXT | Bcrypt hashed password |
| `name` | TEXT | User's display name |
| `avatar_url` | TEXT | Profile picture URL |
| `created_at` | TIMESTAMPTZ | Creation timestamp |

### `user_identities`
Links OAuth providers to users (supports multiple providers per user).

| Column | Type | Description |
|--------|------|-------------|
| `id` | BIGSERIAL | Primary key |
| `user_id` | UUID | Foreign key to users |
| `provider` | TEXT | OAuth provider (google, github) |
| `provider_user_id` | TEXT | User ID from provider |
| `email` | TEXT | Email from provider |
| `created_at` | TIMESTAMPTZ | Creation timestamp |

### `refresh_tokens`
Stores hashed refresh tokens.

| Column | Type | Description |
|--------|------|-------------|
| `token_hash` | TEXT | SHA-256 hash of token (PK) |
| `user_id` | UUID | Foreign key to users |
| `expires_at` | TIMESTAMPTZ | Token expiration |
| `revoked_at` | TIMESTAMPTZ | Revocation timestamp (nullable) |

### `webauthn_credentials`
Stores WebAuthn credentials (for future passkey support).

| Column | Type | Description |
|--------|------|-------------|
| `credential_id` | TEXT | Credential ID (PK) |
| `user_id` | UUID | Foreign key to users |
| `public_key` | BYTEA | Public key |
| `sign_count` | BIGINT | Signature counter |
| `transports` | TEXT | Supported transports |
| `created_at` | TIMESTAMPTZ | Creation timestamp |

## Token Configuration

| Token Type | Default TTL | Description |
|------------|-------------|-------------|
| Access Token | 15 minutes | Short-lived JWT for API access |
| Refresh Token | 14 days | Long-lived token for obtaining new access tokens |

## Security Features

- **Password Hashing**: Uses bcrypt with default cost
- **Refresh Token Storage**: Tokens are hashed (SHA-256) before storage
- **Token Rotation**: Refresh tokens are rotated on each use (old token revoked)
- **CSRF Protection**: OAuth state parameter prevents CSRF attacks
- **HttpOnly Cookies**: OAuth state stored in HttpOnly cookies

## Setting Up OAuth Providers

### Google

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Navigate to **APIs & Services** > **Credentials**
4. Click **Create Credentials** > **OAuth client ID**
5. Select **Web application**
6. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
7. Copy the Client ID and Client Secret to your `.env`

### GitHub

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click **New OAuth App**
3. Set Authorization callback URL: `http://localhost:8080/auth/github/callback`
4. Copy the Client ID and Client Secret to your `.env`

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o auth-service main.go
```

### Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | Yes | - | Secret key for signing JWTs |
| `JWT_ISSUER` | No | - | JWT issuer claim |
| `CORS_ALLOWED_ORIGINS` | No | `*` | Comma-separated allowed origins |
| `GOOGLE_CLIENT_ID` | No* | - | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | No* | - | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | No* | - | Google OAuth callback URL |
| `GITHUB_CLIENT_ID` | No* | - | GitHub OAuth client ID |
| `GITHUB_CLIENT_SECRET` | No* | - | GitHub OAuth client secret |
| `GITHUB_REDIRECT_URL` | No* | - | GitHub OAuth callback URL |
| `WEB_ORIGIN` | No* | - | Frontend origin for OAuth redirects |

*Required if using the respective OAuth provider

## License

This project is open source. See the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

# Authentication & Authorization Implementation Plan

## Overview
This document outlines the plan to add JWT-based authentication and authorization to the Go Dashboard application.

## Architecture

### Backend Components

#### 1. **User Management**
- **User Model**: Store user credentials (ID, username/email, hashed password)
- **Storage**: Start with in-memory storage (map), easily extensible to database later
- **Password Hashing**: Use `golang.org/x/crypto/bcrypt` for secure password hashing
- **Initial User**: Create a default admin user for testing

#### 2. **JWT Token Management**
- **Library**: Use `github.com/golang-jwt/jwt/v5` for JWT handling
- **Token Structure**: 
  - Claims: `user_id`, `username`, `email`, `exp` (expiration)
  - Expiration: 24 hours (configurable)
  - Secret Key: Environment variable or config file
- **Token Generation**: After successful login
- **Token Validation**: Middleware to verify tokens on protected routes

#### 3. **Authentication Endpoints**
- `POST /api/auth/login` - Login with username/email and password
  - Request: `{ "username": "user@example.com", "password": "password" }`
  - Response: `{ "token": "jwt_token_here", "user": { "id": 1, "username": "user@example.com" } }`
- `POST /api/auth/logout` - Logout (optional, mainly for frontend state management)
- `GET /api/auth/me` - Get current user info (protected)

#### 4. **Authorization Middleware**
- **JWT Middleware**: Verify token from `Authorization: Bearer <token>` header
- **Protected Routes**: All `/api/*` endpoints except `/api/auth/login`
- **Page Protection**: All page routes (`/`, `/sales`, `/inventory`) should check authentication
- **Response**: Return 401 Unauthorized if token is missing/invalid

#### 5. **Frontend Integration**
- **Token Storage**: Store JWT in `localStorage` (or httpOnly cookies for better security)
- **Request Interceptor**: Automatically add `Authorization` header to all API requests
- **Login Page**: Create `/login` route with login form
- **Redirect Logic**: 
  - Unauthenticated users → redirect to `/login`
  - Authenticated users accessing `/login` → redirect to `/`
- **Logout**: Clear token and redirect to login

### Frontend Components

#### 1. **Login Page** (`/login`)
- Form with username/email and password fields
- Submit to `/api/auth/login`
- Store token on success
- Display error messages on failure
- Redirect to home page after successful login

#### 2. **Authentication State Management**
- Check for token on page load
- Verify token validity (optional: call `/api/auth/me`)
- Protect client-side navigation

#### 3. **API Request Enhancement**
- Add `Authorization: Bearer <token>` header to all fetch/HTMX requests
- Handle 401 responses by redirecting to login

#### 4. **UI Updates**
- Add logout button to sidebar
- Show current user info in sidebar (replace hardcoded "Admin User")
- Hide/show navigation based on auth state

## Implementation Steps

### Phase 1: Backend Foundation
1. Add dependencies to `go.mod`:
   - `github.com/golang-jwt/jwt/v5`
   - `golang.org/x/crypto/bcrypt`

2. Create `auth/` package:
   - `user.go` - User model and storage
   - `jwt.go` - JWT token generation and validation
   - `middleware.go` - Authentication middleware
   - `handlers.go` - Login/logout handlers

3. Create `models/` package (optional, for shared types)

### Phase 2: Backend Endpoints
1. Implement login handler
2. Implement JWT middleware
3. Protect existing API endpoints
4. Protect page routes

### Phase 3: Frontend Authentication
1. Create login page template
2. Add login route handler
3. Implement token storage and retrieval
4. Add authentication check to page handlers
5. Update API calls to include token

### Phase 4: UI/UX Enhancements
1. Add logout functionality
2. Update sidebar with user info
3. Add authentication state indicators
4. Handle token expiration gracefully

## Security Considerations

1. **Password Security**:
   - Never store plain text passwords
   - Use bcrypt with appropriate cost factor (10-12)
   - Validate password strength (optional)

2. **Token Security**:
   - Use strong secret key (environment variable)
   - Set reasonable expiration times
   - Consider refresh tokens for long sessions (future enhancement)

3. **HTTPS**: 
   - In production, always use HTTPS
   - Tokens should only be sent over secure connections

4. **CORS** (if needed):
   - Configure CORS properly if frontend is on different domain

5. **Rate Limiting** (future):
   - Add rate limiting to login endpoint to prevent brute force

## File Structure After Implementation

```
go-dashboard/
├── api/
│   └── data.go
├── auth/
│   ├── user.go          # User model and storage
│   ├── jwt.go           # JWT utilities
│   ├── middleware.go    # Auth middleware
│   └── handlers.go      # Login/logout handlers
├── templates/
│   ├── layout.html
│   ├── index.html
│   ├── sales.html
│   ├── inventory.html
│   └── login.html       # NEW: Login page
├── main.go
├── go.mod
└── ...
```

## Testing Strategy

1. **Manual Testing**:
   - Test login with valid credentials
   - Test login with invalid credentials
   - Test accessing protected routes without token
   - Test accessing protected routes with invalid token
   - Test token expiration
   - Test logout

2. **Edge Cases**:
   - Empty username/password
   - SQL injection attempts (if using database later)
   - XSS attempts in username
   - Token manipulation

## Future Enhancements

1. **Database Integration**: Replace in-memory storage with database
2. **User Registration**: Add signup endpoint
3. **Password Reset**: Implement password reset flow
4. **Role-Based Access Control (RBAC)**: Add user roles and permissions
5. **Refresh Tokens**: Implement token refresh mechanism
6. **Session Management**: Track active sessions
7. **Multi-Factor Authentication (MFA)**: Add 2FA support

## Configuration

### Environment Variables
- `JWT_SECRET`: Secret key for signing JWT tokens (required in production)
- `JWT_EXPIRATION_HOURS`: Token expiration time (default: 24)

### Default Credentials (for development)
- Username: `admin@dashboard.com`
- Password: `admin123` (should be changed in production)

## Dependencies to Add

```go
require (
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/crypto v0.17.0
)
```


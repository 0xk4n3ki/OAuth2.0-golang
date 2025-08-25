## Features
- Google OAuth 2.0 Authentication - Users can sign in using their Google accounts.
- Automatic Database Setup - Creates the go_oauth database if it doesn’t exist.
- Postgres Integration - Stores user details (first name, last name, email, tokens, timestamps, UUID).
- Pagination Support - List users with page and recordPerPage query params.
- Fetch Single User - Retrieve a user by UUID.
- Token Management - Stores access token and refresh token returned by Google OAuth.
- Error Handling - Validates UUID format, handles DB errors gracefully.

## Tech Stack
- Language: Go (Golang)
- Framework: Gin Gonic(web framework for REST APIs)
- Database: PostgreSQL
- Extensions: pgcrypto extension for generating UUIDs in Postgres
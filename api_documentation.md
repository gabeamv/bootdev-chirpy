# Chirpy Web Api Path Documentation

* /app/... 
Description: Retrieves web assets from fileserver. TBD.

* GET /admin/metrics 
Description: Retirieves the number of server requests for fileserver assets.
Request:
GET /admin/metrics HTTP/1.1
Host: localhost:8080
Response:
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited X times!</p>
  </body>
</html>

* POST /admin/reset 
Description: Deletes all registered users and their related data. Intended for developing and testing purposes.
Request:
POST /admin/reset HTTP/1.1
Host: localhost:8080
Response:
200 OK

* GET /api/healthz 
Description: Retrieves a status code 200 if server is working properly.
Request:
GET /api/healthz HTTP/1.1
Host: localhost:8080
Response:
200 OK

* POST /api/users 
Description: Registers a user to Chirpy.
Request:
POST /api/users HTTP/1.1
Content-Type: application/json

{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
Response: 
HTTP/1.1 201 Created 
Content-Type: application/json 

{
  "id": "3311741c-680c-4546-99f3-fc9efac2036c", 
  "email": "walt@breakingbad.com", 
  "is_chirpy_red": false, 
  "created_at": "2026-02-10T12:34:56Z" 
  "updated_at": "2026-02-10T12:34:56Z" 
} 

* POST /api/chirps
Description: Adds a chirp for a user.
Request:
POST /api/chirps HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "body": "Hello, world!"
}
Response:
201 Created
Content-Type: application/json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "body": "Hello, world!"
}

* GET /api/chirps
Description: Gets all chirps.
Query Parameters: author_id => gets all chirps for specifc user
Request:
GET /api/chirps?author_id=123e4567-e89b-12d3-a456-426614174000 HTTP/1.1
Host: localhost:8080
Response:
200 OK
Content-Type: application/json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "body": "Hello, world!"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "body": "Another chirp!"
  }
]

* GET /api/chirps/{chirpID}
Description: Get a specific chirp by ID
Request:
GET /api/chirps/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: localhost:8080
Response:
200 OK
Content-Type: application/json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "body": "Hello, world!"
}

* POST /api/login
Description: Login with valid credentials.
Request:
POST /api/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secretpassword"
}
Response:
200 OK
Content-Type: application/json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0..."
}

* POST /api/refresh
Description: Get a new access token using a valid refresh token.
Request:
POST /api/refresh HTTP/1.1
Host: localhost:8080
Authorization: Bearer a1b2c3d4e5f6g7h8i9j0...
Response:
200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

* POST /api/revoke
Description: Revoke a refresh token.
Request:
POST /api/revoke HTTP/1.1
Host: localhost:8080
Authorization: Bearer a1b2c3d4e5f6g7h8i9j0...
Response:
204 No Content

* PUT /api/users
Description: Change the email/password for an authorized user.
Request:
PUT /api/users HTTP/1.1
Content-Type: application/json

{
  "email": "walt@breakingbad.com",
  "password": "123456"
}
Response:
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "3311741c-680c-4546-99f3-fc9efac2036c",
  "email": "walt@breakingbad.com",
  "is_chirpy_red": false,
  "created_at": "2026-02-10T12:34:56Z"
  "updated_at": "2026-02-10T12:34:56Z"
}

* DELETE /api/chirps/{chirpID} 
Description: Delete a chirp by ID for an authorized user.
Request:
DELETE /api/chirps/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Response:
HTTP/1.1 204 No Content

* POST /api/polka/webhooks
Description: Webhook that only Polka (payment provider) is authorized to make calls to. Bootdev specific. Requires a valid Polka API key in Authorization header
Request:
POST /api/polka/webhooks HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: ApiKey your-polka-api-key-here

{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
Response:
204 No Content

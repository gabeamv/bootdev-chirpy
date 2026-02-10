# bootdev-chirpy

Chirpy is a guieded Boot-dev web api project inspired by Twitter/X, where users can register an account, login, and make write 'chirps' (tweets). 

## Features

* Authentication
    - Users can register an account and login. Refresh and access tokens are used for authenticated users to make requests to server resources.
* Chirps
    - Users can create and delete their own chirps.

## Requirements

* Go 1.21+ installed
* A terminal to run commands

## Installation

1. Clone the repository.
2. Install and setup PostgreSQL. 
    * On mac, run **brew install postgres@15** or another suitable version.
    * On Linux / WSL (Debian), run **sudo apt update**, then **sudo apt install postgresql postgresql-contrib**
3. Ensure you have the postgres by running **psql --version**
    * If you are on Linux, update postgres password: **sudo passwd postgres**
    * Don't forget the password.
4. Enter the psql shell: **psql postgres**
5. Create a new database called '**chirpy**': **CREATE DATABASE chirpy;**
6. Set the user password (Linux only): **ALTER USER postgres PASSWORD 'postgres';**
7. In the root of the cloned project, we are going to add a '.env' file with some values: 
    1. DB_URL = "postgres://{username}:{password}@localhost:5432/chirpy?sslmode=disable"
        * Replace {username} and {password} with your own configuration
    2. PLATFORM = "dev"
    3. SECRET = "{generated secret}"
        * This secret will be used to sign and verify JWT tokens. To generate a secret, go into your terminal and run **openssl rand -base64 64**
        * Assign the generated value to 'SECRET'.
    4. POLKA_KEY = "f271c81ff7084ee5b99a5091b42d486e" (This value is used for Polka webhooks, tested by boot.dev)

## Run

1. In terminal, run **go build && ./bootdev-chirpy**.
2. To make requests to this api, you can utilize [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client), [curl](https://curl.se/), your
own built client, or any other client. See the paths below for making requests.

## API Paths

* /app/... 
Description: Retrieves web assets from fileserver.

* GET /admin/metrics 
Description: Retirieves the number of server requests for fileserver assets.

* POST /admin/reset 
Description: Deletes all registered users and their related data.

* GET /api/healthz 
Description: Retrieves a status code 200 if server is working properly.

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
201 Created 
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

* GET /api/chirps
Description: Gets all chirps.
Query Parameters: author_id => gets all chirps for specifc user

* GET /api/chirps/{chirpID}
Description: Get a specific chirp by ID

* POST /api/login
Description: Login with valid credentials.

* POST /api/refresh
Description: Get a new access token using a valid refresh token.

* POST /api/revoke
Description: Revoke a refresh token.

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

* POST /api/polka/webhooks
Description: Webhook that only Polka (payment provider) is authorized to make calls to. Bootdev specific.

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

[Docs](./api_documentation.md)
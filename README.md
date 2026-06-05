# Building Management API

I built this REST API using Go-Fiber and SQLBoiler for my coding assessment. It manages buildings and apartments and stores them in PostgreSQL.

## Architecture & Code Design

- **Dependency Injection**: To avoid using messy global state (which is a bad practice), I passed the database connection pool into the handlers using a dependency environment struct (`Env`). The handlers are defined as methods on this struct, making the database access clean and modular.
- **Upsert Logic**: 
  - For buildings: The API accepts POST requests that can either create or update. It tries to find the record by ID (if provided) or by name (since name is unique). If found, it updates the record. Otherwise, it inserts a new one.
  - For apartments: It looks up by ID or by the combination of building ID and apartment number (since they form a unique key). It verifies that the referenced building exists first before inserting or updating.
- **Relationship Preloading**: The `GET /buildings` endpoint accepts a query parameter `?include_apartments=true` to preload all apartments associated with each building in a single clean query.

## Setup Instructions

### 1. Database Setup
1. Create a database in your local PostgreSQL named `bms`.
2. Run the SQL schema in `schema.sql` against it.
3. Configure your Postgres credentials in `sqlboiler.toml` and `main.go`.

### 2. Generate Models
I used the SQLBoiler CLI tool to generate the database mapping models:
```bash
go install github.com/aarondl/sqlboiler/v4@latest
go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest
sqlboiler psql
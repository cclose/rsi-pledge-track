# rsi-pledge-track
 
 Tracks Star Citizen crowdfunding totals by pulling data from the Roberts Space Industries (RSI) website and storing it in a database, then exposing read access via an HTTP API.
 
 ## What’s in this repo
 
 - **Scraper job** (`cmd/rsiPullFunding`)
   - Calls RSI’s crowd-funding stats endpoint.
   - Writes an (intended) hourly datapoint to Postgres.
 - **Read API** (`cmd/rsiAPI`)
   - Echo-based HTTP API.
   - Provides `GET /pledge-data` for consumers (e.g., Google Sheets Apps Script).
 
 ## Architecture (high level)
 
 ```text
 +-------------------+        +---------------------+
 | EventBridge (cron)| -----> | Lambda: rsiPullFunding|
 +-------------------+        +---------------------+
                                      |
                                      v
                              +----------------+
                              | RDS Postgres   |
                              | table PledgeData|
                              +----------------+
                                      ^
                                      |
 +-------------------+        +---------------------+
 | Google Sheets /   | -----> | API Gateway (public)|
 | other consumers   |        +---------------------+
 +-------------------+                   |
                                         v
                                 +-----------------+
                                 | Lambda: rsiAPI   |
                                 +-----------------+
 ```
 
 Production notes (current state):
 
 - **API is public** (no auth).
 - **RDS is private**; Lambdas connect via **VPC**.
 - Example consumer URL pattern:
   - `https://xxxxxxxxxxxx.execute-api.us-west-2.amazonaws.com/pledge-data`
 
 ## Local development (Docker Compose)
 
 ### Prerequisites
 
 - Docker
 - Docker Compose
 
 ### Configure environment
 
 Create a local `.env` file (it’s gitignored):
 
 ```bash
 DB_USER=rsi
 DB_PASS=change-me
 DB_NAME=rsi_pledge_track
 ```
 
 ### Run
 
 ```bash
 docker compose up --build
 ```
 
 The compose file exposes:
 
 - Postgres on `localhost:5432`
 - API on `http://localhost:8085`
 
 ### Verify
 
 ```bash
 curl "http://localhost:8085/pledge-data?format=json"
 ```
 
 ## Configuration (env vars)
 
 The Go services use these env vars:
 
 - **Database**
   - `DB_HOST`
   - `DB_PORT`
   - `DB_USER`
   - `DB_PASS`
   - `DB_NAME`
 - **API service**
   - `PORT` (required in practice; docker-compose sets this)
   - `START_DELAY` (optional; seconds to sleep before DB connect)
 
 ## HTTP API
 
 ### `GET /pledge-data`
 
 Returns pledge datapoints.
 
 Query parameters:
 
 - `format`
   - `json`, `csv`, `html` (also accepts MIME aliases like `application/json`)
   - Default behavior is HTML-ish (placeholder)
 - `timestamp`
   - RFC3339 timestamp to retrieve a specific entry
 - `startingDateTime`
   - RFC3339 timestamp to retrieve entries after a time
 - `offset`
   - Intended timezone offset from UTC
 - `limit`
   - Intended maximum number of rows
 
 Examples:
 
 ```bash
 curl "http://localhost:8085/pledge-data?format=json"
 curl "http://localhost:8085/pledge-data?format=csv" -o pledgeData.csv
 curl "http://localhost:8085/pledge-data?startingDateTime=2016-12-06T19:09:05Z&format=json"
 ```
 
 Response shape (JSON):
 
 ```json
 [
   {
     "ID": 123,
     "TimeStamp": "2026-03-23 14:00:00",
     "Funding": 800000000,
     "Citizens": 5000000,
     "Fleet": 0
   }
 ]
 ```
 
 ## Scraper job (`rsiPullFunding`)
 
 The scraper calls:
 
 - `https://robertsspaceindustries.com/api/stats/getCrowdfundStats`
 
 It is intended to be run **hourly** and insert a row keyed by an hourly UTC timestamp.
 
 The table has a **unique constraint on `TimeStamp`**, which provides basic protection against inserting duplicates.
 
 ## Database schema
 
 - Postgres schema: `database/schema/postgres.sql`
 - MySQL/SQLite schemas are also included for reference.
 
 ## Repo layout
 
 - `cmd/`
   - `rsiAPI/`: HTTP API entrypoint
   - `rsiPullFunding/`: scraper entrypoint
 - `controller/`: Echo handlers
 - `service/`: DB access/services
 - `database/`: DB connection and schemas
 - `model/`: request/response and data structures
 
 ## Suggested follow-ups (not done in this docs pass)
 
 - **Secrets**
   - Migrate production DB credentials (especially `DB_PASS`) from plain Lambda env vars to **AWS Secrets Manager**.
 - **API correctness**
   - Implement `GetByTimestamp` and `GetAfterTimestamp` in `service/pledge_data.go`.
   - Apply `limit`/`offset` consistently in queries.
 - **TLS**
   - Remove `InsecureSkipVerify` from the scraper HTTP client.
 - **IaC + CI/CD**
   - Add infrastructure-as-code (CDK/Terraform/SAM) and GitHub Actions to auto-deploy from `main`.
 
 ## Additional docs
 
 See `docs/`:
 
 - `docs/architecture.md`
 - `docs/api.md`
 - `docs/runbook.md`
 - `docs/google-sheets.md`
 - `docs/deployment-future.md`

# Architecture

This project consists of two Go executables:

- `rsiPullFunding` (scraper/job)
- `rsiAPI` (read-only HTTP API)

## Components

## Scraper job (`cmd/rsiPullFunding`)

- Pulls Star Citizen crowd-funding totals from RSI.
- Intended to run on an hourly schedule and write one row per hour.
- Inserts into Postgres table `PledgeData`.

## Read API (`cmd/rsiAPI`)

- Echo HTTP server.
- Exposes a single route today:
  - `GET /pledge-data`

## Storage (Postgres)

- Primary schema for Postgres is in `database/schema/postgres.sql`.
- The timestamp column is unique to avoid duplicate datapoints.

## Data flow

```text
RSI website
  POST /api/stats/getCrowdfundStats
        |
        v
Scheduled job (rsiPullFunding)
        |
        v
RDS Postgres (PledgeData)
        ^
        |
Read API (rsiAPI) via API Gateway
        ^
        |
Consumers (Google Sheets, browsers, scripts)
```

## Environments

## Local dev

- `docker-compose.yml` brings up:
  - Postgres container (initialized with `database/schema/postgres.sql`)
  - App container (builds and runs `rsiAPI`)

## Production (current)

- `rsiAPI` runs as **AWS Lambda** behind **API Gateway**.
- `rsiPullFunding` runs as **scheduled Lambda** (EventBridge).
- Postgres is **RDS** and is **private**; Lambdas connect via **VPC**.
- API is **public** (no auth).

## Operational notes

- The unique timestamp constraint means re-running the scraper at the same hourly timestamp should fail fast rather than create duplicates.
- If you need to backfill, you will want a controlled process that:
  - Computes the intended hourly timestamps
  - Retries inserts
  - Detects conflicts on the unique timestamp constraint

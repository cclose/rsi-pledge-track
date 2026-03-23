# Runbook

This document describes operational checks for the scraper and API in production.

## Components

- **Scraper**: scheduled Lambda running `rsiPullFunding`
- **API**: Lambda running `rsiAPI` behind API Gateway
- **Database**: private RDS Postgres (Lambdas connect via VPC)

## Scraper: health checks

- Verify EventBridge rule is enabled and triggering on schedule (hourly).
- Check the scraper Lambda logs for:
  - Non-200 responses from RSI
  - JSON parsing errors
  - DB connectivity issues
  - Unique constraint violations on the timestamp (expected only if re-run)

## API: health checks

- Check API Gateway metrics:
  - 4xx/5xx rates
  - latency
- Check API Lambda logs for:
  - DB connect failures (timeouts, DNS/VPC issues)
  - query errors

## Database checks

- Verify RDS is reachable from Lambda VPC subnets/security groups.
- Verify the `PledgeData` table exists and is receiving new rows.
- Verify storage growth and retention expectations.

## Common failure modes

- **VPC networking**
  - Lambda in private subnets without proper routing (e.g., missing NAT for outbound to RSI) will break scraping.
- **RSI endpoint changes**
  - If RSI response shape changes, JSON unmarshal will fail.
- **DB credential/config drift**
  - DB env vars are configured as plain Lambda environment variables; accidental edits can break connectivity.

## Recommended improvements (follow-up)

- Move DB credentials (especially `DB_PASS`) to AWS Secrets Manager.
- Add alarms:
  - No new rows in `PledgeData` for N hours
  - API 5xx above threshold

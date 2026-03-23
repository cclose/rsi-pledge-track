# API

The API is provided by the `rsiAPI` service (see `cmd/rsiAPI`).

## Endpoint

## `GET /pledge-data`

Returns pledge datapoints from the `PledgeData` table.

### Query parameters

- `format`
  - `json`
  - `csv`
  - `html`
  - Also accepts MIME aliases:
    - `application/json` -> `json`
    - `text/csv` -> `csv`
    - `text/html` -> `html`
- `timestamp`
  - RFC3339 timestamp for a specific entry.
  - Example: `2016-12-06T19:09:05Z`
- `startingDateTime`
  - RFC3339 timestamp to retrieve entries strictly after this time.
- `offset`
  - Intended timezone offset from UTC.
- `limit`
  - Intended max number of rows.

### Response formats

## JSON

- Use `?format=json`.
- Returns a JSON array of objects shaped like `model.PledgeData`.

## CSV

- Use `?format=csv`.
- Returns a CSV attachment named `pledgeData.csv`.

## HTML

- Use `?format=html`.
- Current HTML response is a placeholder.

### Examples

```bash
curl "http://localhost:8085/pledge-data?format=json"
curl "http://localhost:8085/pledge-data?format=csv" -o pledgeData.csv
curl "http://localhost:8085/pledge-data?timestamp=2016-12-06T19:09:05Z&format=json"
curl "http://localhost:8085/pledge-data?startingDateTime=2016-12-06T19:09:05Z&format=json"
```

### Production URL pattern

The current public deployment uses an `execute-api` hostname without a visible stage prefix:

- `https://xxxxxxxxxxxx.execute-api.us-west-2.amazonaws.com/pledge-data`

(Use the real API id in place of `xxxxxxxxxxxx`.)

### Notes / known gaps

- The `GetByTimestamp` and `GetAfterTimestamp` service methods are currently stubs.
- `limit` and `offset` are parsed but are not consistently applied in storage queries.

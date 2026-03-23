# Google Sheets consumer notes

This repo does not include the Google Sheets Apps Script, but the API is designed to be easy to call from Apps Script.

## Recommended pattern

- Prefer `format=json` for Apps Script.
- Cache results in the sheet (or in PropertiesService) to avoid repeatedly calling the API.
- Add basic backoff/retry if you hit transient errors.

## Example (Apps Script)

```javascript
function fetchPledgeData() {
  var url = 'https://xxxxxxxxxxxx.execute-api.us-west-2.amazonaws.com/pledge-data?format=json';
  var resp = UrlFetchApp.fetch(url, { muteHttpExceptions: true });

  if (resp.getResponseCode() !== 200) {
    throw new Error('API error ' + resp.getResponseCode() + ': ' + resp.getContentText());
  }

  var data = JSON.parse(resp.getContentText());
  return data;
}
```

Replace `xxxxxxxxxxxx` with the real API id.

## CSV

If you want a downloadable CSV for manual inspection:

- `https://xxxxxxxxxxxx.execute-api.us-west-2.amazonaws.com/pledge-data?format=csv`

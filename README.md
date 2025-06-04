# Email Tracking Service

This is a minimal Go service that generates tracking links with UTM parameters and records when recipients visit your site.

## Usage

Run the server:

```bash
go run ./cmd
```

### Generating a tracking link

```
GET /generate?email=user@example.com&campaign=newsletter
```

The service returns a JSON object with a `url` containing UTM parameters and a unique ID.

### Tracking visits

Users visiting the generated link hit `/track` with the unique ID. The service marks the click and redirects to `https://example.com/`.

## Disclaimer

This is a simple example using an in-memory store. Data will be lost when the application stops.

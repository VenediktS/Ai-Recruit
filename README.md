# Email Tracking Service

This is a minimal Go service that generates tracking links with UTM parameters and records when recipients visit your site.

## Usage

Run the server with a Postgres database available. Configure the connection
string in `config.json`. By default this repository includes a configuration
that points to a local instance:

```json
{
  "database_url": "postgres://user:pass@localhost:5432/dbname"
}
```

If you want to provide a custom configuration file, set the `CONFIG_PATH`
environment variable to its location.

Run the migration:

```bash
psql "postgres://user:pass@localhost:5432/dbname" -f migrations/0001_create_trackings.sql
```

Then start the server:

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

This example now includes a very small admin interface protected with HTTP
Basic authentication. The credentials are configured in `config.json` as
`admin_username` and `admin_password` (defaults are `admin`/`1234`). The API
endpoints used for generating and tracking links remain publicly accessible.

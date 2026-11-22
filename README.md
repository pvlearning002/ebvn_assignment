# Bookmark Management API

A Go application that provides APIs to shorten URLs, redirect to their original destinations, and generate random passwords. Shortened links are stored in Redis, and the HTTP API is built with Gin.

## Features

- **Password generation:** creates a random hexadecimal string.
- **URL shortening:** generates a short code for a valid URL and stores the code/URL pair in Redis.
- **Redirection:** looks up a URL by code and responds with `301 Moved Permanently`.
- **Health check:** returns service information, an instance ID, and the Redis container status.
- **Swagger UI:** provides interactive API documentation.

## API

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/genpass` | Generate a password. |
| `GET` | `/health` | Get the health check response. |
| `POST` | `/v1/links/shorten` | Create a shortened link. |
| `GET` | `/v1/links/redirect/:code` | Redirect to the original URL. |
| `GET` | `/swagger/index.html` | Open Swagger UI. |

### Shorten a URL

The request body must include `expiration_time` of at least `500` and a valid URL:

```json
{
  "expiration_time": 3600,
  "url": "https://example.com/article"
}
```

A successful response looks like this:

```json
{
  "code": "a1b2c3d",
  "message": "Shorten URL generated successfully"
}
```

Request `GET /v1/links/redirect/a1b2c3d` to be redirected to the original URL.

### Generate a password

```http
GET /genpass
```

The response is JSON. By default, the password is returned in the `Password` field.

### Health check

```http
GET /health
```

The response includes a message, service name, instance UUID, and the `is_redis_connected` flag.

## Run the application

You need the Go version specified in `go.mod` and Redis. Start Redis locally with Docker Compose:

```sh
docker compose up -d redis
```

Then start the API:

```sh
go run ./command/api
```

By default, the API listens on port `8080` and Redis is expected at `localhost:6379`. To build and run the full stack with Docker:

```sh
docker build -t ebvn_test:latest .
docker compose up
```

The Compose configuration exposes Redis at `redis:6379` and the API at port `8080`.

## Configuration

Default values are defined in `internal/config/configs.go`:

| Setting | Default | Description |
| --- | --- | --- |
| `APP_PORT` | `8080` | API HTTP port. |
| `PASSWORD_LENGTH_DEFAULT` | `12` | Default generated password length. |
| `KEYCODE_LENGTH_DEFAULT` | `7` | Default shortened-link code length. |
| `REDIS_EXPIRED_TIME` | `3600` | Link lifetime in Redis, in seconds. |
| `REDIS_ADDRESS` | `localhost:6379` | Redis address. |
| `REDIS_PASSWORD` | empty | Redis password. |
| `REDIS_DB` | `0` | Redis database index. |
| `LOG_LEVEL` | `info` | Logging level. |
| `INSTANCE_ID` | generated when empty | UUID identifying this instance. |

Configuration is loaded through `envconfig` with a prefix for each configuration group. For example, the `Config` group uses variables such as `CONFIG_APP_PORT` and `CONFIG_REDIS_EXPIRED_TIME`. The Redis configuration currently uses the `CONFIG_PARAMETER_` prefix (for example, `CONFIG_PARAMETER_REDIS_ADDRESS`). The repository's `.env` file shows deployment values; ensure variable names match the prefixes read by the application.

## Tests

Run all unit and integration tests with:

```sh
go test ./...
```

Integration tests use Miniredis in memory and do not require a separate Redis server.

## Notes

- The request's `expiration_time` is required and validated to be at least `500`, but the actual Redis link lifetime is controlled by the application's `REDIS_EXPIRED_TIME` setting.
- The health check currently runs `docker inspect` for a container named `redis` to determine the container status. It may report `false` if Docker is unavailable or the container has a different name, even when Redis is reachable at another address.
- The API currently has no authentication or authorization. Add an appropriate protection layer before exposing it publicly.

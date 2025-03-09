# Shortener Service

This service generates short URLs and stores them in PostgreSQL.

## Running the service

```bash
go run services/shortener/main.go
```

## Testing the service

```bash
curl -X POST "http://localhost:8080/shorten" \
  -H "Content-Type: application/json" \
  -d '{"long_url":"https://example.com"}'
```

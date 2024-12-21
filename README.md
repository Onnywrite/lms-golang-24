# LMS Golang 24

Yandex lyceum project of the Golang course 2024.

# How to run

```bash
export SERVER_PORT=80
go run cmd/server/main.go
```

If you run `go run cmd/server/main.go` then it will work on port 8080.

# Use cases

Code 200, result -498:
```bash
curl -v --location localhost/api/v1/calculate \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2*(-5)^3"
}'
```

Code 422:
```bash
curl -v --location localhost/api/v1/calculate \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2*(-5)^"
}'
```

Code 422:
```bash
curl -v --location localhost/api/v1/calculate \
--header 'Content-Type: application/json' \
--data '{
  "expression": "wrong json"
'
```

Code 500, panic:
```bash
curl -v -X POST --location localhost/api/v1/panic
```
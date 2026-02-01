# go-job-queue

Production-oriented Golang job queue API scaffold using Gin.

## Run locally

### With Go

```bash
make run
```

### With Docker Compose

```bash
docker compose up --build
```

## 開始開發

### 本機開發流程

```bash
make tidy
make test
make run
```

### 常用指令

- `make run`：啟動 API（預設在 8080）
- `make test`：執行所有測試
- `make build`：產出二進位檔到 `bin/go-job-queue`

## Endpoints

- GET /healthz -> {"status":"ok"}
- GET / -> {"message":"go-job-queue api is running"}

## Configuration

- PORT (default: 8080)

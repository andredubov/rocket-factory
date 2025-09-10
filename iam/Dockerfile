FROM golang:1.24.4-alpine AS builder

COPY .. /github.com/andredubov/rocket-factory
WORKDIR /github.com/andredubov/rocket-factory/iam

RUN go mod download && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/iam ./cmd/main.go

FROM alpine:3.22.1

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
USER appuser

COPY --from=builder /github.com/andredubov/rocket-factory/bin/iam .
COPY --from=builder /github.com/andredubov/rocket-factory/iam/migrations/*.sql ./migrations/

CMD ["./iam"]
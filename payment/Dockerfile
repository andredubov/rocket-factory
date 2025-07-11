FROM golang:1.24.4-alpine AS builder

COPY .. /github.com/andredubov/rocket-factory
WORKDIR /github.com/andredubov/rocket-factory/payment

RUN go mod download && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/payment ./cmd/main.go

FROM alpine:3.20

WORKDIR /root/
COPY --from=builder /github.com/andredubov/rocket-factory/bin/payment .

CMD ["./payment"]
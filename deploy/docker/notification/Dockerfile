# ============================
# Stage 1: Build stage
# ============================

# Используем официальный образ Go на базе Alpine — лёгкий, быстрый и безопасный
FROM golang:1.24.4-alpine AS builder

# Устанавливаем git — нужен для загрузки зависимостей из приватных и публичных репозиториев
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.work и его контрольную сумму — это позволяет Go правильно связать модули
COPY go.work go.work.sum ./

# Копируем go.mod и go.sum всех зависимых модулей — нужно для кэширования go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/
COPY assembly/go.mod assembly/go.sum ./assembly/
COPY notification/go.mod notification/go.sum ./notification/
COPY iam/go.mod notification/go.sum ./iam/

# Загружаем все зависимости, указанные в модульных файлах
RUN go mod download

# Копируем исходный код всех модулей — включая зависимости, т.к. они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY inventory ./inventory
COPY order ./order
COPY payment ./payment
COPY assembly ./assembly
COPY notification ./notification
COPY iam ./iam

# Собираем бинарный файл Notification-сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notification-app ./notification/cmd/main.go


# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.22.1

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Запускаем приложение под non-root пользователем
USER appuser

# Копируем сертификаты из builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/notification-app .

# Копируем .env файл для конфигурации приложения
COPY deploy/compose/notification/.env deploy/compose/notification/.env

# Устанавливаем команду запуска — запускаем наш бинарь
ENTRYPOINT ["./notification-app"]
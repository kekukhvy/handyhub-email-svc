FROM golang:latest AS builder

WORKDIR /app

# Копирование всего проекта
COPY . /app

# Полный обход всех проверок SSL и сборка
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV GOPRIVATE=*
ENV GOINSECURE=*
RUN git config --global http.sslverify false

# Попробуем обновить модули и собрать
RUN go mod tidy -e || true
RUN go build -o /handyhub-email-svc -ldflags "-w -s" cmd/main.go

# Финальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Копирование бинарного файла
COPY --from=builder /handyhub-email-svc .

# Копирование конфигурационных файлов
COPY --from=builder /app/internal/config /app/internal/config

# Создание директории для логов
RUN mkdir -p logs

# Открытие порта
EXPOSE 8008

# Переменные окружения
ENV MONGODB_URL=""
ENV DB_NAME="handyhub"
ENV RABBITMQ_URL=""

# Запуск приложения
CMD ["./handyhub-email-svc"]
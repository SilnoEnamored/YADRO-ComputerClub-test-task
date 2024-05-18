# syntax=docker/dockerfile:1

# Используем официальный образ Golang
FROM golang:1.22

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod файл для установки зависимостей
COPY go.mod ./

# Загружаем зависимости, если есть
RUN go mod download || true

# Копируем остальные исходные файлы
COPY . .

# Компилируем приложение
RUN go build -o task.exe cmd/main.go

# Определяем команду для запуска программы
ENTRYPOINT ["./task.exe"]

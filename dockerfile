
FROM golang:1.23.3

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/

COPY backend/ ./backend/

COPY frontend ./frontend/

WORKDIR /app/backend

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /app.bin

CMD ["/app/backend/app.bin"]
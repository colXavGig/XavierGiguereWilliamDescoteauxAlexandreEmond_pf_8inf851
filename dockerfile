
FROM golang:1.23.3

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/

RUN go mod download

COPY backend/ ./backend/

COPY frontend ./frontend/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app.bin

CMD ["/app/app.bin"]
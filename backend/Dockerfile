FROM golang:1.24-alpine
WORKDIR /app

# Install air for hot reloading and goose for migrations
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
CMD ["air"]

FROM golang:1.25-alpine

WORKDIR /app

# Install only necessary dependencies
RUN apk add --no-cache git curl bash gpgme-dev libassuan-dev libgpg-error-dev build-base btrfs-progs-dev

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Expose the API port
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

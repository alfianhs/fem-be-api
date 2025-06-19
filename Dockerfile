# Stage 1 Build
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install git for go get and build
RUN apk add --no-cache git

# Dependencies caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Minimal Runtime
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=build /app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]

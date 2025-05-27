# Stage 1 Build
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install git for go get and build deps
RUN apk add --no-cache git

# Dependencies Caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app .
# Stage 2: Minimal Runtime
FROM alpine:latest

WORKDIR /app

# Only copy binary
COPY --from=build /app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]

# Build Go BInary
# RUN go build -o app .

# # Stage 2. Run Coyy
# FROM golang:1.24-alpine

# WORKDIR /app

# # Copy from builder stage 1
# COPY --from=build /app/app .

# EXPOSE 8080

# ENTRYPOINT [ "./app" ]

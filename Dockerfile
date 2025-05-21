# Stage 1 Build
FROM golang:1.24-alpine AS build

WORKDIR /app

# Dependencies Caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Go BInary
RUN go build -o app .

# Stage 2. Run Coyy
FROM golang:1.24-alpine

WORKDIR /app

# Copy from builder stage 1
COPY --from=build /app/app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]

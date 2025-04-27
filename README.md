# \ud83d\udce6 Futsal Event Management

> A backend service developed with Go (GIN framework) and MongoDB to manage futsal events and facilitate online ticket sales.

---

## \ud83d\ude80 Tech Used

- [Go (Golang)](https://golang.org/) \u2014 Backend programming language
- [MongoDB](https://www.mongodb.com/) \u2014 NoSQL Database
- [GIN](https://github.com/gin-gonic/gin) — Web framework for building APIs in Go
- **Modified Clean Architecture** — Software architecture pattern that separates concerns into layers (Delivery, Usecase, Repository, Domain)

---

## \ud83d\udcc2 Directory Structure

```
├── app/
│   ├── delivery/
│   │   └── http/          # Handler HTTP (controllers)
│   ├── repository/        # Repository implementation
│   └── usecase/           # Business logic / service layer
├── docs/                  # API documentation with swagger
├── domain/
│   ├── model/             # Entity / model data
│   └── request/           # Request struct
├── helpers/               # Utility
├── go.mod                 # File dependency Go
├── main.go                # Entry point
├── README.md              # Project Documentation
```

---

## \u2699\ufe0f Installation and Setup

1. **Clone this repo**
   ```bash
   git clone https://github.com/alfianhs/fem-be-api.git
   cd fem-be-api
   ```

2. **Copy file `.env`**
   ```bash
   cp .env.example .env
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run app**
   ```bash
   go run main.go
   ```

---

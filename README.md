# Belalai E-Wallet Backend

![badge golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![badge postgresql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![badge redis](https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white)

<img src="/public/belalai-wallet.png"  align="right" width="370px">

Welcome to Belalai E-Wallet! 🐘 The digital wallet application designed to give you fast and easy access to your money, anytime, anywhere. Inspired by the elephant's trunk—a versatile, multi-functional tool—Belalai E-Wallet offers the same effortless reach for all your financial transactions, from bill payments to fund transfers. We’ve built a nimble and robust platform, ensuring your e-wallet experience is as swift and reliable as the movement of a trunk. This project is backend for [Belalai E-Wallet Frontend](https://github.com/FebryanHernanda/Belalai-E-Wallet-Frontend) web application build gin gonic as framework for HTTP API, Go language, PostgreSQL as database, and redis as cache sistem.

## 🔧 Tech Stack

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/install-redis-on-windows/)
- [JWT](https://github.com/golang-jwt/jwt)
- [argon2](https://pkg.go.dev/golang.org/x/crypto/argon2)
- [migrate](https://github.com/golang-migrate/migrate)
- [Docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- [Swagger for API docs](https://swagger.io/) + [Swaggo](https://github.com/swaggo/swag)

## 🗝️ Environment

```bash
# database
DBUSER=<your_database_user>
DBPASS=<your_database_password>
DBNAME=<your_database_name
DBHOST=<your_database_host>
DBPORT=<your_database_port>

# JWT hash
JWT_SECRET=<your_secret_jwt>
JWT_ISSUER=<your_jwt_issuer>

# Redish
RDB_HOST=<your_redis_host>
RDB_PORT=<your_redis_port>
RDB_USER=<your_redis_user>
RDB_PWD=<your_redis_password>

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your_email>
SMTP_PASS=<your_app_password_email>
SMTP_FROM="<aplication-name> <your_email>" # wtih " "
FRONTEND_URL=<your_fronend_url>
```

## ⚙️ Installation

1. Clone the project

```sh
$ https://github.com/FebryanHernanda/Belalai-E-Wallet-Backend.git
```

2. Navigate to project directory

```sh
$ cd Belalai-E-Wallet-Backend
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Setup your [environment](##-environment)

5. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) for DB migration

6. Do the DB Migration

```sh
$ migrate -database YOUR_DATABASE_URL -path ./db/migrations up
```

or if u install Makefile run command

```sh
$ make migrate-createUp
```

7. Run the project

```sh
$ go run ./cmd/main.go
```

## 🚧 API Documentation

| Method | Endpoint                 | Body                                                           | Description                            |
| ------ | ------------------------ | -------------------------------------------------------------- | -------------------------------------- |
| GET    | /img                     |                                                                | Static File                            |
| POST   | /auth                    | email:string, password:string                                  | Login                                  |
| POST   | /auth/register           | email:string, password:string                                  | Register                               |
| DELETE | /auth                    | header: Authorization (token jwt)                              | Logout                                 |
| PATCH  | /auth/update-pin         | header: Authorization (token jwt), body                        | create pin new user                    |
| PATCH  | /auth/change-pin         | header: Authorization (token jwt), body                        | change pin registered user             |
| PATCH  | /auth/change-password    | header: Authorization (token jwt), body                        | change password regitered user         |
| POST   | /auth/forgot-password    | body                                                           | change password with SMTP              |
| POST   | /auth/reset-password     | email:string                                                   | reset password with email              |
| POST   | /auth/forgot-pin         |                                                                | change pin with SMTP                   |
| POST   | /auth/reset-pin          |                                                                | reset pin with email                   |
| POST   | /auth/confirm-pin        |                                                                | verify any transaction                 |
| GET    | /profile                 | header: Authorization (token jwt)                              | get user data                          |
| PATCH  | /profile                 | header: Authorization (token jwt), body                        | update user data                       |
| DELETE | /profile/avatar          | header: Authorization (token jwt)                              | delete user avatar                     |
| GET    | /balance                 | header: Authorization (token jwt)                              | get wallet data a user                 |
| GET    | /chart/:duration         | header: Authorization (token jwt), duration: string            | get statistic data a user              |
| GET    | /transaction/history     | header: Authorization (token jwt)                              | get transaction hsitories data a user  |
| GET    | /transaction/history/all | header: Authorization (token jwt)                              |                                        |
| DELETE | /transaction/:id         | header: Authorization (token jwt), id : integer                | soft delete history transaction        |
| GET    | /transfer                | header: Authorization (token jwt), page:integer, search:string | filter/search user before transfer     |
| POST   | /transfer                | header: Authorization (token jwt), body                        | transfer balance from a user to a user |
| GET    | /topup/method            | header: Authorization (token jwt)                              | get all payment method for top up      |
| POST   | /topup/                  | header: Authorization (token jwt), body                        | Topup wallet a user                    |

## 📄 LICENSE

MIT License

Copyright (c) 2025 Belalai team

## 📧 Contact Info & Contributor

[https://github.com/Darari17](https://github.com/Darari17)

[https://github.com/raihaninkam](https://github.com/raihaninkam)

[https://github.com/M16Yusuf](https://github.com/M16Yusuf)

[https://github.com/FebryanHernanda](https://github.com/FebryanHernanda)

[https://github.com/federus1105](https://github.com/federus1105)

[https://github.com/habibmrizki](https://github.com/habibmrizki)

## 🎯 Related Project

[https://github.com/FebryanHernanda/Belalai-E-Wallet-Frontend](https://github.com/FebryanHernanda/Belalai-E-Wallet-Frontend)

[https://github.com/FebryanHernanda/Belalai-E-Wallet-Backend](https://github.com/FebryanHernanda/Belalai-E-Wallet-Backend)

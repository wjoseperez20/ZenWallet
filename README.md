# ZenWallet

## Overview

This app represents a tech challenge from Storicard. It is a RESTful API for account management developed in Go, incorporating features such as JWT Authentication, rate limiting, Swagger documentation, and database operations through GORM. The application utilizes the Gin Gonic web framework and is containerized using Docker.

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [dbmate](https://github.com/amacneil/dbmate#installation)
- [Docker + Docker compose](https://docs.docker.com/engine/install/)

### Installation

1. Clone the repository

```bash
git clone https://github.com/wjoseperez20/zenwallet.git
```

2. Navigate to the directory

```bash
cd zenwallet
```

3. Build and run the Docker containers

```bash
make setup && make build && make up
```

### Environment Variables

You can set the environment variables in the `.env` file. Here are some important variables:

- `GMAIL_USER`
- `GMAIL_SECRET`
- `POSTGRES_HOST`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `JWT_SECRET`
- `API_SECRET_KEY`
- `AWS_REGION`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### API Documentation

The API is documented using Swagger and can be accessed at:

```
http://localhost:8001/swagger/index.html
```

## Usage

### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

```bash
curl -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8001/api/v1/accounts
```
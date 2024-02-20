# Go POS

## Description

A simple RESTful Point of Sale (POS) web service written in Go programming language. This project is a part of my learning process in understanding [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) in Go.

It uses [Gin](https://gin-gonic.com/) as the HTTP framework and [PostgreSQL](https://www.postgresql.org/) as the database with [pgx](https://github.com/jackc/pgx/) as the driver and [Squirrel](https://github.com/Masterminds/squirrel/) as the query builder. It also utilizes [Redis](https://redis.io/) as the caching layer with [go-redis](https://github.com/redis/go-redis/) as the client.

This project idea was inspired by the [Ide Project untuk Upgrade Portfolio Backend Engineer](https://www.youtube.com/watch?v=uAR1kjyeDtg) video on YouTube by [Asdita Prasetya](https://www.youtube.com/@asditaprasetya), which provided valuable guidance and inspiration for its development.

## Setting up the development environment

Before you start contributing to the project, you will need to set up your development environment. To get started, you should have the following software installed:

- [Go](https://golang.org/) 1.21 or higher
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [PostgreSQL](https://hub.docker.com/_/postgres) container
- [Redis](https://hub.docker.com/_/redis) container

## Getting Started

1. If you do not use devcontainer, ensure you have [Go](https://go.dev/dl/) 1.21 or higher and [Task](https://taskfile.dev/installation/) installed on your machine:

    ```bash
    go version && task --version
    ```

2. Create a copy of the `.env.example` file and rename it to `.env`:

    ```bash
    cp .env.example .env
    ```

    Update configuration values as needed.

3. Install all dependencies, run docker compose, create database schema, and run database migrations:

    ```bash
    task
    ```

4. Run the project in development mode:

    ```bash
    task dev
    ```

## Documentation

For database schema documentation, see [here](https://dbdocs.io/bagashiz/Go-POS/), powered by [dbdocs.io](https://dbdocs.io/).

API documentation can be found in `docs/` directory. To view the documentation, open the browser and go to `http://localhost:8080/docs/index.html`. The documentation is generated using [swaggo](https://github.com/swaggo/swag/) with [gin-swagger](https://github.com/swaggo/gin-swagger/) middleware.

## Learning References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn
- [Ready for changes with Hexagonal Architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749) by Netflix Technology Blog
- [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3) by Matias Varela

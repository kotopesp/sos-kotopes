# SOS-Kotopes Web Application Backend

This repository contains the backend for the SOS-Kotopes web application.
The backend is built using Go and provides RESTful APIs to interact with the system.

## To-Do List

The following tasks are pending to complete the backend development:

- [ ] **Logger**
    - Replace the default Fiber logger with a custom implementation.

- [ ] **CI (Continuous Integration)**
    - Set up a CI pipeline to automate testing, building, and deployment.

- [ ] **Docker Compose**
    - `docker-compose.yml` file to orchestrate the application and its dependencies using Docker.

- [ ] **License**
    - LICENSE file to specify the licensing terms for the project.

- [ ] **README**
    - Write a comprehensive README file.

- [ ] **Automatic Database Migrations**
    - Implement automatic database migrations.
  
- [ ] **Swagger Documentation**
  - Select library for automatic API Swagger generation, add annotations for handlers

## Getting Started

### Prerequisites

- Go 1.22.5+
- PostgreSQL
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:

    ```sh
    git clone https://gitflic.ru/project/spbu-se/sos-kotopes
    cd sos-kotopes/backend
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Running the Application

To run the application:

```sh
make run
```

*Note*: To run command `PG_URL` needs to be updated in the Makefile.

### Running with Docker Compose

Fill the `docker-compose.yml` file (not provided yet) and run:

```sh
docker-compose up --build
```

### Database Migrations

To run database migrations manually:

```sh
make migrate-up
```

*Note*: To run command `PG_URL` needs to be updated in the Makefile.

## Development Practices

### Adding Dependencies

When adding new dependencies to the project:

1. Add the dependency to your `go.mod` file using `go get` or directly editing `go.mod`.
2. Run `go mod vendor` .

### Code Quality Checks

To maintain code quality, follow these practices:

- Run `make lint` before committing changes or submitting a pull request to ensure your code meets the linting standards.


## Acknowledgments

- [Fiber](https://github.com/gofiber/fiber/v2)
- [GORM](https://gorm.io/)
- [zerolog](https://github.com/rs/zerolog)
- [PostgreSQL](https://www.postgresql.org/)


# XYZ Skilltest - By Kholidss

## Getting started

This boilerplate is built on top of [Go Fiber](https://docs.gofiber.io) Golang Framework.

## Dependencies

There is some dependencies that we used in this boilerplate:
- [Go Fiber](https://docs.gofiber.io/) [Go Framework]
- [Viper](https://github.com/spf13/viper) [Go Configuration]
- [Cobra](https://github.com/spf13/cobra) [Go Modern CLI]
- [Logrus Logger](https://github.com/sirupsen/logrus) [Go Logger]
- [Goose Migration](https://github.com/pressly/goose) [Go Migration]
- [Gobreaker](https://github.com/sony/gobreaker) [Go Circuit Breaker]
- [OpenTelemetry](https://pkg.go.dev/go.opentelemetry.io/otel) [OpenTelemetry Tracer]

## Requirement

- Golang version 1.21 or latest
- Database MySQL
- RabbitMQ

## Usage

### Config Initialization
Copy .env.example to .env
```bash
cp .env.example .env
```

### Installation
install required dependencies
```bash
make install
```

### Run HTTP Service
run current http service after all dependencies installed
```bash
make start-http
```

## Database Migration
migration status
```bash
make migration-status
```

create migration table
```bash
make migration-create name={migration_name}

# example
make migration-create name=create_example_table
```

migration up
```bash
make migration-up
```

migration down
```bash
make migration-down
```

to show all migration command
```bash
make migration
```

# Order System – Building Cloud-Native Systems in Go

## Overview

**Order System** is a learning-oriented project that demonstrates how to design and build a modern cloud-native backend using Go.

Instead of introducing technologies independently, every new component is added only when the architecture requires it. The project starts as a simple CRUD application and gradually evolves into a distributed event-driven system composed of multiple services communicating through gRPC and Redis Streams.

The objective of this repository is not simply to build an Order Management System, but to understand **why modern backend architectures are designed the way they are.**

---

# Architecture

```text
                    Client
                       │
                     gRPC
                       │
                       ▼
                Order Service
                       │
          PostgreSQL (ordersdb)
                │         │
                ▼         ▼
             orders    outbox_events
                       │
                       ▼
                   Publisher
                       │
                       ▼
             Redis Stream (stream:orders)
                  ┌──────────────┐
                  ▼              ▼
          Notification      Inventory
                                │
                     PostgreSQL (inventorydb)
                          │              │
                          ▼              ▼
                      orders       outbox_events
                                        │
                                        ▼
                                   Publisher
                                        │
                                        ▼
                          Redis Stream (stream:inventory)
                                        │
                                        ▼
                                    Shipping
```

---

# Features

- gRPC API
- PostgreSQL persistence
- Repository Pattern
- Dependency Injection
- Transport Adapter Pattern
- Outbox Pattern
- Redis Streams
- Consumer Groups
- Event-driven communication
- CQRS-style read model
- Docker multi-stage builds
- Docker Compose orchestration
- Health checks
- Environment-based configuration

---

# Technology Stack

| Technology | Purpose |
|------------|---------|
| Go | Backend services |
| PostgreSQL | Persistent storage |
| pgx | PostgreSQL driver |
| Redis Streams | Event streaming |
| gRPC | Service communication |
| Protocol Buffers | API contracts |
| Docker | Containerization |
| Docker Compose | Local orchestration |

---

# Project Structure

```text
cmd/
├── client/
├── server/
├── publisher/
├── notification/
├── inventory/
└── shipping/

grpc/
├── order.proto
├── order.pb.go
└── order_grpc.pb.go

internal/
├── config/
├── domain/
├── inventory/
├── notification/
├── publisher/
├── redis/
├── repository/
├── service/
└── transport/
```

---

# Folder Responsibilities

## cmd/

Contains all runnable applications.

| Folder | Responsibility |
|---------|----------------|
| server | Order gRPC Server |
| client | Sample gRPC Client |
| publisher | Publishes Outbox events |
| notification | Consumes Order events |
| inventory | Reserves inventory |
| shipping | Consumes Inventory events |

---

## internal/domain

Contains business entities and domain concepts.

Examples:

- Order
- OutboxEvent
- InventoryEvent
- Domain Errors

The domain layer knows nothing about:

- PostgreSQL
- Redis
- gRPC
- Protocol Buffers

---

## internal/service

Contains business logic.

Responsible for:

- Validation
- Business rules
- Event creation
- Orchestration

The service layer never performs SQL or transport-specific operations.

---

## internal/repository

Persistence layer.

Responsible for:

- SQL
- Transactions
- PostgreSQL
- Outbox persistence

Repositories expose business operations rather than CRUD methods.

Examples:

- `CreateOrder()`
- `ReserveInventory()`

instead of generic CRUD methods.

---

## internal/transport

Transport adapters.

Responsible for converting:

```text
protobuf
    ↓
domain
```

and

```text
domain
    ↓
protobuf
```

The business layer remains transport-independent.

---

## internal/publisher

Generic Outbox Publisher.

Responsibilities:

- Poll pending events
- Publish to Redis Streams
- Mark events as published

The same publisher implementation is reused by multiple services.

---

## internal/redis

Redis abstraction.

Provides:

- Stream publishing
- Consumer Groups
- Message acknowledgement
- Group creation

---

## internal/config

Application configuration.

Loads configuration from environment variables.

Supports both:

- Local development
- Docker Compose
- Kubernetes

without changing application code.

---

# Event Flow

## Order Creation

```text
Client
    │
    ▼
gRPC
    │
    ▼
Order Service
    │
    ▼
BEGIN
    │
    ▼
Insert Order
    │
    ▼
Insert Outbox Event
    │
    ▼
COMMIT
```

---

## Event Publication

```text
Outbox
   │
   ▼
Publisher
   │
   ▼
Redis Stream
```

---

## Inventory Processing

```text
OrderCreated
      │
      ▼
Inventory
      │
      ▼
Reserve Stock
      │
      ▼
Update Projection
      │
      ▼
Insert Inventory Outbox
      │
      ▼
COMMIT
```

---

## Shipping

```text
InventoryReserved
        │
        ▼
Shipping
```

---

# Running the Project

## Using Docker Compose (Recommended)

```bash
docker compose up --build
```

This command:

- Builds all application images
- Starts PostgreSQL
- Starts Redis
- Initializes databases
- Starts Order Service
- Starts Publisher
- Starts Notification
- Starts Inventory
- Starts Shipping

---

## Local Development

Start infrastructure:

```bash
docker compose up postgres redis
```

Run services individually:

```bash
go run ./cmd/server
go run ./cmd/publisher
go run ./cmd/notification
go run ./cmd/inventory
go run ./cmd/shipping
```

Run the client:

```bash
go run ./cmd/client
```

---

# Configuration

The application is configured through environment variables.

| Variable | Description |
|----------|-------------|
| DB_HOST | PostgreSQL host |
| DB_PORT | PostgreSQL port |
| DB_USER | PostgreSQL user |
| DB_PASSWORD | PostgreSQL password |
| DB_NAME | Database name |
| REDIS_ADDR | Redis server |
| ORDER_SERVER_ADDR | Order Service gRPC endpoint |

The same binaries work:

- Local development
- Docker Compose
- Kubernetes

Only the environment changes.

---

# Design Patterns

## Repository Pattern

Repositories expose business operations instead of CRUD methods.

---

## Dependency Injection

Services depend on abstractions rather than concrete implementations.

---

## Transport Adapter

Transport-specific models never leak into the business layer.

---

## Outbox Pattern

Business data and events are written within the same transaction, ensuring reliable event publication.

---

## CQRS

Inventory maintains its own read model rather than querying the Order Service directly.

---

## Event-Driven Architecture

Services communicate asynchronously using Redis Streams.

---

## Consumer Groups

Multiple services independently consume the same stream while maintaining their own processing state.

---

# Docker Concepts

This project demonstrates several Docker concepts.

## Multi-stage Builds

A single reusable Dockerfile builds every service using a build argument.

```dockerfile
ARG SERVICE
```

---

## Distroless Runtime Images

The runtime image contains only the compiled Go binary.

No compiler.

No source code.

No package manager.

---

## Docker Compose

Compose orchestrates:

- PostgreSQL
- Redis
- Order Service
- Publisher
- Notification
- Inventory
- Shipping

---

## Service Discovery

Containers communicate using Docker service names.

Examples:

- `postgres`
- `redis`
- `order-server`

instead of `localhost`.

---

## Health Checks

Docker waits until PostgreSQL and Redis are healthy before starting dependent services.

---

## Environment-based Configuration

Applications load configuration from environment variables, allowing the same binaries to run in different environments without modification.

---

# Learning Objectives

This repository demonstrates why modern backend systems adopt layered architecture and asynchronous communication.

Topics covered include:

- Project organization
- Dependency Injection
- Repository Pattern
- PostgreSQL
- Transactions
- gRPC
- Protocol Buffers
- Transport Adapters
- Outbox Pattern
- Redis Streams
- Consumer Groups
- CQRS
- Event-Driven Architecture
- Docker
- Docker Compose
- Health Checks
- Service Discovery

Each technology is introduced to solve a concrete architectural problem rather than as an isolated concept.

---

# Roadmap

## Completed

- ✅ Go Project Structure
- ✅ Dependency Injection
- ✅ Repository Pattern
- ✅ PostgreSQL
- ✅ Docker
- ✅ Docker Compose
- ✅ gRPC
- ✅ Protocol Buffers
- ✅ Transport Adapters
- ✅ Outbox Pattern
- ✅ Redis Streams
- ✅ Consumer Groups
- ✅ Event-Driven Architecture
- ✅ CQRS
- ✅ Notification Service
- ✅ Inventory Service
- ✅ Shipping Service

## Next

- ⬜ Kubernetes
- ⬜ Helm
- ⬜ Prometheus
- ⬜ Grafana
- ⬜ OpenTelemetry
- ⬜ Skaffold
- ⬜ Bazel

---

# Philosophy

The purpose of this project is not to build an Order Management System.

The purpose is to understand the architectural reasoning behind modern distributed systems.

Every technology introduced in this repository solves a real architectural problem created by the previous stage of the application. The result is a gradual evolution from a simple CRUD application into a cloud-native, event-driven backend built with production-oriented design principles.
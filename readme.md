# Order System – Building Cloud-Native Systems in Go

## Overview

This project is a learning-oriented implementation of a modern backend system written in Go.

Rather than focusing only on language features, the project gradually evolves from a simple CRUD application into a cloud-native, event-driven architecture. Every technology is introduced only when the architecture requires it.

The project currently demonstrates:

* Go project organization
* Dependency Injection
* Repository Pattern
* PostgreSQL with `pgx`
* Integration testing
* Docker & Docker Compose
* gRPC with Protocol Buffers
* Transport Adapters
* Outbox Pattern
* Redis integration (work in progress)

Future topics include:

* Redis Streams
* Event-Driven Architecture
* Consumer Groups
* Kubernetes
* Skaffold
* Bazel

---

# Architecture

Current synchronous request flow:

```text
Client
   │
   │ gRPC
   ▼
OrderGRPCServer
   │
   ▼
OrderService
   │
   ▼
OrderRepository
   │
   ▼
PostgreSQL
```

Current persistence flow:

```text
CreateOrder()

↓

BEGIN TRANSACTION

↓

Insert Order

↓

Insert Outbox Event

↓

COMMIT
```

Planned event flow:

```text
PostgreSQL Outbox

↓

Publisher

↓

Redis Stream

↓

Notification
Inventory
Analytics
Fraud Detection
```

---

# Project Structure

```text
cmd/
    client/
    server/
    publisher/

grpc/
    order.proto
    order.pb.go
    order_grpc.pb.go

internal/

    domain/

    repository/

    service/

    transport/
        grpc/

    redis/
```

## Folder Responsibilities

### `cmd/`

Contains runnable executables.

* `server` – gRPC server
* `client` – example gRPC client
* `publisher` – publishes Outbox events

### `internal/domain`

Business entities and domain concepts.

Examples:

* Order
* OutboxEvent
* Domain errors

This package should not know about:

* gRPC
* PostgreSQL
* Redis

---

### `internal/service`

Business logic.

Responsible for:

* validation
* orchestration
* creating business events

Should not contain:

* SQL
* protobuf
* Redis commands

---

### `internal/repository`

Persistence layer.

Responsible for:

* SQL
* transactions
* PostgreSQL

Owns persistence concerns.

---

### `internal/transport/grpc`

Transport Adapter.

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

The transport layer should never leak into the business layer.

---

### `grpc/`

Contains the Protocol Buffer contract and generated code.

The `.proto` file is the API contract between client and server.

---

# Technology Stack

* Go
* PostgreSQL
* pgx
* Docker
* gRPC
* Protocol Buffers
* Redis (in progress)

---

# Design Principles

This project follows a few core principles.

## Separation of Concerns

Each layer owns exactly one responsibility.

```text
Transport

↓

Business

↓

Persistence
```

---

## Depend on Abstractions

Services depend on interfaces, not concrete implementations.

Example:

```go
OrderService
    ↓
OrderRepository
```

instead of:

```go
OrderService
    ↓
PostgresOrderRepository
```

---

## Domain Independence

Business objects remain independent of transport.

```text
orderspb.Order

↓

Transport Adapter

↓

domain.Order
```

Changing gRPC to another communication technology should not affect the business layer.

---

## Transaction Ownership

Repositories own database transactions.

Services express business operations.

---

## Reliable Event Publication

The project implements the Outbox Pattern.

Business operations write both:

* the Order
* the Outbox Event

inside the same database transaction.

A separate publisher later distributes those events.

---

# Running the Project

Start infrastructure:

```bash
docker compose up -d
```

Run the server:

```bash
go run ./cmd/server
```

Run the client:

```bash
go run ./cmd/client
```

Run the publisher (work in progress):

```bash
go run ./cmd/publisher
```

--- 

# Learning Goals

This repository is intentionally educational.

The objective is to understand:

* why repositories exist
* why transport adapters exist
* why protobuf models differ from domain models
* why distributed systems require asynchronous communication
* why the Outbox Pattern exists
* why Redis Streams solve a different problem than queues

Every architectural decision in this repository is intended to answer a "why" question rather than simply demonstrate a framework.

---

# Roadmap

* [x] Go project structure
* [x] Repository Pattern
* [x] PostgreSQL
* [x] Docker
* [x] gRPC
* [x] Outbox Pattern
* [ ] Redis Streams
* [ ] Consumer Groups
* [ ] Event-Driven Architecture
* [ ] Kubernetes
* [ ] Skaffold
* [ ] Bazel

---

# Philosophy

The goal of this project is not to build an Order Service.

The goal is to build a realistic backend system while understanding the architectural reasoning behind every component.

Every new technology should solve a problem introduced by the previous stage of the project.

# Architecture

## Goal

Order System is a distributed event-driven application built in Go to demonstrate
production-grade backend engineering practices.

## Technologies

- Go
- gRPC
- PostgreSQL
- Redis Streams
- Docker
- Kubernetes

## Service Architecture

Client
    │
    ▼
Order Service
    │
    ├── ordersdb
    └── Outbox
            │
            ▼
      Publisher
            │
            ▼
     Redis Stream
      │         │
      ▼         ▼
Notification Inventory
                 │
                 ▼
            inventorydb
                 │
             Outbox
                 │
                 ▼
             Publisher
                 │
                 ▼
          stream:inventory
                 │
                 ▼
             Shipping

## Architectural Patterns

- Repository Pattern
- Dependency Injection
- gRPC
- Outbox Pattern
- CQRS
- Event-driven Architecture
# Aurora-order-system
A distributed order processing system built in Go using an Amazon Aurora–style primary-replica database architecture. Demonstrates read-write splitting, replication, and multi-node deployment across multiple machines.

---

# Overview

This project implements a **distributed order processing system** using a **microservice architecture written in Go**.
The system demonstrates key distributed systems concepts such as:

* Transaction coordination
* Concurrency control
* Service-to-service communication
* Database replication
* Aurora-style writer–reader database architecture

The database layer simulates **Amazon Aurora’s Primary–Replica cluster model**, where:

* The **Primary database handles writes**
* The **Replica database handles read operations**

This improves **scalability, availability, and performance**.

---

# System Architecture

The system is composed of three microservices and a replicated database cluster.

```
Client
   │
   ▼
Order Service (Transaction Coordinator)
   │
   ▼
Inventory Service
   │
   ▼
PostgreSQL Cluster
   ├── Primary Database (Writer)
   └── Replica Database (Reader)
```

---

# Deployment Architecture

The system can run across **three laptops**.

```
Laptop 1
Client

Laptop 2
Order Service
Inventory Service

Laptop 3
PostgreSQL Cluster
   ├── Primary (Port 5432)
   └── Replica (Port 5433)
```

Even when services run on the same machine, they communicate via **HTTP REST APIs**, maintaining a microservice architecture.

---

# Technologies Used

| Technology            | Purpose                           |
| --------------------- | --------------------------------- |
| Go (Golang)           | Backend microservices             |
| PostgreSQL            | Database                          |
| REST APIs             | Inter-service communication       |
| SQL Transactions      | Data consistency                  |
| Streaming Replication | Aurora-style database replication |
| JSON                  | Data exchange format              |

---
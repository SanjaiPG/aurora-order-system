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

How to Run the Project
1. Install Requirements

Install the following:

Go (1.20 or later)

PostgreSQL 14+

Git

Verify installation:

go version
psql --version
2. Setup PostgreSQL Primary

Start PostgreSQL normally.

Verify:

psql -p 5432 -U postgres

Create database:

CREATE DATABASE aurora_orders;

The database schema is provided in schema.sql. Create tables using schema.sql.

3. Setup PostgreSQL Replica

Create replica using:

pg_basebackup -h localhost -U replicator -D C:\aurora-replica -Fp -Xs -P -R

Start replica server:

& "C:\Program Files\PostgreSQL\18\bin\pg_ctl.exe" -D C:\aurora-replica start -l logfile

Edit replica configuration:

C:\aurora-replica\postgresql.conf

Set:

port = 5433

Verify replica:

psql -p 5433 -U postgres

Check:

SELECT pg_is_in_recovery();

If result is:

t

Replication is working.

4. Configure Environment Variables

Create .env file for services.

5. Run Inventory Service
cd inventory-service
go run main.go

Output:

Inventory Service running on port 8081
6. Run Order Service
cd order-service
go run main.go

Output:

Order Service running on port 8080
7. Run Client

Edit the client to point to the Order Service laptop:

http://<ORDER_SERVICE_IP>:8080/place-order

Run:

the index.html file and place the orders

Expected output:

Order Response: 200 OK
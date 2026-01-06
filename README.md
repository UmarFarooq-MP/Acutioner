# Auctioner – Real-Time Auction Service

## Overview

Auctioner is a **real-time auction system** built with **WebSockets** and designed using **Hexagonal (Ports & Adapters) Architecture**.

The primary goals of this project are:
- Clear separation of business logic from infrastructure
- Real-time bid propagation
- Horizontal scalability using Redis Pub/Sub
- Deterministic and testable core logic
- Production-ready Docker-based setup

---

## Key Features

- Real-time bidding via WebSockets  
- Redis Pub/Sub for horizontal scaling  
- Automatic auction timeout handling  
- Explicit bid rejection after auction close  
- Clean Hexagonal architecture  
- Fully Dockerized (one-command startup)  
- Deterministic and testable application layer  

---

## High-Level Architecture

```
Client (Browser / Postman / wscat)
        │
        │ WebSocket
        ▼
WebSocket Adapter (Hub, Clients)
        │
        │ Commands / Events
        ▼
Application Service (Use-cases)
        │
        │ Domain rules
        ▼
Domain (Auction Entity)
        │
        │ Persistence
        ▼
Repository Adapter (In-memory)
```

Redis is used as an **event bus**, not as a source of truth.

---

## Layer Responsibilities

### Domain Layer (`internal/domain`)
**Purpose:** Business rules only

- Core entity: `Auction`
- Enforces invariants:
  - Auction must be open
  - Bid must be higher than current highest
  - Auction closes after end time
- No knowledge of WebSockets, Redis, HTTP, or goroutines

---

### Application Layer (`internal/application`)
**Purpose:** Use-case orchestration

- Coordinates domain logic and infrastructure via ports
- Main use-cases:
  - `PlaceBid`
  - `CloseExpiredAuctions`
- Injected clock for deterministic testing
- Emits domain events through `Broadcaster`

---

### Ports (`internal/ports`)
**Purpose:** Define contracts

- `AuctionRepository`
- `Broadcaster`

These interfaces allow infrastructure to be swapped without touching domain or application code.

---

### Adapters (`internal/adapters`)
**Purpose:** Infrastructure & transport

#### Repository Adapter
- In-memory implementation
- Thread-safe using mutex
- Supports scanning open auctions (for timeouts)

#### WebSocket Adapter
- Manages WebSocket lifecycle
- Per-auction subscriptions
- Translates domain errors into protocol responses
- Sends:
  - Broadcast events (`NEW_BID`, `AUCTION_TIMED_OUT`)
  - Direct responses (`BID_REJECTED`)

#### Redis Adapter
- Redis Pub/Sub publisher
- Redis subscriber that bridges events into local WebSocket hub
- Enables cross-node real-time updates

---

## WebSocket Message Types

### Client → Server

```json
{
  "type": "PLACE_BID",
  "auction_id": "auction-1",
  "amount": 200
}
```

---

### Server → Client (Broadcast)

**New Bid**
```json
{
  "type": "NEW_BID",
  "auction_id": "auction-1",
  "amount": 200
}
```

**Auction Timed Out**
```json
{
  "type": "AUCTION_TIMED_OUT",
  "auction_id": "auction-1",
  "final_bid": 200
}
```

---

### Server → Client (Direct Response)

**Bid Rejected**
```json
{
  "type": "BID_REJECTED",
  "auction_id": "auction-1",
  "reason": "AUCTION_CLOSED"
}
```

Reasons:
- `AUCTION_CLOSED`
- `BID_TOO_LOW`
- `UNKNOWN_ERROR`

---

## Auction Timeout Handling

- A background worker runs every second
- Scans all open auctions
- Automatically closes expired auctions
- Emits `AUCTION_TIMED_OUT` event
- Works across all nodes via Redis Pub/Sub

> Note: Every node currently runs the timeout worker.  
> In production, this would be guarded via leader election or moved to a dedicated worker.

---

## Horizontal Scaling Strategy

- Each application instance:
  - Maintains its own WebSocket connections
  - Subscribes to Redis channels (`auction:*`)
- Events published by one node are delivered to all nodes
- Each node broadcasts only to its local subscribed clients

This keeps WebSocket connections **stateless across nodes**.

---

## Running the Application (Docker)

### Prerequisites
- Docker
- Docker Compose

### Start everything

```bash
docker compose up --build
```

This starts:
- Redis
- Auctioner WebSocket service

The application waits for Redis and fails fast if it is unavailable.

---

## Testing the Application

```bash
wscat -c "ws://localhost:8080/ws?auction_id=auction-1"
```

Place a bid:
```json
{"type":"PLACE_BID","auction_id":"auction-1","amount":200}
```

Timeout test:
- Wait until auction end time
- Clients receive `AUCTION_TIMED_OUT`
- Further bids return `BID_REJECTED`

---

## Entry Point & CLI

The application uses **Cobra CLI**.

```bash
auctioner serve --pub redis
auctioner serve --pub inmemory
```

---


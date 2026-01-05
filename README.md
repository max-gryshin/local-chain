# Local-Chain

A distributed blockchain implementation built with Go, featuring Raft consensus, UTXO-based transactions, and gRPC communication.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Main Entities](#main-entities)
- [Getting Started](#getting-started)

## Overview

Local-Chain is a private blockchain implementation that demonstrates core blockchain concepts including:

- **UTXO (Unspent Transaction Output) Model**: Similar to Bitcoin's transaction model
- **Merkle Trees**: For efficient transaction verification
- **Raft Consensus**: Ensures distributed consensus across multiple nodes
- **ECDSA Signatures**: For transaction authentication and security
- **gRPC API**: For client-server communication
- **LevelDB Storage**: Persistent data storage

## Architecture

The blockchain operates as a distributed system with multiple nodes communicating via Raft consensus protocol. Each node runs:

1. **gRPC Server**: Handles client requests (transactions, balance queries, user management) but executes on leader only
2. **Raft Consensus Layer**: Ensures all nodes agree on the blockchain state
3. **Blockchain Runner**: Periodically creates new blocks (every 10 seconds)
4. **Storage Layer**: LevelDB for persistent storage of blocks, transactions, UTXOs, and users

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Node 1     │◄────►│   Node 2     │◄────►│   Node 3     │
│  (Leader)    │ Raft │  (Follower)  │ Raft │  (Follower)  │
└──────────────┘      └──────────────┘      └──────────────┘
       │                     │                     │
       │                     │                     │
       ▼                     ▼                     ▼
  ┌─────────┐          ┌─────────┐          ┌─────────┐
  │ LevelDB │          │ LevelDB │          │ LevelDB │
  └─────────┘          └─────────┘          └─────────┘
```

## Main Entities

### 1. Block
Represents a container for transactions in the blockchain.

### 2. Transaction
Represents a transfer of value between users using the UTXO model.

### 3. User
Represents a user with cryptographic keys for transaction signing.

## Mains Services

### 1. Transactor

Core service handling transaction creation and validation.

**Responsibilities:**
- Create new transactions with proper UTXO selection
- Verify transaction signatures
- Calculate user balances
- Manage transaction pool (mempool)
- Prevent double-spending

### 2. Blockchain

Manages the blockchain state and block creation.

**Responsibilities:**
- Create new blocks from pending transactions
- Maintain chain integrity
- Interact with Raft for consensus
- Compute Merkle roots for blocks

### 3. Merkle Tree

Data structure for efficient transaction verification.

**Features:**
- Binary tree where each leaf is a transaction hash
- Parent nodes contain hashes of their children
- Root hash represents all transactions in the block

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make

### Running locally with docker

1. **Build and run the project:**
   ```bash
   make docker-up
   ```

2. **Add peer nodes to the cluster:**
   ```bash
   # Add node 2
   ./bin/debug add-peer --server 127.0.0.1:9001 \
     --id 00000000-0000-0000-0000-000000000002 \
     --address 172.25.0.12:8001

   # Add node 3
   ./bin/debug add-peer --server 127.0.0.1:9001 \
     --id 00000000-0000-0000-0000-000000000003 \
     --address 172.25.0.13:8001
   ```

3. **Add voting rights:**
   ```bash
   ./bin/debug add-voter --server 127.0.0.1:9001 \
     --id 00000000-0000-0000-0000-000000000002 \
     --address 172.25.0.12:8001

   ./bin/debug add-voter --server 127.0.0.1:9001 \
     --id 00000000-0000-0000-0000-000000000003 \
     --address 172.25.0.13:8001
   ```


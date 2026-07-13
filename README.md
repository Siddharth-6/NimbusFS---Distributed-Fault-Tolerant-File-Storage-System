# NimbusFS

A distributed, fault-tolerant file storage system built in Go that supports chunk-based storage, replication, automatic failure recovery, and heartbeat-based node monitoring.

NimbusFS is inspired by the architecture of Google File System (GFS) and HDFS. Files are divided into chunks, distributed across multiple storage nodes, replicated for fault tolerance, and automatically repaired when storage node failures occur.

---

## Features

- Chunk-based file storage
- Metadata server architecture
- Multiple independent storage nodes
- Configurable chunk placement
- Chunk replication for fault tolerance
- Automatic download failover
- Heartbeat-based node health monitoring
- Automatic replica recovery (Self-Healing)
- REST APIs for upload, download and metadata management
- Modular architecture for future scalability

---

## Architecture

```
                 Client
                    │
                    ▼
          +------------------+
          | Metadata Server  |
          +------------------+
             │         │
             │         └───────────────+
             │                         │
             ▼                         ▼
      Placement Manager         Heartbeat Monitor
             │                         │
             ▼                         ▼
       Storage Client          Failure Detection
             │
     ┌───────┼────────┐
     ▼       ▼        ▼
 Node1    Node2    Node3
```

---

## System Workflow

### Upload

```
Client

↓

Metadata Server

↓

Split file into chunks

↓

Placement Manager

↓

Replicate chunks

↓

Storage Nodes

↓

Store metadata
```

---

### Download

```
Client

↓

Metadata Server

↓

Lookup metadata

↓

Download chunks

↓

If replica fails

↓

Try another replica

↓

Merge chunks

↓

Return file
```

---

### Failure Recovery

```
Heartbeat detects failure

↓

Repair Service

↓

Download healthy replica

↓

Upload to new node

↓

Update metadata

↓

Replication restored
```

---

## Project Structure

```
cmd/
    metadata-server/
    storage-node/

internal/
    api/
    chunker/
    client/
    metadata/
    models/
    node/
    repair/
    service/
    storage/

data/
    metadata/
    node9001/
    node9002/
    node9003/
```

---

## Tech Stack

- Go
- HTTP REST APIs
- Goroutines
- Mutexes (sync.RWMutex)
- UUID
- JSON Metadata Storage

---

## Distributed Systems Concepts

- Distributed Metadata Management
- Chunk-based Storage
- Replication
- Fault Tolerance
- Failure Detection
- Heartbeat Monitoring
- Automatic Replica Recovery
- Download Failover
- Consistent Chunk Metadata
- Concurrent Network Communication

---

## Future Improvements

- SHA-256 checksum verification
- Parallel chunk uploads/downloads
- Dynamic node registration
- Consistent hashing
- Persistent metadata database
- Authentication & authorization

---

## Running the Project

### Start Storage Nodes

```bash
go run ./cmd/storage-node 9001

go run ./cmd/storage-node 9002

go run ./cmd/storage-node 9003
```

### Start Metadata Server

```bash
go run ./cmd/metadata-server
```

---

## Upload File

```bash
curl -X POST \
-F "file=@example.pdf" \
http://localhost:8080/upload
```

---

## Download File

```bash
curl \
"http://localhost:8080/download?id=<file-id>" \
-o output.pdf
```

---

## Demonstrated Capabilities

- Distributed storage architecture
- Fault-tolerant file storage
- Automatic recovery after node failures
- Replica-based high availability
- Modular and extensible system design

---

## Inspiration

- Google File System (GFS)
- Hadoop Distributed File System (HDFS)
- Ceph

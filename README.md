# distributed_cache_go
A simple distributed caching system using Go.
This project implements a distributed caching system consisting of two main services: a Cache Manager and multiple Data Nodes. These components communicate using a pub/sub mechanism and utilize Go channels for concurrent operations.

## Running the Project

```bash
go run cmd/cache_manager/main.go
```

### CURL Commands

```bash
curl -X GET "http://localhost:8080/?key=12345678"
curl -X POST "http://localhost:8080/?key=123456&value=hello"
curl -X DELETE "http://localhost:8080/?key=123456"
```

> POST is used for both inserting and updating key/value pairs.

## Running the tests

```bash
go test ./...
```

## System Components

1. Cache Manager 
 
   The Cache Manager is responsible for:
   - Handling client requests for data retrieval and storage.
   - Maintaining a mapping of keys to Data Node locations.
   - Coordinating with Data Nodes for data operations.
   - Managing cache invalidation across the system.
2. Data Nodes

   Each Data Node is responsible for:
   - Storing and retrieving multiple key/value pairs.
   - Handling concurrent read/write operations.
   - Implementing a simple eviction policy (Least Recently Used) for its local cache.
   - Communicating with the Cache Manager about its status and operations.

## Requirements
1. Use Go's standard library extensively.
2. Implement pub/sub communication using a simple in-memory solution.
3. Utilize Go channels for handling concurrent operations.
4. Implement basic error handling and logging.
5. Write unit tests for critical components.

## Project Structure

```
distributed_cache_go/
├── cmd/
│   └── cache_manager/
│       └── main.go
├── internal/
│   ├── cache_manager/
│   │   └── cache_manager.go
│   ├── data_node/
│   │   ├── data_node.go
│   │   └── eviction.go
│   └── pub_sub/
│       └── pubsub.go
├── pkg/
│   └── consistent_hash/
│       └── hash_ring.go
├── api/
│   └── http/
│       └── handler.go
├── test/
│   ├── cache_manager_test.go
│   ├── data_node_test.go
│   └── hash_ring_test.go
├── go.mod
└── README.md
```
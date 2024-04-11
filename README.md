# distributed-db

## How to Run

To start the first node and initialize the Raft cluster
```
go run .
```
This will start a node listening on localhost:8080 with a Raft address of localhost:12000 and persistently storing data in the `data/` directory

To run another node
```
go run . -tcpPort [port] -id [node id] -raftAddr [addr] -joinAddr [first node addr] -dataDir [data directory]
```

As an example
```
 go run . -tcpPort 8070 -id node1 -raftAddr localhost:12001 -joinAddr localhost:8080 -dataDir data2/
 ```
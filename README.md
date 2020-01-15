# chord-golang (WIP)

Implementation of the Chord protocol as described in the paper written in Golang: [Chord: A Scalable Peer-to-peer Lookup Protocol for Internet Applications](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf)

## Implementation

The code implements the Chord protocol as defined in the paper, as of now it does not take advantage of successor tables rather only a single successor is used. The code is based on the psuedocode provided in the publication. RPC Backed Virtual Nodes are used to transparently use the same Chord protocol functions (FindSucessor, Notify, etc, as defined in `local.go`). The system is designed such that multiple local worker threads can work together exposing RPC interfaces on different ports while communicating with each via direct method calls rather than using network resources.
```
      ---> LocalVNode: Local Implementation of a VNode. Contains implementation of the Chord Protocol
      |
VNode.VNodeProtcol ---> Generic Interface to the Chord Protocol.
      |
      ---> RemoteVNode: RPC Backed VNode which communicates to the RPC Server (defined in `rpcserver.go`) which calls the remote processes' LocalVNode to fulfil the RPC. 
```

## Try it out.
- Build the program.
```
  go build
```

- Run a single worker thread on **127.0.0.1:8000**.
```
  ./src -host 127.0.0.1:8000
```

- Connect to a remote host.
```
  ./src -mode join -host 127.0.0.1:8001 -rhost 127.0.0.1:8000
```

- Run 8 Local Worker Threads on Randomly Assigned Ports **(0.0.0.0:0)**.
```
  ./src -workers 8
```

## TODO
- successor tables
- Add graceful leaves for VNodes.
- Different RPC implementations.
- Lookup HTTP Server.
- Hostname, ports, polishing.
- blog article
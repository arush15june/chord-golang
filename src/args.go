package main

import (
	"flag"
)

const (
	// DefaultIDBits is the default number of bits held by a Identifier.
	DefaultIDBits = 64

	// DefaultFingers is the default number of fingers held by a VNode.
	DefaultFingers = 1

	// DefaultSuccessors is the default number of successors held by a VNode.
	DefaultSuccessors = 1

	// DefaultVirtualNodes is the default number of virtual nodes to map to the physical node.
	DefaultVirtualNodes = 1
)

var (
	// NodeMode selects the mode of the node to be create or join.
	NodeMode = flag.String("mode", "create", "'join' or 'create'")
	// NodeModeShort = flag.String("m")

	// Workers selects the number of virtual nodes to create on the server.
	Workers = flag.Int("workers", DefaultVirtualNodes, "Number of virtual nodes to start.")
	// WorkersShort = flag.Int("w")

	// HostName is the hostname of the physical chord node.
	HostName = flag.String("host", "127.0.0.1:8000", "Self hostname. Default: :0 (use <hostname>:0 for random port assignment")
	// HostNameShort = flag.String("h")

	// RemoteHost is the hostname of an existing node in the chord network if join mode is selected.
	RemoteHost = flag.String("rhost", "127.0.0.1:8000", "Remote chord node to join.")
	// RemoteHostShort = flag.String("r")
)

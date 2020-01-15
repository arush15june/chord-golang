package rpc

// Pluggable transports for RemoteVNodes.

import (
	VNode "github.com/arush15june/chord-golang/src/pkg/vnode"
)

// ChordProtocolRPC is the interface required by the RemoteVNode,
// A Concrete ChordRPC will resolve the methods on the
// remote node and return the results to the caller.

type ChordProtocolRPC interface {
	// FindSuccessors finds N successors of the VNode.
	// FindSuccessors(int) ([]*VNode.VNodeProtocol, error)

	// FindSuccessor finds the successor for a Key.
	FindSuccessor(uint64) (VNode.VNodeProtocol, error)

	// Notify notifies the VNode of its new predecessor.
	Notify(VNode.VNodeProtocol) error

	// Ping sends a request to a VNode
	Ping() error

	// GetPredecessor returns the predecessor VNode.
	GetPredecessor() (VNode.VNodeProtocol, error)
}

type FindSuccRpcArgs struct {
	ID uint64
}
type FindSuccRpcReply struct {
	Hostname string
}

type NotifyRpcArgs struct {
	Hostname string
}
type NotifyRpcReply struct{}

type PingRpcArgs struct{}
type PingRpcReply struct{}

type GetPredecessorRpcArgs struct{}
type GetPredecessorRpcReply struct {
	Hostname string
}

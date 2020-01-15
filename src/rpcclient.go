package main

import (
	"errors"
	"net/rpc"

	RPC "github.com/arush15june/chord-golang/src/pkg/rpc"
	VNode "github.com/arush15june/chord-golang/src/pkg/vnode"
)

// ChordTCPRPCClient implements RPC for the Chord protocol using Golang net/rpc.
type ChordTCPRPCClient struct {
	client   *rpc.Client
	vnode    VNode.VNodeProtocol
	Hostname string
}

// InitChordTCPRPCClient initializes a ChordTCPRPC object ready to create a server or call a client.
func InitChordTCPRPCClient(HostnameWithPort string, vnode VNode.VNodeProtocol) *ChordTCPRPCClient {
	rpc := &ChordTCPRPCClient{
		Hostname: HostnameWithPort,
		vnode:    vnode,
	}

	return rpc
}

// InitClient creates a connection to a ChordTCPRPC server and returns its client.
func (rpcInstance *ChordTCPRPCClient) InitClient() error {
	if rpcInstance.client != nil {
		return nil
	}

	client, err := rpc.Dial("tcp", rpcInstance.Hostname)
	if err != nil {
		return errors.New("client is dead")
	}
	rpcInstance.client = client

	pingErr := rpcInstance.Ping()
	if pingErr != nil {
		return errors.New("client is dead")
	}

	return nil
}

// FindSuccessor calls FindSuccessorRPC on a remote node and returns the successor.
func (rpc *ChordTCPRPCClient) FindSuccessor(ID uint64) (VNode.VNodeProtocol, error) {
	var err error

	err = rpc.InitClient()
	if err != nil {
		return nil, err
	}

	args := &RPC.FindSuccRpcArgs{ID: ID}
	reply := &RPC.FindSuccRpcReply{}

	err = rpc.client.Call(findSuccRPCName, args, reply)

	if err != nil {
		return nil, err
	}
	return InitRemoteVNode(reply.Hostname), nil
}

// Notify calls NotifyRPC on the remote node and returns the successor node.
func (rpc *ChordTCPRPCClient) Notify(vnode VNode.VNodeProtocol) error {
	var err error

	err = rpc.InitClient()
	if err != nil {
		return err
	}

	args := &RPC.NotifyRpcArgs{Hostname: vnode.Hostname()}
	reply := &RPC.NotifyRpcReply{}

	err = rpc.client.Call(notifyRPCName, args, reply)

	if err != nil {
		return err
	}

	return nil
}

// Ping calls PingRPC on the remote node.
func (rpc *ChordTCPRPCClient) Ping() error {
	var err error

	err = rpc.InitClient()
	if err != nil {
		return err
	}

	args := &RPC.PingRpcArgs{}
	reply := &RPC.PingRpcReply{}

	err = rpc.client.Call(pingRPCName, args, reply)
	if err != nil {
		return err
	}

	return nil
}

// GetPredecessor calls GetPredecessorRPC on the remote node and returns the predecessor.
func (rpc *ChordTCPRPCClient) GetPredecessor() (VNode.VNodeProtocol, error) {
	var err error

	err = rpc.InitClient()
	if err != nil {
		return nil, err
	}

	args := &RPC.GetPredecessorRpcArgs{}
	reply := &RPC.GetPredecessorRpcReply{}

	err = rpc.client.Call(getPredRPCName, args, reply)
	if err != nil {
		return nil, err
	}

	return InitRemoteVNode(reply.Hostname), nil
}

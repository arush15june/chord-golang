package main

import (
	"errors"
	"net"
	"net/rpc"

	RPC "github.com/arush15june/chord-golang/src/pkg/rpc"
	VNode "github.com/arush15june/chord-golang/src/pkg/vnode"
)

const (
	findSuccRPCName = "ChordTCPRPCServer.FindSuccessorRPC"
	notifyRPCName   = "ChordTCPRPCServer.NotifyRPC"
	pingRPCName     = "ChordTCPRPCServer.PingRPC"
	getPredRPCName  = "ChordTCPRPCServer.GetPredecessorRPC"
)

// ChordTCPRPCServer implements RPC for the Chord protocol using Golang net/rpc.
type ChordTCPRPCServer struct {
	client   *rpc.Client
	vnode    VNode.VNodeProtocol
	Hostname string
}

// InitChordTCPRPCServer initializes a ChordTCPRPC object ready to create a server or call a client.
func InitChordTCPRPCServer(HostnameWithPort string, vnode VNode.VNodeProtocol) *ChordTCPRPCServer {
	rpc := &ChordTCPRPCServer{
		Hostname: HostnameWithPort,
		vnode:    vnode,
	}

	return rpc
}

// InitServer starts Chord Protocol TCP-RPC server on Hostname:Port.
func InitServer(rpcInstance *ChordTCPRPCServer) error {
	rpc.Register(rpcInstance)
	l, e := net.Listen("tcp", rpcInstance.Hostname)
	if e != nil {
		return errors.New("failed to start Listen server")
	}

	// Reset address as acquired by Listener.
	address := l.Addr().String()
	rpcInstance.Hostname = address

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	return nil
}

// FindSuccessorRPC implements the method executed the by the RPC server to find successors on local vnode.
func (rpc *ChordTCPRPCServer) FindSuccessorRPC(args *RPC.FindSuccRpcArgs, reply *RPC.FindSuccRpcReply) error {
	successor, err := rpc.vnode.FindSuccessor(args.ID)
	if err != nil {
		return err
	}

	reply.Hostname = successor.Hostname()

	return nil
}

// NotifyRPC implements the method executed by the RPC server to notify local vnode.
func (rpc *ChordTCPRPCServer) NotifyRPC(args *RPC.NotifyRpcArgs, reply *RPC.NotifyRpcReply) error {
	err := rpc.vnode.Notify(InitRemoteVNode(args.Hostname))
	if err != nil {
		return err
	}
	return nil
}

// PingRPC implements the method executed by the RPC server to ping local vnode.
func (rpc *ChordTCPRPCServer) PingRPC(args *RPC.PingRpcArgs, reply *RPC.PingRpcReply) error {
	err := rpc.vnode.Ping()
	if err != nil {
		return errors.New("node died")
	}
	return err
}

// GetPredecessorRPC implements the method executed by the RPC server to get predecessor of local vnode.
func (rpc *ChordTCPRPCServer) GetPredecessorRPC(args *RPC.GetPredecessorRpcArgs, reply *RPC.GetPredecessorRpcReply) error {
	predecessor, err := rpc.vnode.GetPredecessor()
	if err != nil {
		return err
	}

	reply.Hostname = predecessor.Hostname()

	return nil
}

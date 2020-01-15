package main

import (
	Hash "github.com/arush15june/chord-golang/src/pkg/hash"
	RPC "github.com/arush15june/chord-golang/src/pkg/rpc"
	Util "github.com/arush15june/chord-golang/src/pkg/util"
	VNode "github.com/arush15june/chord-golang/src/pkg/vnode"
)

type RemoteVNode struct {
	VNode.VNode
	rpc RPC.ChordProtocolRPC
}

func InitRemoteVNode(Hostname string) *RemoteVNode {
	rvnode := &RemoteVNode{
		VNode: VNode.VNode{Hostname: Hostname},
		rpc:   InitChordTCPRPCClient(Hostname, nil),
	}
	return rvnode
}

func (node *RemoteVNode) FindSuccessors(id int) ([]VNode.VNodeProtocol, error) {
	return nil, nil
}
func (node *RemoteVNode) Notify(vnode VNode.VNodeProtocol) error {
	return node.rpc.Notify(vnode)
}
func (node *RemoteVNode) FindSuccessor(id uint64) (VNode.VNodeProtocol, error) {
	n, err := node.rpc.FindSuccessor(id)

	return n, err
}
func (node *RemoteVNode) Ping() error {
	return node.rpc.Ping()
}
func (node *RemoteVNode) CheckPredecessor() error {
	return nil
}
func (node *RemoteVNode) GetPredecessor() (VNode.VNodeProtocol, error) {
	return node.rpc.GetPredecessor()
}

func (node *RemoteVNode) IsBetweenNodes(vlow VNode.VNodeProtocol, vhigh VNode.VNodeProtocol) bool {
	return Util.IsBetweenID(node.ID(), vlow.ID(), vhigh.ID())
}

func (node *RemoteVNode) Hostname() string {
	return node.VNode.Hostname
}

func (node *RemoteVNode) ID() uint64 {
	return Hash.Sum([]byte(node.Hostname()))
}

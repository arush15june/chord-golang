package vnode

// VNodeProtocol implements the Chord protocol on Vnodes.
// Local VNodes can implement it via method calls.
// Remote VNodes can use RPC to transparently work like Local VNodes.
type VNodeProtocol interface {
	// FindSuccessors finds N successors of the VNode.
	FindSuccessors(int) ([]VNodeProtocol, error)

	// FindSuccessor finds the successor for a Key.
	FindSuccessor(uint64) (VNodeProtocol, error)

	// Notify notifies the VNode of its new predecessor.
	Notify(VNodeProtocol) error

	// Ping sends a request to a VNode
	Ping() error

	// CheckPredecessor checks the aliveness of VNode's predecessor.
	CheckPredecessor() error

	// GetPredecessor returns the predecessor VNode.
	GetPredecessor() (VNodeProtocol, error)

	// IsBetweenNodes
	IsBetweenNodes(VNodeProtocol, VNodeProtocol) bool

	// ID returns the ID of the VNode.
	ID() uint64

	// Hostname returns the hostname of the VNode.
	Hostname() string
}

// VNode is a virtual node running the chord protocol.
type VNode struct {
	// Hostname is the hostname of the VNode.
	Hostname string
}

// Conform VNode to VNodeProtocol.
func (v *VNode) FindSuccessors(int) ([]*VNodeProtocol, error) {
	return nil, nil
}
func (v *VNode) Notify(*VNodeProtocol) error {
	return nil
}
func (v *VNode) FindSuccessor(uint64) (*VNodeProtocol, error) {
	return nil, nil
}
func (v *VNode) Ping() error {
	return nil
}
func (v *VNode) CheckPredecessor() error {
	return nil
}
func (v *VNode) GetPredecessor() (*VNodeProtocol, error) {
	return nil, nil
}
func (v *VNode) IsBetweenNodes(*VNodeProtocol, *VNodeProtocol) bool {
	return true
}

// InitVNode initializes the VNode by computing the hash
// of the supplied hostname and setting it as the id.
func InitVNode(workerHostname string) (*VNode, error) {
	vnode := VNode{Hostname: workerHostname}
	return &vnode, nil
}

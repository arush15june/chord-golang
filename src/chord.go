package main

// Implements Chord Node,
// Initiates the VNodes for the current physical node,
// Initiates the RPC Transport Listener,
// Handles requests from VNodes to the RPC transport,
// Delegates requests to VNodes from the RPC Transport

var (
	workers []*LocalVNode
)

// NewLocalVNodeWithRPC initializes a NewLocalVNode and starts a ChordTCPRPCServer on it.
func NewLocalVNodeWithRPC(hostname string, minStabilizeInterval int, maxStabilizeInterval int, fixFingerInterval int, checkPredInterval int, maxSuccessors int, maxFingers int) (*LocalVNode, error) {
	logger.Printf("Initializing New Local VNode: %s\n", hostname)

	vnode, err := InitLocalVNode(hostname, minStabilizeInterval, maxStabilizeInterval, fixFingerInterval, checkPredInterval, maxSuccessors, maxFingers)
	if err != nil {
		return nil, err
	}

	rpc := InitChordTCPRPCServer(hostname, vnode)
	InitServer(rpc)
	vnode.SetHostname(rpc.Hostname)

	logger.Printf("%s RPC Server Initialized", vnode.Hostname())

	return vnode, nil
}

// JoinRing initializes local vnode workers and joins an existing chord ring.
// It selects a random.
func JoinRing(nWorkers int, hostname string, remoteHostName string, minStabilizeInterval int, maxStabilizeInterval int, fixFingerInterval int, checkPredInterval, maxSuccessors int, maxFingers int) error {
	remoteVNode := InitRemoteVNode(remoteHostName)

	workers := make([]*LocalVNode, nWorkers)
	for i := 0; i < nWorkers; i++ {
		vnode, _ := NewLocalVNodeWithRPC(
			hostname,
			minStabilizeInterval,
			maxStabilizeInterval,
			fixFingerInterval,
			checkPredInterval,
			maxSuccessors,
			maxFingers,
		)

		go func() {
			vnode.Join(remoteVNode)
			vnode.InitializeFingerTables()
			vnode.StartWorker()
		}()
		workers[i] = vnode
	}

	return nil
}

// CreateRing creates a Chord ring in one of the local virtual
// and joins all other local nodes to it.
func CreateRing(nWorkers int, hostname string, minStabilizeInterval int, maxStabilizeInterval int, fixFingerInterval int, checkPredInterval int, maxSuccessors int, maxFingers int) error {
	workers := make([]*LocalVNode, nWorkers)
	for i := 0; i < nWorkers; i++ {
		vnode, _ := NewLocalVNodeWithRPC(
			hostname,
			minStabilizeInterval,
			maxStabilizeInterval,
			fixFingerInterval,
			checkPredInterval,
			maxSuccessors,
			maxFingers,
		)
		workers[i] = vnode
	}

	workers[0].Create()
	workers[0].StartWorker()

	for i := 1; i < nWorkers; i++ {
		go func(i int) {
			workers[i].Join(workers[0])
			workers[i].InitializeFingerTables()
			workers[i].StartWorker()
		}(i)
	}

	return nil
}

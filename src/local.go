package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	Hash "github.com/arush15june/chord-golang/src/pkg/hash"
	Util "github.com/arush15june/chord-golang/src/pkg/util"
	VNode "github.com/arush15june/chord-golang/src/pkg/vnode"
)

// LocalVNode is a VNode communicating using direct method calls rather than RPC.
type LocalVNode struct {
	VNode.VNode

	successors    []VNode.VNodeProtocol
	predecessor   VNode.VNodeProtocol
	maxSuccessors int

	fingers    []VNode.VNodeProtocol
	maxFingers int

	minStabilizeInterval int
	maxStabilizeInterval int
	fixFingerInterval    int
	checkPredInterval    int

	stopStabilizeChan chan bool
	stopFixFingerChan chan bool
	stopCheckPredChan chan bool
}

// InitLocalVNode initializes a local vnode by computing the hash of the hostname string.
func InitLocalVNode(
	Hostname string,
	minStabilizeInterval int,
	maxStabilizeInternval int,
	fixFingerInterval int,
	checkPredInterval int,
	maxSuccessors int,
	maxFingers int,
) (*LocalVNode, error) {
	vnode := &LocalVNode{
		VNode:                VNode.VNode{Hostname: Hostname},
		minStabilizeInterval: minStabilizeInterval,
		maxStabilizeInterval: maxStabilizeInternval,
		fixFingerInterval:    fixFingerInterval,
		checkPredInterval:    checkPredInterval,
		maxSuccessors:        maxSuccessors,
		maxFingers:           maxFingers,
	}

	vnode.initStopChannels()
	vnode.initLists()
	return vnode, nil
}

// StopVNode stops the VNode background operations.
func (node *LocalVNode) StopVNode() {
	node.stopStabilizeChan <- true
	node.stopFixFingerChan <- true
	node.stopCheckPredChan <- true
}

// SetHostname sets a new hostname for the id and recomputes the ID.
func (node *LocalVNode) SetHostname(newHostname string) {
	logger.Printf("Changing Hostname: %s -> %s", node.Hostname(), newHostname)

	node.VNode.Hostname = newHostname
}

// initStopChannels initializes the channels used to stop backgorund goroutines.
func (node *LocalVNode) initStopChannels() {
	node.stopStabilizeChan = make(chan bool)
	node.stopFixFingerChan = make(chan bool)
	node.stopCheckPredChan = make(chan bool)
}

func (node *LocalVNode) initLists() {
	node.successors = make([]VNode.VNodeProtocol, node.maxSuccessors)
	node.fingers = make([]VNode.VNodeProtocol, node.maxFingers)
}

// Hostname returns the hostname of the LocalVNode.
func (node *LocalVNode) Hostname() string {
	return node.VNode.Hostname
}

// ID returns the ID of the LocalVNode.
func (node *LocalVNode) ID() uint64 {
	return Hash.Sum([]byte(node.Hostname()))
}

// Stabilize executes after certain time intervals to fix successors.
// > Each time node n runs Stabilize(), it asks its successor
// > for the successor’s predecessor p, and decides whether p
// > should be n’s successor instead. stabilize() notifies node
// > n’s successor of n’s existence, giving the successor the chance
// > to change its predecessor to n.
func (node *LocalVNode) Stabilize() error {
	logger.Printf("[%s, %d] Stabilizing VNode\n", node.Hostname(), node.ID())

	verifySuccesorNode, _ := node.successors[0].GetPredecessor()
	if verifySuccesorNode != nil && verifySuccesorNode.IsBetweenNodes(node, node.successors[0]) {
		node.successors[0] = verifySuccesorNode
		logger.Printf("[%s, %d] Updated successor: %s\n", node.Hostname(), node.ID(), verifySuccesorNode.Hostname())
	}

	if node.successors[0].ID() != node.ID() {
		err := node.successors[0].Notify(node)
		logger.Printf("[%s, %d] Notified %s of VNode.\n", node.Hostname(), node.ID(), node.successors[0].Hostname())
		return err
	}

	return nil
}

// StabilizeRoutine runs Stabilize() periodically by choosing an interval
// between minStabilizeInterval and maxStabilizeInterval.
func (node *LocalVNode) StabilizeRoutine() error {
	exit := false

	for {
		node.Stabilize()

		interval := Util.GetRandomBetween(node.minStabilizeInterval, node.maxStabilizeInterval)
		timer := time.NewTimer(time.Duration(interval) * time.Second)
		select {
		case <-timer.C:
		case <-node.stopStabilizeChan:
			exit = true
		}
		if exit {
			break
		}
	}
	return nil
}

// FixFinger updates the finger tables. fingerNumber is the n'th finger not the finger list index.
func (node *LocalVNode) FixFinger(fingerNumber int) error {
	logger.Printf("[%s, %d] Fixing Finger %d", node.Hostname(), node.ID(), fingerNumber)
	if fingerNumber < 1 {
		return errors.New("invalid finger number")
	}
	if fingerNumber > node.maxFingers {
		return errors.New("finger number out of bounds")
	}

	fingerIndex := fingerNumber - 1
	fingerID := node.ID() + uint64(math.Exp2(float64(fingerNumber-1)))

	var err error
	node.fingers[fingerIndex], err = node.FindSuccessor(fingerID)

	if node.fingers[fingerIndex] != nil {
		logger.Printf("[%s, %d] Finger[%d], %d: %s", node.Hostname(), node.ID(), fingerNumber, fingerID, node.fingers[0].Hostname())
	}

	return err
}

// FixFingersRoutine periodically fixes the finger indices.
func (node *LocalVNode) FixFingersRoutine() error {
	fingerNumber := 1
	exit := false

	for {
		err := node.FixFinger(fingerNumber)
		if err != nil {
			// Reset fixing process
			fingerNumber = 1
			continue
		}

		fingerTimer := time.NewTimer(time.Duration(node.fixFingerInterval) * time.Second)
		select {
		case <-fingerTimer.C:
		case <-node.stopFixFingerChan:
			exit = true
		}
		if exit {
			break
		}

		if fingerNumber >= node.maxFingers {
			fingerNumber = 1
		} else {
			fingerNumber++
		}
	}

	return nil
}

// Ping returns nil as LocalVNode will always be alive.
func (node *LocalVNode) Ping() error {
	logger.Printf("[%s, %d] Received request for Ping.", node.Hostname(), node.ID())
	return nil
}

// FindSuccessors finds the next N successors on the VNode.
func (node *LocalVNode) FindSuccessors(n int) ([]VNode.VNodeProtocol, error) {
	return nil, nil
}

// FindSuccessor finds the successor for the key id recursively.
func (node *LocalVNode) FindSuccessor(id uint64) (VNode.VNodeProtocol, error) {
	logger.Printf("[%s, %d] Finding Successor: %d\n", node.Hostname(), node.ID(), id)

	if Util.IsBetweenID(id, node.ID(), node.successors[0].ID()) {
		logger.Printf("[%s, %d] %d lies between node[%s, %d] and successor[%s, %d]", node.Hostname(), node.ID(), id, node.Hostname(), node.ID(), node.successors[0].Hostname(), node.successors[0].ID())
		return node.successors[0], nil
	}

	logger.Printf("[%s, %d] %d not in successor, finding closest predecessor.", node.Hostname(), node.ID(), id)
	closestNode := node.ClosestPrecedingNode(id)
	if closestNode.ID() == node.ID() {
		return node, nil
	}

	return closestNode.FindSuccessor(id)
}

// ClosestPrecedingNode finds the closest preceding node to the ID in the FingerTable.
func (node *LocalVNode) ClosestPrecedingNode(id uint64) VNode.VNodeProtocol {
	for _, finger := range node.fingers {
		if finger != nil {
			if Util.IsBetweenID(finger.ID(), node.ID(), id) {
				logger.Printf("[%s, %d] closest preceeding node for %d: %s, %d", node.Hostname(), node.ID(), id, finger.Hostname(), finger.ID())
				return finger
			}
		}
	}

	return node
}

// Notify verifies the notifying node to be its predecessor and updates itself.
func (node *LocalVNode) Notify(notifyingNode VNode.VNodeProtocol) error {
	logger.Printf("[%s, %d] Notification from [%s, %d]\n", node.Hostname(), node.ID(), notifyingNode.Hostname(), notifyingNode.ID())

	if node.predecessor == nil || notifyingNode.IsBetweenNodes(node.predecessor, node) {
		logger.Printf("[%s, %d] Update predecessor to [%s, %d]\n", node.Hostname(), node.ID(), notifyingNode.Hostname(), notifyingNode.ID())
		node.predecessor = notifyingNode
	}

	return nil
}

// CheckPredecessor verifies if the nodes predecessor is alive.
func (node *LocalVNode) CheckPredecessor() error {
	logger.Printf("[%s, %d] Checking liveness of predecessor\n", node.Hostname(), node.ID())
	if node.predecessor == nil {
		logger.Printf("[%s, %d] No predecessor present.", node.Hostname(), node.ID())
		return nil
	}

	err := node.predecessor.Ping()
	if err != nil {
		errLog := fmt.Sprintf("[%s, %d] Predecessor %s dead.", node.Hostname(), node.ID(), node.predecessor.Hostname())
		logger.Println(errLog)
		node.predecessor = nil
		return errors.New(errLog)
	}

	return nil
}

// CheckPredecessorRoutine checks the liveness of predecessor periodically.
func (node *LocalVNode) CheckPredecessorRoutine() error {
	exit := false

	for {
		node.CheckPredecessor()

		interval := node.checkPredInterval
		timer := time.NewTimer(time.Duration(interval) * time.Second)
		select {
		case <-timer.C:
		case <-node.stopCheckPredChan:
			exit = true
		}

		if exit {
			break
		}
	}

	return nil
}

// GetPredecessor returns the predecessor of the VNode.
func (node *LocalVNode) GetPredecessor() (VNode.VNodeProtocol, error) {
	if node.predecessor == nil {
		logger.Printf("[%s, %d] Node has no predecessor.\n", node.Hostname(), node.ID())
		return nil, errors.New("VNode does not have predecessor")
	}

	logger.Printf("[%s, %d] Returning predecessor %s\n", node.Hostname(), node.ID(), node.predecessor.Hostname())
	return node.predecessor, nil
}

func (node *LocalVNode) IsBetweenNodes(vlow VNode.VNodeProtocol, vhigh VNode.VNodeProtocol) bool {
	return Util.IsBetweenID(node.ID(), vlow.ID(), vhigh.ID())
}

// Join joins an existing ChordVNode.
func (node *LocalVNode) Join(chordVNode VNode.VNodeProtocol) error {
	var err error

	logger.Printf("[%s, %d] Joining VNode: [%s, %d]", node.Hostname(), node.ID(), chordVNode.Hostname(), chordVNode.ID())

	node.predecessor = nil
	node.successors[0], err = chordVNode.FindSuccessor(node.ID())
	logger.Printf("[%s, %d] Found first successor [%s, %d]", node.Hostname(), node.ID(), node.successors[0].Hostname(), node.successors[0].ID())

	// for i := 0; i < node.maxFingers; i++ {
	// 	node.fingers[i] = node.successors[0]
	// }

	node.successors[0].Notify(node)

	return err
}

// Create creates a new chord ring.
func (node *LocalVNode) Create() error {
	logger.Printf("[%s, %d] Creating new Chord ring\n", node.Hostname(), node.ID())

	node.predecessor = nil
	node.successors[0] = node
	// for i := 0; i < node.maxFingers; i++ {
	// 	node.fingers[i] = node
	// }

	return nil
}

// Lookup finds the successor of ID.
func (node *LocalVNode) Lookup(Key string) (string, error) {
	ID := Hash.Sum([]byte(Key))

	logger.Printf("[%s, %d] Lookup request for %d\n", node.Hostname(), node.ID(), ID)

	vnode, err := node.FindSuccessor(ID)
	if err != nil {
		logger.Printf("[%s, %d] Error occured: %s\n", node.Hostname(), node.ID(), err)
		return "", err
	}
	logger.Printf("[%s, %d] Lookup: %d -> [%s, %d]\n", node.Hostname(), node.ID(), ID, vnode.Hostname(), vnode.ID())

	return vnode.Hostname(), nil
}

func (node *LocalVNode) InitializeFingerTables() {
	node.FixFinger(1)
}

// StartWorker initializes background operations of the worker.
func (node *LocalVNode) StartWorker() {
	go node.StabilizeRoutine()
	go node.FixFingersRoutine()
	go node.CheckPredecessorRoutine()
}

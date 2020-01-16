package main

import (
	"flag"
)

// Lookup is supported via an HTTP API.

// CreateStrategy is used to create a new Chord ring.
func CreateStrategy() {
	nWorkers := *Workers
	Hostname := *HostName

	logger.Println("Creating New Ring")
	logger.Printf("Hostname: %s", Hostname)
	logger.Printf("VNode Workers: %d", nWorkers)

	minStabilizeInterval := 15
	maxStabilizeInterval := 45
	fixFingerInterval := 15
	checkPredInterval := 15
	maxSuccessors := 1
	maxFingers := 6

	logger.Println("VNode Worker Configuration")
	logger.Printf("Minimum Stabilization Interval: %d", minStabilizeInterval)
	logger.Printf("Maximum Stabilization Interval: %d", maxStabilizeInterval)
	logger.Printf("Finger Fixing Interval: %d", fixFingerInterval)
	logger.Printf("Check Predecessor Interval: %d", checkPredInterval)
	logger.Printf("No of Successors in Successor Table: %d", maxSuccessors)
	logger.Printf("No of Fingers in Finger Table: %d", maxFingers)

	CreateRing(
		nWorkers,
		Hostname,
		minStabilizeInterval,
		maxStabilizeInterval,
		fixFingerInterval,
		checkPredInterval,
		maxSuccessors,
		maxFingers,
	)

}

// JoinStrategy is used to join an existing chord ring.
func JoinStrategy() {
	nWorkers := *Workers
	Hostname := *HostName
	RemoteHost := *RemoteHost

	logger.Println("Joining Existing Ringg")
	logger.Printf("Hostname: %s", Hostname)
	logger.Printf("VNode Workers: %d", nWorkers)
	logger.Printf("Remote Ring: %s", RemoteHost)

	minStabilizeInterval := 15
	maxStabilizeInterval := 45
	fixFingerInterval := 15
	checkPredInterval := 15
	maxSuccessors := 1
	maxFingers := 6

	logger.Println("VNode Worker Configuration")
	logger.Printf("Minimum Stabilization Interval: %d", minStabilizeInterval)
	logger.Printf("Maximum Stabilization Interval: %d", maxStabilizeInterval)
	logger.Printf("Finger Fixing Interval: %d", fixFingerInterval)
	logger.Printf("Check Predecessor Interval: %d", checkPredInterval)
	logger.Printf("No of Successors in Successor Table: %d", maxSuccessors)
	logger.Printf("No of Fingers in Finger Table: %d", maxFingers)

	JoinRing(
		nWorkers,
		Hostname,
		RemoteHost,
		minStabilizeInterval,
		maxStabilizeInterval,
		fixFingerInterval,
		checkPredInterval,
		maxSuccessors,
		maxFingers,
	)

}

func main() {

	InitLogger()
	flag.Parse()

	switch *NodeMode {
	case "create":
		CreateStrategy()
		break
	case "join":
		JoinStrategy()
		break
	default:
		break
	}

	InitHttpServer()
	for {

	}

}

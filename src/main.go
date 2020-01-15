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
	logger.Printf("Hostname: %s\n", Hostname)
	logger.Println("VNode Workers: %s\n", nWorkers)

	minStabilizeInterval := 15
	maxStabilizeInterval := 45
	fixFingerInterval := 15
	checkPredInterval := 15
	maxSuccessors := 1
	maxFingers := 6

	logger.Println("VNode Worker Configuration")
	logger.Printf("Minimum Stabilization Interval: %d\n", minStabilizeInterval)
	logger.Printf("Maximum Stabilization Interval: %d\n", maxStabilizeInterval)
	logger.Printf("Finger Fixing Interval: %d\n", fixFingerInterval)
	logger.Printf("Check Predecessor Interval: %d\n", checkPredInterval)
	logger.Printf("No of Successors in Successor Table: %d\n", maxSuccessors)
	logger.Printf("No of Fingers in Finger Table: %d\n", maxFingers)

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

	for {

	}

}

// JoinStrategy is used to join an existing chord ring.
func JoinStrategy() {
	nWorkers := *Workers
	Hostname := *HostName
	RemoteHost := *RemoteHost

	logger.Println("Joining Existing Ringg")
	logger.Printf("Hostname: %s\n", Hostname)
	logger.Println("VNode Workers: %s\n", nWorkers)
	logger.Println("Remote Ring: %s\n", RemoteHost)

	minStabilizeInterval := 15
	maxStabilizeInterval := 45
	fixFingerInterval := 15
	checkPredInterval := 15
	maxSuccessors := 1
	maxFingers := 6

	logger.Println("VNode Worker Configuration")
	logger.Printf("Minimum Stabilization Interval: %d\n", minStabilizeInterval)
	logger.Printf("Maximum Stabilization Interval: %d\n", maxStabilizeInterval)
	logger.Printf("Finger Fixing Interval: %d\n", fixFingerInterval)
	logger.Printf("Check Predecessor Interval: %d\n", checkPredInterval)
	logger.Printf("No of Successors in Successor Table: %d\n", maxSuccessors)
	logger.Printf("No of Fingers in Finger Table: %d\n", maxFingers)

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

	for {

	}
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

}

package control

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// createUuid
//
//	Will generated a unique uuid for the node upon node creation
func createUuid() string {
	return uuid.New().String()
}

// retrieveFromDataDictionary
// Will retrieve the key from the dictionary if it exists
func retrieveFromDataDictionary(key *string) (data Data) {

	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].DataKey == *key {
			return controller.Data.DataLocations[i]
		}
	}

	return
}

// addToDataDictionary
//
//	Will add the data structure to the dictionary, or update the location
func addToDataDictionary(dataToInsert Data) (new, updated bool) {

	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].DataKey == dataToInsert.DataKey {
			// already exists, so check if the data matches
			controller.Data.DataLocations[i] = dataToInsert

			updated = true
			break
		}
	}

	// if the data doesn't exist already
	if !updated {
		controller.Data.DataLocations = append(controller.Data.DataLocations, dataToInsert)
		new = true
	}

	// Update the collective with this new data - just this data though
	if data, err := json.Marshal(dataToInsert); err != nil {
		log.Println(err)
		return false, false
	} else {
		UpdateCollective(&data)
	}
	return
}

// determineIpAddress
//
//	Will get the ip address that makes this instance recognizable
//	Will determine if this node is inside of a k8s service and will get the service name if applicable
//	Will allow for the ip address to be configurable as well upon service initialization
func determineIpAddress() string {

	// Check if the env variable is set for the IP address
	envIp := os.Getenv("COLLECTIVE_HOST_URL")
	if envIp != "" {
		return envIp
	}

	envIp = os.Getenv("COLLECTIVE_IP")

	resolverFile := "/etc/resolv.conf"
	// Check if resolver file is provided
	envResolverFile := os.Getenv("COLLECTIVE_RESOLVER_FILE")
	if envResolverFile != "" {
		resolverFile = envResolverFile
	}

	// function to return the local ip address for the network
	discoverLocalIp := func() string {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	// determine if this pod is running in k8s
	checkIfK8s := func() (bool, string) {
		readfile, err := os.Open(resolverFile)
		if err != nil {
			return false, ""
		}
		searchRegexp := regexp.MustCompile(`^\s*search\s*(([^\s]+\s*)*)$`)
		fileScanner := bufio.NewScanner(readfile)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {

			match := searchRegexp.FindSubmatch([]byte(fileScanner.Text()))
			if match != nil && strings.Contains(string(match[1]), "svc.cluster.local") {
				matchedStrings := strings.Split(string(match[1]), " ")
				return true, string(matchedStrings[0])
			}
		}
		return false, ""
	}

	// This pod has a search dns route for k8s, compose the dns route for the pod
	if isInK8s, svcValue := checkIfK8s(); isInK8s {

		// set as kubernetes being active
		controller.KubeDeployed = true

		// Need to get the $HOSTNAME env variable here, then create the connection to it
		// 	split by the periods to get the namespace - default.svc.cluster.local
		// 	convert into the dns route for a pod:
		// 		pod-ip-address.my-namespace.pod.cluster-domain.example
		//		10-42-0-180.default.pod.cluster.local

		// discover the local ip address
		localIpAddress := envIp
		if localIpAddress == "" {
			localIpAddress = discoverLocalIp()
		}

		// format the ip and convert svc to pod
		formattedIp := strings.Replace(localIpAddress, ".", "-", -1)
		formattedDnsRoute := strings.Replace(svcValue, ".svc.", ".pod.", -1)

		dnsLocalIp := fmt.Sprintf("%s.%s", formattedIp, formattedDnsRoute)
		return dnsLocalIp
	}

	return discoverLocalIp()
}

// syncData
//
//	is responsible for syncing data between nodes once an application starts up
func syncData() (err error) {
	// TODO: Need to add all of the logic into here
	return nil
}

// removeNode
//
//	when adding a node to an existing replica group, remove a node that holds the data but is not part of the replication group
func removeNode(replicationGroup int) (nodeRemoved Node, err error) {
	// Get the replica group with nodes from collective
	replicatedNodesForGroup := []Node{}
	for i := range controller.CollectiveNodes {
		if controller.CollectiveNodes[i].ReplicaNum == replicationGroup {
			replicatedNodesForGroup = controller.CollectiveNodes[i].ReplicaNodes
			break
		}
	}

	if replicatedNodesForGroup == nil {
		return Node{}, errors.New("replication group doesn't exist")
	}

	found := false
	// Determine node to remove data from
	for i := range controller.Data.DataLocations {
		// Check to see if the replica group number matches the one we are looking for
		if controller.Data.DataLocations[i].ReplicaNodeGroup == replicationGroup {
			// Cycle through and see if the data has a node id that doesn't match with the replica nodes
			for j := range controller.Data.DataLocations[i].ReplicatedNodeIds {
				// Set boolean if it matches the current group
				matchesCurrentGroup := false

				// Go through nodes for the replica group and compare
				for rg := range replicatedNodesForGroup {
					if controller.Data.DataLocations[i].ReplicatedNodeIds[j] == replicatedNodesForGroup[rg].NodeId {
						matchesCurrentGroup = true
					}
				}

				// If the node wasn't discovered, then set it as the discovered node
				if !matchesCurrentGroup {
					nodeRemoved.NodeId = controller.Data.DataLocations[i].ReplicatedNodeIds[j]
				}
			}
		}
		if found {
			break
		}
	}

	if nodeRemoved.NodeId == "" {
		return Node{}, errors.New("there is no node to be removed")
	}

	// Cycle through the collective group to determine which node we are removing
	found = false
	for i := range controller.CollectiveNodes {
		for rg := range controller.CollectiveNodes[i].ReplicaNodes {
			// Does the collective node replica group node match the node id we are attempting to remove
			if controller.CollectiveNodes[i].ReplicaNodes[rg].NodeId == nodeRemoved.NodeId {
				// Update so that we have the IP to send removal of data requests to
				nodeRemoved = controller.CollectiveNodes[i].ReplicaNodes[rg]
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if nodeRemoved.IpAddress == "" {
		// FIXME: In the future, if the node no longer exists then we should do some cleanup in the data dictionary
		return Node{}, errors.New("this node no longer exists")
	}

	var callToRemove = func(data []Data) {
		// TODO: API Call to remove all of these data entries to the IP of the remove node
	}

	// Go through and remove all of the data that this node has for this replica group
	dataToRemove := []Data{}
	for i := range controller.Data.DataLocations {
		// Check if the data is in this replication group and should be removed
		if replicationGroup == controller.Data.DataLocations[i].ReplicaNodeGroup {
			for j := range controller.Data.DataLocations[i].ReplicatedNodeIds {
				if controller.Data.DataLocations[i].ReplicatedNodeIds[j] == nodeRemoved.NodeId {
					dataToRemove = append(dataToRemove, controller.Data.DataLocations[i])
				}
			}
		}
		if len(dataToRemove) > 50 {
			go callToRemove(dataToRemove)
			dataToRemove = []Data{}
		}
	}
	go callToRemove(dataToRemove)

	return nodeRemoved, nil
}

// DetermineReplicas
//
//	Will determine the replicas for this new node
func determineReplicas() (err error) {

	replicaCount := 1
	// Get the environment variable on the wanted replica count
	if rc := os.Getenv("COLLECTIVE_REPLICA_COUNT"); rc != "" {
		// Convert env variable to number
		if replicaCount, err = strconv.Atoi(rc); err != nil {
			return err
		}
	}
	// Scale OUT Algorithm
	//
	// 	OPEN replica group?
	// 		YES - Add to group, pull data, remove data from expired replica node (check data and update data dictionary)
	// 		NO - Create new group
	// 			 Pull replica % from random nodes up to total replica count (rc)
	// 			 	eg., 33% from 3 nodes for replica count of 3
	// 			 Update for data location
	// 			 For each new replica added, remove data from 1 node per replica added
	// 				eg., search through the data dictionary for replica group and remove the replica nodes

	// Cycle through the collective replica groups and determine if there is a new group
	for _, rg := range controller.CollectiveNodes {
		// if it is not a full replica group
		if !rg.FullGroup {
			// Determine if there are more nodes in the group than the replica count
			if len(rg.ReplicaNodes) > replicaCount {
				// We don't want to add to a replica group that is already oversized
				return errors.New("too many nodes in replica currently")
			}

			// Update the node that we are part of the group now
			// apiCall that will go through `ReplicateRequest` on another node

			// Update this controller data with the replication group
			controller.ReplicaNodeId = rg.ReplicaNum
			controller.ReplicaNodes = rg.ReplicaNodes
			controller.ReplicaNodes = append(controller.ReplicaNodes, Node{
				NodeId:    controller.NodeId,
				IpAddress: controller.IpAddress,
			})
			controller.ReplicaNodeIds = append(controller.ReplicaNodeIds, controller.NodeId)

			// remove node not in replication group
			go removeNode(rg.ReplicaNum)

			// start pulling in the data required
			go syncData()

			// do not go through process of pulling in new data
			return
		}
	}

	return nil
}

// terminateReplicas
//
//	Is responsible for alerting when termination is starting, and sending data to another replica group
func terminateReplicas() (err error) {

	// Scale IN Algorithm
	//
	// 	FULL Group?
	// 		YES - Send data to randomized replica group
	// 			  Add nodes in replica group to replica node list in data dictionary for that entry
	// 			  Update Data Dictionary and CollectiveNodes
	// 		NO - Determine the replica group data is being sent to already
	// 			 Send more data to that group
	// 			 Update Data Dictionary and CollectiveNodes

	// replicaCount := 1
	// // Get the environment variable on the wanted replica count
	// if rc := os.Getenv("COLLECTIVE_REPLICA_COUNT"); rc != "" {
	// 	// Convert env variable to number
	// 	if replicaCount, err = strconv.Atoi(rc); err != nil {
	// 		return err
	// 	}
	// }

	// // Set the seed for the random number generator
	// rand.Seed(time.Now().UnixNano())

	// if len(controller.CollectiveNodes) < replicaCount {
	// 	return errors.New("replica count exceeds node count")
	// }

	// // Randomly assign strings from the slice
	// for i := 0; i < replicaCount; i++ {
	// 	// Generate a random index
	// 	randIndex := rand.Intn(len(controller.CollectiveNodes))

	// 	// Print the string at the random index
	// 	controller.ReplicaNodes = append(controller.ReplicaNodes, controller.CollectiveNodes[randIndex])
	// }
	// return nil

	return nil
}

func distributeData(key, bucket *string, data *[]byte) error {

	if *key == "" || *bucket == "" {
		return errors.New("invalid parameters")
	}

	// Add this node to the DataDictionary
	addToDataDictionary(Data{
		ReplicaNodeGroup:  controller.ReplicaNodeId,
		DataKey:           *key,
		Database:          *bucket,
		ReplicatedNodeIds: controller.ReplicaNodeIds,
	})

	// Fire off DataDictionary update process through the collective

	return nil
}

package control

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/insomniadev/collective-db/api/client"
	"github.com/insomniadev/collective-db/api/proto"
)

var (
	replicaCount int
	err          error
)

func init() {
	if replicaCountString := os.Getenv("COLLECTIVE_REPLICA_COUNT"); replicaCountString != "" {
		if replicaCount, err = strconv.Atoi(os.Getenv("COLLECTIVE_REPLICA_COUNT")); err != nil {
			log.Fatal(err)
		}
	} else {
		replicaCount = 1
	}
}

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
func addToDataDictionary(dataToInsert Data) (updateType int) {
	collectiveMemoryMutex.Lock()

	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].DataKey == dataToInsert.DataKey {
			// already exists, so check if the data matches
			controller.Data.DataLocations[i] = dataToInsert

			// Unlock and return
			collectiveMemoryMutex.Unlock()
			return UPDATE
		}
	}

	// if the data doesn't exist already
	controller.Data.DataLocations = append(controller.Data.DataLocations, dataToInsert)

	// Unlock and return
	collectiveMemoryMutex.Unlock()
	return NEW
}

// collectiveUpdate
//
// This will go through and update the collective memory, not touching the actual data
func collectiveUpdate(update *DataUpdate) {
	collectiveMemoryMutex.Lock()

	// If this is a data update
	if update.DataUpdate.Update {
		// Update the data dictionary
		switch update.DataUpdate.UpdateType {
		case NEW:
			// Adds the new element to the end of the array
			controller.Data.DataLocations = append(controller.Data.DataLocations, update.DataUpdate.UpdateData)
		case UPDATE:
			// Updates the element where it is
			for i := range controller.Data.DataLocations {
				if controller.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					controller.Data.DataLocations[i] = update.DataUpdate.UpdateData
					break
				}
			}
		case DELETE:
			// Deletes the element from the array
			for i := range controller.Data.DataLocations {
				if controller.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					controller.Data.DataLocations = removeFromDictionarySlice(controller.Data.DataLocations, i)
					break
				}
			}
		}
	} else if update.ReplicaUpdate.Update {
		// Update the data dictionary
		switch update.ReplicaUpdate.UpdateType {
		case NEW:
			// Adds the new element to the end of the array
			controller.Data.CollectiveNodes = append(controller.Data.CollectiveNodes, update.ReplicaUpdate.UpdateReplica)
		case UPDATE:
			// Updates the element where it is
			for i := range controller.Data.CollectiveNodes {
				if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					controller.Data.CollectiveNodes[i] = update.ReplicaUpdate.UpdateReplica
					break
				}
			}
		case DELETE:
			// Deletes the element from the array
			for i := range controller.Data.CollectiveNodes {
				if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					controller.Data.CollectiveNodes = removeFromDictionarySlice(controller.Data.CollectiveNodes, i)
					break
				}
			}
		}
	}

	collectiveMemoryMutex.Unlock()
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
	// Check if there are more than one replica nodes in this group, since in the previous call we set the node in the group
	// Check that the first member of the group is not this node, if it is then there is no point requesting data
	if len(controller.ReplicaNodes) > 1 && controller.ReplicaNodes[0].NodeId != controller.NodeId {

		// Make a wait group and a channel for the returned data
		dataToStore := make(chan *proto.Data)
		var wg sync.WaitGroup

		wg.Add(1)
		go func(store chan *proto.Data) {
			defer wg.Done()
			for {
				data := <-store
				if data != nil {
					// Store data in the database and do not attempt to distribute
					go storeDataInDatabase(&data.Key, &data.Database, &data.Data, true)
				} else {
					return
				}
			}

		}(dataToStore)

		if active {
			client.SyncDataRequest(&controller.ReplicaNodes[0].IpAddress, dataToStore)
		}
	}
	return nil
}

// removeNode
//
//	when adding a node to an existing replica group, remove a node that holds the data but is not part of the replication group
func removeNode(replicationGroup int) (nodeRemoved Node, err error) {
	// Get the replica group with nodes from collective
	replicatedNodesForGroup := []Node{}
	for i := range controller.Data.CollectiveNodes {
		if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == replicationGroup {
			replicatedNodesForGroup = controller.Data.CollectiveNodes[i].ReplicaNodes
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
			// If length of replicatedNodeIds is greater than the replica count
			if len(controller.Data.DataLocations[i].ReplicatedNodeIds) > replicaCount {
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
						found = true
						break
					}
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
	for i := range controller.Data.CollectiveNodes {
		for rg := range controller.Data.CollectiveNodes[i].ReplicaNodes {
			// Does the collective node replica group node match the node id we are attempting to remove
			if controller.Data.CollectiveNodes[i].ReplicaNodes[rg].NodeId == nodeRemoved.NodeId {
				// Update so that we have the IP to send removal of data requests to
				nodeRemoved = controller.Data.CollectiveNodes[i].ReplicaNodes[rg]
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
		// Remove all of these data entries to the IP of the remove node
		protoData := proto.DataArray{}
		for i := range data {
			protoData.Data = append(protoData.Data, &proto.Data{
				Key:      data[i].DataKey,
				Database: data[i].Database,
			})
		}
		if active {
			client.DeleteData(&nodeRemoved.IpAddress, &protoData)
		}

		// Dictionary update to remove from dictionary groups
		// Create the data update request object
		var wg sync.WaitGroup
		updateDictionary := make(chan *proto.DataUpdates)

		// Call the dictionary function before going through all of the data elements
		// send the update to the first node in the list
		client.DictionaryUpdate(&controller.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, updateDictionary)

		wg.Add(1)
		go func(data []Data) {
			defer wg.Done()
			for i := range data {
				for j := range controller.Data.DataLocations {
					if controller.Data.DataLocations[j].DataKey == data[i].DataKey {
						newListOfIps := []string{}
						for k := range data[i].ReplicatedNodeIds {
							if data[i].ReplicatedNodeIds[k] != nodeRemoved.NodeId {
								newListOfIps = append(newListOfIps, data[i].ReplicatedNodeIds[k])
							}
						}
						updateDictionary <- &proto.DataUpdates{
							CollectiveUpdate: &proto.CollectiveDataUpdate{
								Update:     true,
								UpdateType: UPDATE,
								Data: &proto.CollectiveData{
									ReplicaNodeGroup:  int32(data[i].ReplicaNodeGroup),
									DataKey:           data[i].DataKey,
									Database:          data[i].Database,
									ReplicatedNodeIds: newListOfIps,
								},
							},
						}
					}
				}
			}
			updateDictionary <- nil
		}(data)

		wg.Wait()
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
			if active {
				go callToRemove(dataToRemove)
			}
			dataToRemove = []Data{}
		}
	}
	if active {
		go callToRemove(dataToRemove)
	}

	return nodeRemoved, nil
}

// DetermineReplicas
//
//	Will determine the replicas for this new node
func determineReplicas() (err error) {

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
	for _, rg := range controller.Data.CollectiveNodes {
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
			controller.ReplicaNodeGroup = rg.ReplicaNodeGroup
			controller.ReplicaNodes = rg.ReplicaNodes
			controller.ReplicaNodes = append(controller.ReplicaNodes, Node{
				NodeId:    controller.NodeId,
				IpAddress: controller.IpAddress,
			})
			controller.ReplicaNodeIds = append(controller.ReplicaNodeIds, controller.NodeId)

			// remove node not in replication group
			go removeNode(rg.ReplicaNodeGroup)

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

	newData := Data{
		ReplicaNodeGroup:  controller.ReplicaNodeGroup,
		DataKey:           *key,
		Database:          *bucket,
		ReplicatedNodeIds: controller.ReplicaNodeIds,
	}

	// TODO: API - send data update to ReplicaStoreData - ReplicaUpdate rpc

	// Add this node to the DataDictionary
	updateType := addToDataDictionary(newData)

	// Always start with the first node so that data gets updated consequtively
	updateData := DataUpdate{
		DataUpdate: CollectiveDataUpdate{
			Update:     true,
			UpdateType: updateType,
			UpdateData: newData,
		},
	}
	log.Println(updateData)

	if active {
		// Fire off DataDictionary update process through the collective - DictionaryUpdate rpc
		// Create the data update request object
		updateDictionary := make(chan *proto.DataUpdates)

		// Call the dictionary function before passing the data into the channel
		// send the update to the first node in the list
		client.DictionaryUpdate(&controller.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, updateDictionary)

		updateDictionary <- &proto.DataUpdates{
			CollectiveUpdate: &proto.CollectiveDataUpdate{
				Update:     true,
				UpdateType: int32(updateType),
				Data: &proto.CollectiveData{
					ReplicaNodeGroup:  int32(newData.ReplicaNodeGroup),
					DataKey:           newData.DataKey,
					Database:          newData.Database,
					ReplicatedNodeIds: newData.ReplicatedNodeIds,
				},
			},
		}
		updateDictionary <- nil
	}

	return nil
}

// removeFromDictionarySlice
// removes the specified index from the slice and returns that slice, this does reorder the array by switching out the elements
func removeFromDictionarySlice[T collective](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

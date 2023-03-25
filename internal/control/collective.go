package control

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/insomniadev/collective-db/internal/node"
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/client"
	"github.com/insomniadev/collective-db/internal/types"
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
func retrieveFromDataDictionary(key *string) (data types.Data) {

	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].DataKey == *key {
			return node.Collective.Data.DataLocations[i]
		}
	}

	return
}

// addToDataDictionary
//
//	Will add the data structure to the dictionary, or update the location
func addToDataDictionary(dataToInsert types.Data) (updateType int) {
	node.CollectiveMemoryMutex.Lock()

	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].DataKey == dataToInsert.DataKey {
			// already exists, so check if the data matches
			node.Collective.Data.DataLocations[i] = dataToInsert

			// Unlock and return
			node.CollectiveMemoryMutex.Unlock()
			return types.UPDATE
		}
	}

	// if the data doesn't exist already
	node.Collective.Data.DataLocations = append(node.Collective.Data.DataLocations, dataToInsert)

	// Unlock and return
	node.CollectiveMemoryMutex.Unlock()
	return types.NEW
}

// collectiveUpdate
//
// This will go through and update the collective memory, not touching the actual data
func collectiveUpdate(update *types.DataUpdate) {
	node.CollectiveMemoryMutex.Lock()

	// If this is a data update
	if update.DataUpdate.Update {
		// Update the data dictionary
		switch update.DataUpdate.UpdateType {
		case types.NEW:
			// Adds the new element to the end of the array
			node.Collective.Data.DataLocations = append(node.Collective.Data.DataLocations, update.DataUpdate.UpdateData)
		case types.UPDATE:
			// Updates the element where it is
			for i := range node.Collective.Data.DataLocations {
				if node.Collective.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					node.Collective.Data.DataLocations[i] = update.DataUpdate.UpdateData
					break
				}
			}
		case types.DELETE:
			// Deletes the element from the array
			for i := range node.Collective.Data.DataLocations {
				if node.Collective.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					node.Collective.Data.DataLocations = removeFromDictionarySlice(node.Collective.Data.DataLocations, i)
					break
				}
			}
		}
	} else if update.ReplicaUpdate.Update {
		// Update the data dictionary
		switch update.ReplicaUpdate.UpdateType {
		case types.NEW:
			// Adds the new element to the end of the array
			node.Collective.Data.CollectiveNodes = append(node.Collective.Data.CollectiveNodes, update.ReplicaUpdate.UpdateReplica)
		case types.UPDATE:
			// Updates the element where it is
			for i := range node.Collective.Data.CollectiveNodes {
				if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					node.Collective.Data.CollectiveNodes[i] = update.ReplicaUpdate.UpdateReplica
					break
				}
			}
		case types.DELETE:
			// Deletes the element from the array
			for i := range node.Collective.Data.CollectiveNodes {
				if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					node.Collective.Data.CollectiveNodes = removeFromDictionarySlice(node.Collective.Data.CollectiveNodes, i)
					break
				}
			}
		}
	}

	node.CollectiveMemoryMutex.Unlock()
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

	// FIXME: Can't we just take the provided master node IP and then have that send us back what our IP is?

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
		return localAddr.IP.String() + ":9090"
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
		node.Collective.KubeServiceDns = svcValue + ":9090"
		log.Printf("Discovered to be part of a kubernetes service: %s", node.Collective.KubeServiceDns)

		// set as kubernetes being active
		node.Collective.KubeDeployed = true

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
		return dnsLocalIp + ":9090"
	}

	return discoverLocalIp()
}

// retrieveDataDictionary
//
//	gets the master node and pulls in the data dictionary to start the application up
func retrieveDataDictionary() {

	// extracting the function for pulling and updating the data from the other functions so that there is only one implementation
	storeIncomingData := func(ipAddress *string) (bool, error) {
		dataToStore := make(chan *types.DataUpdate)
		if err := client.SyncCollectiveRequest(ipAddress, dataToStore); err != nil {

			for {
				storeDictionaryData := <-dataToStore
				if storeDictionaryData != nil {
					// Store the collective data now
					collectiveUpdate(storeDictionaryData)
				} else {
					// we are now done processing
					return true, nil
				}
			}
		} else {
			log.Println(err)
			return false, err
		}
	}

	// Create a new collective cluster
	createNewCollective := func() {
		node.Collective.Data.CollectiveNodes = []types.ReplicaGroup{
			{
				ReplicaNodeGroup:   1,
				SecondaryNodeGroup: 0,
				ReplicaNodes: []types.Node{
					{
						NodeId:    node.Collective.NodeId,
						IpAddress: node.Collective.IpAddress,
					},
				},
				FullGroup: false,
			},
		}
		node.Collective.ReplicaNodes = []types.Node{
			{
				NodeId:    node.Collective.NodeId,
				IpAddress: node.Collective.IpAddress,
			},
		}
	}

	collectiveBrokers := []string{}
	collectiveMainBrokers := os.Getenv("COLLECTIVE_MAIN_BROKERS")
	if collectiveMainBrokers != "" {
		collectiveBrokers = strings.Split(collectiveMainBrokers, ",")
	}
	log.Println(collectiveBrokers)

	// if COLLECTIVE_MAIN_BROKERS is populated, then use first
	if len(collectiveBrokers) > 0 {
		// Cycle through brokers to pull in the dictionary data - need to create all of the functions
		successfullyProcessed := false
		for i := range collectiveBrokers {
			if stored, err := storeIncomingData(&collectiveBrokers[i]); err != nil && stored {
				successfullyProcessed = true
				log.Println("Successfully initialized into a Collective cluster")
				break
			}
		}
		// If we weren't able to pull successfully, then let's panic and kill the application
		if !successfullyProcessed {
			panic(fmt.Sprintf("not able to pull data from the COLLECTIVE_MAIN_BROKERS, %s", collectiveMainBrokers))
		}

	} else if node.Collective.KubeDeployed {
		// Pull data from any other pod that is in this current service group
		if stored, err := storeIncomingData(&node.Collective.KubeServiceDns); err != nil && stored {
			log.Println("Successfully initialized into a Kubernetes service group within a cluster")
		} else {
			log.Println("Current Kubernetes service group does not have an existing Collective cluster, creating new Collective")
			createNewCollective()
		}
	} else {
		createNewCollective()
	}
}

// syncData
//
//	is responsible for syncing data between nodes once an application starts up
func syncData() (err error) {
	// Check if there are more than one replica nodes in this group, since in the previous call we set the node in the group
	// Check that the first member of the group is not this node, if it is then there is no point requesting data
	if len(node.Collective.ReplicaNodes) > 1 && node.Collective.ReplicaNodes[0].NodeId != node.Collective.NodeId {

		// Make a wait group and a channel for the returned data
		dataToStore := make(chan *proto.Data)

		go func(store chan *proto.Data) {
			for {
				data := <-store
				if data != nil {
					// Store data in the database and do not attempt to distribute
					go storeDataInDatabase(&data.Key, &data.Database, &data.Data, true, 0)
				} else {
					return
				}
			}

		}(dataToStore)

		if node.Active {
			client.SyncDataRequest(&node.Collective.ReplicaNodes[0].IpAddress, dataToStore)
		}
	}
	return nil
}

// removeDataFromSecondaryNodeGroup
//
//	when a replica group becomes a full group, then remove all of that replica data from the secondaryNodeGroup
func removeDataFromSecondaryNodeGroup(secondaryGroup int) error {

	// get the ipAddress for the leader of the secondary node group
	ipAddress := ""
	for i := range node.Collective.Data.CollectiveNodes {
		if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == secondaryGroup {
			// This ipAddress is very important, if done incorrectly we can suffer massive data loss
			ipAddress = node.Collective.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress
		}
	}

	// create the channel to start deleting the data
	deleteData := make(chan *proto.Data)
	if err := client.DeleteData(&ipAddress, deleteData); err != nil {
		return err
	}

	// cycle through the data for the entries that are for this secondaryGroup
	for i := range node.Collective.Data.DataLocations {
		// IF the data is set for this replicaGroup
		if node.Collective.Data.DataLocations[i].ReplicaNodeGroup == node.Collective.ReplicaNodeGroup {
			// THEN send a deletion command to the leader node for the replicationGroup that is the secondaryNodeGroup
			deleteData <- &proto.Data{
				Key:      node.Collective.Data.DataLocations[i].DataKey,
				Database: node.Collective.Data.DataLocations[i].Database,
			}
		}
	}
	return nil
}

// DetermineReplicas
//
//	Will determine the replicas for this new node
func determineReplicas() (err error) {

	// Scale OUT Algorithm
	//
	// 	OPEN replica group?
	// 		YES - Add to group, pull data, remove data from secondaryNodeGroup if now a full group (check data and update data dictionary)
	// 		NO - Create new group
	// 				Set the secondaryNodeGroup immediately

	// If there is only one node in the cluster, then there is nothing to work through, so just return
	if len(node.Collective.Data.CollectiveNodes) == 1 && len(node.Collective.Data.CollectiveNodes[0].ReplicaNodes) == 1 {
		return nil
	}

	// Cycle through the collective replica groups and determine if there is a new group
	for _, rg := range node.Collective.Data.CollectiveNodes {
		// if it is not a full replica group
		if !rg.FullGroup {
			// Determine if there are more nodes in the group than the replica count
			if len(rg.ReplicaNodes) > replicaCount {
				// We don't want to add to a replica group that is already oversized
				return errors.New("too many nodes in replica currently")
			}

			// Update this controller data with the replication group
			node.Collective.ReplicaNodeGroup = rg.ReplicaNodeGroup
			node.Collective.ReplicaNodes = rg.ReplicaNodes
			node.Collective.ReplicaNodes = append(node.Collective.ReplicaNodes, types.Node{
				NodeId:    node.Collective.NodeId,
				IpAddress: node.Collective.IpAddress,
			})
			node.Collective.ReplicaNodeIds = append(node.Collective.ReplicaNodeIds, node.Collective.NodeId)

			// If we currently match the amount of replicas expected, then let's set this as a full group
			// 		set the secondaryNodeGroup to 0 IF this is a full group, else keep the current group
			fullGroup := false
			if len(node.Collective.ReplicaNodeIds) == replicaCount {
				fullGroup = true

				// Remove the secondaryNodeGroup from this replica group if it is a full group
				rg.SecondaryNodeGroup = 0

				// Clean up - remove this replica group data from the secondaryNodeGroup
				// since we are now a full group, then delete all of the data from the secondary node group
				if err := removeDataFromSecondaryNodeGroup(rg.ReplicaNodeGroup); err != nil {
					return err
				}
			}

			replicaNodesInSync := []*proto.ReplicaNodes{}
			// Assemble the nodes apart from this node that is being removed
			for i := range node.Collective.ReplicaNodes {
				if node.Collective.ReplicaNodes[i].NodeId != node.Collective.NodeId {
					replicaNodesInSync = append(replicaNodesInSync, &proto.ReplicaNodes{
						NodeId:    node.Collective.ReplicaNodes[i].NodeId,
						IpAddress: node.Collective.ReplicaNodes[i].IpAddress,
					})
				}
			}

			// Update the DataDictionary that this node is now part of the collective
			if err := sendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
				ReplicaUpdate: &proto.CollectiveReplicaUpdate{
					Update:     true,
					UpdateType: types.UPDATE,
					UpdateReplica: &proto.UpdateReplica{
						ReplicaNodeGroup:   int32(node.Collective.ReplicaNodeGroup),
						FullGroup:          fullGroup,
						ReplicaNodes:       replicaNodesInSync,
						SecondaryNodeGroup: int32(rg.SecondaryNodeGroup),
					},
				},
			}); err != nil {
				return err
			}

			// start pulling in the data required
			go syncData()

			// do not go through process of pulling in new data
			return
		}
	}

	// If all groups are full then we should immediately create a new group,
	// 		the new group needs to have the secondaryNodeGroup set
	randIndex := rand.Intn(len(node.Collective.Data.CollectiveNodes))
	if err := sendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
		ReplicaUpdate: &proto.CollectiveReplicaUpdate{
			Update:     true,
			UpdateType: types.NEW,
			UpdateReplica: &proto.UpdateReplica{
				ReplicaNodeGroup: int32(node.Collective.ReplicaNodeGroup),
				FullGroup:        false,
				ReplicaNodes: []*proto.ReplicaNodes{{
					NodeId:    node.Collective.NodeId,
					IpAddress: node.Collective.IpAddress,
				}},
				SecondaryNodeGroup: int32(node.Collective.Data.CollectiveNodes[randIndex].ReplicaNodeGroup),
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// terminateReplicas
//
// When this node shuts down, this function will ensure that there is no data loss and will offload data to other nodes if required
func terminateReplicas() (err error) {

	// We only want this functionality to run IF this is currently a full replicaGroup and will no longer be full
	// 		distribute the existing data into the newly created secondaryNodeGroup
	if len(node.Collective.ReplicaNodeIds) == replicaCount {

		// Generate random index to send all of the data to
		randIndex := rand.Intn(len(node.Collective.Data.CollectiveNodes))

		// Insert the data into the intended node so that there is no data loss
		disperseData := make(chan *proto.Data)
		client.DataUpdate(&node.Collective.Data.CollectiveNodes[randIndex].ReplicaNodes[0].IpAddress, disperseData)

		// Create channel to retrieve all of the stored data through the retrieveAllReplicaData function
		replicaData := make(chan *types.StoredData)
		retrieveAllReplicaData(replicaData)
		for {
			if storedData := <-replicaData; storedData != nil {

				// Send the data to the new decided node with the secondaryNodeGroup populated
				disperseData <- &proto.Data{
					Key:                storedData.DataKey,
					Database:           storedData.Database,
					Data:               storedData.Data,
					SecondaryNodeGroup: int32(node.Collective.Data.CollectiveNodes[randIndex].ReplicaNodeGroup),
				}

			} else {
				disperseData <- nil
				break
			}
		}

		replicaNodesInSync := []*proto.ReplicaNodes{}
		// Assemble the nodes apart from this node that is being removed
		for i := range node.Collective.ReplicaNodes {
			if node.Collective.ReplicaNodes[i].NodeId != node.Collective.NodeId {
				replicaNodesInSync = append(replicaNodesInSync, &proto.ReplicaNodes{
					NodeId:    node.Collective.ReplicaNodes[i].NodeId,
					IpAddress: node.Collective.ReplicaNodes[i].IpAddress,
				})
			}
		}

		// Do a collective update to set the secondaryNodeGroup
		// 		Call the dictionary function before passing the data into the channel
		// 		send the update to the first node in the list - the master node
		if err := sendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
			ReplicaUpdate: &proto.CollectiveReplicaUpdate{
				Update:     true,
				UpdateType: types.UPDATE,
				UpdateReplica: &proto.UpdateReplica{
					ReplicaNodeGroup:   int32(node.Collective.ReplicaNodeGroup),
					FullGroup:          false,
					ReplicaNodes:       replicaNodesInSync,
					SecondaryNodeGroup: int32(node.Collective.Data.CollectiveNodes[randIndex].ReplicaNodeGroup),
				},
			},
		}); err != nil {
			return err
		}
	}

	// Determine if this is the last node removed from a replica node group and have all the data belong to the secondaryNodeGroup now
	// 		This should include a Dictionary update
	// Delete this replica group from the collective
	if len(node.Collective.ReplicaNodes) == 1 && node.Collective.ReplicaNodes[0].NodeId == node.Collective.NodeId {
		// Update all of the DataLocations to have the secondaryNodeGroup now
		updateDictionary := make(chan *proto.DataUpdates)
		client.DictionaryUpdate(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, updateDictionary)
		for i := range node.Collective.Data.DataLocations {
			// IF the data is for this replicaNodeGroup
			if node.Collective.Data.DataLocations[i].ReplicaNodeGroup == node.Collective.ReplicaNodeGroup {
				// THEN change the replicaNodeGroup to the secondaryNodeGroup
				updateDictionary <- &proto.DataUpdates{
					CollectiveUpdate: &proto.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.UPDATE,
						Data: &proto.CollectiveData{
							ReplicaNodeGroup: int32(node.Collective.SecondaryNodeGroup),
							DataKey:          node.Collective.Data.DataLocations[i].DataKey,
							Database:         node.Collective.Data.DataLocations[i].Database,
						},
					},
				}
			}
		}
		updateDictionary <- nil

		// Remove this node from the collective database
		if err := sendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
			ReplicaUpdate: &proto.CollectiveReplicaUpdate{
				Update:     true,
				UpdateType: types.DELETE,
				UpdateReplica: &proto.UpdateReplica{
					ReplicaNodeGroup:   int32(node.Collective.ReplicaNodeGroup),
					FullGroup:          false,
					ReplicaNodes:       []*proto.ReplicaNodes{},
					SecondaryNodeGroup: 0,
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func distributeData(key, bucket *string, data *[]byte, secondaryNodeGroup int) error {

	if *key == "" || *bucket == "" {
		return errors.New("invalid parameters")
	}

	newData := types.Data{
		ReplicaNodeGroup: node.Collective.ReplicaNodeGroup,
		DataKey:          *key,
		Database:         *bucket,
	}

	if node.Active {
		// Create the data object to be sent
		dataUpdate := &proto.Data{
			Key:                *key,
			Database:           *bucket,
			Data:               *data,
			ReplicaNodeGroup:   int32(node.Collective.ReplicaNodeGroup),
			SecondaryNodeGroup: int32(secondaryNodeGroup),
		}

		// Send to each replica attached to this replica node group
		for i := range node.Collective.ReplicaNodes {
			if node.Collective.ReplicaNodes[i].NodeId != node.Collective.NodeId {
				updateReplica := make(chan *proto.Data)
				client.ReplicaDataUpdate(&node.Collective.ReplicaNodes[i].IpAddress, updateReplica)
				updateReplica <- dataUpdate
				updateReplica <- nil
			}
		}

		// Double check that the secondaryNodeGroup is 0 before starting to process
		if secondaryNodeGroup != 0 {
			for i := range node.Collective.Data.CollectiveNodes {
				if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == secondaryNodeGroup {
					for j := range node.Collective.Data.CollectiveNodes[i].ReplicaNodes {
						updateReplica := make(chan *proto.Data)
						client.ReplicaDataUpdate(&node.Collective.Data.CollectiveNodes[j].ReplicaNodes[j].IpAddress, updateReplica)
						updateReplica <- dataUpdate
						updateReplica <- nil
					}

					// IF this replicaGroup is not complete and has a secondaryNodeGroup, THEN forward to all nodes in that group as well
					if !node.Collective.Data.CollectiveNodes[i].FullGroup {
						for j := range node.Collective.Data.CollectiveNodes {
							if node.Collective.Data.CollectiveNodes[j].ReplicaNodeGroup == node.Collective.Data.CollectiveNodes[i].SecondaryNodeGroup {
								// Send the update to the first node of that replica to start the update process from there
								dataUpdate.SecondaryNodeGroup = int32(node.Collective.Data.CollectiveNodes[i].SecondaryNodeGroup)

								updateReplica := make(chan *proto.Data)
								client.ReplicaDataUpdate(&node.Collective.Data.CollectiveNodes[j].ReplicaNodes[0].IpAddress, updateReplica)
								updateReplica <- dataUpdate
								updateReplica <- nil
								break
							}
						}
						break
					}
				}
			}

			// Only update the data dictionary with this data if it was not sent here as part of the secondaryNodeGroup
			if secondaryNodeGroup != node.Collective.ReplicaNodeGroup {

				// Add this node to the DataDictionary
				updateType := addToDataDictionary(newData)

				if err := sendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
					CollectiveUpdate: &proto.CollectiveDataUpdate{
						Update:     true,
						UpdateType: int32(updateType),
						Data: &proto.CollectiveData{
							ReplicaNodeGroup: int32(newData.ReplicaNodeGroup),
							DataKey:          newData.DataKey,
							Database:         newData.Database,
						},
					},
				}); err != nil {
					return err
				}
			}
		}

	}

	return nil
}

// removeFromDictionarySlice
// removes the specified index from the slice and returns that slice, this does reorder the array by switching out the elements
func removeFromDictionarySlice[T types.Collective](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// sendClientUpdateDictionaryRequest
//
// extracted function that is used to send the update without all of the additional boilerplate code everywhere
func sendClientUpdateDictionaryRequest(ipAddress *string, update *proto.DataUpdates) error {
	// TODO: Add unit test

	// Create the channel
	updateDictionary := make(chan *proto.DataUpdates)

	// Call the dictionary function before passing the data into the channel
	// send the update to the first node in the list
	go client.DictionaryUpdate(ipAddress, updateDictionary)

	updateDictionary <- update
	updateDictionary <- nil
	return nil
}

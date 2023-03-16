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

		if active {
			client.SyncDataRequest(&controller.ReplicaNodes[0].IpAddress, dataToStore)
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
	for i := range controller.Data.CollectiveNodes {
		if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == secondaryGroup {
			// This ipAddress is very important, if done incorrectly we can suffer massive data loss
			ipAddress = controller.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress
		}
	}

	// create the channel to start deleting the data
	deleteData := make(chan *proto.Data)
	if err := client.DeleteData(&ipAddress, deleteData); err != nil {
		return err
	}

	// cycle through the data for the entries that are for this secondaryGroup
	for i := range controller.Data.DataLocations {
		// IF the data is set for this replicaGroup
		if controller.Data.DataLocations[i].ReplicaNodeGroup == controller.ReplicaNodeGroup {
			// THEN send a deletion command to the leader node for the replicationGroup that is the secondaryNodeGroup
			deleteData <- &proto.Data{
				Key:      controller.Data.DataLocations[i].DataKey,
				Database: controller.Data.DataLocations[i].Database,
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

			// Update this controller data with the replication group
			controller.ReplicaNodeGroup = rg.ReplicaNodeGroup
			controller.ReplicaNodes = rg.ReplicaNodes
			controller.ReplicaNodes = append(controller.ReplicaNodes, Node{
				NodeId:    controller.NodeId,
				IpAddress: controller.IpAddress,
			})
			controller.ReplicaNodeIds = append(controller.ReplicaNodeIds, controller.NodeId)

			// If we currently match the amount of replicas expected, then let's set this as a full group
			// 		set the secondaryNodeGroup to 0 IF this is a full group, else keep the current group
			fullGroup := false
			if len(controller.ReplicaNodeIds) == replicaCount {
				fullGroup = true

				// Remove the secondaryNodeGroup from this replica group if it is a full group
				rg.SecondaryNodeGroup = 0

				// TODO: Clean up - remove this replica group data from the secondaryNodeGroup
				// since we are now a full group, then delete all of the data from the secondary node group
				go removeDataFromSecondaryNodeGroup(rg.ReplicaNodeGroup)
			}

			replicaNodesInSync := []*proto.ReplicaNodes{}
			// Assemble the nodes apart from this node that is being removed
			for i := range controller.ReplicaNodes {
				if controller.ReplicaNodes[i].NodeId != controller.NodeId {
					replicaNodesInSync = append(replicaNodesInSync, &proto.ReplicaNodes{
						NodeId:    controller.ReplicaNodes[i].NodeId,
						IpAddress: controller.ReplicaNodes[i].IpAddress,
					})
				}
			}

			// Update the DataDictionary that this node is now part of the collective
			sendClientUpdateDictionaryRequest(&controller.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
				ReplicaUpdate: &proto.CollectiveReplicaUpdate{
					Update:     true,
					UpdateType: UPDATE,
					UpdateReplica: &proto.UpdateReplica{
						ReplicaNodeGroup:   int32(controller.ReplicaNodeGroup),
						FullGroup:          fullGroup,
						ReplicaNodes:       replicaNodesInSync,
						SecondaryNodeGroup: int32(rg.SecondaryNodeGroup),
					},
				},
			})

			// start pulling in the data required
			go syncData()

			// do not go through process of pulling in new data
			return
		}
	}

	// TODO: If all groups are full then we should immediately create a new group, the new group needs to have the secondaryNodeGroup set

	// TODO: Need to do the NO part of this check, it is currently missing

	return nil
}

// terminateReplicas
//
// When this node shuts down, this function will ensure that there is no data loss and will offload data to other nodes if required
func terminateReplicas() (err error) {

	// TODO: Need to build out this functionality

	// FIXME: In the future we want to have the data evenly spread across the currently active nodes rather than just one

	// We only want this functionality to run IF this is currently a full replicaGroup and will no longer be full
	// 		distribute the existing data into the newly created secondaryNodeGroup
	if len(controller.ReplicaNodeIds) == replicaCount {

		// Generate random index to send all of the data to
		randIndex := rand.Intn(len(controller.Data.CollectiveNodes))

		// Insert the data into the intended node so that there is no data loss
		disperseData := make(chan *proto.Data)
		client.DataUpdate(&controller.Data.CollectiveNodes[randIndex].ReplicaNodes[0].IpAddress, disperseData)

		// Create channel to retrieve all of the stored data through the retrieveAllReplicaData function
		replicaData := make(chan *StoredData)
		retrieveAllReplicaData(replicaData)
		for {
			if storedData := <-replicaData; storedData != nil {

				// Send the data to the new decided node with the secondaryNodeGroup populated
				disperseData <- &proto.Data{
					Key:                storedData.DataKey,
					Database:           storedData.Database,
					Data:               storedData.Data,
					SecondaryNodeGroup: int32(controller.Data.CollectiveNodes[randIndex].ReplicaNodeGroup),
				}

			} else {
				disperseData <- nil
				break
			}
		}

		replicaNodesInSync := []*proto.ReplicaNodes{}
		// Assemble the nodes apart from this node that is being removed
		for i := range controller.ReplicaNodes {
			if controller.ReplicaNodes[i].NodeId != controller.NodeId {
				replicaNodesInSync = append(replicaNodesInSync, &proto.ReplicaNodes{
					NodeId:    controller.ReplicaNodes[i].NodeId,
					IpAddress: controller.ReplicaNodes[i].IpAddress,
				})
			}
		}

		// Do a collective update to set the secondaryNodeGroup
		// 		Call the dictionary function before passing the data into the channel
		// 		send the update to the first node in the list - the master node
		sendClientUpdateDictionaryRequest(&controller.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
			ReplicaUpdate: &proto.CollectiveReplicaUpdate{
				Update:     true,
				UpdateType: UPDATE,
				UpdateReplica: &proto.UpdateReplica{
					ReplicaNodeGroup:   int32(controller.ReplicaNodeGroup),
					FullGroup:          false,
					ReplicaNodes:       replicaNodesInSync,
					SecondaryNodeGroup: int32(controller.Data.CollectiveNodes[randIndex].ReplicaNodeGroup),
				},
			},
		})
	}

	// TODO: Determine if this is the last node removed from a replica node group and have all the data belong to the secondaryNodeGroup now
	// 		This should include a Dictionary update

	return nil
}

func distributeData(key, bucket *string, data *[]byte, secondaryNodeGroup int) error {

	if *key == "" || *bucket == "" {
		return errors.New("invalid parameters")
	}

	newData := Data{
		ReplicaNodeGroup:  controller.ReplicaNodeGroup,
		DataKey:           *key,
		Database:          *bucket,
		ReplicatedNodeIds: controller.ReplicaNodeIds,
	}

	if active {
		// Create the data object to be sent
		dataUpdate := &proto.Data{
			Key:                *key,
			Database:           *bucket,
			Data:               *data,
			ReplicaNodeGroup:   int32(controller.ReplicaNodeGroup),
			SecondaryNodeGroup: int32(secondaryNodeGroup),
		}

		// Send to each replica attached to this replica node group
		for i := range controller.ReplicaNodes {
			if controller.ReplicaNodes[i].NodeId != controller.NodeId {
				updateReplica := make(chan *proto.Data)
				client.ReplicaDataUpdate(&controller.ReplicaNodes[i].IpAddress, updateReplica)
				updateReplica <- dataUpdate
				updateReplica <- nil
			}
		}

		// Double check that the secondaryNodeGroup is 0 before starting to process
		if secondaryNodeGroup != 0 {
			for i := range controller.Data.CollectiveNodes {
				if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == secondaryNodeGroup {
					for j := range controller.Data.CollectiveNodes[i].ReplicaNodes {
						updateReplica := make(chan *proto.Data)
						client.ReplicaDataUpdate(&controller.Data.CollectiveNodes[j].ReplicaNodes[j].IpAddress, updateReplica)
						updateReplica <- dataUpdate
						updateReplica <- nil
					}

					// IF this replicaGroup is not complete and has a secondaryNodeGroup, THEN forward to all nodes in that group as well
					if !controller.Data.CollectiveNodes[i].FullGroup {
						for j := range controller.Data.CollectiveNodes {
							if controller.Data.CollectiveNodes[j].ReplicaNodeGroup == controller.Data.CollectiveNodes[i].SecondaryNodeGroup {
								// Send the update to the first node of that replica to start the update process from there
								dataUpdate.SecondaryNodeGroup = int32(controller.Data.CollectiveNodes[i].SecondaryNodeGroup)

								updateReplica := make(chan *proto.Data)
								client.ReplicaDataUpdate(&controller.Data.CollectiveNodes[j].ReplicaNodes[0].IpAddress, updateReplica)
								updateReplica <- dataUpdate
								updateReplica <- nil
								break
							}
						}
						break
					}
				}
			}
		}

		// Only update the data dictionary with this data if it was not sent here as part of the secondaryNodeGroup
		if secondaryNodeGroup != controller.ReplicaNodeGroup {

			// Add this node to the DataDictionary
			updateType := addToDataDictionary(newData)

			sendClientUpdateDictionaryRequest(&controller.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
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
			})
		}
	}

	return nil
}

// removeFromDictionarySlice
// removes the specified index from the slice and returns that slice, this does reorder the array by switching out the elements
func removeFromDictionarySlice[T collective](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// sendClientUpdateDictionaryRequest
//
// extracted function that is used to send the update without all of the additional boilerplate code everywhere
func sendClientUpdateDictionaryRequest(ipAddress *string, update *proto.DataUpdates) {
	// TODO: Add unit test

	// Create the channel
	updateDictionary := make(chan *proto.DataUpdates)

	// Call the dictionary function before passing the data into the channel
	// send the update to the first node in the list
	client.DictionaryUpdate(ipAddress, updateDictionary)

	updateDictionary <- update
	updateDictionary <- nil
}

package control

// Thoughts: for the node IP it could be <IP_ADDRESS>/node?<NODE_ID>
import (
	"log"
)

// Pull from local database, if doesn't exist then
//
//	Create node id
//	Get IP Address
//	Determine replica nodes
//	Get Node List
func init() {

	// Allow at some point for the node to start back up and begin an update task for the data
	// nodeData := "node"
	// if exists, value := database.Get(&nodeData, &nodeData); exists {
	// 	if err := json.Unmarshal(*value, &controller); err != nil {
	// 		log.Fatal("Failed to parse the configuration data")
	// 	}

	// determine if the replica still
	// update and refresh data

	// return if this is the correct group, if the group no longer exists, then start this as a new collective
	// 	return
	// }

	controller.NodeId = createUuid()

	// Utilizes Environment variables:
	//	COLLECTIVE_HOST_URL - will set this as it's IP address with no additional logic
	// 	COLLECTIVE_IP - will use this IP but still configure for K8S
	// 	COLLECTIVE_RESOLVER_FILE - will override default /etc/resolv.conf file
	controller.IpAddress = determineIpAddress()

	// Will assign replicas to this node
	determineReplicas()

}

// IsActive
//
//	Returns a confirmation on if this node is currently active and processing
//
// THOUGHTS: If this server is up then it should be running, should this be where it has been synced with other nodes?
func IsActive() bool {
	return active
}

// Deactivate
//
//	Will deactivate the node, redistribute leaders, and send data if needed
func Deactivate() bool {
	return false
}

// NodeInfo
//
//	Returns info on this node
func NodeInfo() *Controller {
	return &controller
}

// CollectiveUpdate
//
//	Will update this node with the incoming collective information from the other nodes
func CollectiveUpdate(update *DataUpdate) {
	// Update this node with the incoming information
	// Send the data to the first url in the next replica group
	// Send the new data to all replicas in this replica group

	// If we get data here, then this is the leader of the replica group (or the first node in the array)
	collectiveUpdate(update)

	// Send the data onward
	// Update the replicas in this replica group
	for i := range controller.ReplicaNodes {
		// Let's not send to ourselves here
		if controller.ReplicaNodes[i].NodeId != controller.NodeId {
			// Send the data to the replica node

			// TODO: API - Send the data to the api endpoint for the ReplicaUpdate function - ReplicaUpdate rpc
			log.Println(controller.ReplicaNodes[i].IpAddress)
		}
	}

	// Send to the next replica group in the list
	for i := range controller.Data.CollectiveNodes {
		// Go to where this node group is in the array
		if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == controller.ReplicaNodeGroup {
			// Send to the next replica group in the list
			// 	Check that there is another element in the array
			if len(controller.Data.CollectiveNodes) >= i+2 {

				// TODO: API - Send the data to the next node - DictionaryUpdate rpc
				log.Println(controller.Data.CollectiveNodes[i+1].ReplicaNodes[0].IpAddress)
			}
		}
	}
}

// ReplicateRequest
//
//	Node requesting to join this replica group
func ReplicateRequest() {
	// Respond with success or failure
}

// ReplicaUpdate
//
//	Will update this node with the data coming from another replica related to collective data
//	This update call will not attempt to continue distributing the update
func ReplicaUpdate(update *DataUpdate) {

	// Update the collective with the new information
	collectiveUpdate(update)

}

// ReplicaStoreData
//
//	Will store the data provided from another replica and not update DataDictionary or attempt to replicate
func ReplicaStoreData(key, bucket string, data []byte) {
	storeData(&key, &bucket, &data, true)
}

// storeData
//
//	Will store the provided data on this node
func storeData(key, bucket *string, data *[]byte, replicaStore bool) (bool, *string) {
	return storeDataInDatabase(key, bucket, data, replicaStore)
}

// RetrieveData
//
//	Will return the requested data to the calling application or node
func RetrieveData(key, bucket *string) (bool, *[]byte) {
	return retrieveDataFromDatabase(key, bucket)
}

// DeleteData
//
//	Delete data will remove the data from the database
func DeleteData(key, bucket *string) (bool, error) {
	return deleteDataFromDatabase(key, bucket)
}

// var mu sync.Mutex
// mu.Lock()
// defer mu.Unlock()

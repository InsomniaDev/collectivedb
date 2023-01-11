package control

// Thoughts: for the node IP it could be <IP_ADDRESS>/node?<NODE_ID>
import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	database "github.com/insomniadev/collective-db/internal/database"
)

// Main Controller struct
type Controller struct {
	NodeId         string         `json:"nodeId"`             // UUID of this node within the NodeList
	IpAddress      string         `json:"ipAddress"`          // IpAddress of this node
	KubeDeployed   bool           `json:"kubernetesDeployed"` // This app is deployed in kubernetes
	ReplicaNodeId  int            `json:"replicaNodeId"`      // The replica node id
	ReplicaNodeIds []string       `json:"replicaNodeIds"`     // Replica node ids for distributing traffic
	ReplicaNodes   []Node         `json:"replicaNodes"`       // Replica nodes of this node
	Data           DataDictionary `json:"data"`               // Location of all the keys to nodes
}

// ReplicaGroup
type ReplicaGroup struct {
	ReplicaNodeGroup int    `json:"replicaNodeGroup"`
	ReplicaNodes     []Node `json:"nodes"`
	FullGroup        bool   `json:"fullGroup"`
}

// Node struct
type Node struct {
	NodeId    string `json:"nodeId"`
	IpAddress string `json:"ipAddress"`
}

// DataDictionary struct
type DataDictionary struct {
	DataLocations   []Data         `json:"data"`
	CollectiveNodes []ReplicaGroup `json:"collectiveNodes"` // List of node IPs that are connected to the collective
}

// Data struct
type Data struct {
	ReplicaNodeGroup  int      `json:"replicaNodeGroup"`
	DataKey           string   `json:"dataKey"`
	Database          string   `json:"database"`
	ReplicatedNodeIds []string `json:"replicatedNodes"`
}

var (
	active     bool
	controller Controller
)

// Pull from local database, if doesn't exist then
//
//	Create node id
//	Get IP Address
//	Determine replica nodes
//	Get Node List
func init() {

	nodeData := "node"
	if exists, value := database.Get(&nodeData, &nodeData); exists {
		if err := json.Unmarshal(*value, &controller); err != nil {
			log.Fatal("Failed to parse the configuration data")
		}

		// TODO: determine if the replica still
		// TODO: update and refresh data

		// return if this is the correct group, if the group no longer exists, then start this as a new collective
		return
	}

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

// UpdateCollective
//
//	Update node list with the provided data
func UpdateCollective(data *[]byte) {
	// Get the first IP in the first replication group in the collective list and send the updated information to that one
}

// NodeUpdate
//
//	Will update this node with the incoming information from the other nodes
func NodeUpdate() {
	// Update this node with the incoming information
	// Send the data to the first url in the next replica group
	// Send the new data to all replicas in this replica group
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
func ReplicaUpdate() {

}

// ReplicaStoreData
//
//	Will store the data provided from another replica and not update DataDictionary or attempt to replicate
func ReplicaStoreData(key, bucket string, data []byte) {
	StoreData(&key, &bucket, &data, true)
}

// StoreData
//
//	Will store the provided data on this node
func StoreData(key, bucket *string, data *[]byte, replicaStore bool) (bool, *string) {

	ackLevel := os.Getenv("COLLECTIVE_ACK_LEVEL")
	if ackLevel == "" {
		ackLevel = "NONE"
	}

	// Data is new and doesn't exist
	// Create a unique key and update since this is new data
	if *key == "" {
		newKey := uuid.New().String()
		key = &newKey

		updated, key := database.Update(key, bucket, data)

		if !replicaStore {
			switch ackLevel {
			case "ALL":
				distributeData(key, bucket, data)
			case "NONE":
				go distributeData(key, bucket, data)
			}
		}

		return updated, key
	}

	// This data exists already
	// TODO: Determine what node the data is on, if the data does exist on a node
	dataVolume := retrieveFromDataDictionary(key)

	// If the data doesn't exist yet, but a key was provided
	if dataVolume.DataKey == "" {
		updated, key := database.Update(key, bucket, data)

		if !replicaStore {
			switch ackLevel {
			case "ALL":
				distributeData(key, bucket, data)
			case "NONE":
				go distributeData(key, bucket, data)
			}
		}

		return updated, key
	}

	// TODO: Send update command to all nodes in replica group to store the data
	// TODO: Wait for response of confirmation if that env variable is set

	return false, nil
}

// RetrieveData
//
//	Will return the requested data to the calling application or node
func RetrieveData(key, bucket *string) (bool, *[]byte) {
	if exists, value := database.Get(key, bucket); exists {
		return exists, value
	}

	// The data does not exist on this node
	// TODO: Determine what node the data exists on
	// TODO: Go retrieve the data and then return it here

	// If data is not found on one replica node, it should attempt to pull from at least two nodes before declaring data doesn't exist
	// This should attempt to grab data from the ReplicatedNodes list

	return false, nil
}

// DeleteData
//
//	Delete data will remove the data from the database
func DeleteData(key, bucket *string) (bool, error) {
	if deleted, err := database.Delete(key, bucket); deleted {
		return true, err
	} else if err != nil {
		return false, err
	}

	// The data does not exist on this node
	// TODO: Determine what node the data exists on
	// TODO: Send a delete command to all replicas to delete the data
	// TODO: Wait for response of confirmation if that env variable is set

	return false, nil
}

// UpdateReplicas
//
//	Will update the replicas on the data changes
//
// TODO: This might need to go somewhere else
func UpdateReplicas(urls []string, hashedKey int) {

	// Use a mutex to synchronize access to the list of URLs.
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// for replicaNode := range

	// Send the data to the URL with the next highest hash.
	// url := urlMap[nextHash]
	// sendData(url, data)
}

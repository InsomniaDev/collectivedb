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
	NodeId          string `json:"nodeId"`             // UUID of this node within the NodeList
	IpAddress       string `json:"ipAddress"`          // IpAddress of this node
	KubeDeployed    bool   `json:"kubernetesDeployed"` // This app is deployed in kubernetes
	ReplicaNodes    []Node `json:"replicaNodes"`       // Replica nodes of this node
	CollectiveNodes []Node `json:"collectiveNodes"`    // List of node IPs that are connected to the collective
}

// Node struct
type Node struct {
	Hash       string `json:"hash"`
	IpAddress  string `json:"ipAddress"`
	NodeId     string `json:"nodeId"`
	LeaderNode bool   `json:"leader"` // Do I need to have a leader here?
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
	} else {
		controller.NodeId = createUuid()

		// Utilizes Environment variables:
		//	COLLECTIVE_HOST_URL - will set this as it's IP address with no additional logic
		// 	COLLECTIVE_IP - will use this IP but still configure for K8S
		// 	COLLECTIVE_RESOLVER_FILE - will override default /etc/resolv.conf file
		controller.IpAddress = determineIpAddress()
	}
}

// startNode
//
//	Is responsible for starting the node up, syncing data
func startNode() {

	// Check if there are nodes to sync with
	if controller.CollectiveNodes == nil {
		// Get the node id to sync with and populate
		controller = findNodeLeader()
	}

	// Determine which replica group to fit into
	// Pull data from replica group
	// Update current database if needed
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

// SyncReplicas
//
//	Will determine replicas that need to be synced
func SyncReplicas() {}

// SyncNodeList
//
//	Update node list with the provided data
func SyncNodeList() {}

// StoreData
//
//	Will store the provided data on this node
func StoreData(key, bucket *string, data *[]byte) (bool, *string) {
	ackLevel := os.Getenv("COLLECTIVE_ACK_LEVEL")
	if ackLevel == "" {
		ackLevel = "NONE"
	}

	distributeData := func(key, bucket *string, data *[]byte) {

	}

	// Create a unique key and update since this is new data
	if *key == "" {
		newKey := uuid.New().String()
		key = &newKey

		updated, key := database.Update(key, bucket, data)

		switch ackLevel {
		case "ALL":
			distributeData(key, bucket, data)
		case "NONE":
			go distributeData(key, bucket, data)
		}
		return updated, key
	}

	// This data exists already
	// TODO: Determine what node the data is on
	// TODO: Send update command to all nodes to store the data
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

	return false, nil
}

// DeleteData
//
//	Delete data will remove the data from the database
func DeleteData(key, bucket *string) (bool, error) {
	if deleted, err := database.Delete(key, bucket); deleted {
		return true, err
	}

	// The data does not exist on this node
	// TODO: Determine what node the data exists on
	// TODO: Send a delete command to all replicas to delete the data
	// TODO: Wait for response of confirmation if that env variable is set

	return true, nil
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

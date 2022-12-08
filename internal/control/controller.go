package control

// Thoughts: for the node IP it could be <IP_ADDRESS>/node?<NODE_ID>
import (
	"encoding/json"
	"log"

	database "github.com/insomniadev/collective-db/internal/database"
)

// Main Controller struct
type Controller struct {
	NodeId          string `json:"nodeId"`          // UUID of this node within the NodeList
	IpAddress       string `json:"ipAddress"`       // IpAddress of this node
	ReplicaNodes    []Node `json:"replicaNodes"`    // Replica nodes of this node
	CollectiveNodes []Node `json:"collectiveNodes"` // List of node IPs that are connected to the collective
}

// Node struct
type Node struct {
	IpAddress  string `json:"ipAddress"`
	NodeId     string `json:"nodeId"`
	LeaderNode bool   `json:"leader"`
}

var controller Controller

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

// Active
//
//	Returns a confirmation on if this node is currently active and processing
//
// THOUGHTS: If this server is up then it should be running, should this be where it has been synced with other nodes?
func Active() {}

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
func StoreData() {}

// RetrieveData
//
//	Will return the requested data to the calling application or node
func RetrieveData() {}

// UpdateReplicas
//
//	Will update the replicas on the data changes
//
// TODO: This might need to go somewhere else
func UpdateReplicas() {}

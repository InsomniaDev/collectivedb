package control

// Thoughts: for the node IP it could be <IP_ADDRESS>/node?<NODE_ID>

// Main Controller struct
type Controller struct {
	NodeId       string `json:"nodeId"`       // UUID of this node within the NodeList
	IpAddress    string `json:"ipAddress"`    // IpAddress of this node
	ReplicaNodes string `json:"replicaNodes"` // Replica nodes of this node
	NodeList     []Node `json:"nodeList"`     // List of node IPs that are connected to the collective
}

// Node struct
type Node struct {
	IpAddress      string
	NodeId         string
	ReplicaNodeIds []string
}

// Active
// 		Returns a confirmation on if this node is currently active and processing
func Active() {}

// NodeInfo
// 		Returns info on this node
func NodeInfo() {}

// SyncNodeList
// 		Update node list with the provided data
func SyncNodeList() {}

// StoreData
// 		Will store the provided data on this node
func StoreData() {}

// RetrieveData
// 		Will return the requested data to the calling application or node
func RetrieveData() {}

// UpdateReplicas
// 		Will update the replicas on the data changes
// TODO: This might need to go somewhere else
func UpdateReplicas() {}

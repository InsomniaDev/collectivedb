package control

import "sync"

// Main Controller struct
type Controller struct {
	NodeId           string         `json:"nodeId"`             // UUID of this node within the NodeList
	IpAddress        string         `json:"ipAddress"`          // IpAddress of this node
	KubeDeployed     bool           `json:"kubernetesDeployed"` // This app is deployed in kubernetes
	ReplicaNodeGroup int            `json:"ReplicaNodeGroup"`   // The replica node id
	ReplicaNodeIds   []string       `json:"replicaNodeIds"`     // Replica node ids for distributing traffic
	ReplicaNodes     []Node         `json:"replicaNodes"`       // Replica nodes of this node
	Data             DataDictionary `json:"data"`               // Location of all the keys to nodes
}

type collective interface {
	ReplicaGroup | Data
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

	collectiveMemoryMutex sync.Mutex
)

// Update Data Type
type DataUpdate struct {
	DataUpdate    CollectiveDataUpdate    `json:"dataUpdate"`
	ReplicaUpdate CollectiveReplicaUpdate `json:"replicaUpdate"`
}
type CollectiveDataUpdate struct {
	Update     bool `json:"update"`
	UpdateType int  `json:"updateType"`
	UpdateData Data `json:"data"`
}
type CollectiveReplicaUpdate struct {
	Update        bool         `json:"update"`
	UpdateType    int          `json:"updateType"`
	UpdateReplica ReplicaGroup `json:"data"`
}

// CollectiveUpdateEnumeration
const (
	NEW    = 1
	UPDATE = 2
	DELETE = 3
)

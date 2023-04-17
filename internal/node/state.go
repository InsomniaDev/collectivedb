package node

import (
	"sync"

	"github.com/insomniadev/collective-db/internal/types"
)

// TODO: Set an overall environment level scenario here

var (
	Active bool // sets if this current node is active

	Collective types.Controller // holds the data for the collective

	CollectiveMemoryMutex sync.RWMutex
)

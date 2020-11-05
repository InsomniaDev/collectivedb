package database

// Database is the top level main main for this data
type Database struct {
	TotalDepth int
	memData    data
}

// data is used for each level in the datastructure to determine which hash to go into
type data struct {
	hashKey  string
	depth    int
	values   []entry
	hashData []data
}

// entry is the actual key and value stored in memory
type entry struct {
	key   string
	value string
}

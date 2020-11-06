package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/insomniadev/collective-db/database"
)

func main() {

	// Retrieve the depth that the hash function should extend to
	depth := os.Getenv("DEPTH")
	num, err := strconv.Atoi(depth)
	if err != nil {
		// Non numeric environment variable input, setting to 1
		fmt.Println("Depth incorrectly input, setting to 1")
		num = 1
	}

	fmt.Println(database.FNV32a("test"))

	db := database.Init(num)
	fmt.Println(db.TotalDepth)
}

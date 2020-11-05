package main

import (
	"fmt"

	"github.com/insomniadev/collective-db/database"
)

func main() {
	resp := database.Init(1)
	fmt.Println(resp.TotalDepth)
}

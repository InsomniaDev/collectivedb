package database

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

// How to split up the depth?

// Init will initialize the database and return a Database with the total depth set
func Init(layers int) (db Database) {
	db = Database{
		TotalDepth: layers,
	}
	return
}

// FNV32a will convert the string argument into a numerical hash
func FNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

// hashNum will take the proferred string and hash it
func hashNum(key string, round int, hashS string) (int, string) {
	if round == 0 {
		return 0, hashS
	}
	if hashS != "" {
		hashStr := string([]rune(hashS)[:round])
		hashArrInt, err := strconv.Atoi(hashStr)
		if err != nil {
			fmt.Println(err)
		}
		return hashArrInt, hashS
	}
	hashS = fmt.Sprint(FNV32a(key))
	hashStr := string([]rune(hashS)[:round])
	hashArrInt, err := strconv.Atoi(hashStr)
	if err != nil {
		fmt.Println(err)
	}
	return hashArrInt, hashS
}

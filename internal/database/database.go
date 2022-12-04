package database

import (
	"time"

	"github.com/InsomniaDev/bolt"
)

func Insert(key, value *string) (inserted bool, insertedKey *string) {
	if *key == "" || *value == "" {
		return false, key
	}

	if db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	

	return false, key
}

func Edit(key, value *string) (updated bool, updatedKey *string) {
	if *key == "" || *value == "" {
		return false, key
	}

	return false, key
}

func Delete(key *string) (deleted bool) {
	if *key == "" {
		return false
	}

	return false
}

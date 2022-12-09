package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/InsomniaDev/bolt"
	log "github.com/sirupsen/logrus"
)

var (
	// list of open database connections, each bucket translates to a database
	connections map[*string]*bolt.DB
)

func init() {
	connections = make(map[*string]*bolt.DB)
}

// getDatabase
// 		will return the open database if exists or create and return
func getDatabase(bucket *string) *bolt.DB {
	for connection := range connections {
		if strings.Compare(*connection, *bucket) == 0 {
			return connections[connection]
		}
	}

	connection, err := bolt.Open(fmt.Sprintf("%s.db", *bucket), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	connections[bucket] = connection
	return connection
}

// CloseConnections
// 		Responsible for safely closing all open database connections upon shutdown
func CloseConnections() {
	for connection := range connections {
		connections[connection].Close()
	}
	connections = make(map[*string]*bolt.DB)
}

// Update
// 		Responsible for updating the provided key and value in the connected database
func Update(key, bucket *string, value *[]byte) (updated bool, updatedKey *string) {
	if *key == "" || string(*value) == "" {
		return false, key
	}

	err := getDatabase(bucket).Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(*bucket))
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(*key), *value)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return false, key
	}
	return true, key
}

// Get
// 		Responsible for retrieving values from the database based on specified key and bucket
func Get(key, bucket *string) (exists bool, value *[]byte) {
	if *key == "" {
		return false, nil
	}

	err := getDatabase(bucket).View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*bucket))
		if bucket == nil {
			return errors.New("bucket doesn't exist")
		}
		bv := bucket.Get([]byte(*key))
		value = &bv
		return nil
	})

	if err != nil {
		log.Debug("get bucket err:", bucket, err)
		return false, nil
	}

	if len(*value) == 0 {
		return false, value
	}

	return true, value
}

// Delete
// 		Will remove the specified key from the database
func Delete(key, bucket *string) (deleted bool, err error) {
	if *key == "" {
		return false, errors.New("unable to delete an empty key")
	}

	err = getDatabase(bucket).Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*bucket))
		if bucket == nil {
			log.Error("get bucket err:", bucket, err)
			return errors.New("bucket doesn't exist")
		}

		err = bucket.Delete([]byte(*key))
		return err
	})

	if err != nil {
		log.Error("get bucket err:", bucket, err)
		return false, err
	}
	return true, nil
}

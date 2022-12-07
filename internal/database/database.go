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
	connection  *bolt.DB
	connections map[*string]*bolt.DB
)

func init() {
	connections = make(map[*string]*bolt.DB)
	// TODO: Need to adjust the bolt library to allow for multiple different databases
	// TODO: Add a function to close all libraries on application shutdown
	// connection, _ = bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
}

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

func closeConnections() {
	for connection := range connections {
		connections[connection].Close()
	}
	connections = make(map[*string]*bolt.DB)
}

func Update(key, value, bucket *string) (updated bool, updatedKey *string) {
	if *key == "" || *value == "" {
		return false, key
	}

	err := getDatabase(bucket).Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(*bucket))
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(*key), []byte(*value))
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
		log.Error("get bucket err:", bucket, err)
		return false, nil
	}
	return true, value
}

func Delete(key, bucket *string) (deleted bool, err error) {
	if *key == "" {
		return false, errors.New("unable to delete an empty key")
	}

	err = getDatabase(bucket).Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*bucket))
		if bucket == nil {
			log.Error("get bucket err:", bucket, err)
			return err
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

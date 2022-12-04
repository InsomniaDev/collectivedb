package database

import (
	"errors"
	"time"

	"github.com/InsomniaDev/bolt"
	log "github.com/sirupsen/logrus"
)

var connection *bolt.DB

func init() {
	connection, _ = bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
}

func Update(key, value, bucket *string) (updated bool, updatedKey *string) {
	if *key == "" || *value == "" {
		return false, key
	}

	err := connection.Update(func(tx *bolt.Tx) error {
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

	err := connection.View(func(tx *bolt.Tx) error {
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

	err = connection.Update(func(tx *bolt.Tx) error {
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

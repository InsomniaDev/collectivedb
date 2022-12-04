package database

import (
	"errors"
	"time"

	"github.com/InsomniaDev/bolt"
	log "github.com/sirupsen/logrus"
)

func Update(key, value, bucket *string) (updated bool, updatedKey *string) {
	if *key == "" || *value == "" {
		return false, key
	}

	if db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second}); err == nil {

		err := db.Update(func(tx *bolt.Tx) error {
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
	} else {
		log.Fatal(err)
	}

	return false, key
}

func Get(key, bucket *string) (exists bool, value *[]byte) {
	if *key == "" {
		return false, nil
	}

	if db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second}); err == nil {
		var value []byte
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(*bucket))
			if bucket == nil {
				log.Error("get bucket err:", bucket, err)
				return err
			}

			value = bucket.Get([]byte(*key))
			return nil
		})

		if err != nil {
			log.Error("get bucket err:", bucket, err)
			return false, nil
		}
		return true, &value
	} else {
		log.Fatal(err)
	}

	return false, nil
}

func Delete(key, bucket *string) (deleted bool, err error) {
	if *key == "" {
		return false, errors.New("unable to delete an empty key")
	}

	if db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second}); err == nil {
		err = db.Update(func(tx *bolt.Tx) error {
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
	} else {
		log.Fatal(err)
		return false, err
	}
}

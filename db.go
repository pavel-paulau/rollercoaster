package main

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	dbName     = "data/rollercoaster.db"
	bucketName = "benchmarks"

	db *bolt.DB
)

func open() *bolt.DB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}

	return db
}

func initBucket() {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func put(value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		id, _ := b.NextSequence()
		return b.Put(itob(id), value)
	})
}

func iter(values chan []byte) {
	defer close(values)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		b.ForEach(func(key, value []byte) error {
			values <- value
			return nil
		})
		return nil
	})
}

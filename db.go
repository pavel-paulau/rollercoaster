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

type kvPair struct {
	key   uint64
	value []byte
}

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

func put(id uint64, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		if id == 0 {
			id, _ = b.NextSequence()
		}
		return b.Put(itob(id), value)
	})
}

func iter(values chan kvPair) {
	defer close(values)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		b.ForEach(func(key, value []byte) error {
			values <- kvPair{binary.BigEndian.Uint64(key), value}
			return nil
		})
		return nil
	})
}

func del(key uint64) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		return b.Delete(itob(key))
	})
}

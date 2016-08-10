package main

import (
	"fmt"
)

type benchmark struct {
	Group     string  `json:"group" binding:"required"`
	ID        uint64  `json:"id"`
	Metric    string  `json:"metric" binding:"required"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value" binding:"required"`
}

func greeting() {
	fmt.Println("\n\t.:: Please navigate to http://127.0.0.1:8080/ ::.\n")
}

func main() {
	db = open()
	initBucket()
	defer db.Close()

	greeting()

	httpEngine().Run()
}

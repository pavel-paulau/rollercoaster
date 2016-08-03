package main

type benchmark struct {
	Group     string  `json:"group" binding:"required"`
	Metric    string  `json:"metric" binding:"required"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value" binding:"required"`
}

func main() {
	db = open()
	initBucket()
	defer db.Close()

	httpEngine().Run()
}

package main

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode("release")
}

func httpEngine() *gin.Engine {
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")
	router.Static("/static", "./static")

	v1 := router.Group("/api/v1")

	v1.GET("benchmarks", func(c *gin.Context) {
		var benchmarks []benchmark

		kvPairs := make(chan kvPair, 10)

		go iter(kvPairs)
		for p := range kvPairs {
			var b benchmark
			if err := json.Unmarshal(p.value, &b); err != nil {
				c.IndentedJSON(500, gin.H{"message": err.Error()})
				return
			}
			b.ID = p.key
			benchmarks = append(benchmarks, b)
		}

		c.IndentedJSON(200, benchmarks)
	})

	v1.POST("benchmarks", func(c *gin.Context) {
		var b benchmark
		if err := c.BindJSON(&b); err != nil {
			c.IndentedJSON(400, gin.H{"message": err.Error()})
			return
		}
		if b.Timestamp == 0 {
			b.Timestamp = time.Now().UnixNano()
		}

		value, _ := json.Marshal(b) // Ignoring errors because they are not really possible

		if err := put(b.ID, value); err != nil {
			c.IndentedJSON(500, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(201, gin.H{"message": "ok"})
	})

	v1.DELETE("benchmarks", func(c *gin.Context) {
		var payload struct {
			ID uint64 `json:"id"`
		}
		if err := c.BindJSON(&payload); err != nil {
			c.IndentedJSON(400, gin.H{"message": err.Error()})
			return
		}

		if err := del(payload.ID); err != nil {
			c.IndentedJSON(500, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(200, gin.H{"message": "ok"})
	})

	return router
}

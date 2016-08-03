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

		values := make(chan []byte, 10)

		go iter(values)
		for value := range values {
			var b benchmark
			if err := json.Unmarshal(value, &b); err != nil {
				panic(err)
			}
			benchmarks = append(benchmarks, b)
		}

		c.JSON(200, benchmarks)
	})

	v1.POST("benchmarks", func(c *gin.Context) {
		var b benchmark
		if err := c.BindJSON(&b); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
		}
		b.Timestamp = time.Now().UnixNano()

		value, err := json.Marshal(b)
		if err != nil {
			panic(err)
		}
		put(value)

		c.JSON(200, gin.H{"message": "ok"})
	})

	return router
}

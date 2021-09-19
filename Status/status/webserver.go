package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func configureRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Running": true,
		})
	})
	router.GET("/info", getInfo)

	return router
}

func getInfo(c *gin.Context) {
	jsonMap := make(map[string]interface{})
	jsonFile, err := os.Open(statusFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &jsonMap)
	if err == nil {
		c.JSON(200, jsonMap)
	}
}

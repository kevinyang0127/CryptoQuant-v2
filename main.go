package main

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/quant"
	"CryptoQuant-v2/strategy"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var platform *quant.Platform

func main() {
	log.Println("Start my crypto quant v2!!")
	mongoDB, disconnect, err := db.NewMongoDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	platform = quant.NewPlatform(mongoDB)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/strategy", addStrategy)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func addStrategy(c *gin.Context) {
	type Param struct {
		UserID       string `json:"userID"`
		Exchange     string `json:"exchange"`
		Symbol       string `json:"symbol"`
		Timeframe    string `json:"timeframe"`
		Status       int    `json:"status"`
		StrategyName string `json:"strategyName"`
		Script       string `json:"script"`
	}

	param := Param{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println(err)
		return
	}

	err = platform.AddStrategy(c, param.UserID, param.Exchange, param.Symbol, param.Timeframe,
		strategy.StrategyStatus(param.Status), param.StrategyName, param.Script)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "add strategy success",
	})
}

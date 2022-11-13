package main

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/quant"
	"CryptoQuant-v2/simulation"
	"CryptoQuant-v2/strategy"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var platform *quant.Platform
var MongoDB *db.MongoDB

func main() {
	log.Println("Start my crypto quant v2!!")
	mongoDB, disconnect, err := db.NewMongoDB(db.URI)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	MongoDB = mongoDB
	platform = quant.NewPlatform(mongoDB)
	simulation.RunNewManager(mongoDB)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/strategy", addStrategy)
	r.POST("/backtesting", backtesting)
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

func backtesting(c *gin.Context) {
	type Param struct {
		UserID              string `json:"userID"`
		Exchange            string `json:"exchange"`
		Symbol              string `json:"symbol"`
		Timeframe           string `json:"timeframe"`
		StrategyID          string `json:"strategyID"`
		StartBalance        string `json:"startBalance"`
		Lever               int    `json:"lever"`
		TakerCommissionRate string `json:"takerCommissionRate"`
		MakerCommissionRate string `json:"makerCommissionRate"`
		StartTimeMs         int64  `json:"startTimeMs"`
		EndTimeMs           int64  `json:"endTimeMs"`
	}

	param := Param{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println(err)
		return
	}

	backtestingClient := quant.NewBackTestingClient(MongoDB, param.UserID, param.StrategyID, param.Exchange, param.Symbol,
		param.Timeframe, param.StartBalance, param.Lever, param.TakerCommissionRate, param.MakerCommissionRate,
		param.StartTimeMs, param.EndTimeMs)

	simulationID, err := backtestingClient.Backtest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Backtest error",
			"err": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulationID": simulationID,
		"msg":          "Backtesting run success",
	})
}

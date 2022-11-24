package api

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/quant"
	"CryptoQuant-v2/simulation"
	"CryptoQuant-v2/strategy"
	"CryptoQuant-v2/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	r *gin.Engine
}

func (r *Router) ListenAndServ() {
	r.r.Run()
}

func NewRouter(mongoDB *db.MongoDB, platform *quant.Platform, strategyManager *strategy.Manager,
	exchangeManager *exchange.Manager, userManager *user.Manager, simulationManager *simulation.Manager) *Router {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", getRegisterHandler(userManager))
	r.POST("/binancefuture/key", getUpsertBinanceKeyHandler(userManager))
	r.POST("/strategy", getAddStrategyHandler(platform))
	r.POST("/backtesting", getBacktestingHandler(mongoDB, strategyManager, exchangeManager, userManager, simulationManager))
	return &Router{
		r: r,
	}
}

func getRegisterHandler(userManager *user.Manager) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type Param struct {
			UserID string `json:"userID"`
			Name   string `json:"name"`
		}
		param := Param{}
		err := c.BindJSON(&param)
		if err != nil {
			log.Println(err)
			return
		}

		err = userManager.Register(c, param.UserID, param.Name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "register success",
		})
	}

	return fn
}

func getUpsertBinanceKeyHandler(userManager *user.Manager) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type Param struct {
			UserID    string `json:"userID"`
			ApiKey    string `json:"apiKey"`
			SecretKey string `json:"secretKey"`
		}
		param := Param{}
		err := c.BindJSON(&param)
		if err != nil {
			log.Println(err)
			return
		}

		err = userManager.UpsertBinanceKeys(c, param.UserID, param.ApiKey, param.SecretKey)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "upsert binance key success",
		})
	}
	return fn
}

func getAddStrategyHandler(platform *quant.Platform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		strategyID, err := platform.AddStrategy(c, param.UserID, param.Exchange, param.Symbol, param.Timeframe,
			strategy.StrategyStatus(param.Status), param.StrategyName, param.Script)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"strategyID": strategyID,
			"msg":        "add strategy success",
		})
	}
	return fn
}

func getBacktestingHandler(mongoDB *db.MongoDB, strategyManager *strategy.Manager, exchangeManager *exchange.Manager,
	userManager *user.Manager, simulationManager *simulation.Manager) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type Param struct {
			UserID                string `json:"userID"`
			Exchange              string `json:"exchange"`
			Symbol                string `json:"symbol"`
			Timeframe             string `json:"timeframe"`
			KlineHistoryTimeframe string `json:"klineHistoryTimeframe"`
			StrategyID            string `json:"strategyID"`
			StartBalance          string `json:"startBalance"`
			Lever                 int    `json:"lever"`
			TakerCommissionRate   string `json:"takerCommissionRate"`
			MakerCommissionRate   string `json:"makerCommissionRate"`
			StartTimeMs           int64  `json:"startTimeMs"`
			EndTimeMs             int64  `json:"endTimeMs"`
		}

		param := Param{}
		err := c.BindJSON(&param)
		if err != nil {
			log.Println(err)
			return
		}

		backtestingClient := quant.NewBackTestingClient(mongoDB, strategyManager, exchangeManager, userManager, simulationManager, param.UserID, param.StrategyID, param.Exchange, param.Symbol,
			param.Timeframe, param.KlineHistoryTimeframe, param.StartBalance, param.Lever, param.TakerCommissionRate, param.MakerCommissionRate,
			param.StartTimeMs, param.EndTimeMs)

		simulationID, err := backtestingClient.Backtest(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
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
	return fn
}

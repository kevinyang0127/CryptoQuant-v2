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
	// r.GET("/strategy", getAddStrategyHandler(platform))
	r.POST("/strategy", getAddStrategyHandler(strategyManager))
	r.POST("/strategy/live", getRunStrategyHandler(platform))
	// r.POST("/strategy/stop", getAddStrategyHandler(platform))
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
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "request body error",
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "request body error",
			})
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

func getAddStrategyHandler(strategyManager *strategy.Manager) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type Param struct {
			UserID       string `json:"userID"`
			Exchange     string `json:"exchange"`
			Symbol       string `json:"symbol"`
			Timeframe    string `json:"timeframe"`
			StrategyName string `json:"strategyName"`
			Script       string `json:"script"`
		}

		param := Param{}
		err := c.BindJSON(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "request body error",
			})
			return
		}

		strategyID, err := strategyManager.Add(c, param.UserID, param.Exchange, param.Symbol,
			param.Timeframe, strategy.Draft, param.StrategyName, param.Script)
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
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "request body error",
			})
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

func getRunStrategyHandler(platform *quant.Platform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type Param struct {
			StrategyID string `json:"strategyID"`
		}

		param := Param{}
		err := c.BindJSON(&param)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "request body error",
			})
			return
		}

		err = platform.RunStrategy(c, param.StrategyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
				"msg": "run strategy error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"strategyID": param.StrategyID,
			"msg":        "Start live trade success",
		})
	}
	return fn
}

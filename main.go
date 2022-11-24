package main

import (
	"CryptoQuant-v2/api"
	"log"
)

type TradingQueen struct {
	router *api.Router
}

func NewTradingQueen(router *api.Router) *TradingQueen {
	return &TradingQueen{
		router: router,
	}
}

func main() {
	log.Println("Start my crypto quant v2!!")

	tradingQueen, cleanup, err := InitTradingQueen()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	tradingQueen.router.ListenAndServ() // listen and serve on 0.0.0.0:8080
}

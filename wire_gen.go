// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

// Injectors from wire.go:

func InitTradingQueen() (*TradingQueen, func(), error) {
	dburi := provideDBURI()
	mongoDB, cleanup, err := provideMongoDB(dburi)
	if err != nil {
		return nil, nil, err
	}
	manager := provideUserManager(mongoDB)
	exchangeManager := provideExchangeManager(manager)
	simulationManager := provideSimulationManager(mongoDB)
	luaScriptHandler := provideLuaScriptHandler(exchangeManager, simulationManager)
	strategyManager := provideStrategyManager(mongoDB, luaScriptHandler)
	platform := providePlatform(strategyManager, exchangeManager, manager)
	router := provideRouter(mongoDB, platform, strategyManager, exchangeManager, manager, simulationManager)
	tradingQueen := NewTradingQueen(router)
	return tradingQueen, func() {
		cleanup()
	}, nil
}

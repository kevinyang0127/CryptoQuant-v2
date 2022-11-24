package main

import (
	"CryptoQuant-v2/api"
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/quant"
	"CryptoQuant-v2/script"
	"CryptoQuant-v2/simulation"
	"CryptoQuant-v2/strategy"
	"CryptoQuant-v2/user"

	"github.com/google/wire"
)

var RouterSet = wire.NewSet(provideRouter, provideMongoDB, provideDBURI, providePlatform, provideStrategyManager,
	provideExchangeManager, provideUserManager, provideSimulationManager, provideLuaScriptHandler)

func provideRouter(mongoDB *db.MongoDB, platform *quant.Platform, strategyManager *strategy.Manager,
	exchangeManager *exchange.Manager, userManager *user.Manager, simulationManager *simulation.Manager) *api.Router {
	return api.NewRouter(mongoDB, platform, strategyManager, exchangeManager, userManager, simulationManager)
}

func provideMongoDB(uri db.DBURI) (*db.MongoDB, func(), error) {
	return db.NewMongoDB(uri)
}

func provideDBURI() db.DBURI {
	return db.URI
}

func providePlatform(strategyManager *strategy.Manager, exchangeManager *exchange.Manager, userManager *user.Manager) *quant.Platform {
	return quant.NewPlatform(strategyManager, exchangeManager, userManager)
}

func provideStrategyManager(mongoDB *db.MongoDB, luaScriptHandler *script.LuaScriptHandler) *strategy.Manager {
	return strategy.NewManager(mongoDB, luaScriptHandler)
}

func provideExchangeManager(userManager *user.Manager) *exchange.Manager {
	return exchange.NewExchangeManager(userManager)
}

func provideUserManager(mongoDB *db.MongoDB) *user.Manager {
	return user.NewUserManager(mongoDB)
}

func provideSimulationManager(mongoDB *db.MongoDB) *simulation.Manager {
	return simulation.NewSimulationManager(mongoDB)
}

func provideLuaScriptHandler(simulationManager *simulation.Manager) *script.LuaScriptHandler {
	return script.NewLuaScriptHandler(simulationManager)
}

//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitTradingQueen() (*TradingQueen, func(), error) {
	panic(wire.Build(NewTradingQueen, RouterSet))
}

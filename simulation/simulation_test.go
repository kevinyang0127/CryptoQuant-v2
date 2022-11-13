package simulation

import (
	"CryptoQuant-v2/db"
	"context"
	"log"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEntry(t *testing.T) {
	mongoDB, disconnect, err := db.NewMongoDB(db.LocalURI)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	startBalance := decimal.NewFromFloat(10000)
	lever := decimal.NewFromInt(10)
	takerCommissionRate := decimal.NewFromFloat(0.0004)
	makerCommissionRate := decimal.NewFromFloat(0.0002)

	s := NewSimulation(mongoDB, "kevinyang", startBalance, lever, takerCommissionRate, makerCommissionRate)

	price := decimal.NewFromFloat(1500)
	qty := decimal.NewFromFloat(0.5)
	s.Entry(context.TODO(), price, qty, false, time.Now().Unix())

	assert.Equal(t, price, s.positon.OpenPrice)
	assert.Equal(t, qty, s.positon.Quantity)
}

func TestLong(t *testing.T) {
	mongoDB, disconnect, err := db.NewMongoDB(db.LocalURI)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	startBalance := decimal.NewFromFloat(10000)
	lever := decimal.NewFromInt(10)
	takerCommissionRate := decimal.NewFromFloat(0.0004)
	makerCommissionRate := decimal.NewFromFloat(0.0002)

	//// win case
	s := NewSimulation(mongoDB, "kevinyang", startBalance, lever, takerCommissionRate, makerCommissionRate)

	price := decimal.NewFromFloat(1500)
	qty := decimal.NewFromFloat(0.5)
	s.Entry(context.TODO(), price, qty, false, time.Now().Unix())

	assert.Equal(t, price, s.positon.OpenPrice)
	assert.Equal(t, qty, s.positon.Quantity)

	exitPrice := decimal.NewFromFloat(1600)
	exitQty := decimal.NewFromFloat(-0.5)
	s.Exit(context.TODO(), exitPrice, exitQty, false, time.Now().Unix())

	assert.Empty(t, s.positon)
	assert.Equal(t, decimal.NewFromFloat(10049.938).String(), s.balance.String())

	///// loss case
	s2 := NewSimulation(mongoDB, "kevinyang", startBalance, lever, takerCommissionRate, makerCommissionRate)

	price = decimal.NewFromFloat(1500)
	qty = decimal.NewFromFloat(0.5)
	s2.Entry(context.TODO(), price, qty, false, time.Now().Unix())

	assert.Equal(t, price, s2.positon.OpenPrice)
	assert.Equal(t, qty, s2.positon.Quantity)

	exitPrice = decimal.NewFromFloat(1400)
	exitQty = decimal.NewFromFloat(-0.5)
	s2.Exit(context.TODO(), exitPrice, exitQty, false, time.Now().Unix())

	assert.Empty(t, s2.positon)
	assert.Equal(t, decimal.NewFromFloat(9949.942).String(), s2.balance.String())
}

func TestShort(t *testing.T) {
	mongoDB, disconnect, err := db.NewMongoDB(db.LocalURI)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	startBalance := decimal.NewFromFloat(10000)
	lever := decimal.NewFromInt(10)
	takerCommissionRate := decimal.NewFromFloat(0.0004)
	makerCommissionRate := decimal.NewFromFloat(0.0002)

	//// win case
	s := NewSimulation(mongoDB, "kevinyang", startBalance, lever, takerCommissionRate, makerCommissionRate)

	price := decimal.NewFromFloat(1500)
	qty := decimal.NewFromFloat(-0.5)
	s.Entry(context.TODO(), price, qty, false, time.Now().Unix())

	assert.Equal(t, price, s.positon.OpenPrice)
	assert.Equal(t, qty, s.positon.Quantity)

	exitPrice := decimal.NewFromFloat(1400)
	exitQty := decimal.NewFromFloat(0.5)
	s.Exit(context.TODO(), exitPrice, exitQty, false, time.Now().Unix())

	assert.Empty(t, s.positon)
	assert.Equal(t, decimal.NewFromFloat(10049.942).String(), s.balance.String())

	///// loss case
	s2 := NewSimulation(mongoDB, "kevinyang", startBalance, lever, takerCommissionRate, makerCommissionRate)

	price = decimal.NewFromFloat(1500)
	qty = decimal.NewFromFloat(-0.5)
	s2.Entry(context.TODO(), price, qty, false, time.Now().Unix())

	assert.Equal(t, price, s2.positon.OpenPrice)
	assert.Equal(t, qty, s2.positon.Quantity)

	exitPrice = decimal.NewFromFloat(1600)
	exitQty = decimal.NewFromFloat(0.5)
	s2.Exit(context.TODO(), exitPrice, exitQty, false, time.Now().Unix())

	assert.Empty(t, s2.positon)
	assert.Equal(t, decimal.NewFromFloat(9949.938).String(), s2.balance.String())
}

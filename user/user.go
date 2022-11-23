package user

type User struct {
	UserID           string `bson:"userID"`
	Name             string `bson:"name"`
	BinanceApiKey    string `bson:"binanceApiKey"`
	BinanceSecretKey string `bson:"binanceSecretKey"`
}

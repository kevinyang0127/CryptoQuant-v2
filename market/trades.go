package market

type Trades struct {
	ID            int64  `json:"id" bson:"id"`
	Price         string `json:"price" bson:"price"`
	Quantity      string `json:"qty" bson:"qty"`
	QuoteQuantity string `json:"quote_qty" bson:"quote_qty"`
	TimestampMs   int64  `json:"timestampMs" bson:"timestampMs"`
	IsBuyerMaker  bool   `json:"isBuyerMaker" bson:"isBuyerMaker"`
}

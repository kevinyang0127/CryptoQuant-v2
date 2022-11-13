package user

// type UserInfo struct {
// 	UserID    string `bson:"userID"`
// 	ApiKey    string `bson:"apiKey"`
// 	SecretKey string `bson:"secretKey"`
// }

// func RegisterUser(ctx context.Context, userID string) error {

// }

// func GetUsersApiKey(ctx context.Context, userID string) (apiKey string, secretKey string, err error) {
// 	info := &UserInfo{}
// 	err = m.mongoDB.FindOne(ctx, "cryptoQuantV2", "user", bson.D{{"userID", userID}}, info)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			log.Println("FindOne no result")
// 			log.Println("can't find Strategy by StrategyID = " + strategyID)
// 		} else {
// 			log.Println("FindOne fail")
// 		}
// 		return nil, err
// 	}

// }

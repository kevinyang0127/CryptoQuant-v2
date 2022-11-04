package stream

import "context"

// stream提供pub/sub資料流的介面，方法都要private，應該都交由manager來操作
type Stream interface {
	// 訂閱者要提供唯一key來識別身份，以便退訂閱時知道是誰要退訂
	subscribe(ctx context.Context, subscriberKey string, subscriberCh chan<- []byte) error

	// 訂閱者取消訂閱時要確保不再從通道內讀取資料
	unsubscribe(ctx context.Context, subscriberKey string)

	// 發布訊息給所有訂閱者
	publish(ctx context.Context, data []byte)

	// 資料推送源
	wsConnect(ctx context.Context) error
}

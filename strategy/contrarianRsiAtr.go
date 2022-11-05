package strategy

import (
	"CryptoQuant-v2/indicator"
	"fmt"
)

/*
	策略腳本，未來要移出去讓用戶可以自行新增腳本
*/

type ContrarianRsiAtr struct {
}

func (s *ContrarianRsiAtr) HandleKline(kline indicator.Kline) {
	fmt.Println("handle kline success")
	fmt.Println(kline)
}

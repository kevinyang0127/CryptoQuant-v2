function HandleKline(klines, kline)
    if kline.isFinal then
        local cryptoquant = require("cryptoquant")
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 14, 2)
        local emaArr = cryptoquant.ema(closing, 144)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local supertrend = supertrendArr[#supertrendArr]
        local direction = directionArr[#directionArr]
        local ema = emaArr[#emaArr]
        local macd = macdArr[#macdArr]
        local signal = signalArr[#signalArr]

        cryptoquant.lineNotif({
            msg = "new kline closed",
            ema = ema,
            supertrend = supertrend,
            direction = direction,
            macd = macd,
            signal = signal
        })
    end
end
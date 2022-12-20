function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    if mystatus == nil then
        mystatus = {
            isNotifEarly = false,
            earlyNotifTimeMs = kline.endTime - 900000, -- 15m = 900s = 900000ms
            earlyNotifResetTime = kline.endTime
        }
    end

    if not mystatus.isNotifEarly and cryptoquant.nowTimeMs() >= mystatus.earlyNotifTimeMs then
        table.insert(klines, kline)
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 14, 2)
        local emaArr = cryptoquant.ema(closing, 144)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local closePrice = closing[#closing]
        local prev_closePrice = closing[#closing-1]
        local supertrend = supertrendArr[#supertrendArr]
        local direction = directionArr[#directionArr]
        local prev_direction = directionArr[#directionArr-1]
        local ema = emaArr[#emaArr]
        local prev_ema = emaArr[#emaArr-1]
        local macdEng = macdArr[#macdArr] - signalArr[#signalArr]
        local prev_macdEng = macdArr[#macdArr-1] - signalArr[#signalArr-1]

        local supertrendTurn = direction ~= prev_direction
        local crossEMA = (closePrice - ema)*(prev_closePrice - prev_ema) < 0
        local macdTurn = macdEng * prev_macdEng < 0

        if supertrendTurn or crossEMA or macdTurn then
            cryptoquant.lineNotif({
                supertrendTurn = supertrendTurn,
                crossEMA = crossEMA,
                macdTurn = macdTurn,
            })
        end

        mystatus.isNotifEarly = true
    end

    if mystatus.isNotifEarly and kline.startTime > mystatus.earlyNotifResetTime then
        mystatus.isNotifEarly = false
        mystatus.earlyNotifTimeMs = kline.endTime - 900000
        mystatus.earlyNotifResetTime = kline.endTime
    end

    if kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 14, 2)
        local emaArr = cryptoquant.ema(closing, 144)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local closePrice = closing[#closing]
        local prev_closePrice = closing[#closing-1]
        local supertrend = supertrendArr[#supertrendArr]
        local direction = directionArr[#directionArr]
        local prev_direction = directionArr[#directionArr-1]
        local ema = emaArr[#emaArr]
        local prev_ema = emaArr[#emaArr-1]
        local macdEng = macdArr[#macdArr] - signalArr[#signalArr]
        local prev_macdEng = macdArr[#macdArr-1] - signalArr[#signalArr-1]

        local supertrendTurn = direction ~= prev_direction
        local crossEMA = (closePrice - ema)*(prev_closePrice - prev_ema) < 0
        local macdTurn = macdEng * prev_macdEng < 0

        if supertrendTurn or crossEMA or macdTurn then
            cryptoquant.lineNotif({
                supertrendTurn = supertrendTurn,
                crossEMA = crossEMA,
                macdTurn = macdTurn,
            })
        end
    end

    cryptoquant.saveData(mystatus)
end
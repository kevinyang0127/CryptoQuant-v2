function HandleKline(klines, kline)
    if not kline.isFinal then
        return
    end

    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    if mystatus == nil then
        mystatus = {
            isOpen = false,
            openSide = nil,
            openQty = 0.1,
        }
    end

    local closing = cryptoquant.closedArr(klines)
    local emaArr = cryptoquant.ema(closing, 14)
    local ema = emaArr[#emaArr]
    local prev_ema = emaArr[#emaArr-1]
    local prev_close = closing[#closing-1]

    if not mystatus.isOpen and prev_close < prev_ema and kline.close >= ema then
        cryptoquant.entry(true,mystatus.openQty)
        mystatus.isOpen = true
        mystatus.openSide = true
    elseif not mystatus.isOpen and prev_close > prev_ema and kline.close <= ema then
        cryptoquant.entry(false,mystatus.openQty)
        mystatus.isOpen = true
        mystatus.openSide = false
    end

    if mystatus.isOpen and mystatus.openSide and kline.close < ema then
        cryptoquant.exitAll()
        mystatus.isOpen = false
    elseif mystatus.isOpen and not mystatus.openSide and kline.close > ema then
        cryptoquant.exitAll()
        mystatus.isOpen = false
    end

    cryptoquant.saveData(mystatus)
end



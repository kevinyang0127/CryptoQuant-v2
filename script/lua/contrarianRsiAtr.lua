local function getOpenQuantity(mystatus, openPrice, winProfitRate, commissionRate)
    local nowBet = mystatus.firstBet

    if mystatus.round > 1 and mystatus.round <= mystatus.maxRound then
        local moreBetToWinLoss = (mystatus.totalLoss + mystatus.totalCommission) / winProfitRate
        nowBet = nowBet + moreBetToWinLoss
    end

    local moreBetToWinCommission = nowBet * commissionRate / winProfitRate
    nowBet = nowBet + moreBetToWinCommission

    return nowBet / openPrice
end

function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    
    if mystatus == nil then
        mystatus = {
            isOpen = false,
            openSide = nil,
            waitForDown = false,
            waitForRaise = false,
            openPrice = 0,
            openQty = 0,
            takeProfitPrice = 0,
            stopLossPrice = 0,
            rsiDownTooMuch = false,
            rsiRaiseTooMuch = false,
            profitLossRatio = 1.2,
            takerFeeRate = 0.0004,
            makerFeeRate = 0.0002,
            -- 下注相關變數
            firstBet = 300,
            maxRound = 9,
            totalLoss = 0,
            totalCommission = 0,
            round = 1
        }
    end
    
    --已經開倉且達到停損或停利
    if mystatus.isOpen then
        if (mystatus.openSide and kline.high >= mystatus.takeProfitPrice) or
            (mystatus.openSide and kline.low <= mystatus.stopLossPrice) or
            (not mystatus.openSide and kline.low <= mystatus.takeProfitPrice) or
            (not mystatus.openSide and kline.high >= mystatus.stopLossPrice) then
                
            if (mystatus.openSide and kline.high >= mystatus.takeProfitPrice) or
            (not mystatus.openSide and kline.low <= mystatus.takeProfitPrice) then
                -- 停利
                mystatus.round = 1
                mystatus.totalLoss = 0
                mystatus.totalCommission = 0
                cryptoquant.lineNotif({
                    msg = "Take profit notify"
                })
            else
                -- 停損
                local loss = 0
                if mystatus.openSide then
                    loss = (mystatus.openPrice - mystatus.stopLossPrice) * mystatus.openQty
                else
                    loss = (mystatus.stopLossPrice - mystatus.openPrice) * mystatus.openQty
                end

                if mystatus.round + 1 > mystatus.maxRound then
                    mystatus.round = 1
                    mystatus.totalLoss = 0
                    mystatus.totalCommission = 0
                else
                    mystatus.round = mystatus.round + 1
                    mystatus.totalLoss = mystatus.totalLoss + loss
                    mystatus.totalCommission = mystatus.totalCommission + mystatus.stopLossPrice * mystatus.openQty * mystatus.makerFeeRate
                end
                cryptoquant.lineNotif("[Stop loss notify]")
            end

            mystatus.isOpen = false
            mystatus.waitForDown = false
            mystatus.waitForRaise = false
            cryptoquant.cancelAllOrders()
        end
    end
    
    if kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
		local high = cryptoquant.highArr(klines)
		local low = cryptoquant.lowArr(klines)
        local rsiArr = cryptoquant.rsi(closing, 14)
        local rsi_fast_Arr = cryptoquant.rsi(closing, 4)
        local atrArr = cryptoquant.atr(high, low, closing, 14)
        local rsi = rsiArr[#rsiArr]
        local rsi_fast = rsi_fast_Arr[#rsi_fast_Arr]
        local atr = atrArr[#atrArr]

        if mystatus.rsiDownTooMuch and rsi >= 45 then
            mystatus.rsiDownTooMuch = false
        end

        if mystatus.rsiRaiseTooMuch and rsi <= 55 then
            mystatus.rsiRaiseTooMuch = false
        end

        if mystatus.rsiDownTooMuch or rsi < 20 then
            mystatus.rsiDownTooMuch = true
            mystatus.waitForDown = false
            mystatus.waitForRaise = false
            cryptoquant.saveData(mystatus)
            return
        end

        if mystatus.rsiRaiseTooMuch or rsi > 80 then
            mystatus.rsiRaiseTooMuch = true
            mystatus.waitForDown = false
            mystatus.waitForRaise = false
            cryptoquant.saveData(mystatus)
            return
        end

        if not mystatus.isOpen and rsi <= 30 then
            mystatus.waitForRaise = true
        end

        if not mystatus.isOpen and rsi >= 70 then
            mystatus.waitForDown = true
        end

        --開多倉
        if not mystatus.isOpen and mystatus.waitForRaise and rsi_fast > rsi then
            local stopLossPrice = kline.close - (atr * 2)
            local takeProfitPrice = kline.close + (atr * 2 * mystatus.profitLossRatio)
            local openFee = kline.close * mystatus.takerFeeRate
            local closeFee = takeProfitPrice * mystatus.makerFeeRate
            local profit = atr * 2 * mystatus.profitLossRatio - openFee - closeFee
            if profit <= 0 then
                mystatus.waitForRaise = false
                cryptoquant.saveData(mystatus)
                return
            end

            local winProfitRate = (takeProfitPrice - kline.close) / kline.close
            local qty = getOpenQuantity(mystatus, kline.close, winProfitRate, mystatus.takerFeeRate)
            cryptoquant.entry(true, qty) --市價開倉
            cryptoquant.stopLossOrder(false, stopLossPrice, qty, stopLossPrice) --限價停損單
            cryptoquant.takeProfitOrder(false, takeProfitPrice, qty, takeProfitPrice) --限價停利單

            local msg = string.format("[Market entry notify]\nmartinRound: %d\ntotalLoss: %f\nside: %s\nprice: %f\nqty: %f\nstopLossPrice: %f\ntakeProfitPrice: %f\n",mystatus.round, mystatus.totalLoss, true, kline.close, qty, stopLossPrice, takeProfitPrice)
            cryptoquant.lineNotif(msg)
            
            mystatus.isOpen = true
            mystatus.openSide = true
            mystatus.openPrice = kline.close
            mystatus.openQty = qty
            mystatus.stopLossPrice = stopLossPrice
            mystatus.takeProfitPrice = takeProfitPrice
            mystatus.waitForRaise = false
            cryptoquant.saveData(mystatus)
            return
        end

        --開空倉
        if not mystatus.isOpen and mystatus.waitForDown and rsi_fast < rsi then
            local stopLossPrice = kline.close + (atr * 2)
            local takeProfitPrice = kline.close - (atr * 2 * mystatus.profitLossRatio)
            local openFee = kline.close * mystatus.takerFeeRate
            local closeFee = takeProfitPrice * mystatus.makerFeeRate
            local profit = atr * 2 * mystatus.profitLossRatio - openFee - closeFee
            if profit <= 0 then
                mystatus.waitForDown = false
                cryptoquant.saveData(mystatus)
                return
            end

            local winProfitRate = (kline.close - takeProfitPrice) / kline.close
            local qty = getOpenQuantity(mystatus, kline.close, winProfitRate, mystatus.takerFeeRate)
            cryptoquant.entry(false, qty) --市價開倉
            cryptoquant.stopLossOrder(true, stopLossPrice, qty, stopLossPrice) --限價停損單
            cryptoquant.takeProfitOrder(true, takeProfitPrice, qty, takeProfitPrice) --限價停利單

            local msg = string.format("[Market entry notify]\nmartinRound: %d\ntotalLoss: %f\nside: %s\nprice: %f\nqty: %f\nstopLossPrice: %f\ntakeProfitPrice: %f\n",mystatus.round, mystatus.totalLoss, false, kline.close, qty, stopLossPrice, takeProfitPrice)
            cryptoquant.lineNotif(msg)
            
            mystatus.isOpen = true
            mystatus.openSide = false
            mystatus.openPrice = kline.close
            mystatus.openQty = qty
            mystatus.stopLossPrice = stopLossPrice
            mystatus.takeProfitPrice = takeProfitPrice
            mystatus.waitForDown = false
            cryptoquant.saveData(mystatus)
            return
        end
    end

    cryptoquant.saveData(mystatus)
end
function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    
    if mystatus == nil then
        mystatus = {
            isOpen = false,
            principal = 30000, --本金
            targetRate = 0.05, --目標漲迭幅
            totalStage = 5, --總共多少關卡 (0 ~ n)
            nowStage = 0, -- 目前達標關卡
            firstStopLossRate = 0.01, --開倉時的停損趴數
            openSide = nil,
            openPrice = 0,
            takeProfitPrice = 0,
            stopLossPrice = 0
        }
    end

    -- 檢查supertrend是否轉向
    if mystatus.isOpen and kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 14, 3)
        local direction = directionArr[#directionArr]
        if direction ~= mystatus.openSide then
            cryptoquant.exitAll()
            mystatus.isOpen = false
            mystatus.nowStage = 0
            cryptoquant.cancelAllOrders()
        end
    end

    -- 已經開倉，檢查停損停利或加倉
    if mystatus.isOpen then
        -- 達到停利點
        if (mystatus.openSide and kline.close >= mystatus.takeProfitPrice) or
        (not mystatus.openSide and kline.close <= mystatus.takeProfitPrice) then
            cryptoquant.exitAll()
            mystatus.isOpen = false
            mystatus.nowStage = 0
            cryptoquant.cancelAllOrders()

        -- 達到停損點
        elseif (mystatus.openSide and kline.close <= mystatus.stopLossPrice) or
        (not mystatus.openSide and kline.close >= mystatus.stopLossPrice) then
            cryptoquant.exitAll()
            mystatus.isOpen = false
            mystatus.nowStage = 0
            cryptoquant.cancelAllOrders()

        -- 檢查是否有達到新關卡 多單 
        elseif mystatus.openSide then
            local stageRate = mystatus.targetRate / mystatus.totalStage
            for i=mystatus.nowStage+1,(mystatus.totalStage-1) do
                local nextStagePrice = mystatus.openPrice + mystatus.openPrice * stageRate * i
                -- 加倉
                if kline.close >= nextStagePrice then
                    local stageBet = mystatus.principal / mystatus.totalStage
                    local addQty = stageBet / nextStagePrice
                    cryptoquant.entry(true, addQty)
                    mystatus.nowStage = i
                else
                    break
                end
            end

            --移動停損價位
            if mystatus.nowStage > 0 then
                local stopLossStage = mystatus.nowStage / 2
                if mystatus.nowStage >= (mystatus.totalStage / 2) then
                    stopLossStage = mystatus.nowStage / 1.5
                end
                mystatus.stopLossPrice = mystatus.openPrice + mystatus.openPrice * stopLossStage * stageRate
            end

        -- 檢查是否有達到新關卡 空單 
        elseif not mystatus.openSide then
            local stageRate = mystatus.targetRate / mystatus.totalStage
            for i=mystatus.nowStage+1,(mystatus.totalStage-1) do
                local nextStagePrice = mystatus.openPrice - mystatus.openPrice * stageRate * i
                -- 加倉
                if kline.close <= nextStagePrice then
                    local stageBet = mystatus.principal / mystatus.totalStage
                    local addQty = stageBet / nextStagePrice
                    cryptoquant.entry(false, addQty)
                    mystatus.nowStage = i
                else
                    break
                end
            end

            --移動停損價位
            if mystatus.nowStage > 0 then
                local stopLossStage = mystatus.nowStage / 2
                if mystatus.nowStage >= (mystatus.totalStage / 2) then
                    stopLossStage = mystatus.nowStage / 1.5
                end
                mystatus.stopLossPrice = mystatus.openPrice - mystatus.openPrice * stopLossStage * stageRate
            end
        end
    end
    
    -- 檢查開倉條件
    if not mystatus.isOpen and kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 14, 3)
        local emaArr = cryptoquant.ema(closing, 200)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local closePrice = closing[#closing]
        local supertrend = supertrendArr[#supertrendArr]
        local direction = directionArr[#directionArr]
        local prev_direction = directionArr[#directionArr-1]
        local ema = emaArr[#emaArr]
        local macd = macdArr[#macdArr]
        local signal = signalArr[#signalArr]

        -- 開多倉條件
        if direction == true and prev_direction == false and closePrice > ema and macd > 0 and signal > 0 then
            local stageBet = mystatus.principal / mystatus.totalStage
            local qty = stageBet / closePrice
            cryptoquant.entry(true, qty) --市價開倉
            
            mystatus.stopLossPrice = closePrice - closePrice * mystatus.firstStopLossRate            
            mystatus.takeProfitPrice = closePrice + closePrice * mystatus.targetRate
            mystatus.isOpen = true
            mystatus.openSide = true
            mystatus.openPrice = closePrice
        end
        
        -- 開空倉條件
        if direction == false and prev_direction == true and closePrice < ema and macd < 0 and signal < 0 then
            local stageBet = mystatus.principal / mystatus.totalStage
            local qty = stageBet / closePrice
            cryptoquant.entry(false, qty) --市價開倉
            
            mystatus.stopLossPrice = closePrice + closePrice * mystatus.firstStopLossRate
            mystatus.takeProfitPrice = closePrice - closePrice * mystatus.targetRate
            mystatus.isOpen = true
            mystatus.openSide = false
            mystatus.openPrice = closePrice
        end
        
    end
    
    cryptoquant.saveData(mystatus)
end
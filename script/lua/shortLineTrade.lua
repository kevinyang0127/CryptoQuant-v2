function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    
    if mystatus == nil then
        mystatus = {
            isOpen = false,
            targetRate = 0.015, --目標漲迭幅
            stageRate = 0.005,
            totalStage = 3,
            nowStage = 1, -- 目前達標關卡
            firstStopLossRate = 0.01, --開倉時的停損趴數
            openSide = nil,
            openPrice = 0,
            takeProfitPrice = 0,
            stopLossPrice = 0
        }
    end

    -- 檢查supertrend是否轉向或macd變向
    if mystatus.isOpen and kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 5, 5)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local direction = directionArr[#directionArr]
        local macd = macdArr[#macdArr]
        local signal = signalArr[#signalArr]
        local macdCross = false
        if mystatus.openSide then
            macdCross = macd - signal < 0
        else
            macdCross = macd - signal > 0
        end
        if direction ~= mystatus.openSide or macdCross then
            cryptoquant.exitAll()
            mystatus.isOpen = false
            mystatus.nowStage = 0
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

        -- 達到停損點
        elseif (mystatus.openSide and kline.close <= mystatus.stopLossPrice) or
        (not mystatus.openSide and kline.close >= mystatus.stopLossPrice) then
            cryptoquant.exitAll()
            mystatus.isOpen = false
            mystatus.nowStage = 0

        -- 檢查是否有達到新關卡 多單 
        elseif mystatus.openSide then
            for i=mystatus.nowStage,(mystatus.totalStage) do
                local nextStagePrice = mystatus.openPrice + mystatus.openPrice * mystatus.stageRate * i
                -- 加倉
                if kline.close >= nextStagePrice then
                    local principal = cryptoquant.getBalance()
                    local stageBet = principal / mystatus.totalStage
                    local addQty = stageBet / nextStagePrice
                    cryptoquant.entry(true, addQty)
                    mystatus.nowStage = i+1
                else
                    break
                end
            end

            --移動停損價位
            if mystatus.nowStage > 1 then
                local stopLossStage = mystatus.nowStage - 1.5
                mystatus.stopLossPrice = mystatus.openPrice + mystatus.openPrice * stopLossStage * mystatus.stageRate
            end

        -- 檢查是否有達到新關卡 空單 
        elseif not mystatus.openSide then
            for i=mystatus.nowStage,(mystatus.totalStage) do
                local nextStagePrice = mystatus.openPrice - mystatus.openPrice * mystatus.stageRate * i
                -- 加倉
                if kline.close <= nextStagePrice then
                    local principal = cryptoquant.getBalance()
                    local stageBet = principal / mystatus.totalStage
                    local addQty = stageBet / nextStagePrice
                    cryptoquant.entry(false, addQty)
                    mystatus.nowStage = i+1
                else
                    break
                end
            end

            --移動停損價位
            if mystatus.nowStage > 1 then
                local stopLossStage = mystatus.nowStage - 1.5
                mystatus.stopLossPrice = mystatus.openPrice - mystatus.openPrice * stopLossStage * mystatus.stageRate
            end
        end
    end
    
    -- 檢查開倉條件
    if not mystatus.isOpen and kline.isFinal then
        local closing = cryptoquant.closedArr(klines)
        local high = cryptoquant.highArr(klines)
        local low = cryptoquant.lowArr(klines)
        local supertrendArr, directionArr = cryptoquant.supertrend(high, low, closing, 5, 5)
        local emaArr = cryptoquant.ema(closing, 144)
        local macdArr, signalArr = cryptoquant.macd1226(closing)
        
        local closePrice = closing[#closing]
        local direction = directionArr[#directionArr]
        local ema = emaArr[#emaArr]
        local macd = macdArr[#macdArr]
        local signal = signalArr[#signalArr]

        -- 開多倉條件
        if direction == true and closePrice > ema and macd-signal > 0 then
            local principal = cryptoquant.getBalance()
            local stageBet = principal / mystatus.totalStage
            local qty = stageBet / closePrice
            cryptoquant.entry(true, qty) --市價開倉
            
            mystatus.stopLossPrice = closePrice - closePrice * mystatus.firstStopLossRate
            mystatus.takeProfitPrice = closePrice + closePrice * mystatus.targetRate
            mystatus.isOpen = true
            mystatus.openSide = true
            mystatus.openPrice = closePrice
            mystatus.nowStage = 1
        end
        
        -- 開空倉條件
        if direction == false and closePrice < ema and macd-signal < 0 then
            local principal = cryptoquant.getBalance()
            local stageBet = principal / mystatus.totalStage
            local qty = stageBet / closePrice
            cryptoquant.entry(false, qty) --市價開倉
            
            mystatus.stopLossPrice = closePrice + closePrice * mystatus.firstStopLossRate
            mystatus.takeProfitPrice = closePrice - closePrice * mystatus.targetRate
            mystatus.isOpen = true
            mystatus.openSide = false
            mystatus.openPrice = closePrice
            mystatus.nowStage = 1
        end
        
    end
    
    cryptoquant.saveData(mystatus)
end
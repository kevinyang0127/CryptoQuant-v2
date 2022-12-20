function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    if mystatus == nil then
        mystatus = {
            missionDone = false,
            targetPrice = 79.27,
            side = true,
            entryUSDT= 1390,
        }
    end

    if mystatus.missionDone then
        return
    end

    if (mystatus.side and kline.close >= mystatus.targetPrice) or
        (not mystatus.side and kline.close <= mystatus.targetPrice) then
        local qty = mystatus.entryUSDT / kline.close
        cryptoquant.entry(true,  qty)
        
        mystatus.missionDone = true
        cryptoquant.lineNotif({
            msg = "mission done, stop me!"
        })
    end

    cryptoquant.saveData(mystatus)
end
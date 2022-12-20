function HandleKline(klines, kline)
    local cryptoquant = require("cryptoquant")
    local mystatus = cryptoquant.getData()
    if mystatus == nil then
        mystatus = {
            missionDone = false,
            targetPrice = 75.66,
            side = true
        }
    end

    if mystatus.missionDone then
        return
    end

    if (mystatus.side and kline.close >= mystatus.targetPrice) or
        (not mystatus.side and kline.close <= mystatus.targetPrice) then
        cryptoquant.lineNotif({
            msg = "target price notif",
            targetPrice = mystatus.targetPrice,
            nowPrice = kline.close
        })
        cryptoquant.lineNotif({
            msg = "mission done, stop me!"
        })
        mystatus.missionDone = true
        cryptoquant.saveData(mystatus)
    end
end
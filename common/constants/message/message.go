package message

// exchange trade events
const TopicAggTrade = "ExchangeAggTrade"

// an order was completed
const TopicCompletedOrder = "OrderCompleted"

// request candles
const TopicCandleDataRequest = "GetSomeCandles"

// when a notification is generated
const TopicNotification = "Notification"

// used to notify other services that an engine was started
const TopicEngineStart = "EngineStart"

// when the user deletes a key we must tell the plan service to kill the plans belonging to that account id
const TopicAccountDeleted = "KillAccountPlans"

// engine sends this out to fill a binance order
const TopicFillBinanceOrder = "FillBinanceOrder"

// this will start a new indicator
const TopicStartIndicator = "TopicStartIndicator"

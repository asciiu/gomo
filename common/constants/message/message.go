package message

// TopicNewKey when a new key is added
const TopicNewKey = "NewKey"

// TopicKeyVerified when a key has been validated
const TopicKeyVerified = "KeyVerified"

// TopicBalanceUpdate update balance from exchange has arrived
const TopicBalanceUpdate = "BalanceUpdated"

// TopicAggTrade emitted after trade event
const TopicAggTrade = "ExchangeAggTrade"

// this is used to tell the plan service that an order was aborted successfully
const TopicAbortedOrder = "AbortedOrderSuccess"

// new order event
const TopicNewPlan = "NewPlan"

// TopicOrderFilled when order has filled
const TopicCompletedOrder = "OrderCompleted"

// used to notifiy exchange services to execute an order
const TopicTriggeredOrder = "TriggeredOrder"

const TopicCandleDataRequest = "GetSomeCandles"

// TopicNotification when a notification is generated
const TopicNotification = "Notification"

// used to notify other services that an engine was started
const TopicEngineStart = "EngineStart"

// when the user deletes a key we must tell the plan service to kill the plans belonging to that account id
const TopicAccountDeleted = "KillAccountPlans"

const TopicFillBinanceOrder = "FillBinanceOrder"

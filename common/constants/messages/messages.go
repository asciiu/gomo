package messages

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
const TopicNewOrder = "NewOrder"

// TopicOrderFilled when order has filled
const TopicCompletedOrder = "OrderCompleted"

// used to notifiy exchange services to execute an order
const TopicTriggeredOrder = "TriggeredOrder"

// TopicNotification when a notification is generated
const TopicNotification = "Notification"

// used to notify other services that an engine was started
const TopicEngineStart = "EngineStart"

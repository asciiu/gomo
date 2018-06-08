package messages

// TopicNewKey when a new key is added
const TopicNewKey = "NewKey"

// TopicKeyVerified when a key has been validated
const TopicKeyVerified = "KeyVerified"

// TopicBalanceUpdate update balance from exchange has arrived
const TopicBalanceUpdate = "BalanceUpdated"

// TopicAggTrade emitted after trade event
const TopicAggTrade = "ExchangeAggTrade"

// TopicNewBuyOrder new orders
const TopicNewBuyOrder = "NewBuyOrder"

// TopicNewSellOrder new sell order
const TopicNewSellOrder = "NewSellOrder"

// TopicOrderFilled when order has filled
const TopicOrderFilled = "OrderFilled"

// used to notifiy exchange services to execute an order
const TopicFillOrder = "FillOrder"

// TopicNotification when a notification is generated
const TopicNotification = "Notification"

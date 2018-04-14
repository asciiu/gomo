package models

import "github.com/asciiu/gomo/common/enums"

type ExchangeKey struct {
	ApiKeyId string
	ApiKey   string
	Secret   string
	Exchange enums.Exchange
}

const TopicNewKey = "new.key"

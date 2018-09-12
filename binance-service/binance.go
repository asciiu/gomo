package main

type Symbol struct {
	Symbol             string                    `json:"symbol"`
	Status             string                    `json:"status"`
	BaseAsset          string                    `json:"baseAsset"`
	BaseAssetPrecision int16                     `json:"baseAssetPrecision"`
	QuoteAsset         string                    `json:"quoteAsset"`
	QuotePrecision     int16                     `json:"quotePrecision"`
	OrderTypes         []string                  `json:"orderTypes"`
	IcebergAllowed     bool                      `json:"icebergAllowed"`
	Filters            []*map[string]interface{} `json:"filters"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType`
	Interval      string `json:"interval"`
	Limit         int32  `json:"limit"`
}

type ExchangeInfo struct {
	TimeZone   string       `json:"timezone"`
	ServerTime int64        `json:"serverTime"`
	RateLimits []*RateLimit `json:"rateLimits"`
	Symbols    []*Symbol    `json:"symbols"`
}

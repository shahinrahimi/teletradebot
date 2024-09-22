package bitmex

import "time"

type SideType string

const (
	OrderTypeStop            string = "Stop"
	OrderTypeMarket          string = "Market"
	OrderTypeMarketIfTouched string = "MarketIfTouched"

	OrderStatusTypeNew             string = "New"
	OrderStatusTypePartiallyFilled string = "PartiallyFilled"
	OrderStatusTypeFilled          string = "Filled"
	OrderStatusTypeCanceled        string = "Canceled"

	SideTypeBuy  SideType = "Buy"
	SideTypeSell SideType = "Sell"
)

type MarginData struct {
	Account            int64     `json:"account"`
	Currency           string    `json:"currency"`
	RiskLimit          int64     `json:"riskLimit"`
	Amount             int64     `json:"amount"`
	GrossComm          int64     `json:"grossComm"`
	GrossOpenCost      int64     `json:"grossOpenCost"`
	GrossOpenPremium   int64     `json:"grossOpenPremium"`
	GrossMarkValue     int64     `json:"grossMarkValue"`
	RiskValue          int64     `json:"riskValue"`
	InitMargin         int64     `json:"initMargin"`
	MaintMargin        int64     `json:"maintMargin"`
	TargetExcessMargin int64     `json:"targetExcessMargin"`
	RealisedPnl        int64     `json:"realisedPnl"`
	UnrealisedPnl      int64     `json:"unrealisedPnl"`
	WalletBalance      int64     `json:"walletBalance"`
	MarginBalance      int64     `json:"marginBalance"`
	MarginLeverage     float64   `json:"marginLeverage"`
	MarginUsedPcnt     float64   `json:"marginUsedPcnt"`
	ExcessMargin       int64     `json:"excessMargin"`
	AvailableMargin    int64     `json:"availableMargin"`
	WithdrawableMargin int64     `json:"withdrawableMargin"`
	Timestamp          time.Time `json:"timestamp"`
}

type MarginTable struct {
	Table  string       `json:"table"`
	Action string       `json:"action"`
	Keys   []string     `json:"keys"`
	Types  interface{}  `json:"types"`
	Filter Filter       `json:"filter"`
	Data   []MarginData `json:"data"`
}

type Filter struct {
	Account int64 `json:"account"`
}

type OrderData struct {
	OrderID          string    `json:"orderID"`
	ClOrdID          string    `json:"clOrdID"`
	ClOrdLinkID      string    `json:"clOrdLinkID"`
	Account          int64     `json:"account"`
	Symbol           string    `json:"symbol"`
	Side             string    `json:"side"`
	OrderQty         int64     `json:"orderQty"`
	Price            float64   `json:"price"`
	DisplayQty       int64     `json:"displayQty"`
	StopPx           float64   `json:"stopPx"`
	PegOffsetValue   float64   `json:"pegOffsetValue"`
	PegPriceType     string    `json:"pegPriceType"`
	Currency         string    `json:"currency"`
	SettlCurrency    string    `json:"settlCurrency"`
	OrdType          string    `json:"ordType"`
	TimeInForce      string    `json:"timeInForce"`
	ExecInst         string    `json:"execInst"`
	ContingencyType  string    `json:"contingencyType"`
	OrdStatus        string    `json:"ordStatus"`
	Triggered        string    `json:"triggered"`
	WorkingIndicator bool      `json:"workingIndicator"`
	OrdRejReason     string    `json:"ordRejReason"`
	LeavesQty        int64     `json:"leavesQty"`
	CumQty           int64     `json:"cumQty"`
	AvgPx            float64   `json:"avgPx"`
	Text             string    `json:"text"`
	TransactTime     time.Time `json:"transactTime"`
	Timestamp        time.Time `json:"timestamp"`
}

type OrderTable struct {
	Table  string      `json:"table"`
	Action string      `json:"action"`
	Keys   []string    `json:"keys"`
	Types  interface{} `json:"types"`
	Filter Filter      `json:"filter"`
	Data   []OrderData `json:"data"`
}

type ExecutionData struct {
	ExecID           string    `json:"execID"`
	OrderID          string    `json:"orderID"`
	ClOrdID          string    `json:"clOrdID"`
	ClOrdLinkID      string    `json:"clOrdLinkID"`
	Account          int64     `json:"account"`
	Symbol           string    `json:"symbol"`
	Side             string    `json:"side"`
	LastQty          int64     `json:"lastQty"`
	LastPx           float64   `json:"lastPx"`
	LastLiquidityInd string    `json:"lastLiquidityInd"`
	OrderQty         int64     `json:"orderQty"`
	Price            float64   `json:"price"`
	DisplayQty       int64     `json:"displayQty"`
	StopPx           float64   `json:"stopPx"`
	PegOffsetValue   float64   `json:"pegOffsetValue"`
	PegPriceType     string    `json:"pegPriceType"`
	Currency         string    `json:"currency"`
	SettlCurrency    string    `json:"settlCurrency"`
	ExecType         string    `json:"execType"`
	OrdType          string    `json:"ordType"`
	TimeInForce      string    `json:"timeInForce"`
	ExecInst         string    `json:"execInst"`
	ContingencyType  string    `json:"contingencyType"`
	OrdStatus        string    `json:"ordStatus"`
	Triggered        string    `json:"triggered"`
	WorkingIndicator bool      `json:"workingIndicator"`
	OrdRejReason     string    `json:"ordRejReason"`
	LeavesQty        int64     `json:"leavesQty"`
	CumQty           int64     `json:"cumQty"`
	AvgPx            float64   `json:"avgPx"`
	Commission       float64   `json:"commission"`
	FeeType          string    `json:"feeType"`
	TradePublishInd  string    `json:"tradePublishIndicator"`
	TrdMatchID       string    `json:"trdMatchID"`
	ExecCost         int64     `json:"execCost"`
	ExecComm         int64     `json:"execComm"`
	HomeNotional     float64   `json:"homeNotional"`
	ForeignNotional  float64   `json:"foreignNotional"`
	TransactTime     time.Time `json:"transactTime"`
	Timestamp        time.Time `json:"timestamp"`
}

type ExecutionTable struct {
	Table  string          `json:"table"`
	Action string          `json:"action"`
	Keys   []string        `json:"keys"`
	Types  interface{}     `json:"types"`
	Filter Filter          `json:"filter"`
	Data   []ExecutionData `json:"data"`
}

type InstrumentTable struct {
	Table  string           `json:"table"`
	Action string           `json:"action"`
	Data   []InstrumentData `json:"data"`
}

type InstrumentData struct {
	Symbol    string    `json:"symbol"`
	OpenValue int64     `json:"openValue"`
	FairPrice float64   `json:"fairPrice"`
	MarkPrice float64   `json:"markPrice"`
	Timestamp time.Time `json:"timestamp"` // You may use time.Time with proper unmarshaling
}

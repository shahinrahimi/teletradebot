package types

type OrderIDType string
type ExecutionType string

type BotMessage struct {
	ChatID int64
	MsgStr string
}

type BotError struct {
	Msg string
}

func (b *BotError) Error() string {
	return b.Msg
}

const (
	ACCOUNT_B string = `Binance`
	ACCOUNT_M string = `Bitmex`

	SIDE_L string = `LONG`
	SIDE_S string = `SHORT`

	STATE_IDLE      string = `idle`
	STATE_CANCELED  string = `canceled`
	STATE_PLACED    string = `placed`
	STATE_FILLED    string = `filled`
	STATE_STOPPED   string = `stopped`
	STATE_PROFITED  string = `profited`
	STATE_CLOSED    string = `closed`
	STATE_EXPIRED   string = `expired`
	STATE_REVERTING string = `reverting`

	OrderIDTypeMain              OrderIDType = `main`
	OrderIDTypeStopLoss          OrderIDType = `stop-loss`
	OrderIDTypeTakeProfit        OrderIDType = `take-profit`
	OrderIDTypeReverseMain       OrderIDType = `reverse-main`
	OrderIDTypeReverseStopLoss   OrderIDType = `reverse-stop-loss`
	OrderIDTypeReverseTakeProfit OrderIDType = `reverse-take-profit`
	OrderIDTypeNone              OrderIDType = `none`

	GetOrderExecution             ExecutionType = "get-order"
	CancelOrderExecution          ExecutionType = "cancel-order"
	StopPriceExecution            ExecutionType = "stop-price"
	StopLossExecution             ExecutionType = "stop-loss"
	TakeProfitExecution           ExecutionType = "take-profit"
	ClosePositionExecution        ExecutionType = "close-position"
	ReverseStopPriceExecution     ExecutionType = "reverse-stop-price"
	ReverseStopLossExecution      ExecutionType = "reverse-stop-loss"
	ReverseTakeProfitExecution    ExecutionType = "reverse-take-profit"
	ClosePositionReverseExecution ExecutionType = "close-position-reverse"
)

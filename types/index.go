package types

type OrderIDType string

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
	OrderIDTypeMain              OrderIDType = `main`
	OrderIDTypeStopLoss          OrderIDType = `stop-loss`
	OrderIDTypeTakeProfit        OrderIDType = `take-profit`
	OrderIDTypeReverseMain       OrderIDType = `reverse-main`
	OrderIDTypeReverseStopLoss   OrderIDType = `reverse-stop-loss`
	OrderIDTypeReverseTakeProfit OrderIDType = `reverse-take-profit`
	OrderIDTypeNone              OrderIDType = `none`

	ACCOUNT_B string = `Binance`
	ACCOUNT_M string = `Bitmex`

	SIDE_L string = `LONG`
	SIDE_S string = `SHORT`

	STATE_IDLE      string = `idle`
	STATE_CANCELED  string = `canceled`
	STATE_PLACED    string = `placed`
	STATE_FILLED    string = `filled`
	STATE_REVERTING string = `reverting`
	STATE_STOPPED   string = `stopped`
	STATE_PROFITED  string = `profited`
	STATE_CLOSED    string = `closed`
)

package types

type SideType string
type StateType string
type ExchangeType string
type OrderIDType string
type ExecutionType string
type AccountType string
type OrderTitleType string
type VerbType string

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
	ExchangeBinance ExchangeType = `binance`
	ExchangeBitmex  ExchangeType = `bitmex`

	SideLong  SideType = `LONG`
	SideShort SideType = `SHORT`

	StateIdle      StateType = `idle`
	StatePlaced    StateType = `placed`
	StateCanceled  StateType = `canceled`
	StateFilled    StateType = `filled`
	StateStopped   StateType = `stopped`
	StateProfited  StateType = `profited`
	StateClosed    StateType = `closed`
	StateExpired   StateType = `expired`
	StateReverting StateType = `reverting`

	OrderTitleMain              OrderTitleType = `main`
	OrderTitleStopLoss          OrderTitleType = `stop-loss`
	OrderTitleTakeProfit        OrderTitleType = `take-profit`
	OrderTitleReverseMain       OrderTitleType = `reverse-main`
	OrderTitleReverseStopLoss   OrderTitleType = `reverse-stop-loss`
	OrderTitleReverseTakeProfit OrderTitleType = `reverse-take-profit`
	OrderTitleNone              OrderTitleType = `none`

	VerbPlaced   VerbType = `placed`
	VerbReplaced VerbType = `replaced`
	VerbCanceled VerbType = `canceled`
	VerbClosed   VerbType = `closed`
	VerbFilled   VerbType = `filled`
	VerbExecuted VerbType = `executed`

	OrderIDTypeMain              OrderIDType = `main`
	OrderIDTypeStopLoss          OrderIDType = `stop-loss`
	OrderIDTypeTakeProfit        OrderIDType = `take-profit`
	OrderIDTypeReverseMain       OrderIDType = `reverse-main`
	OrderIDTypeReverseStopLoss   OrderIDType = `reverse-stop-loss`
	OrderIDTypeReverseTakeProfit OrderIDType = `reverse-take-profit`
	OrderIDTypeNone              OrderIDType = `none`

	ExecutionGetOrder               ExecutionType = "get-order"
	ExecutionCancelOrder            ExecutionType = "cancel-order"
	ExecutionCloseMainOrder         ExecutionType = "close-main-order"
	ExecutionCloseReverseMainOrder  ExecutionType = "close-reverse-main-order"
	ExecutionTakeProfitOrder        ExecutionType = "take-profit-order"
	ExecutionStopLossOrder          ExecutionType = "stop-loss-order"
	ExecutionEntryMainOrder         ExecutionType = "entry-main-order"
	ExecutionEntryReverseMainOrder  ExecutionType = "entry-reverse-main-order"
	ExecutionStopLossReverseOrder   ExecutionType = "stop-loss-reverse-order"
	ExecutionTakeProfitReverseOrder ExecutionType = "take-profit-reverse-order"

	ExecutionNone ExecutionType = "none"
)

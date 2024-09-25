package bot

import (
	"fmt"

	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) getMessagePlacedOrder(orderTitle types.OrderTitleType, verb types.VerbType, tradeID int64, orderIDorOrderRes interface{}) string {
	orderIDStr := utils.ExtractOrderIDStr(orderIDorOrderRes)
	msg := fmt.Sprintf("%s order %s successfully.\n\nOrder ID: %s\nTrade ID: %d", utils.Capitalize(string(orderTitle)), string(verb), orderIDStr, tradeID)
	return msg
}

func (b *Bot) getMessageStopped(orderTitle types.OrderTitleType, verb types.VerbType, tradeID int64) string {
	msg := fmt.Sprintf("🛑 %s order %s successfully.\n\nTrade ID: %d", utils.Capitalize(string(orderTitle)), string(verb), tradeID)
	return msg
}

func (b *Bot) getMessageProfited(orderTitle types.OrderTitleType, verb types.VerbType, tradeID int64) string {
	msg := fmt.Sprintf("✅ %s order %s successfully.\n\nTrade ID: %d", utils.Capitalize(string(orderTitle)), string(verb), tradeID)
	return msg
}

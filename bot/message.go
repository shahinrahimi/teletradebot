package bot

import (
	"fmt"
	"strings"

	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (b *Bot) getMessagePlacedOrder(orderTitle types.OrderTitleType, verb types.VerbType, tradeID int64, orderIDorOrderRes interface{}) string {
	orderIDStr := utils.ExtractOrderIDStr(orderIDorOrderRes)
	msg := fmt.Sprintf("%s %s successfully.\n\nOrder ID: %s\nTrade ID: %d", strings.ToUpper(string(orderTitle)), string(verb), orderIDStr, tradeID)
	return msg
}

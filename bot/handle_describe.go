package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
)

func (b *Bot) HandleDescribe(u *tgbotapi.Update, ctx context.Context) error {
	t := ctx.Value(models.KeyTrade{}).(models.Trade)
	userID := u.Message.From.ID
	if t.Account == types.ACCOUNT_B {
		// if trade state is idle or placed it should get latest describer
		if !((t.State == types.STATE_IDLE) || (t.State == types.STATE_PLACED)) {
			// try to read describer from memory
			d, exist := models.GetDescriber(t.ID)
			if exist {
				b.MsgChan <- types.BotMessage{
					ChatID: userID,
					MsgStr: d.ToString(&t),
				}
			} else {
				d, err := b.bc.FetchDescriber(ctx, &t)
				if err != nil {
					b.l.Printf("error fetching describer")
					return err
				}
				b.MsgChan <- types.BotMessage{
					ChatID: userID,
					MsgStr: d.ToString(&t),
				}
			}
		} else {
			d, err := b.bc.FetchDescriber(ctx, &t)
			if err != nil {
				b.l.Printf("error fetching describer")
				return err
			}
			b.MsgChan <- types.BotMessage{
				ChatID: userID,
				MsgStr: d.ToString(&t),
			}
		}

	} else {
		d, err := b.mc.FetchDescriber(ctx, &t)
		if err != nil {
			b.l.Printf("error fetching describer")
			return err
		}
		b.MsgChan <- types.BotMessage{
			ChatID: userID,
			MsgStr: d.ToString(&t),
		}
	}
	return nil
}

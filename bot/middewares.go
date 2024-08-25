package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gihub.com/shahinrahimi/teletradebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Logger(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		if u.Message.Command() == "" {
			b.l.Printf("Received message: %s", u.Message.Text)
		} else {
			b.l.Printf("Received command: %s\t args: %s", u.Message.Command(), u.Message.CommandArguments())
		}
		next(u, ctx)
	}
}

func (b *Bot) ProvideAddOrder(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		args := strings.Split(u.Message.CommandArguments(), " ")
		o, err := ParseOrder(args)
		if err != nil {
			b.SendMessage(u.Message.From.ID, err.Error())
			return
		}
		ctx = context.WithValue(ctx, models.KeyOrder{}, *o)
		next(u, ctx)
	}
}

func ParseOrder(args []string) (*models.Order, error) {
	var o models.Order
	if len(args) < 9 {
		return nil, fmt.Errorf("the length of args is not sufficient for parsing")
	}
	// account it should be string
	// m for bitmex
	// b for binance
	part1 := strings.TrimSpace(strings.ToLower(args[0]))
	if len(part1) > 1 || (part1 != "m" && part1 != "b") {
		return nil, fmt.Errorf("the valid value for account should be 'm' => bitmex, 'b' => binance")
	} else if part1 == "m" {
		o.Account = models.ACCOUNT_M
	} else if part1 == "b" {
		o.Account = models.ACCOUNT_B
	} else {
		// should never happen
		return nil, fmt.Errorf("internal error")
	}
	// pair
	// TODO maybe add check if pair exist on the tickers
	part2 := strings.TrimSpace(strings.ToUpper(args[1]))
	o.Pair = part2
	// side
	part3 := strings.TrimSpace(strings.ToUpper(args[2]))
	if part3 != models.SIDE_L && part3 != models.SIDE_S {
		return nil, fmt.Errorf("the valid value for side should be 'long' or 'short'")
	} else {
		o.Side = part3
	}
	// pair
	part4 := strings.TrimSpace(strings.ToLower(args[3]))
	if part4 != models.CANDLE_15MIN && part4 != models.CANDLE_30MIN && part4 != models.CANDLE_1H && part4 != models.CANDLE_4H {
		return nil, fmt.Errorf("the valid value for candle should be '15min', '30min', '1h' or '4h'")
	} else {
		o.Pair = part4
	}
	// offset
	part5 := strings.TrimSpace(args[4])
	offset, err := strconv.ParseFloat(part5, 64)
	if err != nil {
		return nil, fmt.Errorf("the valid value for offset_entry should be amount (float or integer)")
	} else {
		o.Offset = float32(offset)
	}
	// size percent
	part6 := strings.TrimSpace(args[5])
	size_percent, err := strconv.Atoi(part6)
	if err != nil {
		return nil, fmt.Errorf("the valid value for size should be amount in percent (e.g 5)")
	} else if size_percent <= 0 || size_percent > 50 {
		return nil, fmt.Errorf("the valid value for size should be a non-zero none-negative number (max: 50)")
	} else {
		o.SLPercent = size_percent
	}

	// stop-loss percent
	part7 := strings.TrimSpace(args[6])
	stop_percent, err := strconv.Atoi(part7)
	if err != nil {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be amount in percent (e.g 105)")
	} else if stop_percent < 100 {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be a non-zero none-negative number (min: 100)")
	} else {
		o.SLPercent = stop_percent
	}

	// target-point percent
	part8 := strings.TrimSpace(args[7])
	target_percent, err := strconv.Atoi(part8)
	if err != nil {
		return nil, fmt.Errorf("the valid value for target-point percent should be amount in percent (e.g 105)")
	} else if target_percent < 100 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (min: 100)")
	} else {
		o.TPPercent = target_percent
	}

	// reverse-multiplier
	part9 := strings.TrimSpace(args[8])
	reverse_multiplier, err := strconv.Atoi(part9)
	if err != nil {
		return nil, fmt.Errorf("the valid value for reverse_multiplier should be number (1 or 2)")
	} else if reverse_multiplier <= 0 || reverse_multiplier > 2 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (1 or 2)")
	} else {
		o.ReverseMultiplier = reverse_multiplier
	}

	o.State = models.STATE_IDLE

	return &o, nil
}

package bot

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gihub.com/shahinrahimi/teletradebot/config"
	"gihub.com/shahinrahimi/teletradebot/models"
	"gihub.com/shahinrahimi/teletradebot/types"
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

func (b *Bot) ProvideAddTrade(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		var tradeArgs []string
		args := strings.Split(u.Message.CommandArguments(), " ")
		// check shortcuts
		if len(args) == 1 {
			value, exist := config.Shortcuts[args[0]]
			if !exist {
				b.SendMessage(u.Message.From.ID, fmt.Sprintf("shortcut is not found: '%s'", args[0]))
				return
			}
			tradeArgs = strings.Split(value, " ")
		} else {
			tradeArgs = args
		}
		o, err := ParseTrade(tradeArgs)
		if err != nil {
			b.SendMessage(u.Message.From.ID, err.Error())
			return
		}

		var isAvailable bool = false
		switch o.Account {
		case types.ACCOUNT_B:
			// check pair for binance
			for _, s := range b.bc.Symbols {
				if s.Symbol == strings.ToUpper(o.Pair) {
					isAvailable = true
					break
				}
			}
		case types.ACCOUNT_M:
			// check pair for bitmex
			// TODO implement a method to check if pair is available in bitmex
			c := http.Client{}
			resp, err := c.Get(fmt.Sprintf("https://www.bitmex.com/api/udf/symbols?symbol=%s", o.Pair))
			if resp.StatusCode == 200 || err != nil {
				isAvailable = true
			}
		default:
			// should never happen

			b.l.Panicf("error checking pair availability pair: %s", o.Pair)

		}

		if !isAvailable {
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("the pair '%s' is not available", o.Pair))
			return
		}

		ctx = context.WithValue(ctx, models.KeyTrade{}, *o)
		next(u, ctx)
	}
}

func (b *Bot) ProvideTradeByID(next Handler) Handler {
	return func(u *tgbotapi.Update, ctx context.Context) {
		args := strings.Split(u.Message.CommandArguments(), " ")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("the id is not valid: %s", args[0]))
			return
		}
		o, err := b.s.GetTrade(id)
		if err != nil {
			b.SendMessage(u.Message.From.ID, fmt.Sprintf("the id is not valid: %d", id))
			return
		}
		ctx = context.WithValue(ctx, models.KeyTrade{}, *o)
		next(u, ctx)
	}
}

func ParseTrade(tradeArgs []string) (*models.Trade, error) {
	var t models.Trade
	if len(tradeArgs) < 9 {
		return nil, fmt.Errorf("the length of args is not sufficient for parsing")
	}
	// account it should be string
	// m for bitmex
	// b for binance
	part1 := strings.TrimSpace(strings.ToLower(tradeArgs[0]))
	if len(part1) > 1 || (part1 != "m" && part1 != "b") {
		return nil, fmt.Errorf("the valid value for account should be 'm' => bitmex, 'b' => binance")
	} else if part1 == "m" {
		t.Account = types.ACCOUNT_M
	} else if part1 == "b" {
		t.Account = types.ACCOUNT_B
	} else {
		// should never happen
		return nil, fmt.Errorf("internal error")
	}
	// pair
	// TODO maybe add check if pair exist on the tickers
	part2 := strings.TrimSpace(strings.ToUpper(tradeArgs[1]))
	t.Pair = part2
	// side
	part3 := strings.TrimSpace(strings.ToUpper(tradeArgs[2]))
	if part3 != types.SIDE_L && part3 != types.SIDE_S {
		return nil, fmt.Errorf("the valid value for side should be 'long' or 'short'")
	} else {
		t.Side = part3
	}
	// candle
	part4 := strings.TrimSpace(tradeArgs[3])
	if !types.IsValidCandle(part4) {
		return nil, fmt.Errorf("the valid value for candle should be %s", types.GetValidCandlesString())
	} else {
		t.Candle = part4
	}
	// offset
	part5 := strings.TrimSpace(tradeArgs[4])
	offset, err := strconv.ParseFloat(part5, 64)
	if err != nil {
		return nil, fmt.Errorf("the valid value for offset_entry should be amount (float or integer)")
	} else {
		t.Offset = float32(offset)
	}
	// size percent
	part6 := strings.TrimSpace(tradeArgs[5])
	size_percent, err := strconv.Atoi(part6)
	fmt.Println(part6, size_percent)
	if err != nil {
		return nil, fmt.Errorf("the valid value for size should be amount in percent (e.g 5)")
	} else if size_percent <= 0 || size_percent > 50 {
		return nil, fmt.Errorf("the valid value for size should be a non-zero none-negative number (max: 50)")
	} else {
		t.SizePercent = size_percent
	}

	// stop-loss percent
	part7 := strings.TrimSpace(tradeArgs[6])
	stop_percent, err := strconv.Atoi(part7)
	if err != nil {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be amount in percent (e.g 105)")
	} else if stop_percent < 100 {
		return nil, fmt.Errorf("the valid value for stop-loss percent should be a non-zero none-negative number (min: 100)")
	} else {
		t.SLPercent = stop_percent
	}

	// target-point percent
	part8 := strings.TrimSpace(tradeArgs[7])
	target_percent, err := strconv.Atoi(part8)
	if err != nil {
		return nil, fmt.Errorf("the valid value for target-point percent should be amount in percent (e.g 105)")
	} else if target_percent < 100 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (min: 100)")
	} else {
		t.TPPercent = target_percent
	}

	// reverse-multiplier
	part9 := strings.TrimSpace(tradeArgs[8])
	reverse_multiplier, err := strconv.Atoi(part9)
	if err != nil {
		return nil, fmt.Errorf("the valid value for reverse_multiplier should be number (1 or 2)")
	} else if reverse_multiplier <= 0 || reverse_multiplier > 2 {
		return nil, fmt.Errorf("the valid value for target-point percent should be a non-zero none-negative number (1 or 2)")
	} else {
		t.ReverseMultiplier = reverse_multiplier
	}

	t.State = types.STATE_IDLE

	return &t, nil
}

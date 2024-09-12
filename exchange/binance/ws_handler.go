package binance

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shahinrahimi/teletradebot/models"
	"github.com/shahinrahimi/teletradebot/types"
	"github.com/shahinrahimi/teletradebot/utils"
)

func (bc *BinanceClient) handleOrderTradeUpdate(f futures.WsOrderTradeUpdate) {
	switch f.Status {
	case futures.OrderStatusTypeCanceled:
		bc.l.Println("Order was canceled.")
		// b.HandleCanceled(f)
	case futures.OrderStatusTypeFilled:
		bc.l.Println("Order filled successfully.")
		go bc.handleFilled(f)
	case futures.OrderStatusTypeRejected:
		bc.l.Println("Order was rejected.")
	case futures.OrderStatusTypeNew:
		bc.l.Println("New order received.")
	case futures.OrderStatusTypeExpired:
		bc.l.Println("Order has expired.")
	case futures.OrderStatusTypePartiallyFilled:
		bc.l.Println("Order partially filled.")
	default:
		bc.l.Println("Unknown order status received.")
	}
}

func (bc *BinanceClient) handleFilled(f futures.WsOrderTradeUpdate) {
	orderID := utils.ConvertBinanceOrderID(f.ID)
	var t *models.Trade
	var err error
	// sleep a little bit to make sure the store is updated for early filled orders
	// TODO maybe change the logic in future for better handling
	time.Sleep(time.Second)
	// check if order related to a trade
	t, err = bc.s.GetTradeByOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			bc.l.Panic("Internal error while retrieving trade:", err)
		}
	}
	// handle new filled order
	if t != nil {
		if err := bc.s.UpdateTradeFilled(t); err != nil {
			bc.l.Printf("Error updating trade state to FILLED: %v", types.STATE_FILLED)
			return
		}
		msg := fmt.Sprintf("Order filled successfully.\n\nTrade ID: %d", t.ID)
		bc.MsgChan <- types.BotMessage{
			ChatID: t.UserID,
			MsgStr: msg,
		}
		bc.handleNewFilled(t, &f)
		return
	}
	// check if order is for stop loss
	t, err = bc.s.GetTradeBySLOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			bc.l.Panic("Internal error while retrieving stop-loss trade:", err)
		}
	}
	if t != nil {
		if err := bc.s.UpdateTradeStopped(t); err != nil {
			bc.l.Printf("Error updating trade state to STOPPED: %v", types.STATE_STOPPED)
			return
		}
		bc.handleSLFilled(t, &f)
		return
	}
	// check if order is for take profit
	t, err = bc.s.GetTradeByTPOrderID(orderID)
	if err != nil {
		if err != sql.ErrNoRows {
			bc.l.Panic("Internal error while retrieving take-profit trade:", err)
		}
	}
	if t != nil {
		if err := bc.s.UpdateTradeProfited(t); err != nil {
			bc.l.Printf("Error updating trade state to PROFITED: %v", types.STATE_PROFITED)
			return
		}
		bc.handleTPFilled(t, &f)
		return
	}
}

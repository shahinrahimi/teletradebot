package models

import (
	"fmt"
	"time"
)

var (
	describers = map[int64]*Describer{}
)

type Describer struct {
	From  time.Time
	Till  time.Time
	Open  string
	High  string
	Low   string
	Close string
	SP    string // strop price or entry
	TP    string // take-profit price
	SL    string // take-loss price
}

func GetDescriber(tradeID int64) (*Describer, bool) {
	if d, exist := describers[tradeID]; exist {
		return d, true
	}
	return nil, false
}

func SetDescriber(d *Describer, tradeID int64) {
	describers[tradeID] = d
}

func UpdateDescriberSL(tradeID int64, sl string) {
	if _, exist := describers[tradeID]; exist {
		describers[tradeID].SL = sl
	}
}

func UpdateDescriberTP(tradeID int64, tp string) {
	if _, exist := describers[tradeID]; exist {
		describers[tradeID].TP = tp
	}
}

func DeleteDescriber(tradeID int64) {
	delete(describers, tradeID)
}

func (d *Describer) ToString(t *Trade) string {
	sizeStr := fmt.Sprintf("%.1f%%", float64(t.Size))
	slStr := fmt.Sprintf("%.1f%%", float64((t.StopLoss - 100)))
	tpStr := fmt.Sprintf("%.1f%%", float64((t.TakeProfit - 100)))

	format := "2006-01-02 15:04:05"
	FromStr := d.From.Local().Format(format)
	TillStr := d.Till.Local().Format(format)

	msg := fmt.Sprintf("Trade ID %d\n\n", t.ID)
	msg = fmt.Sprintf("%s From:  %s\n Till:  %s\n Open:  %s\n High:  %s\n Low:  %s\n Close:  %s\n\n", msg, FromStr, TillStr, d.Open, d.High, d.Low, d.Close)
	msg = fmt.Sprintf("%sTrading:\n", msg)
	msg = fmt.Sprintf("%s Entry %s at %s with %s of balance.\n", msg, t.Side, d.SP, sizeStr)
	msg = fmt.Sprintf("%s TP at %s with %s.\n", msg, d.TP, tpStr)
	msg = fmt.Sprintf("%s SL at %s with %s.\n", msg, d.SL, slStr)
	return msg
}

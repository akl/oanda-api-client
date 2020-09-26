package main

import (
	"fmt"
	"github.com/yuki-inoue-eng/oanda-api-client"
	"log"
	"time"
)

// example
func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	gtdTime := time.Now().AddDate(0, 0, 2).UTC()
	order := oanda.Order{
		TakeProfitOnFill:       nil,
		StopLossOnFill:         nil,
		TrailingStopLossOnFill: nil,
		CreateTime:             nil,
		ID:                     "12",
		Instrument:             oanda.InstrumentUSDJPY,
		GtdTime:                &gtdTime,
		PartialFill:            "DEFAULT_FILL",
		PositionFill:           "DEFAULT",
		Price:                  oanda.Price(107.000),
		State:                  "PENDING",
		TimeInForce:            oanda.TimeInForceGTD,
		TriggerCondition:       "DEFAULT",
		Type:                   oanda.OrderTypeMarketIfTouched,
		Units:                  oanda.Unit(-1),
	}
	if err := client.UpdateOrder(order); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("order updated")
}

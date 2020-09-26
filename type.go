package oanda

import "strconv"

const (
	SideBuy                     = side("buy")
	SideSell                    = side("sell")
	OrderTypeMarket             = orderType("MARKET")
	OrderTypeLimit              = orderType("LIMIT")
	OrderTypeStop               = orderType("STOP")
	OrderTypeMarketIfTouched    = orderType("MARKET_IF_TOUCHED")
	OrderTypeTakeProfit         = orderType("TAKE_PROFIT")
	OrderTypeStopLoss           = orderType("STOP_LOSS")
	OrderTypeGuaranteedStopLoss = orderType("GUARANTEED_STOP_LOSS")
	OrderTypeTrailingStopLoss   = orderType("TRAILING_STOP_LOSS")
	OrderFixedPrice             = orderType("FIXED_PRICE")
	TimeInForceGTC              = timeInForce("GTC")
	TimeInForceGTD              = timeInForce("GTD")
	TimeInForceGFD              = timeInForce("GFD")
	TimeInForceFOK              = timeInForce("FOK")
	TimeInForceIOC              = timeInForce("IOC")
	InstrumentUSDJPY            = instrument("USD_JPY")
	InstrumentEURJPY            = instrument("EUR_JPY")
	InstrumentEURUSD            = instrument("EUR_USD")
)

type Price float64
func (p *Price)String()string{
	return strconv.FormatFloat(float64(*p), 'f', 7, 64)
}

type Pips float64 // valid up to the first minority

type side string

type orderType string

type timeInForce string

type orderID string

type instrument string

type Unit int

func (p *Pips)PipsToPrice(instrument string) Price {
	if instrument == "USD_JPY" {
		return Price(float64(*p) * 0.01)
	}
	if instrument == "EUR_JPY" {
		return Price(float64(*p) * 0.01)
	}
	if instrument == "EUR_USD" {
		return Price(float64(*p) * 0.0001)
	}
	return 0
}
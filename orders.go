package oanda_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	SideBuy                  = side("buy")
	SideSell                 = side("sell")
	OrderTypeLimit           = orderType("limit")
	OrderTypeStop            = orderType("stop")
	OrderTypeMarketIfTouched = orderType("marketIfTouched")
	OrderTypeMarket          = orderType("market")
)

type Price float64

type Pips float64 // valid up to the first minority

type side string

type orderType string

type orderID int64

// Order defines order.
type Order struct {
	ID           orderID
	Instrument   string
	Units        int
	Side         side
	OrderType    orderType
	Expiry       *time.Time
	Price        Price
	LowerBound   Price
	UpperBound   Price
	StopLoss     Price
	TakeProfit   Price
	TrailingStop Pips
}

// PostOrder posts new order, and returns order created by this posts.
func (c *Client) PostOrder(accountID int64, order Order) (*Order, error) {
	body, err := buildPostOrderBody(order)
	if err != nil {
		return nil, fmt.Errorf("failed to build request body to create new order: %v", err)
	}
	respBody, err := c.createOrder(accountID, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}
	o := struct {
		Instrument  string     `json:"instrument"`
		Time        *time.Time `json:"time"`
		Price       Price      `json:"price"`
		TradeOpened struct {
			ID           orderID `json:"id"`
			Units        int     `json:"units"`
			Side         side    `json:"side"`
			TakeProfit   Price   `json:"takeProfit"`
			StopLoss     Price   `json:"stopLoss"`
			TrailingStop Pips    `json:"trailingStop"`
		} `json:"tradeOpened"`
	}{}
	if err := json.Unmarshal(respBody, &o); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return &Order{
		ID:           o.TradeOpened.ID,
		Instrument:   o.Instrument,
		Units:        o.TradeOpened.Units,
		Side:         o.TradeOpened.Side,
		OrderType:    order.OrderType,
		Expiry:       order.Expiry,
		Price:        o.Price,
		LowerBound:   order.LowerBound,
		UpperBound:   order.UpperBound,
		StopLoss:     o.TradeOpened.StopLoss,
		TakeProfit:   o.TradeOpened.TakeProfit,
		TrailingStop: o.TradeOpened.TrailingStop,
	}, nil
}

func buildPostOrderBody(order Order) (string, error) {

	// check required items
	if len(order.Instrument) == 0 {
		return "", errors.New("instrument is required")
	}
	if len(order.Side) == 0 {
		return "", errors.New("side is required")
	}
	if len(order.OrderType) == 0 {
		return "", errors.New("order type is required")
	}
	if order.Expiry == nil && order.OrderType != OrderTypeMarket {
		return "", errors.New("expiry is required")
	}
	if order.Price <= Price(0) && order.OrderType != OrderTypeMarket {
		return "", errors.New("price is required")
	}

	// build form
	form := url.Values{}
	form.Add("instrument", order.Instrument)
	form.Add("side", string(order.Side))
	form.Add("type", string(order.OrderType))
	if order.OrderType != OrderTypeMarket {
		form.Add("expiry", order.Expiry.UTC().Format(time.RFC3339))
		form.Add("price", fmt.Sprintf(strconv.FormatFloat(float64(order.Price), 'f', -1, 32)))
	}
	if order.LowerBound > Price(0) {
		form.Add("lowerBound", fmt.Sprintf(strconv.FormatFloat(float64(order.LowerBound), 'f', -1, 32)))
	}
	if order.UpperBound > Price(0) {
		form.Add("upperBound", fmt.Sprintf(strconv.FormatFloat(float64(order.UpperBound), 'f', -1, 32)))
	}
	if order.StopLoss > Price(0) {
		form.Add("stopLoss", fmt.Sprintf(strconv.FormatFloat(float64(order.StopLoss), 'f', -1, 32)))
	}
	if order.TakeProfit > Price(0) {
		form.Add("takeProfit", fmt.Sprintf(strconv.FormatFloat(float64(order.TakeProfit), 'f', -1, 32)))
	}
	if order.TrailingStop > Pips(0) {
		form.Add("trailingStop", fmt.Sprintf(strconv.FormatFloat(float64(order.TrailingStop), 'f', -1, 32)))
	}
	return form.Encode(), nil
}

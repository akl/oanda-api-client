package oanda

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type orderInfo struct {
	ClientExtensions *struct {
		Comment string `json:"comment,omitempty"`
		ID      string `json:"id,omitempty"`
		Tag     string `json:"tag,omitempty"`
	} `json:"clientExtensions,omitempty"`
	TakeProfitOnFill *struct {
		Price       string `json:"price"`
		TimeInForce string `json:"timeInForce"`
	} `json:"takeProfitOnFill"`
	StopLossOnFill *struct {
		Price       string `json:"price"`
		TimeInForce string `json:"timeInForce"`
	} `json:"stopLossOnFill"`
	TrailingStopLossOnFill *struct {
		Distance    string `json:"distance"`
		TimeInForce string `json:"timeInForce"`
	} `json:"trailingStopLossOnFill"`
	CreateTime       *time.Time `json:"createTime"`
	ID               string     `json:"id"`
	Instrument       string     `json:"instrument,omitempty"`
	PartialFill      string     `json:"partialFill"`
	PositionFill     string     `json:"positionFill"`
	Price            string     `json:"price"`
	ReplacesOrderID  string     `json:"replacesOrderID,omitempty"`
	State            string     `json:"state"`
	TimeInForce      string     `json:"timeInForce"`
	GtdTime          *time.Time `json:"gtdTime"`
	TriggerCondition string     `json:"triggerCondition"`
	Type             string     `json:"type"`
	Units            string     `json:"units,omitempty"`
}

type retrievedOrders struct {
	LastTransactionID string      `json:"lastTransactionID"`
	Orders            []orderInfo `json:"orders,omitempty"`
}

type OrderPayloadBody struct {
	Units            int        `json:"units"`
	Instrument       string     `json:"instrument"`
	TimeInForce      string     `json:"timeInForce"`
	GtdTime          *time.Time `json:"gtdTime"`
	Type             string     `json:"type"`
	PositionFill     string     `json:"positionFill,omitempty"`
	Price            string     `json:"price,omitempty"`
	TakeProfitOnFill *onFillStr `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill   *onFillStr `json:"stopLossOnFill,omitempty"`
}

type OrderPayload struct {
	Order OrderPayloadBody `json:"order"`
}

type Order struct {
	TakeProfitOnFill       *onFill
	StopLossOnFill         *onFill
	TrailingStopLossOnFill *onFill
	CreateTime             *time.Time
	ID                     orderID
	Instrument             instrument
	PartialFill            string
	PositionFill           string
	Price                  Price
	State                  string
	TimeInForce            timeInForce
	GtdTime                *time.Time
	TriggerCondition       string
	Type                   orderType
	Units                  Unit
}

type onFill struct {
	Price       Price
	TimeInForce timeInForce
}

type onFillStr struct {
	TimeInForce string `json:"timeInForce,omitempty"`
	Price       string `json:"price,omitempty"` // must be a string for float precision
}

func (o *onFill) ToOnFillStr() *onFillStr {
	if o == nil {
		return nil
	}
	return &onFillStr{
		Price:       strconv.FormatFloat(float64(o.Price), 'f', 7, 64),
		TimeInForce: string(o.TimeInForce),
	}
}

func (o *Order) toOrderPayload() OrderPayload {
	return OrderPayload{
		OrderPayloadBody{
			Units:            int(o.Units),
			Instrument:       string(o.Instrument),
			TimeInForce:      string(o.TimeInForce),
			GtdTime:          o.GtdTime,
			Type:             string(o.Type),
			PositionFill:     o.PositionFill,
			Price:            o.Price.String(),
			TakeProfitOnFill: o.TakeProfitOnFill.ToOnFillStr(),
			StopLossOnFill:   o.StopLossOnFill.ToOnFillStr(),
		},
	}
}

func (r *retrievedOrders) toOrders() ([]Order, error) {
	var orders []Order
	for _, o := range r.Orders {
		var tpOnFill onFill
		var slOnFill onFill
		var trOnFill onFill
		if o.TakeProfitOnFill != nil {
			p, err := strconv.ParseFloat(o.TakeProfitOnFill.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse take plofit on fill price to float64: %v", err)
			}
			tpOnFill = onFill{
				Price(p),
				timeInForce(o.TakeProfitOnFill.TimeInForce),
			}
		}
		if o.StopLossOnFill != nil {
			p, err := strconv.ParseFloat(o.StopLossOnFill.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse take stop loss on fill price to float64: %v", err)
			}
			slOnFill = onFill{
				Price(p),
				timeInForce(o.StopLossOnFill.TimeInForce),
			}
		}
		if o.TrailingStopLossOnFill != nil {
			dp, err := strconv.ParseFloat(o.TrailingStopLossOnFill.Distance, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse trailing stop loss on fill distance to float64: %v", err)
			}
			trOnFill = onFill{
				Price(dp),
				timeInForce(o.TrailingStopLossOnFill.TimeInForce),
			}
		}
		p, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse order price to float64: %v", err)
		}
		u, err := strconv.Atoi(o.Units)
		if err != nil {
			return nil, fmt.Errorf("failed to parse units to int")
		}
		orders = append(orders, Order{
			&tpOnFill,
			&slOnFill,
			&trOnFill,
			o.CreateTime,
			orderID(o.ID),
			instrument(o.Instrument),
			o.PartialFill,
			o.PositionFill,
			Price(p),
			o.State,
			timeInForce(o.TimeInForce),
			o.GtdTime,
			o.TriggerCondition,
			orderType(o.Type),
			Unit(u),
		})
	}
	return orders, nil
}

func (c *Client) FetchOrders() ([]Order, error) {
	body, err := c.fetchOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	var ro retrievedOrders
	if err := json.Unmarshal(body, &ro); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %v", err)
	}
	o, err := ro.toOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to convert retrieved orders to type of Orders")
	}
	return o, nil
}

func (c *Client) FetchOrdersJSON() ([]byte, error) {
	return c.fetchOrders()
}

func (c *Client) UpdateOrder(order Order) error {
	body, err := json.Marshal(order.toOrderPayload())
	if err != nil {
		return fmt.Errorf("failed to marshal order payload to json: %v", err)
	}
	if err = c.updateOrder(order.ID, body); err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}
	return nil
}

func (c *Client) CreateOrder(order Order) error {
	body, err := json.Marshal(order.toOrderPayload())
	if err != nil {
		return fmt.Errorf("failed to marshal order payload to json: %v", err)
	}
	if err = c.createOrder(body); err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}
	return nil
}

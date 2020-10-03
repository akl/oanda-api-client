package oanda

import (
	"encoding/json"
	"fmt"
	"time"
)

type receivedTrades struct {
	LastTransactionID string `json:"lastTransactionID"`
	Trades            []struct {
		CurrentUnits string    `json:"currentUnits"`
		Financing    string    `json:"financing"`
		ID           string    `json:"id"`
		InitialUnits string    `json:"initialUnits"`
		Instrument   string    `json:"instrument"`
		OpenTime     time.Time `json:"openTime"`
		Price        string    `json:"price"`
		RealizedPL   string    `json:"realizedPL"`
		State        string    `json:"state"`
		UnrealizedPL string    `json:"unrealizedPL"`
	} `json:"trades"`
}

type Trade struct {
	ID       tradeID
	OpenTime *time.Time
}

func (r *receivedTrades) toTrades() []Trade {
	var trades []Trade
	for _, t := range r.Trades {
		trades = append(trades, Trade{
			ID:       tradeID(t.ID),
			OpenTime: &t.OpenTime,
		})
	}
	return trades
}

func (c *Client) FetchOpenTrades() ([]Trade, error) {
	body, err := c.fetchOpenTrades()
	if err != nil {
		return nil, fmt.Errorf("fariled to fetch open trades: %v", err)
	}
	var rt receivedTrades
	if err := json.Unmarshal(body, &rt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	return rt.toTrades(), nil
}

func (c *Client) FetchOpenTradesJSON() ([]byte, error){
	return c.fetchOpenTrades()
}

func (c *Client) CloseOpenTrade(id tradeID) error {
	body, err := json.Marshal(struct {
		Units string `json:"units"`
	}{Units: "ALL"})
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}
	if err := c.reduceTradeSize(id, body); err != nil {
		return fmt.Errorf("failed to close trade (id=%s): %v", string(id), err)
	}
	return nil
}

package oanda

import (
	"encoding/json"
	"fmt"
)

type AccountID int64

type Account struct {
	ID           AccountID `json:"accountId"`
	Name         string    `json:"accountName"`
	Balance      int64     `json:"balance"`
	UnrealizedPl int64     `json:"unrealizedPl"`
	RealizedPl   int64     `json:"realizedPl"`
	MarginUsed   int64     `json:"marginUsed"`
	MarginAvail  int64     `json:"marginAvail"`
	OpenTrades   int       `json:"openTrades"`
	OpenOrders   int       `json:"openOrders"`
	MarginRate   float64   `json:"marginRate"`
	Currency     string    `json:"accountCurrency"`
}

func (c *Client) ListAccountNames() ([]string, error) {
	body, err := c.fetchAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accout data: %v", err)
	}
	a := struct {
		Accounts []struct {
			ID   AccountID `json:"accountId"`
			Name string    `json:"accountName"`
		} `json:"accounts"`
	}{}
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	var result []string
	for _, account := range a.Accounts {
		result = append(result, account.Name)
	}
	return result, nil
}

// FetchAccountID fetches account id by account name.
// Returns nil if no match account found.
func (c *Client) FetchAccountID(accountName string) (*AccountID, error) {
	body, err := c.fetchAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accout data: %v", err)
	}
	a := struct {
		Accounts []struct {
			ID   AccountID `json:"accountId"`
			Name string    `json:"accountName"`
		} `json:"accounts"`
	}{}
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	for _, account := range a.Accounts {
		if accountName == accountName {
			return &account.ID, nil
		}
	}
	return nil, nil
}

// FetchAccount fetches account.
func (c *Client) FetchAccount(accountID AccountID) (*Account, error) {
	body, err := c.fetchAccountInfo(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account info: %v", err)
	}
	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return &account, nil
}

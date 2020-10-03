package oanda

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const authorizationPrefix = "Bearer "

// Client implements operations trade of oanda through OANDA API.
type Client struct {
	accountID       string
	client          *http.Client
	endpoint        string
	requiredHeaders http.Header
}

// NewClient constructs OANDA API client objects.
func NewClient(accountID, apiKey string, environment string) *Client {
	requiredHeaders := http.Header{}
	requiredHeaders.Add("Authorization", authorizationPrefix+apiKey)
	requiredHeaders.Add("Content-Type", "application/json")
	var endpoint string
	if environment == "Trade" {
		endpoint = "https://api-fxtrade.oanda.com"
	} else {
		endpoint = "https://api-fxpractice.oanda.com"
	}
	return &Client{
		accountID:       accountID,
		client:          &http.Client{},
		endpoint:        endpoint,
		requiredHeaders: requiredHeaders,
	}
}

func (c *Client) fetchOpenTrades() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint+"/v3/accounts/"+c.accountID+"/openTrades", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *Client) reduceTradeSize(id tradeID, body []byte) error {
	req, err := http.NewRequest(
		http.MethodPut,
		c.endpoint+"/v3/accounts/"+c.accountID+"/trades/"+string(id)+"/close",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}
	req.Header = c.requiredHeaders
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s: %s", resp.Status, respBody)
	}
	return nil
}

func (c *Client) fetchOrders() ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		c.endpoint+"/v3/accounts/"+c.accountID+"/orders",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *Client) updateOrder(orderID orderID, body []byte) error {
	req, err := http.NewRequest(
		http.MethodPut,
		c.endpoint+"/v3/accounts/"+c.accountID+"/orders/"+string(orderID),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("HTTP %s: %s", resp.Status, respBody)
	}
	return nil
}

func (c *Client) createOrder(body []byte) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.endpoint+"/v3/accounts/"+c.accountID+"/orders",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("HTTP %s: %s", resp.Status, respBody)
	}
	return nil
}

func (c *Client) cancelOrder(orderID orderID) error {
	req, err := http.NewRequest(
		http.MethodPut,
		c.endpoint+"/v3/accounts/"+c.accountID+"/orders/"+string(orderID)+"/cancel",
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	req.URL.RawQuery = req.URL.Query().Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s: %s", resp.Status, respBody)
	}
	return nil
}

func (c *Client) fetchOrderBook(instrument instrument, dateTime *time.Time) ([]byte, error) {
	url := c.endpoint + "/v3/instruments/" + string(instrument) + "/orderBook"
	if dateTime != nil {
		url = c.endpoint + "/v3/instruments/" + string(instrument) + "/orderBook?time=" + dateTime.UTC().Format(time.RFC3339Nano)
	}
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return body, nil
}

func safeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close: %v", err)
		}
	}
}

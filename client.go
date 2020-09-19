package oanda

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const authorizationPrefix = "Bearer "

// Client implements operations trade of oanda through OANDA API.
type Client struct {
	client          *http.Client
	endpoint        string
	requiredHeaders http.Header
}

// NewClient constructs OANDA API client objects.
func NewClient() *Client {
	requiredHeaders := http.Header{}
	requiredHeaders.Add("Authorization", authorizationPrefix+ParamOandaAPIKey.FetchValue())
	requiredHeaders.Add("Content-Type", "application/x-www-form-urlencoded")
	return &Client{
		client:          &http.Client{},
		endpoint:        ParamOandaEndpoint.FetchValue(),
		requiredHeaders: requiredHeaders,
	}
}

func (c *Client) fetchAccounts() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint+"/v1/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
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

func (c *Client) fetchAccountInfo(accountID AccountID) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint+"/v1/accounts/"+strconv.FormatInt(int64(accountID), 10), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
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

func (c *Client) fetchOrders(accountID int64, count int, instrument string) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		c.endpoint+"/v1/accounts/"+strconv.FormatInt(accountID, 10)+"/orders",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	req.Header = c.requiredHeaders
	req.URL.Query().Add("count", strconv.Itoa(count))
	req.URL.Query().Add("instrument", instrument)
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

func (c *Client) createOrder(accountID int64, reqBody string) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		c.endpoint+"/v1/accounts/"+strconv.FormatInt(accountID, 10)+"/orders",
		strings.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
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

func (c *Client) fetchOrderInfo(accountID int64, orderID int64) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		c.endpoint+"/v1/accounts/"+strconv.FormatInt(accountID, 10)+"/orders/"+strconv.FormatInt(orderID, 10),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
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

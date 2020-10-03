package oanda

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)



type book struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []bucket  `json:"buckets"`
}

type retrievedBook struct {
	Book book `json:"orderBook"`
}

type bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}
type OrderBook struct {
	Instrument instrument
	Time       time.Time
	Price      Price
	Buckets    []OrderBookBucket
}
type OrderBookBucket struct {
	Price             Price
	LongCountPercent  float64
	ShortCountPercent float64
}

func (b *book) toOrderBook() (*OrderBook, error) {
	price, err := strconv.ParseFloat(b.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price to float64: %v", err)
	}
	var buckets []OrderBookBucket
	for _, bu := range b.Buckets {
		p, err := strconv.ParseFloat(bu.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bucket price to type of float64: %v", err)
		}
		l, err := strconv.ParseFloat(bu.LongCountPercent, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse long count percent to float64: %v", err)
		}
		s, err := strconv.ParseFloat(bu.ShortCountPercent, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse short count percent to float64: %v", err)
		}
		buckets = append(buckets, OrderBookBucket{
			Price(p),
			l,
			s,
		})
	}
	return &OrderBook{
		instrument(b.Instrument),
		b.Time,
		Price(price),
		buckets,
	}, nil
}

func (o *OrderBook) ExtractBucketVicinityOfPrice(price Price, n int) (short, long []OrderBookBucket, err error) {
	var lowerBuckets []OrderBookBucket
	var higherBuckets []OrderBookBucket

	fmt.Printf("buckets.len = %d\n",len(o.Buckets))

	for i, b := range o.Buckets {
		if b.Price > price {
			lowerBuckets = o.Buckets[:i-1]
			higherBuckets = o.Buckets[i-1:]
			break
		}
	}

	fmt.Printf("lower.len = %d\n",len(lowerBuckets))
	fmt.Printf("higher.len = %d\n",len(higherBuckets))

	for i, j := 0, len(lowerBuckets)-1; i < j; i, j = i+1, j-1 {
		lowerBuckets[i], lowerBuckets[j] = lowerBuckets[j], lowerBuckets[i]
	}

	fmt.Printf("lower.len = %d\n",len(lowerBuckets))
	fmt.Printf("higher.len = %d\n",len(higherBuckets))

	if len(lowerBuckets[:n]) < n {
		return nil, nil, fmt.Errorf("price is too low: lowerBuckets[%d] is not exist", n-1)
	}
	if len(higherBuckets[:n]) < n {
		return nil, nil, fmt.Errorf("price is too high: higherBuckets[%d] is not exist", n-1)
	}
	return lowerBuckets[:n], higherBuckets[:n], nil
}

func (c *Client) FetchOrderBook(instrument instrument, dateTime *time.Time) (*OrderBook, error) {
	body, err := c.fetchOrderBook(instrument, dateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order book: %v", err)
	}
	var rb retrievedBook
	if err := json.Unmarshal(body, &rb); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %v", err)
	}
	ob, err := rb.Book.toOrderBook()
	if err != nil {
		return nil, fmt.Errorf("failed to convert book to order book: %v", err)
	}
	return ob, nil
}

func (c *Client) FetchOrderBookJSON(instrument instrument, dateTime *time.Time) ([]byte, error) {
	return c.fetchOrderBook(instrument, dateTime)
}

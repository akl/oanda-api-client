package main

import (
	"fmt"
	"log"

	"github.com/yuki-inoue-eng/oanda-api-client"
)

func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	book, err := client.FetchOrderBook(oanda.InstrumentUSDJPY, nil)
	if err != nil {
		log.Printf("failed to fetch order book: %v", err)
		return
	}
	s, l, err := book.ExtractBucketVicinityOfPrice(106, 3)
	if err != nil {
		log.Printf("failed to get vop from order book: %v", err)
		return
	}
	fmt.Printf("s.len = %d, l.len = %d\n", len(s), len(l))
	fmt.Printf("s=%v\n\n", s)
	fmt.Printf("l=%v\n\n", l)
}

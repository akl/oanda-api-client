package main

import (
	"fmt"
	"log"

	"github.com/yuki-inoue-eng/oanda-api-client"
)

func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	bytes, err := client.FetchOpenTradesJSON()
	if err != nil {
		log.Printf("failed to close trade: %v", err)
		return
	}
	fmt.Print(string(bytes))
}
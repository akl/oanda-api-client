package main

import (
	"log"

	"github.com/yuki-inoue-eng/oanda-api-client"
)

func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	err := client.CloseOpenTrade("1")
	if err != nil {
		log.Printf("failed to close trade: %v", err)
		return
	}
}

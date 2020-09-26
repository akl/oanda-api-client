package main

import (
	"fmt"
	"github.com/yuki-inoue-eng/oanda-api-client"
	"log"
)



func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	bytes, err  := client.FetchOrdersJSON()
	if err != nil {
		log.Printf("failed to fetch orders: %v", err)
		return
	}
	fmt.Print(string(bytes))
}

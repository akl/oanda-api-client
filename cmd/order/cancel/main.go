package main

import (
	"fmt"
	"github.com/yuki-inoue-eng/oanda-api-client"
	"log"
)

// example
func main() {
	client := oanda.NewClient(oanda.ParamOandaAccountID.FetchValue(), oanda.ParamOandaAPIKey.FetchValue(), "Practice")
	if err := client.CancelOrder("21"); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("order canceled")
}
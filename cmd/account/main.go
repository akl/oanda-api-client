package account

import (
	"fmt"
	"github.com/yuki-inoue-eng/oanda-api-client"
	"log"
)

func main() {
	client := oanda.NewClient()
	accountNames, err := client.ListAccountNames()
	if err != nil {
		log.Printf("failed to fetch account names: %v", err)
		return
	}
	fmt.Printf("account name list: %v", accountNames)
}

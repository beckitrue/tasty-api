package jsondecode

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type Accounts struct {
	Data struct {
		Items []struct {
			Account struct {
				AccountNumber        string `json:"account-number"`
				OpenedAt             string `json:"opened-at"`
				Nickname             string `json:"nickname"`
				AccountTypeName      string `json:"account-type-name"`
				DayTraderStatus      bool   `json:"day-trader-status"`
				IsClosed             bool   `json:"is-closed"`
				IsFirmError          bool   `json:"is-firm-error"`
				IsFirmProprietary    bool   `json:"is-firm-proprietary"`
				IsFuturesApproved    bool   `json:"is-futures-approved"`
				IsTestDrive          bool   `json:"is-test-drive"`
				MarginOrCash         string `json:"margin-or-cash"`
				IsForeign            bool   `json:"is-foreign"`
				InvestmentObjective  string `json:"investment-objective"`
				SuitableOptionsLevel string `json:"suitable-options-level"`
				CreatedAt            string `json:"created-at"`
			} `json:"account"`
			AuthorityLevel string `json:"authority-level"`
		} `json:"items"`
	} `json:"data"`
	Context string `json:"context"`
}

func PrintDataAccounts(data string) {
	dec := json.NewDecoder(strings.NewReader(data))
	for {
		var accounts Accounts
		if err := dec.Decode(&accounts); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// loop through the items and print the account numbers
		// TODO: include more information in the output
		for i := range accounts.Data.Items {
			fmt.Printf("Account %d: %s\n", i, accounts.Data.Items[i].Account.AccountNumber)
		}
	}
}

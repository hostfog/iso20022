package iso20022_test

import (
	"fmt"
	"time"

	iso "github.com/hostfog/iso20022"
	"github.com/hostfog/iso20022/pacs008"
)

func Example() {
	now := time.Now()

	msg := &pacs008.Message{
		Header: pacs008.GroupHeader{
			MessageID:        "HFPAY-1",
			CreatedAt:        iso.DateTime{Time: now},
			SettlementMethod: "INDA",
		},
		Transactions: []pacs008.Transaction{{
			EndToEndID:       "E2E-1",
			UETR:             iso.NewUETR(),
			SettlementAmount: iso.Amount{Currency: "INR", Value: "27613.02"},
			InstructedAmount: &iso.Amount{Currency: "GBP", Value: "250.00"},
			ExchangeRate:     "110.4521",
			ChargeBearer:     "SHAR",
			Debtor:           iso.Party{Name: "Param Singh"},
			DebtorAgent:      iso.Agent{BIC: "HFPYGB2L"},
			CreditorAgent:    iso.Agent{BIC: "HDFCINBB"},
			Creditor:         iso.Party{Name: "R Kaur"},
			CreditorAccount:  &iso.Account{Other: "50100123456789"},
		}},
	}

	_, err := pacs008.Marshal(msg)
	fmt.Println(err)
	// Output: <nil>
}

# iso20022

Build, parse and validate ISO 20022 payment messages in Go, with readable field names and no code generation.

```sh
go get github.com/hostfog/iso20022
```

| Package   | Message            | Use                                         |
|-----------|--------------------|---------------------------------------------|
| `pain001` | pain.001.001.09    | Customer → bank: initiate credit transfers  |
| `pacs008` | pacs.008.001.08    | Bank → bank: customer credit transfer (FX)  |

## Example

```go
msg := &pain001.Message{
	Header: pain001.GroupHeader{
		MessageID:       "MSG-0001",
		CreatedAt:       iso.DateTime{Time: time.Now()},
		InitiatingParty: iso.Party{Name: "HostFog Ltd"},
	},
	PaymentInfos: []pain001.PaymentInfo{{
		ID:            "BATCH-1",
		ExecutionDate: iso.Date{Time: time.Now()},
		Debtor:        iso.Party{Name: "HostFog Ltd"},
		DebtorAccount: iso.Account{IBAN: "GB82WEST12345698765432"},
		DebtorAgent:   iso.Agent{BIC: "WESTGB2L"},
		Transactions: []pain001.Transaction{{
			EndToEndID:      "INV-100",
			Amount:          iso.Amount{Currency: "GBP", Value: "250.00"},
			Creditor:        iso.Party{Name: "Acme Ltd"},
			CreditorAccount: iso.Account{Other: "40051512345678"},
			CreditorAgent:   &iso.Agent{ClearingSystem: "GBDSC", MemberID: "400515"},
		}},
	}},
}

xmlBytes, err := pain001.Marshal(msg)
```

`Marshal` fills in `NbOfTxs`, `CtrlSum` and `PmtMtd`, validates, and returns the XML. It never modifies your struct.

`Parse` decodes a document and checks the namespace. It does not validate, so call `Validate()` on input you don't trust.

## Validation

Errors are joined `*iso20022.FieldError` values with ISO element paths:

```
PmtInf[0].DbtrAcct.Id.IBAN: IBAN checksum failed
PmtInf[0].CdtTrfTxInf[0].Cdtr.Nm: is required
```

Checks include IBAN (format, country length, mod-97), BIC, UETR, currency codes, amount format (max 18 digits, 5 decimals), text lengths, code lists, and count and sum consistency.

## Helpers

- `iso20022.NewUETR()`: a UUIDv4 for SWIFT gpi tracking
- `iso20022.ValidateIBAN`, `iso20022.ValidateBIC`

## Scope

v0.1 covers the commonly used elements of each message. Market practice rules (CBPR+, SEPA, UK Faster Payments) add their own constraints, so check your scheme's guidelines before going live.

## License

MIT
# iso20022

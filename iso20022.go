// Package iso20022 holds the building blocks shared by the message packages
// pain001 and pacs008: parties, accounts, agents, amounts and dates.
//
// Go field names are readable; XML tags use the ISO 20022 element names, so
// validation errors report paths such as "PmtInf[0].Dbtr.Nm".
package iso20022

import (
	"regexp"
	"strings"
	"time"

	"github.com/hostfog/iso20022/internal/check"
)

// FieldError is the type of each error returned by a Validate method.
// Use errors.As to inspect it.
type FieldError = check.FieldError

var (
	countryRe  = regexp.MustCompile(`^[A-Z]{2}$`)
	currencyRe = regexp.MustCompile(`^[A-Z]{3}$`)
	amountRe   = regexp.MustCompile(`^[0-9]+(\.[0-9]{1,5})?$`)
)

type Party struct {
	Name    string         `xml:"Nm,omitempty"`
	Address *PostalAddress `xml:"PstlAdr,omitempty"`
}

func (p Party) Validate() error {
	var e check.Errors
	e.Text("Nm", p.Name, true, 140)
	if p.Address != nil {
		e.Nest("PstlAdr", p.Address.Validate())
	}
	return e.Err()
}

type PostalAddress struct {
	Street         string   `xml:"StrtNm,omitempty"`
	BuildingNumber string   `xml:"BldgNb,omitempty"`
	PostCode       string   `xml:"PstCd,omitempty"`
	Town           string   `xml:"TwnNm,omitempty"`
	Country        string   `xml:"Ctry,omitempty"`
	Lines          []string `xml:"AdrLine,omitempty"`
}

func (a PostalAddress) Validate() error {
	var e check.Errors
	e.Text("StrtNm", a.Street, false, 70)
	e.Text("BldgNb", a.BuildingNumber, false, 16)
	e.Text("PstCd", a.PostCode, false, 16)
	e.Text("TwnNm", a.Town, false, 35)
	if a.Country != "" && !countryRe.MatchString(a.Country) {
		e.Fail("Ctry", "must be a two-letter ISO 3166 code")
	}
	if len(a.Lines) > 7 {
		e.Fail("AdrLine", "at most 7 lines allowed")
	}
	for _, l := range a.Lines {
		e.Text("AdrLine", l, false, 70)
	}
	return e.Err()
}

// Account is identified by an IBAN or, where IBANs aren't used, another
// identifier such as a UK sort code and account number.
type Account struct {
	IBAN     string
	Other    string
	Currency string
}

func (a Account) Validate() error {
	var e check.Errors
	switch {
	case a.IBAN != "" && a.Other != "":
		e.Fail("Id", "set either IBAN or Othr, not both")
	case a.IBAN != "":
		if err := ValidateIBAN(a.IBAN); err != nil {
			e.Fail("Id.IBAN", "%s", err)
		}
	case a.Other != "":
		e.Text("Id.Othr.Id", a.Other, true, 34)
	default:
		e.Fail("Id", "is required")
	}
	if a.Currency != "" && !currencyRe.MatchString(a.Currency) {
		e.Fail("Ccy", "must be a three-letter ISO 4217 code")
	}
	return e.Err()
}

// Agent is a financial institution, identified by BIC or clearing system
// membership (e.g. ClearingSystem "GBDSC" with a sort code as MemberID).
type Agent struct {
	BIC            string
	ClearingSystem string
	MemberID       string
	Name           string
}

func (a Agent) Validate() error {
	var e check.Errors
	if a.BIC == "" && a.MemberID == "" {
		e.Fail("FinInstnId", "BICFI or ClrSysMmbId is required")
	}
	if a.BIC != "" {
		if err := ValidateBIC(a.BIC); err != nil {
			e.Fail("FinInstnId.BICFI", "%s", err)
		}
	}
	if a.ClearingSystem != "" && a.MemberID == "" {
		e.Fail("FinInstnId.ClrSysMmbId.MmbId", "is required with ClrSysId")
	}
	e.Text("FinInstnId.ClrSysMmbId.ClrSysId.Cd", a.ClearingSystem, false, 5)
	e.Text("FinInstnId.ClrSysMmbId.MmbId", a.MemberID, false, 35)
	e.Text("FinInstnId.Nm", a.Name, false, 140)
	return e.Err()
}

// Amount is a decimal value with its currency, e.g. {"GBP", "250.00"}.
type Amount struct {
	Currency string `xml:"Ccy,attr"`
	Value    string `xml:",chardata"`
}

func (a Amount) Validate() error {
	var e check.Errors
	if !currencyRe.MatchString(a.Currency) {
		e.Fail("Ccy", "must be a three-letter ISO 4217 code")
	}
	switch {
	case !amountRe.MatchString(a.Value):
		e.Fail("Amt", "invalid amount %q", a.Value)
	case len(strings.Replace(a.Value, ".", "", 1)) > 18:
		e.Fail("Amt", "more than 18 digits")
	case isZero(a.Value):
		e.Fail("Amt", "must be greater than zero")
	}
	return e.Err()
}

// Date is an ISO date (2006-01-02).
type Date struct{ time.Time }

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.Format(time.DateOnly)), nil
}

func (d *Date) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.DateOnly, string(b))
	d.Time = t
	return err
}

// DateTime is an ISO date-time. It is written with its UTC offset.
type DateTime struct{ time.Time }

func (d DateTime) MarshalText() ([]byte, error) {
	return []byte(d.Format(time.RFC3339)), nil
}

func (d *DateTime) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999", string(b))
	}
	d.Time = t
	return err
}

func isZero(v string) bool {
	for _, r := range v {
		if r != '0' && r != '.' {
			return false
		}
	}
	return true
}

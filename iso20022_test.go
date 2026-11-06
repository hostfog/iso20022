package iso20022

import (
	"errors"
	"testing"
)

func TestValidateIBAN(t *testing.T) {
	valid := []string{
		"GB82WEST12345698765432",
		"DE89370400440532013000",
		"FR1420041010050500013M02606",
		"NL91ABNA0417164300",
	}
	for _, s := range valid {
		if err := ValidateIBAN(s); err != nil {
			t.Errorf("ValidateIBAN(%q) = %v", s, err)
		}
	}

	invalid := []string{
		"GB82WEST12345698765431",   // checksum
		"GB82WEST1234569876543",    // length
		"GB82 WEST 1234 5698 7654", // spaces
		"gb82west12345698765432",   // case
		"",
	}
	for _, s := range invalid {
		if ValidateIBAN(s) == nil {
			t.Errorf("ValidateIBAN(%q) should fail", s)
		}
	}
}

func TestValidateBIC(t *testing.T) {
	for _, s := range []string{"BARCGB22", "HBUKGB4BXXX", "DEUTDEFF500"} {
		if err := ValidateBIC(s); err != nil {
			t.Errorf("ValidateBIC(%q) = %v", s, err)
		}
	}
	for _, s := range []string{"BARCGB2", "barcgb22", "BARC1222", "BARCGB22XX"} {
		if ValidateBIC(s) == nil {
			t.Errorf("ValidateBIC(%q) should fail", s)
		}
	}
}

func TestNewUETR(t *testing.T) {
	seen := map[string]bool{}
	for range 100 {
		u := NewUETR()
		if !ValidUETR(u) {
			t.Fatalf("invalid UETR %q", u)
		}
		if seen[u] {
			t.Fatalf("duplicate UETR %q", u)
		}
		seen[u] = true
	}
}

func TestAmountValidate(t *testing.T) {
	tests := []struct {
		amt Amount
		ok  bool
	}{
		{Amount{"GBP", "250.00"}, true},
		{Amount{"GBP", "1"}, true},
		{Amount{"GBP", "0.00001"}, true},
		{Amount{"GBP", "0.00"}, false},
		{Amount{"GBP", "-1"}, false},
		{Amount{"GBP", "1.123456"}, false},
		{Amount{"GBP", "1,000"}, false},
		{Amount{"gbp", "1"}, false},
		{Amount{"GBP", "1234567890123456789"}, false},
	}
	for _, tt := range tests {
		if err := tt.amt.Validate(); (err == nil) != tt.ok {
			t.Errorf("%+v: err = %v, want ok %v", tt.amt, err, tt.ok)
		}
	}
}

func TestAgentAndAccount(t *testing.T) {
	sortCode := Agent{ClearingSystem: "GBDSC", MemberID: "400515"}
	if err := sortCode.Validate(); err != nil {
		t.Errorf("sort code agent: %v", err)
	}
	if (Agent{}).Validate() == nil {
		t.Error("empty agent should fail")
	}
	if (Account{IBAN: "GB82WEST12345698765432", Other: "1"}).Validate() == nil {
		t.Error("account with IBAN and Othr should fail")
	}

	err := Party{Address: &PostalAddress{Country: "gb"}}.Validate()
	var fe *FieldError
	if !errors.As(err, &fe) || fe.Path != "Nm" {
		t.Errorf("want FieldError on Nm, got %v", err)
	}
}

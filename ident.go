package iso20022

import (
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
)

var (
	ibanRe = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{1,30}$`)
	bicRe  = regexp.MustCompile(`^[A-Z0-9]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	uetrRe = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
)

// IBAN lengths for countries this package checks strictly. Other countries
// are checked for format and checksum only.
var ibanLength = map[string]int{
	"AE": 23, "AT": 20, "BE": 16, "CH": 21, "CY": 28, "CZ": 24, "DE": 22,
	"DK": 18, "EE": 20, "ES": 24, "FI": 18, "FR": 27, "GB": 22, "GR": 27,
	"HR": 21, "HU": 28, "IE": 22, "IT": 27, "LT": 20, "LU": 20, "LV": 21,
	"MT": 31, "NL": 18, "NO": 15, "PK": 24, "PL": 28, "PT": 25, "QA": 29,
	"RO": 24, "SA": 24, "SE": 24, "SI": 19, "SK": 24, "TR": 26,
}

// ValidateIBAN checks format, known country length and the mod-97 checksum.
// The IBAN must be in electronic form: upper case, no spaces.
func ValidateIBAN(iban string) error {
	if !ibanRe.MatchString(iban) {
		return errors.New("invalid IBAN format")
	}
	if n, ok := ibanLength[iban[:2]]; ok && len(iban) != n {
		return fmt.Errorf("%s IBAN must be %d characters", iban[:2], n)
	}

	rearranged := iban[4:] + iban[:4]
	rem := 0
	for _, r := range rearranged {
		if r >= 'A' {
			rem = (rem*100 + int(r-'A'+10)) % 97
		} else {
			rem = (rem*10 + int(r-'0')) % 97
		}
	}
	if rem != 1 {
		return errors.New("IBAN checksum failed")
	}
	return nil
}

// ValidateBIC checks an 8 or 11 character BIC.
func ValidateBIC(bic string) error {
	if !bicRe.MatchString(bic) {
		return fmt.Errorf("invalid BIC %q", bic)
	}
	return nil
}

// NewUETR returns a random UUIDv4 for end-to-end payment tracking (SWIFT gpi).
func NewUETR() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ValidUETR reports whether s is a well-formed UETR.
func ValidUETR(s string) bool {
	return uetrRe.MatchString(s)
}

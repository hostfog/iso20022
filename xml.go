package iso20022

import "encoding/xml"

// encoding/xml writes empty parents for "a>b,omitempty" tags, so optional
// groups are encoded through these explicit shapes instead.

type accountXML struct {
	ID struct {
		IBAN  string   `xml:"IBAN,omitempty"`
		Other *otherID `xml:"Othr,omitempty"`
	} `xml:"Id"`
	Currency string `xml:"Ccy,omitempty"`
}

type otherID struct {
	ID string `xml:"Id"`
}

func (a Account) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var w accountXML
	w.ID.IBAN = a.IBAN
	if a.Other != "" {
		w.ID.Other = &otherID{ID: a.Other}
	}
	w.Currency = a.Currency
	return e.EncodeElement(w, start)
}

func (a *Account) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var w accountXML
	if err := d.DecodeElement(&w, &start); err != nil {
		return err
	}
	*a = Account{IBAN: w.ID.IBAN, Currency: w.Currency}
	if w.ID.Other != nil {
		a.Other = w.ID.Other.ID
	}
	return nil
}

type agentXML struct {
	FinInstnID struct {
		BIC      string     `xml:"BICFI,omitempty"`
		Clearing *memberXML `xml:"ClrSysMmbId,omitempty"`
		Name     string     `xml:"Nm,omitempty"`
	} `xml:"FinInstnId"`
}

type memberXML struct {
	System   Code   `xml:"ClrSysId,omitempty"`
	MemberID string `xml:"MmbId"`
}

func (a Agent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var w agentXML
	w.FinInstnID.BIC = a.BIC
	w.FinInstnID.Name = a.Name
	if a.MemberID != "" {
		w.FinInstnID.Clearing = &memberXML{System: Code(a.ClearingSystem), MemberID: a.MemberID}
	}
	return e.EncodeElement(w, start)
}

func (a *Agent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var w agentXML
	if err := d.DecodeElement(&w, &start); err != nil {
		return err
	}
	*a = Agent{BIC: w.FinInstnID.BIC, Name: w.FinInstnID.Name}
	if c := w.FinInstnID.Clearing; c != nil {
		a.ClearingSystem = string(c.System)
		a.MemberID = c.MemberID
	}
	return nil
}

// Code is an external code list value, encoded as <Cd>value</Cd>.
type Code string

func (c Code) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		Cd string `xml:"Cd"`
	}{string(c)}, start)
}

func (c *Code) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var w struct {
		Cd string `xml:"Cd"`
	}
	err := d.DecodeElement(&w, &start)
	*c = Code(w.Cd)
	return err
}

// Remittance is unstructured remittance text, encoded as <Ustrd>.
type Remittance string

func (r Remittance) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		Ustrd string `xml:"Ustrd"`
	}{string(r)}, start)
}

func (r *Remittance) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var w struct {
		Ustrd string `xml:"Ustrd"`
	}
	err := d.DecodeElement(&w, &start)
	*r = Remittance(w.Ustrd)
	return err
}

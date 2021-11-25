package main

// DCC (Digital Covid Certificate) Top Level CBOR structure
// Section 3.3.1 at https://ec.europa.eu/health/sites/default/files/ehealth/docs/digital-green-certificates_v1_en.pdf
type DCC struct {
	ExpirationTime    int    `cbor:"4,keyasint,omitempty"`
	IssuedAt          int    `cbor:"6,keyasint,omitempty"`
	Issuer            string `cbor:"1,keyasint,omitempty"`
	HealthCertificate HCert  `cbor:"-260,keyasint,omitempty"`
}

type HCert struct {
	DGC DCCPayload `cbor:"1,keyasint,omitempty"`
}

type DCCPayload struct {
	Nam Nam    `json:"nam"`
	Dob string `json:"dob"` // Date of birth in range of 1900-0-01 to 2099-12-31
	V   []V    `json:"v"`   // Vaccination group
	Ver string `json:"ver"`
}

type Nam struct {
	Fn  string `json:"fn"`  // Surname
	Fnt string `json:"fnt"` // Standardised Surname
	Gn  string `json:"gn"`  // Forename
	Gnt string `json:"gnt"` // Standardised Forename
}

// DCC Payload for Vaccination group type of HCERT
type V struct {
	Dn int    `json:"dn"` // Number in a series of doses
	Sd int    `json:"sd"` // The overall number of doses
	Mp string `json:"mp"` // Vaccine Product
	Dt string `json:"dt"` // Date of vaccination
	Tg string `json:"tg"` // Disease or agent targeted
	Vp string `json:"vp"` // Vaccine or Prophylaxis
	Ma string `json:"ma"` // Vaccine marketing authorisation holder or manufacturer
	Co string `json:"co"` // Member state which administered the vaccine
	Is string `json:"is"` // Certificate Issuer
	Ci string `json:"ci"` // Unique certificate identifier
}

var (
	GreenPassVersion       = "1.3.0"
	SNOMEDCovidCode        = "840539006"
	VaccineProduct         = "EU/1/20/1528"
	VaccineType            = "1119349007"    // SARS-CoV-2 mRNA vaccine
	MarketingAuthorisation = "ORG-100030215" // Biontech Manufacturing GmbH
)

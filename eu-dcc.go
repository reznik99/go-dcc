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
	Dob string `json:"dob"`
	V   []V    `json:"v"`
	Ver string `json:"ver"`
}

type Nam struct {
	Fn  string `json:"fn"`
	Fnt string `json:"fnt"`
	Gn  string `json:"gn"`
	Gnt string `json:"gnt"`
}

// DCC Payload for Vaccination type of HCERT
type V struct {
	Dn int    `json:"dn"`
	Sd int    `json:"sd"`
	Mp string `json:"mp"`
	Dt string `json:"dt"`
	Tg string `json:"tg"`
	Vp string `json:"vp"`
	Ma string `json:"ma"`
	Co string `json:"co"`
	Is string `json:"is"`
	Ci string `json:"ci"`
}

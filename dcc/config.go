package dcc

// DCC Payload data to be read from file
type Config struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Dob             string `json:"dob"`
	IssuerCountry   string `json:"issuerCountry"`
	Issuer          string `json:"issuer"`
	VaccinationDate string `json:"vaccinationDate"`
	Doses           int    `json:"doses"`
}

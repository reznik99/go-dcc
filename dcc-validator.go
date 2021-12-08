package dcc

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"

	"github.com/veraison/go-cose"
)

var (
	API_BASE_URL = "https://get.dgc.gov.it/v1/dgc"
)

// Verify verifies a Vaccine Passport, it reads the file of type fileType
func Verify(filePath string) (valid bool, err error) {
	// Parse certificate from raw or QR code, returning raw cose and parsed dcc payload
	_, raw, err := Parse(filePath)
	if err != nil {
		return
	}

	// Fetch KIDs to compare
	kidsList, err := fetchValidKIDs()
	if err != nil {
		return
	}

	// Write PEMs into KIDs map
	kids, err := fetchCerts(kidsList)
	if err != nil {
		return
	}

	// Extract KID from Passports Header
	kidBytes, err := extractHeaderBytes(raw, "kid")
	if err != nil {
		return
	}
	var kid = base64.StdEncoding.EncodeToString(kidBytes)

	// Check KID is in trusted list
	if _, ok := kids[kid]; !ok {
		return
	}

	// Parse Signer certificate belonging to KID. Extract public key
	pemCertificate := kids[kid]
	block, _ := pem.Decode([]byte(pemCertificate))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}

	publicKey := cert.PublicKey.(crypto.PublicKey)

	// Print extra information for debugging
	toBeSigned, _ := raw.SigStructure(nil)
	digest := sha256.Sum256(toBeSigned)

	verifier := cose.Verifier{
		PublicKey: publicKey,
		Alg:       cose.ES256,
	}

	// Verify the Vaccine Passport's Signature
	err = verifier.Verify(digest[:], raw.Signature)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// fetchCerts fetches all the Signer Certificates to be used to validate International Vaccine Passports
func fetchCerts(kids []string) (map[string]string, error) {
	var kidsMap = map[string]string{}
	var token = "0"

	for range kids {
		// Generate request wit headers
		req, err := http.NewRequest(http.MethodGet, API_BASE_URL+"/signercertificate/update", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("x-resume-token", token)

		// Dispatch request TODO: Multi-thread this or cache it
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		// Read response
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Store PEM with appended Headers for x509 parsing
		token = resp.Header.Get("x-resume-token")
		respKID := resp.Header.Get("x-kid")
		kidsMap[respKID] = "-----BEGIN CERTIFICATE-----\n" + string(bodyBytes) + "\n-----END CERTIFICATE-----"
	}

	return kidsMap, nil
}

// fetchValidKIDs fetches all KIDs from Italian DGC site to compare against Greenpass header
func fetchValidKIDs() ([]string, error) {
	resp, err := http.Get(API_BASE_URL + "/signercertificate/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kids = []string{}
	err = json.Unmarshal(bodyBytes, &kids)
	if err != nil {
		return nil, err
	}

	return kids, err
}

// extractHeaderBytes extracts header bytes with given tag name, from Protected or Unprotected header in cose signature object
func extractHeaderBytes(raw *cose.Sign1Message, headerName string) ([]byte, error) {
	var tag, err = cose.GetCommonHeaderTag(headerName)
	if err != nil {
		return nil, err
	}

	var dccKid = raw.Headers.Protected[tag]
	if _, ok := dccKid.([]byte); !ok {
		dccKid = raw.Headers.Unprotected[tag]
		if _, ok = dccKid.([]byte); !ok {
			return nil, errors.New("ERROR: Couldn't extract KID from Vaccine Passport")
		}
	}

	return dccKid.([]byte), nil
}

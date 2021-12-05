package dcc

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/veraison/go-cose"
)

var (
	API_BASE_URL = "https://get.dgc.gov.it/v1/dgc"
)

func VerifyGreenpass(filePath string, fileType int) {

	dcc, raw, err := ParseGreenpass(filePath, fileType)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Fetch KIDs to compare
	kidsList := fetchValidKIDs()
	fmt.Printf("Fetched %d KIDs from %s\n", len(kidsList), API_BASE_URL)

	// Write PEMs into KIDs map
	// kids := map[string]string{}
	kids := fetchCerts(kidsList)

	// Extract KID from Passports Header
	tag, err := cose.GetCommonHeaderTag("kid")
	if err != nil {
		panic(err)
	}

	dccKid := raw.Headers.Protected[tag]
	if _, ok := dccKid.([]byte); !ok {
		dccKid = raw.Headers.Unprotected[tag]
		if _, ok = dccKid.([]byte); !ok {
			panic(errors.New("ERROR: Couldn't extract KID from Vaccine Passport"))
		}
	}
	var dccKidName = base64.StdEncoding.EncodeToString(dccKid.([]byte))

	// Check KID is in trusted list
	if _, ok := kids[dccKidName]; !ok {
		fmt.Printf("vaccine Pass KID %s not found on trusted list at: %s\n", dccKidName, API_BASE_URL)
	}
	fmt.Printf("Vaccine Pass KID '%s' is trusted\n", dccKidName)

	// Parse appropriate PEM of KID
	pemCertificate := kids[dccKidName]
	fmt.Printf("Vaccine Pass Signer Certificate PEM: \n%s\n", pemCertificate)
	block, _ := pem.Decode([]byte(pemCertificate))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	// Extract Public Key from Signer Certificate
	switch cert.PublicKeyAlgorithm {
	case x509.RSA:
		fmt.Println("Signer certificate that signed this Vaccine Passport is of type: RSA")
		return
	case x509.ECDSA:
		fmt.Println("Signer certificate that signed this Vaccine Passport is of type: ECDSA")
	}
	ecdsaPublic := cert.PublicKey.(*ecdsa.PublicKey)

	// Extra information for debugging
	toBeSigned, _ := raw.SigStructure(nil)
	digest := sha256.Sum256(toBeSigned)
	fmt.Printf("Digest: %s\n", hexutil.Encode(digest[:]))
	fmt.Printf("Signature: %s\n", hexutil.Encode(raw.Signature))

	// Work around to get Algorithm struct
	// (Curve is not exported and therefore cannot be initialised outside the go-cose library)
	signer, _ := cose.NewSigner(cose.ES256, nil)

	// Verify the Vaccine Passport's Signature
	verifier := cose.Verifier{
		PublicKey: ecdsaPublic,
		Alg:       signer.GetAlg(),
	}
	err = raw.Verify(nil, verifier)
	if err != nil {
		fmt.Printf("Verification FAILED with err: %s\n", err.Error())
		return
	}
	fmt.Printf("%s's Vaccine Passport has Signature Validated succesfully\n", dcc.HealthCertificate.DGC.Nam.Gn)
}

// fetchCerts fetches all the Signer Certificates to be used to validate International Vaccine Passports
func fetchCerts(kids []string) (kidsMap map[string]string) {
	kidsMap = map[string]string{}
	var token = "0"
	for idx := range kids {
		// Generate request wit headers
		req, err := http.NewRequest(http.MethodGet, API_BASE_URL+"/signercertificate/update", nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("x-resume-token", token)

		// Dispatch request TODO: Multi-thread this or cache it
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		// Read response
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// Store PEM with appended Headers for x509 parsing
		token = resp.Header.Get("x-resume-token")
		respKID := resp.Header.Get("x-kid")
		kidsMap[respKID] = "-----BEGIN CERTIFICATE-----\n" + string(bodyBytes) + "\n-----END CERTIFICATE-----"

		fmt.Printf("\rFetched %d/%d Signer Certificates", idx, len(kids))
	}
	fmt.Printf("Fetched %d/%d Signer Certificates", len(kids), len(kids))
	return
}

// fetchValidKIDs fetches all KIDs from Italian DGC site to compare against Greenpass header
func fetchValidKIDs() []string {
	resp, err := http.Get(API_BASE_URL + "/signercertificate/status")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	kids := []string{}
	err = json.Unmarshal(bodyBytes, &kids)
	if err != nil {
		panic(err)
	}

	return kids
}

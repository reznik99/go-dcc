package dcc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	API_BASE_URL = "https://get.dgc.gov.it/v1/dgc"
)

func VerifyGreenpass(filePath string, fileType int) {

	_, err := ParseGreenpass(filePath, fileType)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Fetch KIDs to compare
	kidsList := fetchValidKIDs()
	fmt.Printf("Fetched %d KIDs from %s\n", len(kidsList), API_BASE_URL)

	// Write PEMs into KIDs map
	// kids := fetchCerts(kidsList)

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

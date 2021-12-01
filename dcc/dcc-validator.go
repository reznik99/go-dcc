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

	// Validation logic

	kids := fetchValidKIDs()

	// TODO: Fetch certificates/signers from EU-DCC API and check against KID in Protected Header

}

// fetchValidKIDs fetches all KIDs from Italian DGC site to compare against Greenpass header
func fetchValidKIDs() map[string]bool {
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

	kidsMap := map[string]bool{}
	for _, kid := range kids {
		fmt.Println(kid)
		kidsMap[kid] = true
	}

	return kidsMap
}

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

type KIDsList struct {
	KIDs []string
}

func VerifyGreenpass(filePath string, fileType int) {

	_, err := ParseGreenpass(filePath, fileType)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Validation logic
	// TODO: Fetch certificates/signers from EU-DCC API and check against KID in Protected Header

	kids := fetchValidKIDs()
	for _, kid := range kids.KIDs {
		fmt.Println(kid)
	}
}

func fetchValidKIDs() *KIDsList {
	resp, err := http.Get(API_BASE_URL + "/signercertificate/status")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	kids := new(KIDsList)
	err = json.Unmarshal(bodyBytes, kids)
	if err != nil {
		panic(err)
	}
	return kids
}

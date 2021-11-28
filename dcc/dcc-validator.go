package dcc

import "fmt"

func VerifyGreenpass(filePath string, fileType int) {

	_, err := ParseGreenpass(filePath, fileType)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Validation logic
	// TODO: Fetch certificates/signers from EU-DCC API and check against KID in Protected Header

}

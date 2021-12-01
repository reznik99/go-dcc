package dcc

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dasio/base45"
	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

func ParseGreenpass(path string, fileType int) (payload *DCC, headers *cose.Headers, messageRaw *cose.Sign1Message, signature []byte, err error) {

	if fileType == TypeQRCode {
		err = fmt.Errorf("QRCode parsing not yet implemented")
		return
	}

	// Read vaccine pass
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// remove magic header : HC1: (Health Certificate Version 1)
	dccBase45 := string(fileBytes)
	if !strings.HasPrefix(dccBase45, "HC1:") {
		err = fmt.Errorf("data does not appear to be EU DCC / Vaccine Passport")
		return
	}
	dccBase45 = strings.Split(dccBase45, "HC1:")[1]

	// Decode base45
	dccCOSECompressed, err := base45.DecodeString(dccBase45)
	if err != nil {
		return
	}

	// Decompress Binary COSE Data with zlib
	dccCOSE, err := zlibDecompress(dccCOSECompressed)
	if err != nil {
		return
	}

	// Unmarshal COSE
	msg := cose.NewSign1Message()
	err = msg.UnmarshalCBOR(dccCOSE)
	if err != nil {
		return
	}

	// Read CBOR DCC Payload into JSON Struct
	dcc := &DCC{}
	err = cbor.Unmarshal(msg.Payload, dcc)
	if err != nil {
		return
	}

	// We got the goods! Return them
	payload = dcc
	headers = msg.Headers
	signature = msg.Signature
	messageRaw = msg

	return
}

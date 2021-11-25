package dcc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dasio/base45"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

func ParseGreenpass(path string, fileType int) (*DCC, error) {

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if fileType == TypeQRCode {
		return nil, fmt.Errorf("QRCode parsing not yet implemented")
	}

	dccBase45 := string(fileBytes)
	if !strings.HasPrefix(dccBase45, "HC1:") {
		return nil, fmt.Errorf("data does not appear to be EU DCC / Vaccine Passport")
	}

	// remove magic header : HC1: (Health Certificate Version 1)
	dccBase45 = strings.Split(dccBase45, "HC1:")[1]

	// Decode base45
	dccCOSECompressed, err := base45.DecodeString(dccBase45)
	if err != nil {
		return nil, err
	}

	// Decompress Binary COSE Data with zlib
	dccCOSE, err := zlibDecompress(dccCOSECompressed)
	if err != nil {
		return nil, err
	}

	// Unmarshal COSE
	msg := cose.NewSign1Message()
	err = msg.UnmarshalCBOR(dccCOSE)
	if err != nil {
		return nil, err
	}

	// Read CBOR DCC Payload into JSON Struct
	dcc := &DCC{}
	err = cbor.Unmarshal(msg.Payload, dcc)
	if err != nil {
		return nil, err
	}

	// Marshal to JSON to pretty-print
	jsonDCC, err := json.MarshalIndent(dcc, "", "	")
	if err != nil {
		return nil, err
	}

	// We got the goods!

	fmt.Printf("Headers: %v\n", msg.Headers)

	fmt.Printf("Payload: %s\n", string(jsonDCC))

	fmt.Printf("Signature: %s\n", hexutil.Encode(msg.Signature))

	return nil, nil
}

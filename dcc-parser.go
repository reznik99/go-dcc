package dcc

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/dasio/base45"
	"github.com/fxamacker/cbor/v2"
	"github.com/veraison/go-cose"
)

var dccPrefix = "HC1:"

// ParseQR parses a Vaccine Passport, it reads the image file at 'path', decoding the QR code and returns the Pass Payload and Raw COSE message containing headers, payload and signatures
func ParseQR(path string) (payload *DCC, messageRaw *cose.Sign1Message, err error) {
	err = errors.New("QRCode parsing not yet implemented")
	return
}

// Parse parses a Vaccine Passport, it reads the file at 'path' and returns the Pass Payload and Raw COSE message containing headers, payload and signatures
func Parse(path string) (payload *DCC, messageRaw *cose.Sign1Message, err error) {

	// Read vaccine pass
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// remove magic header : HC1: (Health Certificate Version 1)
	dccBase45 := string(fileBytes)
	if !strings.HasPrefix(dccBase45, dccPrefix) {
		err = errors.New("data does not appear to be EU DCC / Vaccine Passport")
		return
	}
	dccBase45 = strings.Split(dccBase45, dccPrefix)[1]

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
	messageRaw = msg

	return
}

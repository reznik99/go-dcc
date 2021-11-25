package main

import (
	"bytes"
	"crypto/rand"
	"fmt"

	_ "crypto/sha256"

	"compress/zlib"

	"github.com/dasio/base45"
	"github.com/fxamacker/cbor/v2"
	"github.com/skip2/go-qrcode"

	cose "github.com/veraison/go-cose"
)

func main() {

	// Generate JSON eu-dcc structure
	dccJson := generateDCC("FIRSTNAME", "LASTNAME", "1996-04-17", "NZ", "Ministry of Health, NZ", "2021-11-25", 2)

	// JSON to CBOR
	dccCBORBytes, err := cbor.Marshal(dccJson)
	if err != nil {
		panic(err)
	}

	// Sign CBOR with COSE
	dccCOSESignMsg, err := coseSign(dccCBORBytes)
	if err != nil {
		panic(err)
	}
	dccCOSE, err := dccCOSESignMsg.MarshalCBOR()
	if err != nil {
		panic(err)
	}

	// Compress Binary COSE Data with zlib
	dccCOSEcompressed := zlibCompress(dccCOSE)

	// Encode zlib compressed cose to base45
	dccBase45 := base45.EncodeToString(dccCOSEcompressed)

	// Prepend magic HC1 (Health Certificate Version 1)
	dccBase45 = fmt.Sprintf("HC1:%s\n", dccBase45)

	fmt.Println(dccBase45)

	// Convert to QR Code
	err = qrcode.WriteFile(dccBase45, qrcode.Medium, 256, "vaccinePass.png")
	if err != nil {
		panic(err)
	}
}

func coseSign(data []byte) (*cose.Sign1Message, error) {

	// create a signer with a new private key
	// ES256 (algId: -7), i.e.: ECDSA w/ SHA-256 from RFC8152
	signer, err := cose.NewSigner(cose.ES256, nil)
	if err != nil {
		return nil, err
	}

	msg := cose.NewSign1Message()
	msg.Headers.Protected["alg"] = "ES256"        // ECDSA w/ SHA-256
	msg.Headers.Protected["kid"] = "dy8HnMQYOrE=" // Temporary hardcoded Key ID
	msg.Payload = data

	err = msg.Sign(rand.Reader, nil, *signer)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func zlibCompress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

package dcc

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dasio/base45"
	"github.com/fxamacker/cbor/v2"
	"github.com/skip2/go-qrcode"
	"github.com/veraison/go-cose"
)

// TODO: Temporary hardcoded Key ID
var newZealandKID4 = "dy8HnMQYOrE="

func GenerateGreenpass(inputPath string, outputPath string, outputType int) {

	conf, err := readRaw(inputPath)
	if err != nil {
		panic(err)
	}

	// Generate JSON eu-dcc structure
	dccJson := generateDCCStruct(conf.Name, conf.Surname, conf.Dob, conf.IssuerCountry, conf.Issuer, conf.VaccinationDate, conf.Doses)

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
	dccBase45 = fmt.Sprintf("HC1:%s", dccBase45)

	fmt.Println(dccBase45)

	// Convert to QR Code or raw txt
	switch outputType {
	case TypeRAWGreepass:
		err = ioutil.WriteFile(outputPath, []byte(dccBase45), 0644)
		if err != nil {
			panic(err)
		}
	case TypeQRCode:
		err = qrcode.WriteFile(dccBase45, qrcode.Medium, 256, outputPath)
		if err != nil {
			panic(err)
		}
	}
}

func coseSign(data []byte) (*cose.Sign1Message, error) {

	// create a signer with a new private key
	// TODO: This should be initiated from existing keypair
	signer, err := cose.NewSigner(cose.ES256, nil)
	if err != nil {
		return nil, err
	}

	kid, err := base64.StdEncoding.DecodeString(newZealandKID4)
	if err != nil {
		return nil, err
	}

	msg := cose.NewSign1Message()
	msg.Headers.Protected["alg"] = "ES256" // ECDSA w/ SHA-256
	msg.Headers.Protected["kid"] = kid     // KID is first 8 bytes of signer certificate
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

func zlibDecompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	output, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// readRaw reads raw config to parse data for Payload of DCC/Greenpass
func readRaw(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

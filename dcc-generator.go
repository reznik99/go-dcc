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

// Generates Vaccine Passport with json data read from 'qrPath' and outputs the Pass as a QR code at 'outputPath'
func GenerateQR(qrPath string, outputPath string) error {

	raw, err := Generate(qrPath)
	if err != nil {
		return err
	}

	// Convert to QR Code or raw txt
	err = qrcode.WriteFile(raw, qrcode.Medium, 256, outputPath)
	if err != nil {
		return err
	}

	return nil
}

// Generates Vaccine Passport with json data read from 'dataPath' and returns raw Pass string as `HC1:...`
func Generate(dataPath string) (dccBase45 string, err error) {

	conf, err := readRaw(dataPath)
	if err != nil {
		return
	}

	// Generate JSON eu-dcc structure
	dccJson := generateDCCStruct(conf.Name, conf.Surname, conf.Dob, conf.IssuerCountry, conf.Issuer, conf.VaccinationDate, conf.Doses)

	// JSON to CBOR
	dccCBORBytes, err := cbor.Marshal(dccJson)
	if err != nil {
		return
	}

	// Sign CBOR with COSE
	dccCOSESignMsg, err := coseSign(dccCBORBytes)
	if err != nil {
		return
	}
	dccCOSE, err := dccCOSESignMsg.MarshalCBOR()
	if err != nil {
		return
	}

	// Compress Binary COSE Data with zlib
	dccCOSEcompressed := zlibCompress(dccCOSE)

	// Encode zlib compressed cose to base45
	dccBase45 = base45.EncodeToString(dccCOSEcompressed)

	// Prepend magic HC1 (Health Certificate Version 1)
	dccBase45 = fmt.Sprintf("%s%s", dccPrefix, dccBase45)

	return
}

func coseSign(data []byte) (*cose.Sign1Message, error) {

	// create a signer with a new private key
	// TODO: This should be initiated from existing keypair
	signer, err := cose.NewSigner(cose.AlgorithmES256, nil)
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

	err = msg.Sign(rand.Reader, nil, signer)
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

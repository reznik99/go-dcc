package dcc

import (
	"bytes"
	"compress/zlib"
	"crypto"
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

// Generates Vaccine Passport with json data read from 'dataPath' and writes the Vaccine Pass as a QR code file at 'outputPath'
func GenerateQR(key crypto.Signer, kid, dataPath, outputPath string) error {

	// Generate Vaccine Pass using given key and data
	raw, err := Generate(key, kid, dataPath)
	if err != nil {
		return err
	}

	// Write raw Vaccine Pass to QR code file
	err = qrcode.WriteFile(raw, qrcode.Large, 500, outputPath)
	if err != nil {
		return err
	}

	return nil
}

// Generates Vaccine Passport with json data read from 'dataPath' and returns raw Vaccine Pass string as `HC1:...`
func Generate(key crypto.Signer, kid, dataPath string) (dccBase45 string, err error) {

	conf, err := readRaw(dataPath)
	if err != nil {
		return
	}

	// Generate JSON eu-dcc structure
	payload := generateDCCStruct(conf.Name, conf.Surname, conf.Dob, conf.IssuerCountry, conf.Issuer, conf.VaccinationDate, conf.Doses)

	// JSON struct to CBOR
	dccCBORBytes, err := cbor.Marshal(payload)
	if err != nil {
		return
	}

	// Sign CBOR with COSE
	dccCOSESignMsg, err := coseSign(dccCBORBytes, key, kid)
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

// coseSign Signs the given CBOR payload with the given signer key
func coseSign(payload []byte, key crypto.Signer, kid string) (*cose.Sign1Message, error) {

	// create a signer with a new private key
	signer, err := cose.NewSigner(cose.AlgorithmES256, key)
	if err != nil {
		return nil, err
	}

	kidBytes, err := base64.StdEncoding.DecodeString(kid)
	if err != nil {
		return nil, err
	}

	msg := cose.NewSign1Message()
	msg.Headers.Protected[cose.HeaderLabelKeyID] = kidBytes // KID is first 8 bytes of Signer Certificate
	msg.Payload = payload

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

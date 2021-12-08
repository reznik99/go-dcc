package main

import (
	_ "crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/reznik99/go-dcc/v1/dcc"
)

func main() {

	generate := flag.Bool("gen", false, " Generate DCC/Greenpass and save as QR Code")
	verify := flag.Bool("verify", false, " Verify DCC/Greenpass by reading QR Code")
	info := flag.Bool("info", false, " Inspect DCC/Greenpass by reading QR Code")

	raw := flag.Bool("raw", false, " raw txt file with HC1: greenpass contents instead of QR Code")

	inputPath := flag.String("in", "", "Path to input file")
	outputPath := flag.String("out", "", "Path to output file")
	fileType := dcc.TypeQRCode

	flag.Parse()

	if *inputPath == "" || (!*generate && !*verify && !*info) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *raw {
		fileType = dcc.TypeRAWGreepass
	}

	if *generate {
		dcc.GenerateGreenpass(*inputPath, *outputPath, fileType)
	}
	if *verify {
		dcc.VerifyGreenpass(*inputPath, fileType)
	}
	if *info {
		dcc, raw, err := dcc.ParseGreenpass(*inputPath, fileType)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		prettyDCC, err := json.MarshalIndent(dcc, "", "	")
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		fmt.Printf("Headers:   %v\n", raw.Headers)
		fmt.Printf("Payload:   %s\n", prettyDCC)
		fmt.Printf("Signature: %s\n", hexutil.Encode(raw.Signature))

	}
}

package main

import (
	_ "crypto/sha256"
	"flag"
	"go-dcc/v1/dcc"
	"os"
)

var (
	typeQRCode      = 1
	typeRAWGreepass = 2
)

func main() {

	generate := *flag.Bool("gen", false, " Generate DCC/Greenpass and save as QR Code")
	verify := *flag.Bool("verify", false, " Verify DCC/Greenpass by reading QR Code")
	info := *flag.Bool("info", false, " Inspect DCC/Greenpass by reading QR Code")

	raw := *flag.Bool("raw", false, " raw txt file with HC1: greenpass contents instead of QR Code")

	filePath := *flag.String("in", "", "Path to input file")
	fileType := typeQRCode

	flag.Parse()

	if !generate && !verify && !info {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if raw {
		fileType = typeRAWGreepass
		fileType += 1
	}

	if generate {
		dcc.GenerateGreenpass(filePath)
	}
	if verify {
		// dcc.VerifyGreenpass()
	}
	if info {
		// dcc.ParseGreenpass()
	}

}

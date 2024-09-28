package main

import (
	"fmt"
	"os"

	"github.com/digitorus/pdfsign/verify"
)

func Valid(filepath string) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	pdfInfo, err := verify.File(f)
	if err != nil {
		panic(err)
	}
	for _, signer := range pdfInfo.Signers {
		fmt.Printf("%+v\n", signer)
	}
}

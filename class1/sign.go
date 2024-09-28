package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/digitorus/pdf"
	"github.com/digitorus/pdfsign/revocation"
	"github.com/digitorus/pdfsign/sign"
)

func Sign(inputPDFPath, outputPDFPath string) {
	//First step
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		panic(err)
	}
	// Second step
	x509RootCertificate := &x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Organization:  []string{"Demo"},
			Country:       []string{"CN"},
			Province:      []string{"guangdong"},
			Locality:      []string{"shenzheng"},
			StreetAddress: []string{"fandoudajie 18号"},
			PostalCode:    []string{"010027"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	rootCertificateBytes, err := x509.CreateCertificate(rand.Reader, x509RootCertificate, x509RootCertificate, &privateKey.PublicKey, privateKey)

	if err != nil {
		panic(err)
	}

	rootCertificate, err := x509.ParseCertificate(rootCertificateBytes)

	if err != nil {
		panic(err)
	}

	roots := x509.NewCertPool()

	roots.AddCert(rootCertificate)

	certificateChain, err := rootCertificate.Verify(x509.VerifyOptions{
		Roots: roots,
	})

	if err != nil {
		panic(err)
	}
	// Third step

	outputFile, err := os.Create(outputPDFPath)

	if err != nil {
		panic(err)
	}

	defer func(outputFile *os.File) {
		err = outputFile.Close()

		if err != nil {
			log.Println(err)
		}

		fmt.Println("output file closed")
	}(outputFile)

	inputFile, err := os.OpenFile(inputPDFPath, os.O_RDONLY, 0755)

	if err != nil {
		panic(err)
	}

	defer func(inputFile *os.File) {
		err = inputFile.Close()

		if err != nil {
			log.Println(err)
		}

		fmt.Println("input file closed")
	}(inputFile)

	fileInfo, err := inputFile.Stat()

	if err != nil {
		panic(err)
	}

	size := fileInfo.Size()

	pdfReader, err := pdf.NewReader(inputFile, size)

	if err != nil {
		panic(err)
	}

	// Fourth step

	err = sign.Sign(inputFile, outputFile, pdfReader, size, sign.SignData{
		Signature: sign.SignDataSignature{
			Info: sign.SignDataSignatureInfo{
				Name:        "Su",
				Location:    "CN",
				Reason:      "PDF SIGN",
				ContactInfo: "TEST",
				Date:        time.Now().Local(),
			},
			CertType:   sign.CertificationSignature,
			DocMDPPerm: sign.AllowFillingExistingFormFieldsAndSignaturesPerms,
		},
		Signer:            privateKey,       // crypto.Signer
		DigestAlgorithm:   crypto.SHA256,    // hash algorithm for the digest creation
		Certificate:       rootCertificate,  // x509.Certificate
		CertificateChains: certificateChain, // x509.Certificate.Verify()
		TSA: sign.TSA{
			URL:      "",
			Username: "",
			Password: "",
		},

		// The follow options are likely to change in a future release
		//
		// cache revocation data when bulk signing
		RevocationData: revocation.InfoArchival{},
		// custom revocation lookup
		RevocationFunction: sign.DefaultEmbedRevocationStatusFunction,
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("pdf signed")
	}
}

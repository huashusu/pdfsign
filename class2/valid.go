package main

import "github.com/pdfcpu/pdfcpu/pkg/api"

func Valid(path string) error {
	return api.ValidateFile(path, nil)
}

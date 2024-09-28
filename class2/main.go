package main

import "log"

func main() {
	err := Valid("./class1/output.pdf")
	if err != nil {
		panic(err)
	} else {
		log.Printf("sign success\n")
	}
}

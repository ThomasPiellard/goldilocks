package main

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/field/generator"
	"github.com/consensys/gnark-crypto/field/generator/config"
)

func assertNoError(err error) {
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(-1)
	}
}

//go:generate go run main.go
func main() {

	modulus := "18446744069414584321"
	gold, err := config.NewFieldConfig("fr", "Element", modulus, true)
	assertNoError(err)

	err = generator.GenerateFF(gold, "../fr/")
	assertNoError(err)
}

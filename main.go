package main

import (
	"fmt"
	"os"

	"github.com/gradienthealth/dicom-protos/parse"
)

func main() {
	af, err := os.Open("dicom-standard/standard/attributes.json")
	if err != nil {
		panic(err)
	}
	mf, err := os.Open("dicom-standard/standard/modules.json")
	if err != nil {
		panic(err)
	}

	mta, err := os.Open("dicom-standard/standard/module_to_attributes.json")

	modules, err := parse.Parse(af, mf, mta)
	for _, item := range modules {
		fmt.Println(item)
		fmt.Println()
	}
}

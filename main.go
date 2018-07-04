package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gradienthealth/dicom-protos/gen"
	"github.com/gradienthealth/dicom-protos/parse"
)

func moduleNameToFilename(name string) string {
	s := strings.Replace(name, " ", "", -1)
	s2 := strings.Replace(s, "/", "", -1)
	return fmt.Sprintf("protos/%s.proto", s2)
}

func main() {
	af, err := os.Open("dicom-standard/standard/attributes.json")
	if err != nil {
		panic(err)
	}
	defer af.Close()

	mf, err := os.Open("dicom-standard/standard/modules.json")
	if err != nil {
		panic(err)
	}
	defer mf.Close()

	mta, err := os.Open("dicom-standard/standard/module_to_attributes.json")
	if err != nil {
		panic(err)
	}
	defer mta.Close()

	modules, err := parse.Parse(af, mf, mta)
	for _, m := range modules {
		out, err := os.Create(moduleNameToFilename(m.Name))
		if err != nil {
			log.Printf("Couldn't create file %s \n", moduleNameToFilename(m.Name))
			log.Print(err)
		}
		gen.ModuleProto(m, out) // Generate module proto
		out.Close()
	}

}

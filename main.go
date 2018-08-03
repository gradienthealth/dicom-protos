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
	s = strings.Replace(s, "/", "", -1)
	return fmt.Sprintf("protos/%s.proto", s)
}

func generateFile(name string) *os.File {
	out, err := os.Create(name)
	if err != nil {
		log.Printf("Couldn't create file %s \n", name)
		log.Panic(err)
	}
	return out
}

func main() {
	af, err := os.Open("dicom-standard/standard/attributes.json")
	if err != nil {
		log.Panic(err)
	}
	defer af.Close()

	mf, err := os.Open("dicom-standard/standard/modules.json")
	if err != nil {
		log.Panic(err)
	}
	defer mf.Close()

	mta, err := os.Open("dicom-standard/standard/module_to_attributes.json")
	if err != nil {
		log.Panic(err)
	}
	defer mta.Close()

	// Generate Sequences
	seqOut := generateFile("protos/Sequences.proto")
	fmt.Fprintln(seqOut, gen.ProtoHeader)
	fmt.Fprintln(seqOut, "import \"Attributes.proto\";")
	fmt.Fprintln(seqOut, "")
	defer seqOut.Close()

	attrOut := generateFile("protos/Attributes.proto")
	fmt.Fprintln(attrOut, gen.ProtoHeader)
	fmt.Fprintln(attrOut, "")
	defer attrOut.Close()

	modules, err := parse.Parse(af, mf, mta)
	if err != nil {
		log.Panic()
	}

	numModules, err := gen.ProtosToFile("protos", modules)
	if err != nil {
		panic(err)
	}

	log.Printf("Complete. Generated protos for %d modules.", numModules)
}

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

func generateFile(name string) *os.File {
	out, err := os.Create(name)
	if err != nil {
		log.Printf("Couldn't create file %s \n", name)
		log.Print(err)
		panic(err)
	}
	return out
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

	counter := 0
	modules, err := parse.Parse(af, mf, mta)
	for _, m := range modules {
		out := generateFile(moduleNameToFilename(m.Name))
		gen.ModuleProto(m, out) // Generate module proto
		out.Close()
		counter++
	}

	// Generate Sequences
	seqOut := generateFile("protos/Sequences.proto")
	fmt.Fprintln(seqOut, gen.ProtoHeader)
	fmt.Fprintln(seqOut, "import \"Attributes.proto\";")
	fmt.Fprintln(seqOut, "")

	attrOut := generateFile("protos/Attributes.proto")
	fmt.Fprintln(attrOut, gen.ProtoHeader)
	fmt.Fprintln(attrOut, "")
	fmt.Fprintln(attrOut, "")

	for _, a := range gen.SQMap {
		gen.SequenceAttrProto(a, seqOut, attrOut, map[string]*parse.Attribute{})
	}

	attrOut.Close()
	seqOut.Close()

	log.Printf("Complete. Generated protos for %d modules.", counter)
}

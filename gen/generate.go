package gen

import (
	"io"

	"bytes"
	"fmt"

	"log"

	"github.com/gradienthealth/dicom-protos/parse"
)

const NOTMAPPED = "NOTMAPPED"

// VRToProto maps a DICOM VR to the corresponding proto type (Part 5 6.2)
var VRToProto = map[string]string{
	"UT":       "string",
	"US":       "uint32", // maps 16 bit to 32 bit uint
	"UN":       "bytes",
	"UL":       "uint32",
	"UI":       "string", //TODO(suyash) create a UID proto message
	"TM":       "int64",
	"ST":       "string",
	"SS":       "int32",   // 16 bit -> 32 bit
	"SQ":       NOTMAPPED, //TODO(suyash) figure out how to map sequences best. just bytes w/ word_size?
	"SL":       "int32",
	"SH":       "string",
	"PN":       "string",
	"OW":       "bytes",  // 16 bit words needs to be specified somehow
	"OF":       "string", //TODO(suyash) think about just parsing as float32
	"OD":       "string",
	"OB":       "string",
	"LT":       "string",
	"LO":       "string",
	"IS":       "int32",
	"FD":       "double",
	"FL":       "float",
	"DT":       "int64", //TODO(suyash) consider using protobuf time lib
	"DS":       "float",
	"DA":       "string",
	"CS":       "string",
	"AT":       "string", // attribute tag, consider handing specially
	"AS":       "string",
	"AE":       "string",
	"UC":       "string",
	"UR":       "string", // consider url types later
	"":         "bytes",  // assuming this is unknown VR
	"US or SS": "int64",
	"US or OW": "bytes",
	"OB or OW": "bytes",
}

const ProtoHeader = "syntax = \"proto3\";"

func AttributeProto(a *parse.Attribute) string {
	b := bytes.NewBufferString("")

	messageName := a.Keyword
	if a.Keyword == "" {
		if a.Name == "" && a.ValueRepresentation == "" && a.ValueMultiplicity == "" {
			return ""
		}
		messageName = parse.TagToKey(a.Tag)
		log.Printf("Blank keyword, assigning tag. Attribute: %+v", a)

	}

	fmt.Fprintf(b, "message %s {\n", messageName)

	if a.ValueMultiplicity != "1" {
		fmt.Fprint(b, "\trepeated ")
	} else {
		fmt.Fprint(b, "\t")
	}

	var fieldName string
	if val, ok := VRToProto[a.ValueRepresentation]; ok {
		fieldName = val
	} else {
		log.Printf("Issue mapping VR. Offending VR: %s", a.ValueRepresentation)
	}
	fmt.Fprintf(b, "%s value = 1;\n", fieldName)

	// Add word size for OW VRs. Clients can choose to fill this in
	if a.ValueRepresentation == "OW" {
		fmt.Fprint(b, "int32 word_size = 2;\n")
	}
	fmt.Fprint(b, "}")
	return b.String()
}

// ModuleProto generates the corresponding protocol buffer for the module and writes it to w
func ModuleProto(module *parse.Module, w io.Writer) error {
	fmt.Fprintln(w, ProtoHeader)
	fmt.Fprintln(w, "")

	for _, a := range module.Attributes {
		if a.Retired {
			continue
		}
		ap := AttributeProto(&a)
		fmt.Fprintf(w, ap)
		fmt.Fprintln(w, "")
	}

	return nil
}

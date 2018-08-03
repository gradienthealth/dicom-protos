package gen

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gradienthealth/dicom-protos/parse"
)

// VRToProto maps a DICOM VR to the corresponding proto type (Part 5 6.2)
var VRToProto = map[string]string{
	"UT":       "string",
	"US":       "uint32", // maps 16 bit to 32 bit uint
	"UN":       "bytes",
	"UL":       "uint32",
	"UI":       "string", //TODO(suyash) create a UID proto message?
	"TM":       "int64",
	"ST":       "string",
	"SS":       "int32", // 16 bit -> 32 bit
	"SQ":       "bytes", //TODO(suyash) figure out how to map sequences best. Just bytes fallback now
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
	"":         "bytes",  // assuming this is unknown VR, default to bytes
	"OL":       "int32",
	"US or SS": "int64",
	"US or OW": "bytes",
	"OB or OW": "bytes",
}

// SetComplete holds on to attributes that have been generated already
var SetComplete = map[string]*parse.Attribute{}

// ProtoHeader represents the protocol buffer header for all protos
const ProtoHeader = "syntax = \"proto3\";"

// AttributeProto generates the protocol message for an attribute and returns it as a string
func AttributeProto(a *parse.Attribute) string {

	b := bytes.NewBufferString("")

	// Determine name of proto message:
	messageName := a.Keyword
	if a.Keyword == "" {
		messageName = string(a.TagKey())
		log.Printf("Blank keyword, assigning tag. Attribute: %+v", a)
	}

	fmt.Fprintf(b, "// DICOM Tag: %s\n", a.Tag)
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

	// Add word size for OW VRs. Clients can choose to fill this in, should be 16bits
	if a.ValueRepresentation == "OW" {
		fmt.Fprint(b, "\tint32 word_size = 2;\n")
	}

	if a.ValueRepresentation == "SQ" {
		fmt.Fprint(b, "\t // This attribute has a SQ VR, no proto mapping yet\n")
	}

	fmt.Fprint(b, "}")
	return b.String()
}

// ModuleProto generates the corresponding protocol buffer for the module and writes it to w
func ModuleProto(module *parse.Module, w io.Writer) {
	fmt.Fprintf(w, "// %s\n// Link to standard: %s\n", module.Name, module.LinkToStandard)
	fmt.Fprintln(w, ProtoHeader)
	fmt.Fprintln(w, "option go_package = \"github.com/gradienthealth/dicom-protos/protos\";")
	fmt.Fprintln(w, "import \"Sequences.proto\";")
	fmt.Fprintln(w, "import \"Attributes.proto\";")
	fmt.Fprintln(w, "")

	// Write module proto.
	moduleName := normalizeString(module.Name)
	fmt.Fprintf(w, "message %sModule {\n", moduleName)

	i := 1
	attrs := getSortedAttributes(module.Attributes)
	for _, a := range attrs {
		if !(a.Retired || a.IsEmpty()) {
			fmt.Fprintf(w, "\t%s %s = %d;\n", a.Keyword, getFieldName(a.Name), i)
			i++
		}
	}
	fmt.Fprintln(w, "}")
}

// SequenceAttrProto generates two protocol buffer files for Sequences and Leaf Attributes.
func SequenceAttrProto(a *parse.Attribute, wSeq, wAttr io.Writer) error {
	if val, ok := SetComplete[a.Tag]; ok {
		if a.ValueRepresentation != val.ValueRepresentation {
			log.Print("WARN: ValueRepresentation for same tag does not match\n")
		}
	} else if len(a.SubAttributes) == 0 && !(a.Retired || a.IsEmpty()) {
		// Write the attribute out if it is a leaf node.
		SetComplete[a.Tag] = a
		ap := AttributeProto(a)
		fmt.Fprintf(wAttr, ap)
		fmt.Fprintln(wAttr, "")
	} else {
		// This attribute is not a leaf Attribute
		SetComplete[a.Tag] = a

		fmt.Fprintf(wSeq, "// DICOM Tag: %s\n", a.Tag)
		fmt.Fprintf(wSeq, "message %s {\n ", a.Keyword)

		// Write out the sequence protobuf message.
		i := 1
		attrs := getSortedAttributes(a.SubAttributes)
		for _, attr := range attrs {
			fmt.Fprintf(wSeq, "\t%s %s = %d;\n", attr.Keyword, getFieldName(attr.Name), i)
			i++
		}
		fmt.Fprint(wSeq, "} \n\n")

		// Recursively call for children.
		for _, s := range attrs {
			SequenceAttrProto(s, wSeq, wAttr)
		}

	}
	return nil
}

func getSortedKeys(m map[parse.TagKey]*parse.Attribute) []string {
	keys := make([]string, len(m))
	for k := range m {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	return keys
}

func getSortedAttributes(m map[parse.TagKey]*parse.Attribute) []*parse.Attribute {
	keys := getSortedKeys(m)
	sort.Strings(keys)
	attrs := make([]*parse.Attribute, 0, len(m))
	for _, key := range keys {
		if key == "" {
			continue //TODO: investigate why we have empty string keys
		}
		attrs = append(attrs, m[parse.TagKey(key)])

	}
	return attrs
}

func moduleProtoToFile(dirPath string, m *parse.Module) error {
	out, err := os.Create(dirPath + "/" + moduleNameToFilename(m.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	ModuleProto(m, out) // Generate module proto
	return nil
}

// ProtosToFile generates ALL protos from the provided modules to the directory path specified
func ProtosToFile(dirPath string, modules []*parse.Module) (counter int, err error) {
	// Create Sequences and Attributes file
	seqOut, err := os.Create(dirPath + "/Sequences.proto")
	if err != nil {
		return 0, err
	}
	fmt.Fprintln(seqOut, ProtoHeader)
	fmt.Fprintln(seqOut, "import \"Attributes.proto\";")
	fmt.Fprintln(seqOut, "")
	defer seqOut.Close()

	attrOut, err := os.Create(dirPath + "/Attributes.proto")
	fmt.Fprintln(attrOut, ProtoHeader)
	fmt.Fprintln(attrOut, "")
	defer attrOut.Close()

	// Begin to generate Module protos and populate Sequences.proto and Attributes.proto
	sort.Sort(parse.ModulesByID(modules))
	counter = 0
	for _, m := range modules {
		// Generate Module proto
		err := moduleProtoToFile(dirPath, m)
		if err != nil {
			return counter, err
		}

		// Generate messages for this module's attributes
		attrs := getSortedAttributes(m.Attributes)
		for _, attr := range attrs {
			SequenceAttrProto(attr, seqOut, attrOut)
		}

		counter++
	}

	return counter, nil
}

func normalizeString(s string) string {
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.Replace(s, "'", "", -1)
	s = strings.Replace(s, "&", "_and_", -1)
	s = strings.Replace(s, "0", "zero", -1)
	s = strings.Replace(s, "1", "one", -1)
	s = strings.Replace(s, "2", "two", -1)
	s = strings.Replace(s, "3", "three", -1)
	s = strings.Replace(s, "4", "four", -1)
	s = strings.Replace(s, "5", "five", -1)
	s = strings.Replace(s, "6", "six", -1)
	s = strings.Replace(s, "7", "seven", -1)
	s = strings.Replace(s, "8", "eight", -1)
	s = strings.Replace(s, "9", "nine", -1)
	// Only allow ascii characters into the FieldName
	re := regexp.MustCompile("[[:^ascii:]]")
	s = re.ReplaceAllLiteralString(s, "")
	return s
}

func getFieldName(name string) string {
	s := strings.ToLower(name)
	s = normalizeString(s)
	return s
}

func moduleNameToFilename(name string) string {
	s := strings.Replace(name, " ", "", -1)
	s = strings.Replace(s, "/", "", -1)
	return fmt.Sprintf("%s.proto", s)
}

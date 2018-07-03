package parse

import (
	"encoding/json"
	"io"
	"strings"
)

type Attribute struct {
	Tag                 string `json:"tag"`
	Name                string `json:"name"`
	Keyword             string `json:"keyword"`
	ValueRepresentation string `json:"valueRepresentation"`
	valueMultiplicity   string `json:"valueMultiplicity"`
	Retired             bool   `json:"retired"`
}

type AttributeMap map[string]Attribute

type Module struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	LinkToStandard string `json:"linkToStandard"`
	Attributes     []Attribute
}

func (m *Module) AddAttribute(a Attribute) {
	m.Attributes = append(m.Attributes, a)
}

type ModuleMap map[string]Module
type ModulePointerMap map[string]*Module

type ModuleToAttributeItem struct {
	ModuleID string `json:"module"`
	Tag      string `json:"tag"`
	Type     string `json:"type"`
}

func tagToKey(tag string) string {
	s := strings.Replace(tag, ",", "", -1)
	s2 := strings.Replace(s, "(", "", -1)
	s3 := strings.Replace(s2, ")", "", -1)

	return s3
}

func convertToPointerMap(m ModuleMap) ModulePointerMap {
	var mp ModulePointerMap
	mp = make(ModulePointerMap)
	for k, v := range m {
		var module Module
		module = v
		mp[k] = &module
	}
	return mp
}

func Parse(attributeFile, moduleFile, moduleToAttributesFile io.Reader) ([]*Module, error) {

	aFileDec := json.NewDecoder(attributeFile)
	var aMap AttributeMap
	err := aFileDec.Decode(&aMap)
	if err != nil {
		return []*Module{}, err
	}

	mFileDec := json.NewDecoder(moduleFile)
	var mMap ModuleMap
	err = mFileDec.Decode(&mMap)
	if err != nil {
		return []*Module{}, err
	}
	mpMap := convertToPointerMap(mMap)

	moduleToAttributesDec := json.NewDecoder(moduleToAttributesFile)
	var moduleToAttributes []ModuleToAttributeItem
	err = moduleToAttributesDec.Decode(&moduleToAttributes)

	for _, mapping := range moduleToAttributes {
		tagKey := tagToKey(mapping.Tag)
		objA := mpMap[mapping.ModuleID]
		objA.AddAttribute(aMap[tagKey])
	}

	moduleArray := make([]*Module, len(mpMap))

	i := 0

	for _, v := range mpMap {
		moduleArray[i] = v
		i++
	}

	return moduleArray, nil

}

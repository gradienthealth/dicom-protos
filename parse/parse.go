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
	ValueMultiplicity   string `json:"valueMultiplicity"`
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

func TagToKey(tag string) string {
	s := strings.Replace(tag, ",", "", -1)
	s2 := strings.Replace(s, "(", "", -1)
	s3 := strings.Replace(s2, ")", "", -1)

	return s3
}

func (m ModuleMap) ToPointerMap() ModulePointerMap {
	mp := make(ModulePointerMap)
	for k, v := range m {
		var module Module
		module = v
		mp[k] = &module
	}
	return mp
}

func Parse(attributeFile, moduleFile, moduleToAttributesFile io.Reader) ([]*Module, error) {

	dec := json.NewDecoder(attributeFile)
	var aMap AttributeMap
	err := dec.Decode(&aMap)
	if err != nil {
		return []*Module{}, err
	}

	dec = json.NewDecoder(moduleFile)
	var mMap ModuleMap
	err = dec.Decode(&mMap)
	if err != nil {
		return []*Module{}, err
	}
	mpMap := mMap.ToPointerMap()

	dec = json.NewDecoder(moduleToAttributesFile)
	var moduleToAttributes []ModuleToAttributeItem
	err = dec.Decode(&moduleToAttributes)
	if err != nil {
		panic(err)
	}

	for _, mapping := range moduleToAttributes {
		tagKey := TagToKey(mapping.Tag)
		module := mpMap[mapping.ModuleID]
		module.AddAttribute(aMap[tagKey])
	}

	moduleArray := make([]*Module, len(mpMap))

	i := 0

	for _, v := range mpMap {
		moduleArray[i] = v
		i++
	}

	return moduleArray, nil

}

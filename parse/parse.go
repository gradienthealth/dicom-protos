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

// ToPointerMap converts a ModuleMap to a ModulePointerMap
func (m ModuleMap) ToPointerMap() ModulePointerMap {
	mp := make(ModulePointerMap)
	for k, v := range m {
		module := v // allocate a new module (copy v)
		mp[k] = &module
	}
	return mp
}

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

func Parse(attributeFile, moduleFile, moduleToAttributesFile io.Reader) ([]*Module, error) {
	// Parse AttributeMap
	dec := json.NewDecoder(attributeFile)
	var attrMap AttributeMap
	err := dec.Decode(&attrMap)
	if err != nil {
		return []*Module{}, err
	}

	// Parse ModuleMap
	dec = json.NewDecoder(moduleFile)
	var mMap ModuleMap
	err = dec.Decode(&mMap)
	if err != nil {
		return []*Module{}, err
	}
	moduleMap := mMap.ToPointerMap()

	// Parse ModuleToAttributeItem mappings
	dec = json.NewDecoder(moduleToAttributesFile)
	var moduleToAttributes []ModuleToAttributeItem
	err = dec.Decode(&moduleToAttributes)
	if err != nil {
		return []*Module{}, err
	}

	// Associate all modules with their sets of attributes
	for _, mapping := range moduleToAttributes {
		tagKey := TagToKey(mapping.Tag)
		module := moduleMap[mapping.ModuleID]
		module.AddAttribute(attrMap[tagKey])
	}

	// Create slice of modules to return
	i := 0
	modules := make([]*Module, len(moduleMap))
	for _, v := range moduleMap {
		modules[i] = v
		i++
	}

	return modules, nil
}

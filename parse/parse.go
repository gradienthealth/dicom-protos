package parse

import (
	"encoding/json"
	"io"
	"strings"
)

// TagKey represents a DICOM Tag as a string "key" used in dicom-standard. Tag (aaaa, bbbb) => aaaabbbb as a key
type TagKey string

// Attribute represents a DICOM attribute, with ability to point to "sub-attributes" that may be in a SQ
type Attribute struct {
	Tag                 string `json:"tag"`
	Name                string `json:"name"`
	Keyword             string `json:"keyword"`
	ValueRepresentation string `json:"valueRepresentation"`
	ValueMultiplicity   string `json:"valueMultiplicity"`
	Retired             bool   `json:"retired"`
	tagKey              TagKey
	SubAttributes       map[TagKey]*Attribute
}

// IsEmpty returns true if the Attribute is effectively empty
func (a *Attribute) IsEmpty() bool {
	return a.Name == "" && a.Keyword == "" &&
		a.ValueRepresentation == "" &&
		a.ValueMultiplicity == "" && a.Tag == ""
}

// Returns the corresponding TagKey of an Attribute
func (a *Attribute) TagKey() TagKey {
	if a.tagKey != "" {
		return a.tagKey
	}
	a.tagKey = tagToKey(a.Tag)
	return a.tagKey
}

// AddSubAttribute adds a subattribute in the right place to an attribute "tree" TODO: incoporate into parse flow
func (a *Attribute) AddSubAttribute(at *Attribute, keys []TagKey) {
	if len(keys) == 1 {
		a.SubAttributes[keys[0]] = at
	}

	// If map hasn't been made yet, do so
	if a.SubAttributes == nil {
		a.SubAttributes = make(map[TagKey]*Attribute)
	}

	a.SubAttributes[keys[1]].AddSubAttribute(at, keys[1:])
}

// AttributeMap represents a mapping from TagKey to Attribute entities
type AttributeMap map[TagKey]Attribute

// Module represents a DICOM module
type Module struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	LinkToStandard string `json:"linkToStandard"`
	Attributes     map[TagKey]*Attribute
}

// AddAttribute adds a top level attribute to a module
func (m *Module) AddAttribute(a Attribute, path string) {
	if m.Attributes == nil {
		m.Attributes = make(map[TagKey]*Attribute)
	}

	m.Attributes[a.TagKey()] = &a
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
	Path     string `json:"path"`
}

// Parse consumes data from innolitics' dicom-standard format and returns a slice of parsed Modules.
func Parse(attributeSource, moduleSource, moduleToAttributesSource io.Reader) ([]*Module, error) {
	// Parse AttributeMap
	dec := json.NewDecoder(attributeSource)
	var attrMap AttributeMap
	err := dec.Decode(&attrMap)
	if err != nil {
		return []*Module{}, err
	}

	// Parse ModuleMap
	dec = json.NewDecoder(moduleSource)
	var mMap ModuleMap
	err = dec.Decode(&mMap)
	if err != nil {
		return []*Module{}, err
	}
	moduleMap := mMap.ToPointerMap()

	// Parse ModuleToAttributeItem mappings
	dec = json.NewDecoder(moduleToAttributesSource)
	var moduleToAttributes []ModuleToAttributeItem
	err = dec.Decode(&moduleToAttributes)
	if err != nil {
		return []*Module{}, err
	}

	// Associate all modules with their sets of attributes
	for _, mapping := range moduleToAttributes {
		tagKey := tagToKey(mapping.Tag)
		module := moduleMap[mapping.ModuleID]
		module.AddAttribute(attrMap[tagKey], mapping.Path)
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

// Useful functions
func tagToKey(tag string) TagKey {
	s := strings.Replace(tag, ",", "", -1)
	s2 := strings.Replace(s, "(", "", -1)
	s3 := strings.Replace(s2, ")", "", -1)

	return TagKey(s3)
}

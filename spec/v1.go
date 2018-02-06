// Copyright © 2018 Horst Gutmann <zerok@zerokspot.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package spec

import (
	"encoding/json"
	"fmt"
)

const stringType = "string"
const integerType = "integer"

type validatorV1 struct {
	err  error
	spec *envelopV1
}

func (v *validatorV1) checkValidPropTypes() {
	if v.err != nil {
		return
	}
	for path, sec := range v.spec.RawSpec.RawSecrets {
	propLoop:
		for key, prop := range sec.RawProperties {
			switch prop.RawTypeName {
			case stringType:
				continue propLoop
			case integerType:
				continue propLoop
			default:
				v.err = ErrSpec1PropInvalidType{Path: path, Key: key, Type: prop.RawTypeName}
				return
			}
		}
	}
}

func isValidString(data interface{}) bool {
	switch data.(type) {
	case string:
		return true
	}
	return false
}

func isValidInteger(data interface{}) bool {
	switch data.(type) {
	case json.Number, int8, int16, int32, int64, uint8, uint16, uint32, uint64, int, uint:
		return true
	}
	return false
}

func (v *validatorV1) checkDefaultsOfRightType() {
	if v.err != nil {
		return
	}
	for path, sec := range v.spec.RawSpec.RawSecrets {
	propLoop:
		for key, prop := range sec.RawProperties {
			if prop.RawDefault == nil {
				continue propLoop
			}
			switch prop.RawTypeName {
			case stringType:
				if isValidString(prop.RawDefault) {
					continue propLoop
				}
				break
			case integerType:
				if isValidInteger(prop.RawDefault) {
					continue propLoop
				}
				break
			}
			v.err = ErrSpec1PropInvalidDefault{Path: path, Key: key, Default: prop.RawDefault}
		}
	}
}

func (v *validatorV1) checkKnownInputs() {
	if v.err != nil {
		return
	}
	for path, sec := range v.spec.RawSpec.RawSecrets {
		for key, prop := range sec.RawProperties {
			if prop.RawInput != "" && prop.RawInput != "default" && prop.RawInput != "hidden" {
				v.err = ErrSpec1PropInvalidInput{Path: path, Key: key}
				return
			}
		}
	}
}

type envelopV1 struct {
	RawVersion string `yaml:"version"`
	RawSpec    specV1 `yaml:"spec"`
}

func (e *envelopV1) Validate() error {
	v := validatorV1{spec: e}
	v.checkValidPropTypes()
	v.checkDefaultsOfRightType()
	v.checkKnownInputs()
	return v.err
}

func (e *envelopV1) Version() string {
	return e.RawVersion
}

func (e *envelopV1) SecretPaths() []string {
	result := make([]string, 0, len(e.RawSpec.RawSecrets))
	for path := range e.RawSpec.RawSecrets {
		result = append(result, path)
	}
	return result
}

func (e *envelopV1) Secret(path string) Secret {
	if sec, found := e.RawSpec.RawSecrets[path]; found {
		return &sec
	}
	return nil
}

type specV1 struct {
	RawSecrets map[string]secretV1 `yaml:"secrets"`
}

type secretV1 struct {
	RawPath       string                `yaml:"-"`
	RawLabel      string                `yaml:"label"`
	RawProperties map[string]propertyV1 `yaml:"properties"`
}

func (s *secretV1) Label() string {
	if s.RawLabel == "" {
		return s.RawPath
	}
	return s.RawLabel
}

func (s *secretV1) PropertyNames() []string {
	result := make([]string, 0, len(s.RawProperties))
	for key := range s.RawProperties {
		result = append(result, key)
	}
	return result
}

func (s *secretV1) Property(key string) SecretProperty {
	if prop, found := s.RawProperties[key]; found {
		return &prop
	}
	return nil
}

type propertyV1 struct {
	RawLabel    string      `yaml:"label"`
	RawTypeName string      `yaml:"type"`
	RawDefault  interface{} `yaml:"default"`
	RawHelp     string      `yaml:"help"`
	RawInput    string      `yaml:"input"`
	RawName     string      `yaml:"-"`
}

func (p *propertyV1) String() string {
	return fmt.Sprintf("<label=`%s` type=`%s` default=`%s` help=`%s` input=`%s`>", p.RawLabel, p.RawTypeName, p.RawDefault, p.RawHelp, p.RawInput)
}

func (p *propertyV1) IsValidData(data interface{}) error {
	switch p.RawTypeName {
	case stringType:
		if !isValidString(data) {
			return fmt.Errorf("Invalid data %v for string type", data)
		}
		break
	case integerType:
		if !isValidInteger(data) {
			return fmt.Errorf("Invalid data %v for integer type", data)
		}
		break
	}
	return nil
}

func (p *propertyV1) Type() string {
	return p.RawTypeName
}

func (p *propertyV1) Default() interface{} {
	return p.RawDefault
}

func (p *propertyV1) Help() string {
	return p.RawHelp
}

func (p *propertyV1) Input() string {
	return p.RawInput
}

func (p *propertyV1) Label() string {
	if p.RawLabel == "" {
		return p.RawName
	}
	return p.RawLabel
}

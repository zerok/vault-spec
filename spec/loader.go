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
	"fmt"
	"io/ioutil"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// FromPath parses file under the provided path into a specification object.
func FromPath(path string) (Specification, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var versioned versionedEnvelop
	if err := yaml.Unmarshal(raw, &versioned); err != nil {
		return nil, err
	}
	switch versioned.Version {
	case "1":
		result := envelopV1{}
		if err := yaml.Unmarshal(raw, &result); err != nil {
			return nil, err
		}
		if err := normalizeSpec(&result); err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, fmt.Errorf("Unsupported version '%v'", versioned.Version)
	}
}

func normalizeSpec(env *envelopV1) error {
	for secretPath, secretSpec := range env.RawSpec.RawSecrets {
		secretSpec.RawPath = secretPath
		for propName, propSpec := range secretSpec.RawProperties {
			propSpec.RawName = propName
			if propSpec.RawDefault != nil {
				converted, err := convertDefaultValue(propSpec.RawDefault, propSpec.RawTypeName)
				if err != nil {
					return err
				}
				propSpec.RawDefault = converted
			}
			if propSpec.RawInput == "" {
				propSpec.RawInput = "default"
			}
			secretSpec.RawProperties[propName] = propSpec
		}
		env.RawSpec.RawSecrets[secretPath] = secretSpec
	}
	return nil
}

func convertDefaultValue(v interface{}, typeName string) (interface{}, error) {
	switch typeName {
	case "string":
		return v, nil
	case "integer":
		return strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
	default:
		return nil, fmt.Errorf("%s not supported", typeName)
	}
}

package spec

import (
	"fmt"
	"io/ioutil"

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
		return &result, nil
	default:
		return nil, fmt.Errorf("Unsupported version '%v'", versioned.Version)
	}
}

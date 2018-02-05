package spec

import "fmt"

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
			case "string":
			case "integer":
				continue propLoop
			default:
				v.err = ErrSpec1PropInvalidType{Path: path, Key: key, Type: prop.RawTypeName}
				return
			}
		}
	}
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
			case "string":
				switch prop.RawDefault.(type) {
				case string:
					continue propLoop
				}
				break
			case "integer":
				switch prop.RawDefault.(type) {
				case int32:
				case int64:
				case uint32:
				case uint64:
				case int:
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
}

func (p *propertyV1) String() string {
	return fmt.Sprintf("<label=`%s` type=`%s` default=`%s` help=`%s` input=`%s`>", p.RawLabel, p.RawTypeName, p.RawDefault, p.RawHelp, p.RawInput)
}

func (p *propertyV1) IsValidData(data interface{}) error {
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
	return p.RawLabel
}

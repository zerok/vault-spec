package spec

type Specification interface {
	Version() string
	SecretPaths() []string
	Secret(path string) Secret
	Validate() error
}

type Labelled interface {
	Label() string
}

type Secret interface {
	PropertyNames() []string
	Property(name string) SecretProperty
	Labelled
}

type SecretProperty interface {
	Type() string
	Input() string
	Help() string
	Default() interface{}
	IsValidData(interface{}) error
	Labelled
}

type versionedEnvelop struct {
	Version string `yaml:"version"`
}

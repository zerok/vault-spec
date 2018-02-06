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

// Specification abstracts a Vault specification across multiple versions as
// best as possible.
type Specification interface {
	// Version returns the version string of the specification.
	Version() string
	// SecretPaths returns a list of the paths of all specified secrets.
	SecretPaths() []string
	// Secret returns the specification of a secret or nil.
	Secret(path string) Secret
	// Validate checks the specification itself for validity.
	Validate() error
}

// Labelled implementations provide a Label() method for display to the user.
type Labelled interface {
	// Label provides a string representation can can be directly presented to
	// the user.
	Label() string
}

// Secret is the specification of an item inside Vault's key-value store.
type Secret interface {
	// PropertyNames returns a list of the names of all specified properties.
	PropertyNames() []string
	// Property returns a property specification or nil if not present.
	Property(name string) SecretProperty
	Labelled
}

// SecretProperty is the specification of a property inside a secret.
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

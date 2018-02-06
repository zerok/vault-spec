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

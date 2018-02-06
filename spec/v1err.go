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
	"reflect"
)

// ErrSpec1PropInvalidInput is returned if the specification contains an
// unknown input.
type ErrSpec1PropInvalidInput struct {
	Path string
	Key  string
}

func (err ErrSpec1PropInvalidInput) Error() string {
	return fmt.Sprintf("invalid input defined for %s/%s", err.Path, err.Key)
}

// ErrSpec1PropInvalidType is returned if the specification refers to an
// unknown type.
type ErrSpec1PropInvalidType struct {
	Path string
	Key  string
	Type string
}

func (err ErrSpec1PropInvalidType) Error() string {
	return fmt.Sprintf("invalid type defined for %s/%s: %s", err.Path, err.Key, err.Type)
}

// ErrSpec1PropInvalidDefault is returned if the specification uses a default
// value that doesn't match the specified type.
type ErrSpec1PropInvalidDefault struct {
	Path    string
	Key     string
	Default interface{}
}

func (err ErrSpec1PropInvalidDefault) Error() string {
	return fmt.Sprintf("invalid default defined for %s/%s: `%v` (type=%s)", err.Path, err.Key, err.Default, reflect.TypeOf(err.Default))
}

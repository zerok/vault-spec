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
	"testing"
)

func TestV1Validation(t *testing.T) {
	e := envelopV1{}
	e.RawSpec.RawSecrets = map[string]secretV1{
		"secret1": secretV1{
			RawProperties: map[string]propertyV1{
				"intProp": propertyV1{
					RawTypeName: "integer",
					RawDefault:  "0-with-string",
				},
			},
		},
	}
	tests := []struct {
		rawDefault interface{}
		valid      bool
	}{
		{0, true},
		{"0-with-string", false},
		{int64(0), true},
	}

	for _, test := range tests {
		secretSpec := e.RawSpec.RawSecrets["secret1"]
		propSpec := secretSpec.RawProperties["intProp"]
		propSpec.RawDefault = test.rawDefault
		secretSpec.RawProperties["intProp"] = propSpec
		e.RawSpec.RawSecrets["secret1"] = secretSpec

		err := e.Validate()
		if err == nil && !test.valid {
			t.Fatal("Setting %v (type=%s) as default value for an integer property should have produced an error", test.rawDefault, reflect.TypeOf(test.rawDefault))
		}
		if err != nil && test.valid {
			t.Fatalf("%v (type=%s) as default value for an integer property should have been valid. Got %s instead.", test.rawDefault, reflect.TypeOf(test.rawDefault), err)
		}
	}
}

type dataValidationTestCase struct {
	spec  propertyV1
	data  interface{}
	valid bool
}

func (d dataValidationTestCase) String() string {
	return fmt.Sprintf("<spec=%s data=%v valid=%v>", d.spec, d.data, d.valid)
}

func TestV1DataValidation(t *testing.T) {
	tests := []dataValidationTestCase{
		{
			spec: propertyV1{
				RawTypeName: "string",
			},
			data:  1,
			valid: false,
		}, {
			spec: propertyV1{
				RawTypeName: "string",
			},
			data:  "",
			valid: true,
		}, {
			spec: propertyV1{
				RawTypeName: "string",
			},
			data:  "something",
			valid: true,
		}, {
			spec: propertyV1{
				RawTypeName: "integer",
			},
			data:  "",
			valid: false,
		}, {
			spec: propertyV1{
				RawTypeName: "integer",
			},
			data:  "something",
			valid: false,
		}, {
			spec: propertyV1{
				RawTypeName: "integer",
			},
			data:  1,
			valid: true,
		}, {
			spec: propertyV1{
				RawTypeName: "integer",
			},
			data:  1.2,
			valid: false,
		}, {
			spec: propertyV1{
				RawTypeName: "integer",
			},
			data:  -1,
			valid: true,
		},
	}

	for _, test := range tests {
		err := test.spec.IsValidData(test.data)
		if err == nil && !test.valid {
			t.Errorf("No error received for %s", test)
		} else if err != nil && test.valid {
			t.Errorf("Unexpected error received for %s: %s", test, err)
		}
	}
}

package spec

import (
	"fmt"
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
	err := e.Validate()
	if err == nil {
		t.Fatal("Setting a string as default value for an integer property should have produced an error")
	}
	e.RawSpec.RawSecrets["secret1"].RawProperties["intProp"] = propertyV1{
		RawTypeName: "integer",
		RawDefault:  0,
	}
	err = e.Validate()
	if err != nil {
		t.Fatalf("\"0\" as default value for an integer property should have been valid. Got %s instead.", err)
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

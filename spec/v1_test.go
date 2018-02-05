package spec

import "testing"

func TestV1Validation(t *testing.T) {
	e := envelopV1{}
	e.RawSpec.RawSecrets = map[string]secretV1{
		"secret1": secretV1{
			RawProperties: map[string]propertyV1{
				"stringProp": propertyV1{
					RawTypeName: "string",
					RawDefault:  0,
				},
			},
		},
	}
	err := e.Validate()
	if err == nil {
		t.Fatal("Setting 0 as default value for a string property should have produced an error")
	}
	e.RawSpec.RawSecrets["secret1"].RawProperties["stringProp"] = propertyV1{
		RawTypeName: "string",
		RawDefault:  "stringvalue",
	}
	err = e.Validate()
	if err != nil {
		t.Fatalf("\"stringvalue\" as default value for a string property should have been valid. Got %s instead.", err)
	}
}

func TestV1DataValidation(t *testing.T) {
	tests := []struct {
		spec  propertyV1
		data  interface{}
		valid bool
	}{
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

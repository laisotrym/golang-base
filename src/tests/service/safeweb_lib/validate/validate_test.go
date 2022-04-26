package validate_test

import (
	"fmt"
	"testing"
)

type Test1 struct {
	Name   string `validate:"required,lte=50"`
	Amount int64  `validate:"required"`
	Key    string `validate:"lte=3"`
}

func TestValidateVariable(t *testing.T) {
	tests := []struct {
		input  interface{}
		rule   string
		result bool
	}{
		{"order-code", "code", true},
		{"order3#G#G-code", "code", false},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			ok, _ := safeweb_lib_validate.Variable(test.input, test.rule)
			if ok != test.result {
				t.Errorf("expect: %v got: %v", test.result, ok)
			}
		})
	}
}

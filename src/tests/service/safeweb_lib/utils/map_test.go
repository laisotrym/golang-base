package utils_test

import (
	"fmt"
	"testing"

	safeweb_lib_utils "safeweb.app/service/safeweb_lib/utils"
)

func TestBuildMapEncode(t *testing.T) {
	input := map[string]string{
		"name":        "SafeWEb",
		"id":          "123434",
		"phoneNumber": "0981222222",
		"url":         "https://admin.safeweb.app",
		"description": "",
	}
	test := []struct {
		input  map[string]string
		expect string
	}{
		{input: input, expect: "description=&id=123434&name=SafeWeb&phoneNumber=0981222222&url=https://admin.safeweb.app"},
	}
	for i, test := range test {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := safeweb_lib_utils.MapToQueryParam(test.input)
			if result != test.expect {
				t.Errorf("want %v; got %v", test.expect, result)
			}
		})
	}
}

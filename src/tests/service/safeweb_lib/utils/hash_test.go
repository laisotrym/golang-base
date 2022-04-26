package utils_test

import (
	"fmt"
	"testing"

	safeweb_lib_utils "safeweb.app/service/safeweb_lib/utils"
)

func TestSHA256(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"abc", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"SafeWeb@12434", "678d1e13b2577df2f1a9834709e1a45e22b007b253974856ca9db0507c7cec19"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := safeweb_lib_utils.HashSHA256(tt.in)
			if result != tt.out {
				t.Errorf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out string
	}{
		{"abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := safeweb_lib_utils.HashMD5(tt.in)
			if result != tt.out {
				t.Errorf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestHash(t *testing.T) {
	for i, tt := range []struct {
		hashType string
		input    string
		expect   string
	}{
		{"SHA256", "abc", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{"MD5", "abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"", "abc", ""},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := safeweb_lib_utils.Hash(tt.input, tt.hashType)
			if result != tt.expect {
				t.Errorf("want %v; got %v", tt.expect, result)
			}
		})
	}
}

package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	safeweb_lib_utils "safeweb.app/service/safeweb_lib/utils"
)

func TestSliceHelper_FindFirstDuplicateString(t *testing.T) {
	// Arrange
	testCases := []struct {
		name string
		arr  []string
		exp  struct {
			duplicated bool
			dupVal     string
		}
	}{
		{
			name: "Empty array",
			arr:  make([]string, 0),
			exp: struct {
				duplicated bool
				dupVal     string
			}{
				duplicated: false,
				dupVal:     "",
			},
		},
		{
			name: "No duplicated array",
			arr:  []string{"a", "abc", "abcd", "test", "thanh"},
			exp: struct {
				duplicated bool
				dupVal     string
			}{
				duplicated: false,
				dupVal:     "",
			},
		},
		{
			name: "Duplicated empty item in array",
			arr:  []string{"a", "abc", "", "test", ""},
			exp: struct {
				duplicated bool
				dupVal     string
			}{
				duplicated: true,
				dupVal:     "",
			},
		},
		{
			name: "Non empty duplicated item in array",
			arr:  []string{"x", "y", "a", "x", "abc", "", "test", ""},
			exp: struct {
				duplicated bool
				dupVal     string
			}{
				duplicated: true,
				dupVal:     "x",
			},
		},
		{
			name: "Multiple duplicated items in array",
			arr:  []string{"ps", "y", "a", "x", "abc", "ps", "test", "b", "b", "safeweb", "safeweb"},
			exp: struct {
				duplicated bool
				dupVal     string
			}{
				duplicated: true,
				dupVal:     "ps",
			},
		},
	}

	// Act & Assert
	for _, tc := range testCases {
		dv, d := safeweb_lib_utils.FindFirstDuplicateString(tc.arr)
		assert.Equal(t, tc.exp.duplicated, d)
		assert.Equal(t, tc.exp.dupVal, dv)
	}
}

func TestSliceUtils_FindMissingString(t *testing.T) {
	// Arrange
	testCases := []struct {
		name string
		arr  []string
		keys []string
		exp  []string
	}{
		{
			name: "Empty array & keys",
			arr:  make([]string, 0),
			keys: make([]string, 0),
			exp:  make([]string, 0),
		},
		{
			name: "Empty keys",
			arr:  []string{"a", "b", "c"},
			keys: make([]string, 0),
			exp:  []string{"a", "b", "c"},
		},
		{
			name: "Miss all",
			arr:  []string{"a", "b", "c"},
			keys: make([]string, 0),
			exp:  []string{"a", "b", "c"},
		},
		{
			name: "Miss several",
			arr:  []string{"a", "b", "c"},
			keys: []string{"a"},
			exp:  []string{"b", "c"},
		},
	}

	// Act & Assert
	for _, tc := range testCases {
		miss := safeweb_lib_utils.FindMissingString(tc.arr, tc.keys...)
		assert.Equal(t, tc.exp, miss)
	}
}

func TestSliceUtils_FindOutsideKeys(t *testing.T) {
	testCases := []struct {
		name     string
		source   []string
		keys     []string
		expected []string
	}{
		{
			name:     "No diff",
			source:   []string{"a", "b"},
			keys:     []string{"a", "a", "b"},
			expected: make([]string, 0),
		},
		{
			name:     "Has diff",
			source:   []string{"a", "b", "a"},
			keys:     []string{"a", "c"},
			expected: []string{"c"},
		},
		{
			name:     "Has diff complex",
			source:   []string{"aa", "bb", "cc", "dd"},
			keys:     []string{"ff", "aa", "cc", "ff"},
			expected: []string{"ff", "ff"},
		},
	}

	for _, tc := range testCases {
		fmt.Println(tc.name)
		diffs := safeweb_lib_utils.FindOutsideKeys(tc.source, tc.keys...)
		assert.Equal(t, tc.expected, diffs)
	}
}

func TestRemoveDuplicatedInt64(t *testing.T) {
	testCases := []struct {
		input  []int64
		output []int64
	}{
		{
			input:  []int64{},
			output: []int64{},
		},
		{
			input:  []int64{1, 2},
			output: []int64{1, 2},
		},
		{
			input:  []int64{1, 2, 1, 1, 1},
			output: []int64{1, 2},
		},
		{
			input:  []int64{4, 3, 5, 5, 4, 2, 3, 4},
			output: []int64{4, 3, 5, 2},
		},
	}

	for _, tc := range testCases {
		actual := safeweb_lib_utils.RemoveDuplicatedInt64(tc.input)

		assert.Equal(t, tc.output, actual)
	}
}

func TestGetComplementSliceInt64(t *testing.T) {
	testCases := []struct {
		arr1   []int64
		arr2   []int64
		output []int64
	}{
		{
			arr1:   []int64{},
			arr2:   []int64{},
			output: []int64{},
		},
		{
			arr1:   []int64{2, 3, 1},
			arr2:   []int64{},
			output: []int64{2, 3, 1},
		},
		{
			arr1:   []int64{2, 3, 1},
			arr2:   []int64{69},
			output: []int64{2, 3, 1},
		},
		{
			arr1:   []int64{1, 2, 3, 4, 5},
			arr2:   []int64{3, 5},
			output: []int64{1, 2, 4},
		},
		{
			arr1:   []int64{1, 2, 3, 4, 5},
			arr2:   []int64{2, 5, 10, 20},
			output: []int64{1, 3, 4},
		},
	}

	for _, tc := range testCases {
		actual := safeweb_lib_utils.GetComplementSliceInt64(tc.arr1, tc.arr2)

		assert.Equal(t, tc.output, actual)
	}
}

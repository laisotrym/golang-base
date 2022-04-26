package utils_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"

    "safeweb.app/service/safeweb_lib/utils"
)

type Child struct {
    StringField string
    IntField    int
}

type TrimTest struct {
    StringField string
    IntField int
    Child    Child
}

func TestTrim_String(t *testing.T) {
    s := " abc  "
    actual := safeweb_lib_utils.Trim(s).(string)
    assert.Equal(t, actual, "abc")
}

func TestTrim_NotString(t *testing.T) {
    number := 1
    actual := safeweb_lib_utils.Trim(number).(int)
    assert.Equal(t, actual, 1)
}

func TestTrim_Struct_noNest(t *testing.T) {
    child := Child{
        StringField: "  abc",
        IntField:    1,
    }
    actual := safeweb_lib_utils.Trim(child).(Child)
    assert.Equal(t, actual.StringField, "abc")
    assert.Equal(t, actual.IntField, 1)
}

func TestTrim_Struct_hasNest(t *testing.T) {
    child := Child{
        StringField: "  abc",
        IntField:    1,
    }
    trimObj := TrimTest{
        StringField: "  abc",
        IntField:    1,
        Child:       child,
    }
    actual := safeweb_lib_utils.Trim(trimObj).(TrimTest)
    assert.Equal(t, actual.StringField, "abc")
    assert.Equal(t, actual.IntField, 1)
    assert.Equal(t, actual.Child.StringField, "abc")
    assert.Equal(t, actual.Child.IntField, 1)
}
func TestTrim_Pointer_noNest(t *testing.T) {
    child := &Child{
        StringField: "  abc",
        IntField:    1,
    }
    actual := safeweb_lib_utils.Trim(child).(*Child)
    assert.Equal(t, actual.StringField, "abc")
    assert.Equal(t, actual.IntField, 1)
}

func TestTrim_Pointer_hasNest(t *testing.T) {
    child := Child{
        StringField: "  abc",
        IntField:    1,
    }
    trimObj := &TrimTest{
        StringField: "  abc",
        IntField:    1,
        Child:       child,
    }
    actual := safeweb_lib_utils.Trim(trimObj).(*TrimTest)
    assert.Equal(t, actual.StringField, "abc")
    assert.Equal(t, actual.IntField, 1)
    assert.Equal(t, actual.Child.StringField, "abc")
    assert.Equal(t, actual.Child.IntField, 1)
}

func TestTrim_Slice(t *testing.T) {
    child := Child{
        StringField: "  abc",
        IntField:    1,
    }
    trimObj := TrimTest{
        StringField: "  abc",
        IntField:    1,
        Child:       child,
    }
    actual := safeweb_lib_utils.Trim([]TrimTest{trimObj}).([]TrimTest)
    for _, v := range actual {
        assert.Equal(t, v.StringField, "abc")
        assert.Equal(t, v.IntField, 1)
        assert.Equal(t, v.Child.StringField, "abc")
        assert.Equal(t, v.Child.IntField, 1)
    }
}

func TestTrim_Map(t *testing.T) {
    child := Child{
        StringField: "  abc ",
        IntField:    1,
    }
    trimObj := &TrimTest{
        StringField: "  abc  ",
        IntField:    1,
        Child:       child,
    }
    m := map[string]*TrimTest{"key": trimObj}
    actual := safeweb_lib_utils.Trim(m).(map[string]*TrimTest)
    for _, v := range actual {
        assert.Equal(t, v.StringField, "abc")
        assert.Equal(t, v.IntField, 1)
        assert.Equal(t, v.Child.StringField, "abc")
        assert.Equal(t, v.Child.IntField, 1)
    }
}

func TestToSnakeCase(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"", ""},
        {"already_snake", "already_snake"},
        {"A", "a"},
        {"AA", "aa"},
        {"AaAa", "aa_aa"},
        {"HTTPRequest", "http_request"},
        {"BatteryLifeValue", "battery_life_value"},
        {"Id0Value", "id0_value"},
        {"ID0Value", "id0_value"},
    }
    for _, test := range tests {
        have := safeweb_lib_utils.StringToSnakeCase(test.input)
        if have != test.want {
            t.Errorf("input=%q:\nhave: %q\nwant: %q", test.input, have, test.want)
        }
    }
}

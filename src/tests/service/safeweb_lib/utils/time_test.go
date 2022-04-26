package utils_test

import (
    "fmt"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"

    "safeweb.app/service/safeweb_lib/utils"
)

func TestToVnTime(t *testing.T) {
    myTime := time.Date(2020, 04, 22, 15, 00, 00, 00, time.UTC)
    expected := "20200422220000"
    tests := []struct {
        input    time.Time
        expected string
    }{
        {myTime, expected},
    }
    for i, test := range tests {
        t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
            time, err := safeweb_lib_utils.TimeVNYmdHis(myTime)
            assert.NoError(t, err)
            if time != expected {
                t.Errorf("expect: %v got %v", test.expected, test.input)
            }
        })
    }
}

func BenchmarkTimeVNYmdHis(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        safeweb_lib_utils.TimeVNYmdHis(time.Now().UTC())
    }
}

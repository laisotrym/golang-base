package utils_test

import (
    "fmt"
    "testing"

    "safeweb.app/service/safeweb_lib/utils"
)

func TestNew(t *testing.T) {
    const total int = 10
    for i := 0; i < total; i++ {
        t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
            u1 := safeweb_lib_utils.NewUUID()
            u2 := safeweb_lib_utils.NewUUID()
            if u1 == u2 {
                t.Errorf("uuid same: %s", u1)
            }
        })
    }
}

package promo_error_test

import (
    "testing"
    
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

func Test_FindWithCode(t *testing.T) {
    t.Run("found", func(t *testing.T) {
        expect := &promo_error.Error{
            RpcStatus: 13,
            Code:      "99",
            Success:   false,
            Message:   "Lỗi hệ thống",
        }
        actual := promo_error.FindGeneralWithCode("99")
        assert.Equal(t, expect, actual)
    })
    t.Run("nil", func(t *testing.T) {
        actual := promo_error.FindGeneralWithCode("abcdef")
        assert.Nil(t, actual)
    })
}

func Test_ErrorString(t *testing.T) {
    status := &promo_error.Error{
        RpcStatus: 4,
        Code:      "100",
        Message:   "There are some error",
        Success:   false,
    }
    expected := "code: 100, message: There are some error, success: false, rpc code: 4"
    t.Run("1", func(t *testing.T) {
        assert.Equal(t, expected, status.Error())
    })
}

func Test_MissingRequiredParam(t *testing.T) {
    expect := &promo_error.Error{
        RpcStatus: 3,
        Code:      "02",
        Message:   "test",
        Success:   false,
    }
    actual, _ := promo_error.GetGeneralStatus()
    assert.Equal(t, expect, actual.MissingRequiredParam.ReplaceMessage("test"))
}

func Test_DataInvalid(t *testing.T) {
    expect := &promo_error.Error{
        RpcStatus: 3,
        Code:      "04",
        Message:   "test",
        Success:   false,
    }
    actual, _ := promo_error.GetGeneralStatus()
    assert.Equal(t, expect, actual.DataInvalid.ReplaceMessage("test"))
}

func Benchmark_UnWrapError(b *testing.B) {
    stt := &promo_error.Error{
        RpcStatus: 0,
        Code:      "01",
        Message:   "message",
        Success:   false,
    }
    stt1 := errors.Wrap(stt, "wrap1")
    stt2 := errors.Wrap(stt1, "wrap1")
    stt3 := errors.Wrap(stt2, "wrap1")
    for i := 0; i < b.N; i++ {
        _ = promo_error.UnWrapError(stt3)
    }
}

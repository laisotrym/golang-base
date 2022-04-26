package promo_error_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    gStatus "google.golang.org/grpc/status"
    "rpc.tekoapis.com/rpc/payment"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

func Test_IsError(t *testing.T) {
    stt := &promo_error.Error{
        RpcStatus: 4,
        Code:      "100",
        Message:   "There are some error",
        Success:   false,
    }
    t.Run("1", func(t *testing.T) {
        if promo_error.IsError(stt) == false {
            t.Error("expect true get false")
        }
    })
}

func TestErrorPromo_ToRPCError(t *testing.T) {
    stt := &promo_error.Error{
        RpcStatus: 3,
        Code:      "04",
        Message:   "test",
        Success:   false,
    }
    st := gStatus.New(3, "test")
    br := &payment.PaymentError{}
    code := &payment.PaymentErrorMetadata{
        Key:   "code",
        Value: "04",
    }
    br.Metadatas = append(br.Metadatas, code)
    gs, err := st.WithDetails(br)
    if err != nil {
        t.Error("error")
    }
    expect := gs.Err()
    if err != nil {
        t.Error("error add detail to status")
    }
    actual := stt.ToRPCError()
    assert.Equal(t, expect.Error(), actual.Error())
    
}

func TestErrorPromo_ToRPCErrorWithTrace(t *testing.T) {
    traceID := "EGEGEG-GRGRG"
    stt := &promo_error.Error{
        RpcStatus: 3,
        Code:      "04",
        Message:   "test",
        Success:   false,
    }
    st := gStatus.New(3, "test")
    br := &payment.PaymentError{}
    code := &payment.PaymentErrorMetadata{
        Key:   "code",
        Value: "04",
    }
    br.Metadatas = append(br.Metadatas, code)
    trace := &payment.PaymentErrorMetadata{
        Key:   "trace_id",
        Value: traceID,
    }
    
    br.Metadatas = append(br.Metadatas, trace)
    gs, err := st.WithDetails(br)
    if err != nil {
        t.Error("error")
    }
    expect := gs.Err()
    if err != nil {
        t.Error("error add detail to status")
    }
    actual := stt.ToRPCErrorWithTraceID(traceID)
    assert.Equal(t, expect.Error(), actual.Error())
    
}

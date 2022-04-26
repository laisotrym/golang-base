package status_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    gStatus "google.golang.org/grpc/status"
    "rpc.tekoapis.com/rpc/payment"
    
    promo_status "safeweb.app/service/safeweb_lib/status"
)

func Test_IsStatus(t *testing.T) {
    stt := promo_status.DefaultSuccessTemplate.With("something")
    t.Run("1", func(t *testing.T) {
        if promo_status.IsTrace(stt) == false {
            t.Error("expect true get false")
        }
    })
}

func TestStatus_ToRPCError(t *testing.T) {
    stt := promo_status.DefaultInvalidFormatTemplate.With("test", "abcxyz")
    st := gStatus.New(3, "test sai định dạng, định dạng đúng là abcxyz")
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

func TestStatus_ToRPCErrorWithTrace(t *testing.T) {
    traceID := "EGEGEG-GRGRG"
    stt := promo_status.DefaultInvalidFormatTemplate.With("test", "abcxyz")
    st := gStatus.New(3, "test sai định dạng, định dạng đúng là abcxyz")
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
    actual := stt.WithTraceID(traceID).ToRPCError()
    assert.Equal(t, expect.Error(), actual.Error())
}

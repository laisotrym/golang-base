package safeweb_lib_error

import (
    "errors"
    "fmt"
    "reflect"
    
    "github.com/gogo/protobuf/types"
    "github.com/gogo/status"
    "google.golang.org/grpc/codes"
    
    "safeweb.app/rpc/safeweb_lib"
)

type (
    Error struct {
        RpcStatus int32  `yaml:"rpcStatus"`
        Code      string `yaml:"code"`
        Message   string `yaml:"message"`
        Success   bool   `yaml:"success"`
    }
)

func NewError(rpcStatus int32, code string, msg string, success bool) *Error {
    return &Error{
        RpcStatus: rpcStatus,
        Code:      code,
        Message:   msg,
        Success:   success,
    }
}

func (s *Error) Error() string {
    return fmt.Sprintf("code: %v, message: %v, success: %v, rpc code: %v", s.Code, s.Message, s.Success, s.RpcStatus)
}

func (s *Error) ToRPCErrorWithTraceID(traceID string) error {
    st := status.New(codes.Code(s.RpcStatus), s.Message)
    detail := &safeweb_lib.Error{}
    if s.Code != "" {
        code := &safeweb_lib.ErrorMetadata{
            Key:   "error_code",
            Value: s.Code,
        }
        detail.Metadatas = append(detail.Metadatas, code)
    }
    trace := &safeweb_lib.ErrorMetadata{
        Key:   "trace_id",
        Value: traceID,
    }
    detail.Metadatas = append(detail.Metadatas, trace)
    st, err := st.WithDetails(detail)
    if err != nil {
        return errors.New(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
    }
    return st.Err()
}

func (s *Error) ToRPCError() error {
    st := status.New(codes.Code(s.RpcStatus), s.Message)
    detail := &safeweb_lib.Error{}
    code := &safeweb_lib.ErrorMetadata{
        Key:   "error_code",
        Value: s.Code,
    }
    detail.Metadatas = append(detail.Metadatas, code)
    st, err := st.WithDetails(detail)
    if err != nil {
        return errors.New(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
    }
    return st.Err()
}

func (s Error) F(v ...interface{}) *Error {
    return &Error{
        RpcStatus: s.RpcStatus,
        Code:      s.Code,
        Message:   fmt.Sprintf(s.Message, v...),
        Success:   s.Success,
    }
}

func (s Error) ReplaceMessage(m string) *Error {
    return &Error{
        RpcStatus: s.RpcStatus,
        Code:      s.Code,
        Message:   m,
        Success:   s.Success,
    }
}

func IsError(err error) bool {
    e := &Error{}
    return errors.As(err, &e)
}

func UnWrapError(err error) *Error {
    stt := &Error{}
    for {
        if reflect.TypeOf(err) != reflect.TypeOf(stt) {
            err = errors.Unwrap(err)
        } else {
            break
        }
    }
    return err.(*Error)
}

func CastToError(s *status.Status) (promoError *Error) {
    for _, v := range s.Proto().GetDetails() {
        var detail safeweb_lib.Error
        if types.Is(v, &detail) {
            if err := types.UnmarshalAny(v, &detail); err == nil {
                for _, it := range detail.GetMetadatas() {
                    if it.Key == "error_code" {
                        promoError = &Error{
                            RpcStatus: int32(s.Code()),
                            Code:      it.Value,
                            Message:   s.Message(),
                            Success:   false,
                        }
                        break
                    }
                }
                break
            }
        }
    }
    
    return
}

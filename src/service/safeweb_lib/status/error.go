package safeweb_lib_status

import (
    "errors"
    "fmt"
    "reflect"
    
    "github.com/gogo/status"
    "google.golang.org/grpc/codes"
    
    "safeweb.app/rpc/safeweb_lib"
)

type Trace struct {
    Code    Code   `json:"code"`
    Message string `json:"message"`
    TraceID string `json:"traceId,omitempty"`
}

// Implements error interface
func (p *Trace) Error() string {
    return fmt.Sprintf("code: %v, message: %v, traceID: %v", p.Code, p.Message, p.TraceID)
}

func (p *Trace) WithTraceID(traceID string) *Trace {
    p.TraceID = traceID
    return p
}

func (p *Trace) ToRPCError() error {
    rpcCode, ok := codeToRpcCode[p.Code]
    if !ok {
        return errors.New("perror: codeToRpcCode does not contain " + p.Code.String())
    }
    
    detail := &safeweb_lib.Error{}
    
    codeMetadata := &safeweb_lib.ErrorMetadata{
        Key:   "error_code",
        Value: p.Code.String(),
    }
    detail.Metadatas = append(detail.Metadatas, codeMetadata)
    
    if p.TraceID != "" {
        traceMetadata := &safeweb_lib.ErrorMetadata{
            Key:   "trace_id",
            Value: p.TraceID,
        }
        detail.Metadatas = append(detail.Metadatas, traceMetadata)
    }
    
    rpcStatus, err := status.New(rpcCode, p.Message).WithDetails(detail)
    if err != nil {
        return errors.New(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
    }
    
    return rpcStatus.Err()
}

func (p *Trace) RPCCode() codes.Code {
    return codeToRpcCode[p.Code]
}

func IsTrace(err error) bool {
    e := &Trace{}
    return errors.As(err, &e)
}

func UnwrapTrace(err error) *Trace {
    stt := &Trace{}
    for {
        if err != nil && reflect.TypeOf(err) != reflect.TypeOf(stt) {
            err = errors.Unwrap(err)
        } else {
            break
        }
    }
    
    ps, ok := err.(*Trace)
    if !ok {
        return nil
    }
    return ps
}

func NewErrorTrace(c Code, message string) *Trace {
    return &Trace{Code: c, Message: message}
}

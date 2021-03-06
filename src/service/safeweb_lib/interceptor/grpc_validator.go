package safeweb_lib_interceptor

import (
    "context"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    
    promo_error "safeweb.app/service/safeweb_lib/error"
)

type validator interface {
    Validate() error
}

// UnaryValidatorServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func UnaryValidatorServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if v, ok := req.(validator); ok {
            if err := v.Validate(); err != nil {
                return nil, promo_error.NewError(int32(codes.InvalidArgument), "", err.Error(), false)
            }
        }
        return handler(ctx, req)
    }
}

// StreamValidatorServerInterceptor returns a new streaming server interceptor that validates incoming messages.
//
// The stage at which invalid messages will be rejected with `InvalidArgument` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any userspace
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamValidatorServerInterceptor() grpc.StreamServerInterceptor {
    return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        wrapper := &recvWrapper{stream}
        return handler(srv, wrapper)
    }
}

type recvWrapper struct {
    grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
    if err := s.ServerStream.RecvMsg(m); err != nil {
        return err
    }
    if v, ok := m.(validator); ok {
        if err := v.Validate(); err != nil {
            return promo_error.NewError(int32(codes.InvalidArgument), "", err.Error(), false)
        }
    }
    return nil
}

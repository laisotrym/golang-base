package safeweb_lib_interceptor

import (
	"context"
	"reflect"

	"github.com/gogo/status"
	"github.com/golang/protobuf/proto"
	"safeweb.app/log/level"
	"google.golang.org/grpc"

	"safeweb.app/rpc/safeweb_lib"
	promo_error "safeweb.app/service/safeweb_lib/error"
	safeweb_lib_utils "safeweb.app/service/safeweb_lib/utils"
	promo_status "safeweb.app/service/safeweb_lib/status"
	"safeweb.app/constant"
)

func UnaryErrorTraceInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		traceID := safeweb_lib_utils.NewUUID()
		ctx = context.WithValue(ctx, safeweb_lib_constant.RequestTraceContextKey, traceID)
		resp, err := handler(ctx, req)
		if err != nil {
			// Keeps compatible (remove later)
			if promo_error.IsError(err) {
				level.Error(ctx).F("Trace %#v, err: %#v", traceID, err)
				stt := promo_error.UnWrapError(err)
				return nil, stt.ToRPCErrorWithTraceID(traceID)
			}

			// NewErrorTrace error implementation
			if promo_status.IsTrace(err) {
				stt := promo_status.UnwrapTrace(err)
				return nil, stt.WithTraceID(traceID).ToRPCError()
			}
			s, ok := status.FromError(err)
			if ok {
				promoError := promo_error.CastToError(s)
				if promoError != nil {
					return nil, promoError.ToRPCErrorWithTraceID(traceID)
				}
			}

			// Internal server error
			level.Error(ctx).F("Trace %#v, err: %#v", traceID, err)
			return nil, promo_status.DefaultInternalErrorTemplate.With().WithTraceID(traceID).ToRPCError()
		}

		// wrapper := wrapperRespProcess(resp)

		return resp, nil
	}
}

// type wrapperResp struct {
//     Code       string       `json:"code,omitempty"`
//     Message    string       `json:"message,omitempty"`
//     Result     interface{}  `json:"result,omitempty"`
//     TypeResult reflect.Type `json:"-"`
// }

// func (m *wrapperResp) Reset()         { *m = wrapperResp{} }
// func (m *wrapperResp) String() string { return proto.CompactTextString(m) }
// func (*wrapperResp) ProtoMessage()    {}

func wrapperRespProcess(ukn interface{}) interface{} {
	respStruct := reflect.Indirect(reflect.ValueOf(ukn)).Interface()
	respStructType := reflect.TypeOf(respStruct)
	if _, ok := respStructType.FieldByName("Message"); ok {
		return ukn
	}
	if _, ok := respStructType.FieldByName("Code"); ok {
		return ukn
	}
	// respStructName := respStructType.String()
	// fmt.Printf("Object Name: %s", respStructName)

	value, err := ukn.(proto.Marshaler).Marshal()
	if err != nil {
		return ukn
	}
	return &safeweb_lib.WrapperResp{
		Code:    "00",
		Message: "Thành công",
		Result:  value,
	}

	// return &wrapperResp{
	//     Code:       "00",
	//     Message:    "Thành công",
	//     Result:     respStruct,
	//     TypeResult: respStructType,
	// }
}

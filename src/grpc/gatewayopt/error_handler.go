package gatewayopt

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	errorpb "safeweb.app/grpc/gatewayopt/internal/error"
)

// HTTPError convert GRPC error response to HTTP error response to conform standardized response format.
// To use this error handler in grpc-gateway v1, write this line in your code `runtime.GlobalHTTPErrorHandler = HTTPError`. This needs to be set early before the handler code.
func HTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	contentType := marshaler.ContentType()
	w.Header().Set("Content-Type", contentType)

	st := runtime.HTTPStatusFromCode(s.Code())

	var localized bool
	for _, detail := range s.Details() {
		if _, ok := detail.(*errdetails.LocalizedMessage); ok {
			localized = true
			break
		}
	}
	if !localized {
		msg := defaultLocalizedMessage(st, "vi-VN")
		s, _ = s.WithDetails(msg)
	}

	body := &errorpb.Error{
		Code:    int32(st),
		Message: s.Message(),
		Details: s.Proto().Details,
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

}

const localeVN = "vi-VN"
const localeEN = "en-US"
const defaultLocale = "vi-VN"

var fallbackLocalizedMessage = map[string]map[int]string{
	localeEN: {
		http.StatusBadRequest:          "Client Bad Request",
		http.StatusNotFound:            "Resource Not Found",
		http.StatusInternalServerError: "Something went wrong...",
	},
	localeVN: {
		http.StatusBadRequest:          "Lỗi client",
		http.StatusNotFound:            "Không thấy mục bạn tìm kiếm",
		http.StatusInternalServerError: "Có vấn đề xảy ra trong hệ thống của chúng tôi. Xin lỗi vì sự bất tiện...",
	},
}

func defaultLocalizedMessage(st int, locale string) *errdetails.LocalizedMessage {
	m, ok := fallbackLocalizedMessage[locale]
	if !ok {
		m = fallbackLocalizedMessage[defaultLocale]
	}
	if text, ok := m[st]; ok {
		return &errdetails.LocalizedMessage{
			Locale:  locale,
			Message: text,
		}
	}
	return &errdetails.LocalizedMessage{
		Locale:  locale,
		Message: "Lỗi không xác định...",
	}
}

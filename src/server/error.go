package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"

	"github.com/gogo/status"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"

	"safeweb.app/rpc/safeweb_lib"
	promo_error "safeweb.app/service/safeweb_lib/error"
)

type CustomContentTypeMarshaler interface {
	// ContentTypeFromMessage returns the Content-Type this marshaler produces from the provided message
	ContentTypeFromMessage(v interface{}) string
}

type ErrorBody struct {
	Code    string `json:"resCode,omitempty"`
	Message string `json:"resDesc,omitempty"`
}

func CustomHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	w.Header().Del("Trailer")

	contentType := marshaler.ContentType()
	// Check marshaler on run time in order to keep backwards compatibility
	// An interface param needs to be added to the ContentType() function on
	// the Marshal interface to be able to remove this check
	pb := s.Proto()
	if typeMarshaler, ok := marshaler.(CustomContentTypeMarshaler); ok {
		contentType = typeMarshaler.ContentTypeFromMessage(pb)
	}
	w.Header().Set("Content-Type", contentType)

	// var body interface{} = &safeweb_lib.UnaryError{
	//     Error:   s.Message(),
	//     Message: s.Message(),
	//     Code:    int32(s.Code()),
	//     Details: s.Proto().GetDetails(),
	// }
	var body interface{}
	promoError := promo_error.CastToError(s)
	if promoError != nil {
		body = &safeweb_lib.ErrorResponse{
			Code:    promoError.Code,
			Message: promoError.Message,
		}
	} else {
		body = &safeweb_lib.ErrorResponse{
			Code:    fmt.Sprintf("%02d", int32(s.Code())),
			Message: s.Message(),
		}
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

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, mux, md)
	handleForwardResponseTrailerHeader(w, md)
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	handleForwardResponseTrailer(w, md)
}

func outgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		if h, ok := outgoingHeaderMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

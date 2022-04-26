package safeweb_lib_interceptor

import (
	"context"
	"log"

	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"safeweb.app/service/safeweb_admin/auth"
)

// Auth is a server interceptor for authentication and authorization
type Auth struct {
	jwtManager      *auth.JWTManager
	accessibleRoles map[string][]string
}

// NewInterceptor returns a new auth interceptor
func NewInterceptor(jwtManager *auth.JWTManager, accessibleRoles map[string][]string) *Auth {
	return &Auth{jwtManager, accessibleRoles}
}

// UnaryAuthInterceptor returns a server interceptor function to authenticate and authorize unary RPC
func (i *Auth) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> safeweb_lib unary auth interceptor: ", info.FullMethod)

		err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamAuthInterceptor returns a server interceptor function to authenticate and authorize stream RPC
func (i *Auth) StreamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> safeweb_lib stream auth interceptor: ", info.FullMethod)

		err := i.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (i *Auth) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := i.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

package optimus

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type User struct {
	Username string
	Project  string
}

func getAuthUserFromContext(ctx context.Context) User {
	v := ctx.Value("authorized-user")
	if v != nil {
		return v.(User)
	}
	return User{}
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	allowed_endpoints := map[string]bool{
	// "/OptimusOptimus/CreatePoint":         true,
	}
	if allow, ok := allowed_endpoints[fullMethodName]; allow && ok {
		return ctx, nil
	}

	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Error processing client certificate")
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	if len(tlsInfo.State.PeerCertificates) < 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "Error processing client certificate")
	}

	cert := tlsInfo.State.PeerCertificates[0]
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Error processing client certificate")
	}

	user := User{
		Username: cert.Subject.CommonName,
		Project:  cert.Subject.Organization[0],
	}
	return context.WithValue(ctx, "authorized-user", user), nil
}

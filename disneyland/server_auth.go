package disneyland

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"strings"
)

func getAuthUserFromContext(ctx context.Context) User {
	v := ctx.Value("authorized-user")
	if v != nil {
		return v.(User)
	}
	return User{}
}

func parseCertificateFields(field string) (projectAccess string, kindAccess string, err error) {
	fieldCopy := strings.Split(field, ".")
	if len(fieldCopy) != 2 {
		return "", "", grpc.Errorf(codes.DataLoss, "Error processing Organization Name")
	}
	return fieldCopy[0], fieldCopy[1], nil
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	allowedEndpoints := map[string]bool{}
	if allow, ok := allowedEndpoints[fullMethodName]; allow && ok {
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

	projectAccess, kindAccess, err := parseCertificateFields(cert.Subject.Organization[0])
	if err != nil {
		return nil, err
	}

	user := User{
		Username:      cert.Subject.CommonName,
		KindAccess:    kindAccess,
		ProjectAccess: projectAccess,
	}
	return context.WithValue(ctx, "authorized-user", user), nil
}

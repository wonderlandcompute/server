package disneyland

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"strings"
)

type User struct {
	Username       string
	Project        string
	Project_access string
	Kind_access    string
}

func getAuthUserFromContext(ctx context.Context) User {
	v := ctx.Value("authorized-user")
	if v != nil {
		return v.(User)
	}
	return User{}
}
func parseCertificateFields(field string) (project string, project_access string, kind_access string, err error) {
	fieldCopy := strings.Split(field, ".")
	if len(fieldCopy) != 3 {
		return "", "", "", grpc.Errorf(codes.DataLoss, "Error processing Organization Name")
	}
	return fieldCopy[0], fieldCopy[1], fieldCopy[2], nil

}
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	allowed_endpoints := map[string]bool{}
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

	project, project_access, kind_access, err := parseCertificateFields(cert.Subject.Organization[0])
	if err != nil {
		return nil, err
	}

	user := User{
		Username:       cert.Subject.CommonName,
		Project:        project,
		Kind_access:    kind_access,
		Project_access: project_access,
	}
	return context.WithValue(ctx, "authorized-user", user), nil
}

package disneyland

import (
	"crypto/x509"
	"encoding/asn1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
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
func checkExistance(certificate *x509.Certificate, oid asn1.ObjectIdentifier) (bool, int) {
	atv := certificate.Subject.Names
	index := 0
	for _, a := range atv {
		if a.Type.Equal(oid) {
			return true, index
		}
		index += 1
	}
	return false, 0
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
	project_type := asn1.ObjectIdentifier{1, 2, 3, 4}
	kind_type := asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6}
	projAccessExist, projAccessInd := checkExistance(cert, project_type)
	kindAccessExist, kindAccessInd := checkExistance(cert, kind_type)
	project := "default"
	kind := "default"
	if projAccessExist {
		project = cert.Subject.Names[projAccessInd].Value.(string)
	} else {
		return nil, grpc.Errorf(codes.DataLoss, "Error processing project_access field")

	}
	if kindAccessExist {
		kind = cert.Subject.Names[kindAccessInd].Value.(string)
	} else {
		return nil, grpc.Errorf(codes.DataLoss, "Error processing kind_access field")

	}

	user := User{
		Username:       cert.Subject.CommonName,
		Project:        cert.Subject.Organization[0],
		Kind_access:    kind,
		Project_access: project,
	}
	return context.WithValue(ctx, "authorized-user", user), nil
}

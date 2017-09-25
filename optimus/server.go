package optimus

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Server struct {
	Storage   *OptimusStorage
	SecretKey []byte
}

func detailedInternalError(err error) error {
	return grpc.Errorf(
		codes.Internal,
		fmt.Sprintf("Error creating point: %v", err),
	)
}

func (s *Server) Init() {
}

func (s *Server) CreatePoint(ctx context.Context, in *Point) (*Point, error) {
	user := getAuthUserFromContext(ctx)
	in.Project = user.Project

	created_point, err := s.Storage.CreatePoint(in, user)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return created_point, nil
}

func (s *Server) GetPoint(ctx context.Context, in *RequestWithId) (*Point, error) {
	user := getAuthUserFromContext(ctx)

	point, err := s.Storage.GetPoint(in.Id, user.Project)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return point, nil
}

func (s *Server) ListPoints(ctx context.Context, in *ListPointsRequest) (*ListOfPoints, error) {
	user := getAuthUserFromContext(ctx)
	ret, err := s.Storage.ListPoints(user.Project)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) ModifyPoint(ctx context.Context, in *Point) (*Point, error) {
	user := getAuthUserFromContext(ctx)
	if user.Project != in.Project {
		return nil, grpc.Errorf(codes.PermissionDenied, "point.Project ≠ user.Project")
	}

	ret, err := s.Storage.UpdatePoint(in)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return ret, nil
}

func (s *Server) PullPendingPoints(ctx context.Context, in *ListPointsRequest) (*ListOfPoints, error) {
	return nil, nil
}

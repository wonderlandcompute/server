package optimus

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"fmt"
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

	created_point, err := s.Storage.CreatePoint(in)
	if err != nil {
		return nil, detailedInternalError(err)
	}

	return created_point, nil
}

func (s *Server) GetPoint(ctx context.Context, in *RequestWithId) (*Point, error) {
	return nil, nil
}

func (s *Server) ListPoints(ctx context.Context, in *ListPointsRequest) (*ListOfPoints, error) {
	return nil, nil
}

func (s *Server) ModifyPoint(ctx context.Context, in *Point) (*Point, error) {
	return nil, nil
}
